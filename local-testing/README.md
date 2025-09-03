# Local Testing for KubeAgentic

This directory contains all the tools and scripts needed to test KubeAgentic locally in various configurations.

## 🚀 Quick Start

```bash
# From the project root, run the main test runner
./local-testing/test-local.sh help
```

## 📁 Directory Structure

```
local-testing/
├── test-local.sh              # Main test runner script
├── env.example                # Environment variables template
├── docs/
│   └── TESTING.md            # Comprehensive testing documentation
├── scripts/
│   ├── local-deploy.sh       # Local Kubernetes deployment
│   └── test-basic.sh         # Basic functionality tests
└── docker/
    ├── docker-compose.yml    # Multi-provider testing setup
    ├── Dockerfile.mock-vllm  # Mock vLLM server dockerfile
    └── mock-vllm/           # Mock vLLM server implementation
        ├── app.py
        └── requirements.txt
```

## 🧪 Testing Options

### 1. Standalone Python Agent (Fastest)
Test just the Python agent application:

```bash
# Set your API key
export OPENAI_API_KEY="sk-your-key-here"

# Run standalone test
./local-testing/test-local.sh standalone
```

**What it does:**
- Sets up Python virtual environment
- Starts agent in background
- Runs basic functionality tests
- Cleans up automatically

### 2. Docker Compose Multi-Provider (Recommended)
Test all LLM providers simultaneously:

```bash
# Copy and configure environment
cp local-testing/env.example .env
# Edit .env with your API keys

# Run Docker tests
./local-testing/test-local.sh docker
```

**What it includes:**
- OpenAI agent (port 8081)
- Claude agent (port 8082)
- Gemini agent (port 8083)
- Mock vLLM server (port 8084)
- vLLM agent (port 8085)

### 3. Local Kubernetes Deployment (Complete)
Deploy to a local Kubernetes cluster:

```bash
# Requires kind, minikube, or k3d
kind create cluster --name kubeagentic-test

# Deploy full system
./local-testing/test-local.sh kubernetes
```

**What it deploys:**
- Complete operator with CRDs
- RBAC and service accounts
- Test agents with secrets
- Health checks and monitoring

## 🛠️ Individual Test Scripts

### Basic Functionality Tests
```bash
./local-testing/scripts/test-basic.sh
```
Validates:
- Component health checks
- API endpoint responses  
- Basic chat functionality
- Configuration loading

### Local Kubernetes Deployment
```bash
./local-testing/scripts/local-deploy.sh
```
Features:
- Auto-detects cluster type (kind, minikube, k3d)
- Builds and loads local images
- Interactive agent deployment
- Comprehensive validation

## 🐳 Docker Testing

### Manual Docker Compose
```bash
cd local-testing
docker-compose -f docker/docker-compose.yml up -d --build

# Test individual agents
curl -X POST http://localhost:8081/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello OpenAI!"}'

curl -X POST http://localhost:8082/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello Claude!"}'

# Cleanup
docker-compose -f docker/docker-compose.yml down
```

### Mock vLLM Server
The included mock vLLM server provides:
- OpenAI-compatible API endpoints
- Realistic response simulation
- No actual model weights required
- Perfect for CI/CD testing

## 🔧 Configuration

### Environment Variables
Copy `env.example` to create your `.env` file:

```bash
# API Keys
OPENAI_API_KEY=sk-your-openai-key
CLAUDE_API_KEY=sk-ant-your-claude-key
GEMINI_API_KEY=your-gemini-key

# Optional vLLM settings
VLLM_ENDPOINT=http://localhost:8000/v1
VLLM_API_KEY=your-vllm-token
```

### Kubernetes Configuration
For local Kubernetes testing, ensure you have:
- `kubectl` installed and configured
- Local cluster running (kind/minikube/k3d)
- Docker for building images

## 🐛 Troubleshooting

### Common Issues

**Python virtual environment errors:**
```bash
# Clean up and retry
rm -rf ../agent/venv
./local-testing/test-local.sh standalone
```

**Docker build failures:**
```bash
# Clean Docker cache
docker system prune -a
./local-testing/test-local.sh docker
```

**Kubernetes deployment issues:**
```bash
# Check cluster status
kubectl cluster-info
kubectl get nodes

# Verify images are loaded (for kind)
docker exec -it kubeagentic-test-control-plane crictl images | grep kubeagentic
```

### Debug Mode
Enable verbose logging:
```bash
# For Python agent
export LOG_LEVEL=DEBUG
./local-testing/test-local.sh standalone

# For operator
export OPERATOR_DEBUG=true
./local-testing/scripts/local-deploy.sh
```

## 🧹 Cleanup

Clean up all test resources:
```bash
./local-testing/test-local.sh clean
```

This removes:
- Docker containers and images
- Kubernetes test resources
- Python virtual environments  
- Background processes

## 📚 Additional Resources

- [Comprehensive Testing Guide](docs/TESTING.md) - Detailed testing documentation
- [Main README](../README.md) - Project overview and production deployment
- [Examples](../examples/) - Sample agent configurations

## 💡 Tips for Effective Testing

1. **Start Simple**: Begin with standalone testing to validate basic functionality
2. **Use Mock Services**: Leverage the mock vLLM server for consistent testing  
3. **Test Incrementally**: Validate each component before moving to integration tests
4. **Monitor Resources**: Use `kubectl top` and `docker stats` to monitor resource usage
5. **Clean Between Tests**: Always clean up resources between different test runs

## 🤝 Contributing to Tests

When adding new features, please:
1. Add corresponding tests to the appropriate test script
2. Update this README with new testing procedures
3. Ensure all test modes (standalone, docker, kubernetes) work
4. Add example configurations if needed

This testing framework ensures KubeAgentic works reliably across different deployment scenarios!

