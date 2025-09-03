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
        self.tools_count = int(os.getenv("AGENT_TOOLS_COUNT", "0"))
        
        if not self.api_key:
            logger.error("AGENT_API_KEY environment variable is not set.")
            raise ValueError("AGENT_API_KEY environment variable is required")
        
        logger.info(f"Agent configured with provider: {self.provider}, model: {self.model}")

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

# --- FastAPI Application ---

# Global instances
try:
    agent_config = AgentConfig()
    llm_provider = LLMProvider(agent_config)
except Exception as e:
    logger.critical(f"Failed to initialize agent: {e}", exc_info=True)
    # In a real application, you might want to exit here if the agent can't be initialized.
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
        response_text = await llm_provider.chat(
            message=request.message,
            conversation_id=request.conversation_id
        )
        
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

