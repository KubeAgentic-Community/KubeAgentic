#!/usr/bin/env python3
"""
KubeAgentic Agent Application

This is the core agent application that runs in Kubernetes pods.
It reads configuration from environment variables and provides 
an API endpoint for LLM interactions.
"""

import os
import json
import logging
from typing import Dict, List, Optional, Any
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import uvicorn
from datetime import datetime
import backoff
import httpx

# Import LLM providers
import openai
from anthropic import Anthropic
import google.generativeai as genai

# Configure structured logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

# Import LangGraph components (optional)
try:
    from langgraph.graph import StateGraph, END
    from langgraph.prebuilt import create_react_agent
    from langchain_openai import ChatOpenAI
    from langchain_anthropic import ChatAnthropic
    from langchain_google_genai import ChatGoogleGenerativeAI
    from langchain.schema import HumanMessage, SystemMessage
    from langchain.tools import BaseTool
    from typing import Type
    LANGGRAPH_AVAILABLE = True
    logger.info("LangGraph dependencies loaded successfully")
except ImportError as e:
    LANGGRAPH_AVAILABLE = False
    logger.warning(f"LangGraph dependencies not available: {e}")
    logger.warning("Agent will only support direct framework mode")

app = FastAPI(title="KubeAgentic Agent", version="1.0.0")

# --- Pydantic Models for API Requests and Responses ---

class ChatRequest(BaseModel):
    """Request model for the /chat endpoint."""
    message: str
    conversation_id: Optional[str] = None
    context: Optional[Dict[str, Any]] = None

class ChatResponse(BaseModel):
    """Response model for the /chat endpoint."""
    response: str
    conversation_id: str
    timestamp: datetime
    provider: str
    model: str

class HealthResponse(BaseModel):
    """Response model for the /health and /ready endpoints."""
    status: str
    provider: str
    model: str
    timestamp: datetime

# --- Agent Configuration ---

class AgentConfig:
    """Loads and holds the agent's configuration from environment variables."""
    def __init__(self):
        self.provider = os.getenv("AGENT_PROVIDER", "openai")
        self.model = os.getenv("AGENT_MODEL", "gpt-3.5-turbo")
        self.system_prompt = os.getenv("AGENT_SYSTEM_PROMPT", "You are a helpful AI assistant.")
        self.api_key = os.getenv("AGENT_API_KEY")
        self.endpoint = os.getenv("AGENT_ENDPOINT")
        self.framework = os.getenv("AGENT_FRAMEWORK", "direct")
        self.tools_count = int(os.getenv("AGENT_TOOLS_COUNT", "0"))
        
        # Load LangGraph configuration if framework is langgraph
        self.langgraph_config = None
        if self.framework == "langgraph":
            langgraph_config_str = os.getenv("AGENT_LANGGRAPH_CONFIG")
            if langgraph_config_str:
                try:
                    self.langgraph_config = json.loads(langgraph_config_str)
                    logger.info("LangGraph configuration loaded")
                except json.JSONDecodeError as e:
                    logger.error(f"Invalid LangGraph configuration JSON: {e}")
                    raise ValueError(f"Invalid LangGraph configuration: {e}")
            else:
                logger.warning("Framework set to 'langgraph' but no AGENT_LANGGRAPH_CONFIG provided")
        
        if not self.api_key:
            logger.error("AGENT_API_KEY environment variable is not set.")
            raise ValueError("AGENT_API_KEY environment variable is required")
        
        # Validate framework
        if self.framework not in ["direct", "langgraph"]:
            logger.error(f"Invalid framework: {self.framework}")
            raise ValueError(f"Framework must be 'direct' or 'langgraph', got: {self.framework}")
        
        # Check LangGraph availability if needed
        if self.framework == "langgraph" and not LANGGRAPH_AVAILABLE:
            logger.error("Framework set to 'langgraph' but LangGraph dependencies not available")
            raise ValueError("LangGraph dependencies required for 'langgraph' framework")
        
        logger.info(f"Agent configured with provider: {self.provider}, model: {self.model}, framework: {self.framework}")

# --- LLM Provider Logic ---

class LLMProvider:
    """Handles the interaction with the underlying LLM provider."""
    def __init__(self, config: AgentConfig):
        self.config = config
        self.client = None
        self._initialize_client()
    
    def _initialize_client(self):
        """Initializes the appropriate LLM client based on the configured provider."""
        try:
            if self.config.provider == "openai":
                self.client = openai.OpenAI(
                    api_key=self.config.api_key,
                    base_url=self.config.endpoint
                ) if self.config.endpoint else openai.OpenAI(api_key=self.config.api_key)
            
            elif self.config.provider == "claude":
                self.client = Anthropic(api_key=self.config.api_key)
            
            elif self.config.provider == "gemini":
                genai.configure(api_key=self.config.api_key)
                self.client = genai.GenerativeModel(self.config.model)
            
            elif self.config.provider == "vllm":
                if not self.config.endpoint:
                    raise ValueError("Endpoint is required for the vLLM provider")
                self.client = openai.OpenAI(
                    api_key=self.config.api_key,
                    base_url=self.config.endpoint
                )
            
            else:
                raise ValueError(f"Unsupported provider: {self.config.provider}")
                
            logger.info(f"Successfully initialized client for provider: {self.config.provider}")
            
        except Exception as e:
            logger.error(f"Failed to initialize LLM client: {e}", exc_info=True)
            raise

    @backoff.on_exception(backoff.expo, (httpx.RequestError, openai.RateLimitError), max_tries=3)
    async def chat(self, message: str, conversation_id: Optional[str] = None) -> str:
        """
        Sends a chat message to the LLM and returns the response.
        Includes retry logic for transient network errors and rate limiting.
        """
        try:
            if self.config.provider in ["openai", "vllm"]:
                response = self.client.chat.completions.create(
                    model=self.config.model,
                    messages=[
                        {"role": "system", "content": self.config.system_prompt},
                        {"role": "user", "content": message}
                    ],
                    temperature=0.7,
                    max_tokens=2000
                )
                return response.choices[0].message.content
            
            elif self.config.provider == "claude":
                response = self.client.messages.create(
                    model=self.config.model,
                    max_tokens=2000,
                    system=self.config.system_prompt,
                    messages=[{"role": "user", "content": message}]
                )
                return response.content[0].text
            
            elif self.config.provider == "gemini":
                full_prompt = f"System: {self.config.system_prompt}\n\nUser: {message}"
                response = self.client.generate_content(full_prompt)
                return response.text
                
        except (httpx.RequestError, openai.RateLimitError) as e:
            logger.warning(f"A transient error occurred: {e}. Retrying...")
            raise
        except Exception as e:
            logger.error(f"An unexpected error occurred in chat completion: {e}", exc_info=True)
            raise HTTPException(status_code=500, detail=f"LLM request failed: {str(e)}")

# --- LangGraph Provider ---

class LangGraphProvider:
    """Handles LangGraph-based agent workflows."""
    def __init__(self, config: AgentConfig):
        if not LANGGRAPH_AVAILABLE:
            raise ValueError("LangGraph dependencies not available")
            
        self.config = config
        self.llm = None
        self.workflow = None
        self.sessions = {}  # Simple in-memory session storage
        self._initialize_llm()
        if config.langgraph_config:
            self._build_workflow()
    
    def _initialize_llm(self):
        """Initialize the appropriate LangChain LLM."""
        if self.config.provider == "openai":
            self.llm = ChatOpenAI(
                model=self.config.model,
                openai_api_key=self.config.api_key,
                base_url=self.config.endpoint
            )
        elif self.config.provider == "claude":
            self.llm = ChatAnthropic(
                model=self.config.model,
                anthropic_api_key=self.config.api_key
            )
        elif self.config.provider == "gemini":
            self.llm = ChatGoogleGenerativeAI(
                model=self.config.model,
                google_api_key=self.config.api_key
            )
        else:
            raise ValueError(f"LangGraph provider {self.config.provider} not supported yet")
    
    def _build_workflow(self):
        """Build the LangGraph workflow from configuration."""
        from langgraph.graph import StateGraph, END
        
        workflow = StateGraph(dict)  # Using dict as state for simplicity
        
        # Add nodes
        for node in self.config.langgraph_config.get("nodes", []):
            if node["type"] == "llm":
                workflow.add_node(node["name"], self._create_llm_node(node))
            elif node["type"] == "tool":
                workflow.add_node(node["name"], self._create_tool_node(node))
        
        # Add edges
        for edge in self.config.langgraph_config.get("edges", []):
            if edge.get("condition"):
                workflow.add_conditional_edges(
                    edge["from"],
                    lambda state, e=edge: edge["to"] if self._evaluate_condition(e["condition"], state) else END
                )
            else:
                workflow.add_edge(edge["from"], edge["to"])
        
        # Set entry point
        entrypoint = self.config.langgraph_config.get("entrypoint", "start")
        workflow.set_entry_point(entrypoint)
        
        self.workflow = workflow.compile()
        logger.info("LangGraph workflow compiled successfully")
    
    def _create_llm_node(self, node_config):
        """Create an LLM node function."""
        def llm_node(state):
            prompt = node_config.get("prompt", "{user_input}")
            # Simple template substitution
            try:
                formatted_prompt = prompt.format(**state)
            except KeyError as e:
                logger.warning(f"Template variable missing in state: {e}")
                formatted_prompt = prompt
            
            messages = [
                SystemMessage(content=self.config.system_prompt),
                HumanMessage(content=formatted_prompt)
            ]
            
            response = self.llm.invoke(messages)
            
            # Update state with outputs
            outputs = node_config.get("outputs", ["response"])
            for output in outputs:
                state[output] = response.content
            
            return state
        
        return llm_node
    
    def _create_tool_node(self, node_config):
        """Create a tool node function."""
        def tool_node(state):
            tool_name = node_config.get("tool", "unknown")
            inputs = {key: state.get(key) for key in node_config.get("inputs", [])}
            
            # Mock tool execution - in production, integrate with actual tools
            result = f"Tool {tool_name} executed with inputs: {inputs}"
            
            # Update state with outputs
            outputs = node_config.get("outputs", ["tool_result"])
            for output in outputs:
                state[output] = result
            
            return state
        
        return tool_node
    
    def _evaluate_condition(self, condition, state):
        """Simple condition evaluation."""
        try:
            # Very basic condition evaluation - in production use a safer evaluator
            return eval(condition, {"__builtins__": {}}, state)
        except Exception as e:
            logger.warning(f"Condition evaluation failed: {e}")
            return False
    
    async def chat(self, message: str, conversation_id: str = "default"):
        """Process a chat message through the LangGraph workflow."""
        if not self.workflow:
            raise ValueError("LangGraph workflow not configured")
        
        # Initialize or get existing session state
        if conversation_id not in self.sessions:
            self.sessions[conversation_id] = {"conversation_id": conversation_id}
        
        state = self.sessions[conversation_id].copy()
        state.update({"user_input": message})
        
        # Execute the workflow
        try:
            result = self.workflow.invoke(state)
            
            # Save updated state
            self.sessions[conversation_id] = result
            
            # Return the final response
            return result.get("response", "No response generated from workflow")
        
        except Exception as e:
            logger.error(f"LangGraph workflow execution failed: {e}")
            raise

# --- FastAPI Application ---

# Global instances
try:
    agent_config = AgentConfig()
    
    if agent_config.framework == "direct":
        llm_provider = LLMProvider(agent_config)
        langgraph_provider = None
        logger.info("Initialized with direct framework")
    elif agent_config.framework == "langgraph":
        llm_provider = None
        langgraph_provider = LangGraphProvider(agent_config)
        logger.info("Initialized with LangGraph framework")
    else:
        raise ValueError(f"Unknown framework: {agent_config.framework}")
        
except Exception as e:
    logger.critical(f"Failed to initialize agent: {e}", exc_info=True)
    raise

@app.get("/health", response_model=HealthResponse)
async def health_check():
    """Health check endpoint for Kubernetes liveness probe."""
    return HealthResponse(
        status="healthy",
        provider=agent_config.provider,
        model=agent_config.model,
        timestamp=datetime.now()
    )

@app.get("/ready", response_model=HealthResponse)
async def readiness_check():
    """Readiness check endpoint for Kubernetes readiness probe."""
    if llm_provider.client is None:
        raise HTTPException(status_code=503, detail="LLM client not initialized")
    
    return HealthResponse(
        status="ready",
        provider=agent_config.provider,
        model=agent_config.model,
        timestamp=datetime.now()
    )

@app.post("/chat", response_model=ChatResponse)
async def chat(request: ChatRequest):
    """Main chat endpoint for interacting with the agent."""
    try:
        if agent_config.framework == "direct":
            response_text = await llm_provider.chat(
                message=request.message,
                conversation_id=request.conversation_id
            )
        elif agent_config.framework == "langgraph":
            response_text = await langgraph_provider.chat(
                message=request.message,
                conversation_id=request.conversation_id or "single-turn"
            )
        else:
            raise HTTPException(status_code=500, detail=f"Unknown framework: {agent_config.framework}")
        
        return ChatResponse(
            response=response_text,
            conversation_id=request.conversation_id or "single-turn",
            timestamp=datetime.now(),
            provider=agent_config.provider,
            model=agent_config.model
        )
    
    except HTTPException:
        # Re-raise HTTPException to let FastAPI handle it
        raise
    except Exception as e:
        logger.error(f"Chat request failed: {e}", exc_info=True)
        raise HTTPException(status_code=500, detail="An internal error occurred during the chat request.")

@app.get("/config")
async def get_config():
    """Returns the current agent configuration, excluding sensitive data."""
    return {
        "provider": agent_config.provider,
        "model": agent_config.model,
        "system_prompt_summary": agent_config.system_prompt[:100] + "..." if len(agent_config.system_prompt) > 100 else agent_config.system_prompt,
        "tools_count": agent_config.tools_count,
        "has_custom_endpoint": bool(agent_config.endpoint)
    }

@app.get("/")
async def root():
    """Root endpoint with basic information about the agent."""
    return {
        "name": "KubeAgentic Agent",
        "version": "1.0.0",
        "status": "running",
        "provider": agent_config.provider,
        "model": agent_config.model
    }

if __name__ == "__main__":
    port = int(os.getenv("PORT", "8080"))
    uvicorn.run(
        app,
        host="0.0.0.0",
        port=port,
        log_level="info"
    )

