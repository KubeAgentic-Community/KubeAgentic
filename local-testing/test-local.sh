#!/bin/bash

# KubeAgentic Local Test Runner
# Provides easy commands for testing KubeAgentic locally

set -e

show_help() {
    echo "KubeAgentic Local Test Runner"
    echo "=============================="
    echo ""
    echo "Usage: ./test-local.sh [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  standalone    Test agent in standalone mode (Python only)"
    echo "  docker        Test with Docker Compose"
    echo "  kubernetes    Test with local Kubernetes"
    echo "  basic         Run basic functionality tests"
    echo "  clean         Clean up test resources"
    echo "  help          Show this help message"
    echo ""
    echo "Examples:"
    echo "  ./test-local.sh standalone     # Quick Python test"
    echo "  ./test-local.sh docker         # Full multi-provider test"
    echo "  ./test-local.sh kubernetes     # Deploy to local cluster"
}

test_standalone() {
    echo "🐍 Testing KubeAgentic Standalone"
    echo "================================="
    echo ""
    
    # Check if we have API key
    if [ -z "$OPENAI_API_KEY" ]; then
        echo "⚠️  OPENAI_API_KEY not set. Please set it or create a .env file:"
        echo "   export OPENAI_API_KEY=sk-your-key-here"
        echo "   # OR"
        echo "   cp local-testing/env.example .env  # then edit .env"
        echo ""
        exit 1
    fi
    
    echo "Setting up Python environment..."
    cd ../agent/
    
    if [ ! -d "venv" ]; then
        python -m venv venv
    fi
    
    source venv/bin/activate
    pip install -r requirements.txt
    
    echo "Starting agent in background..."
    export AGENT_PROVIDER=openai
    export AGENT_MODEL=gpt-3.5-turbo
    export AGENT_SYSTEM_PROMPT="You are a test assistant for KubeAgentic."
    export PORT=8080
    
    python main.py &
    AGENT_PID=$!
    
    echo "Waiting for agent to start..."
    sleep 3
    
    echo "Testing agent..."
    cd ../local-testing
    ./scripts/test-basic.sh
    
    echo "Stopping agent..."
    kill $AGENT_PID 2>/dev/null || true
    
    echo "✅ Standalone test complete!"
}

test_docker() {
    echo "🐳 Testing KubeAgentic with Docker Compose"
    echo "=========================================="
    echo ""
    
    # Check if .env exists
    if [ ! -f ".env" ]; then
        echo "Creating .env file from template..."
        cp env.example .env
        echo "⚠️  Please edit .env with your API keys, then run this command again."
        exit 1
    fi
    
    echo "Starting services with Docker Compose..."
    docker-compose -f docker/docker-compose.yml up -d --build
    
    echo "Waiting for services to be ready..."
    sleep 10
    
    echo "Testing services..."
    
    # Test each service
    services=("openai-agent:8081" "claude-agent:8082" "gemini-agent:8083" "vllm-agent:8085")
    
    for service in "${services[@]}"; do
        name=$(echo $service | cut -d: -f1)
        port=$(echo $service | cut -d: -f2)
        
        echo -n "Testing $name... "
        if curl -s http://localhost:$port/health | grep -q "healthy\|ready"; then
            echo "✅"
        else
            echo "❌"
        fi
    done
    
    echo ""
    echo "🧪 Running interactive tests..."
    echo "You can now test the agents:"
    echo ""
    echo "# OpenAI Agent (port 8081)"
    echo 'curl -X POST http://localhost:8081/chat -H "Content-Type: application/json" -d '"'"'{"message": "Hello from OpenAI!"}'"'"''
    echo ""
    echo "# Claude Agent (port 8082)"  
    echo 'curl -X POST http://localhost:8082/chat -H "Content-Type: application/json" -d '"'"'{"message": "Hello from Claude!"}'"'"''
    echo ""
    echo "# vLLM Agent (port 8085)"
    echo 'curl -X POST http://localhost:8085/chat -H "Content-Type: application/json" -d '"'"'{"message": "Hello from vLLM!"}'"'"''
    echo ""
    echo "Press Enter to stop services..."
    read -r
    
    docker-compose -f docker/docker-compose.yml down
    echo "✅ Docker test complete!"
}

test_kubernetes() {
    echo "☸️  Testing KubeAgentic with Kubernetes"
    echo "======================================"
    echo ""
    
    ./scripts/local-deploy.sh
    
    echo ""
    echo "Running basic tests..."
    ./scripts/test-basic.sh
    
    echo "✅ Kubernetes test complete!"
}

run_basic_tests() {
    echo "🧪 Running Basic Tests"
    echo "====================="
    echo ""
    
    ./scripts/test-basic.sh
}

clean_up() {
    echo "🧹 Cleaning up test resources"
    echo "============================="
    echo ""
    
    # Stop Docker Compose if running
    if docker-compose -f docker/docker-compose.yml ps -q > /dev/null 2>&1; then
        echo "Stopping Docker Compose services..."
        docker-compose -f docker/docker-compose.yml down
    fi
    
    # Clean up Kubernetes resources
    if kubectl cluster-info &> /dev/null; then
        echo "Cleaning up Kubernetes resources..."
        kubectl delete agents --all --ignore-not-found=true
        kubectl delete secrets openai-secret claude-secret vllm-secret --ignore-not-found=true
        kubectl delete deployment mock-vllm --ignore-not-found=true
        kubectl delete service mock-vllm --ignore-not-found=true
    fi
    
    # Stop any background Python processes
    pkill -f "python.*main.py" 2>/dev/null || true
    
    # Clean up Python virtual environments
    rm -rf ../agent/venv 2>/dev/null || true
    
    echo "✅ Cleanup complete!"
}

# Main command handling
case "$1" in
    standalone)
        test_standalone
        ;;
    docker)
        test_docker
        ;;
    kubernetes)
        test_kubernetes
        ;;
    basic)
        run_basic_tests
        ;;
    clean)
        clean_up
        ;;
    help|--help|-h|"")
        show_help
        ;;
    *)
        echo "Unknown command: $1"
        echo "Run './test-local.sh help' for available commands"
        exit 1
        ;;
esac
