#!/usr/bin/env python3
"""
Mock vLLM Server for Testing

This creates a simple OpenAI-compatible API server that mimics vLLM's behavior
for testing purposes without needing actual model weights.
"""

import os
import json
import time
import random
from typing import Dict, List
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import uvicorn
from datetime import datetime

app = FastAPI(title="Mock vLLM Server", version="1.0.0")

class ChatMessage(BaseModel):
    role: str
    content: str

class ChatCompletionRequest(BaseModel):
    model: str
    messages: List[ChatMessage]
    temperature: float = 0.7
    max_tokens: int = 2000
    stream: bool = False

class ChatCompletionResponse(BaseModel):
    id: str
    object: str = "chat.completion"
    created: int
    model: str
    choices: List[Dict]
    usage: Dict[str, int]

# Mock responses for different message patterns
MOCK_RESPONSES = [
    "Hello! I'm a mock LLM response. This is simulating a local vLLM server.",
    "I understand you're testing the KubeAgentic system. Everything looks good!",
    "This is a simulated response from a locally hosted language model.",
    "Mock response: I'm helping you test multi-provider AI agent deployment.",
    "Testing successful! This response comes from the mock vLLM server.",
]

@app.get("/health")
async def health_check():
    """Health check endpoint."""
    return {
        "status": "healthy",
        "model": os.getenv("MODEL_NAME", "llama2-7b-chat"),
        "timestamp": datetime.now().isoformat()
    }

@app.get("/v1/models")
async def list_models():
    """List available models (OpenAI-compatible)."""
    model_name = os.getenv("MODEL_NAME", "llama2-7b-chat")
    return {
        "object": "list",
        "data": [
            {
                "id": model_name,
                "object": "model",
                "created": int(time.time()),
                "owned_by": "mock-vllm",
            }
        ]
    }

@app.post("/v1/chat/completions")
async def create_chat_completion(request: ChatCompletionRequest):
    """Create chat completion (OpenAI-compatible)."""
    
    # Simulate some processing time
    await asyncio.sleep(0.1 + random.uniform(0, 0.3))
    
    # Get the user's message
    user_message = ""
    for msg in request.messages:
        if msg.role == "user":
            user_message = msg.content
            break
    
    # Generate a mock response
    mock_response = random.choice(MOCK_RESPONSES)
    if user_message.lower() in ["hello", "hi", "hey"]:
        mock_response = f"Hello! I received your message: '{user_message}'. This is a mock response from vLLM."
    elif "test" in user_message.lower():
        mock_response = f"Test confirmed! Your message was: '{user_message}'. Mock vLLM is working correctly."
    elif "help" in user_message.lower():
        mock_response = "I'm a mock vLLM server for testing KubeAgentic. I can simulate responses to help you test your AI agent deployment."
    
    response = ChatCompletionResponse(
        id=f"chatcmpl-mock-{int(time.time())}{random.randint(1000, 9999)}",
        created=int(time.time()),
        model=request.model,
        choices=[
            {
                "index": 0,
                "message": {
                    "role": "assistant",
                    "content": mock_response
                },
                "finish_reason": "stop"
            }
        ],
        usage={
            "prompt_tokens": len(user_message.split()),
            "completion_tokens": len(mock_response.split()),
            "total_tokens": len(user_message.split()) + len(mock_response.split())
        }
    )
    
    return response

@app.get("/")
async def root():
    """Root endpoint with server info."""
    return {
        "name": "Mock vLLM Server",
        "version": "1.0.0",
        "model": os.getenv("MODEL_NAME", "llama2-7b-chat"),
        "status": "running",
        "description": "Mock OpenAI-compatible API server for testing KubeAgentic"
    }

# Add asyncio import
import asyncio

if __name__ == "__main__":
    port = int(os.getenv("PORT", "8000"))
    uvicorn.run(
        app,
        host="0.0.0.0",
        port=port,
        log_level="info"
    )
