---
layout: page
title: Local Testing
description: Comprehensive guide for testing KubeAgentic locally and in development environmentspermalink: /local-testing/
---

# KubeAgentic Local Testing Guide

This comprehensive guide covers everything you need to test KubeAgentic locally, from simple standalone testing to full Kubernetes deployments.

## 🚀 Quick Start

The fastest way to get started with local testing:

```bash
# Clone the repository
git clone https://github.com/sudeshmu/kubeagentic.git
cd kubeagentic

# Quick test with Docker
./local-testing/test-local.sh docker
```

## 📋 Prerequisites

### Required Tools

- **Docker**: For containerized testing
- **kubectl**: Kubernetes command-line tool
- **Go 1.21+**: For building the operator
- **Python 3.8+**: For agent development

### Optional Tools

- **kind/minikube/k3d**: Local Kubernetes cluster
- **make**: Build automation
- **jq**: JSON processing (for test scripts)

### API Keys

You'll need at least one AI provider API key:
- **OpenAI**: `sk-...` format
- **Anthropic Claude**: `sk-ant-...` format  
- **Google Gemini**: Standard API key
- **vLLM**: Self-hosted endpoint (optional)

## 🧪 Testing Methods

### 1. Standalone Python Agent

**Best for**: Quick development and debugging

```bash
# Set up environment
export OPENAI_API_KEY="sk-your-key-here"

# Run standalone test
./local-testing/test-local.sh standalone
```

**What happens**:
1. Creates Python virtual environment
2. Installs dependencies
3. Starts agent server
4. Runs functionality tests
5. Cleans up automatically

**Output**:
```
🤖 Starting KubeAgentic Standalone Tests...
✅ Virtual environment created
✅ Dependencies installed
✅ Agent server started (PID: 12345)
✅ Health check passed
✅ Chat functionality working
✅ Cleanup completed
```

### 2. Docker Compose Multi-Provider

**Best for**: Testing multiple providers simultaneously

```bash
# Configure environment
cp local-testing/env.example .env
# Edit .env with your API keys

# Run Docker tests
./local-testing/test-local.sh docker
```

**Services included**:
- OpenAI agent (http://localhost:8081)
- Claude agent (http://localhost:8082) 
- Gemini agent (http://localhost:8083)
- Mock vLLM server (http://localhost:8084)
- vLLM agent (http://localhost:8085)

**Testing each service**:
```bash
# Test OpenAI agent
curl -X POST http://localhost:8081/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello from OpenAI!"}'

# Test Claude agent
curl -X POST http://localhost:8082/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello from Claude!"}'

# Test mock vLLM
curl -X POST http://localhost:8084/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "llama2-7b-chat",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

### 3. Local Kubernetes Deployment

**Best for**: Complete integration testing

```bash
# Create local cluster
kind create cluster --name kubeagentic-test

# Deploy full system
./local-testing/test-local.sh kubernetes
```

**What gets deployed**:
- Custom Resource Definitions (CRDs)
- RBAC and ServiceAccounts
- KubeAgentic operator
- Sample agents with secrets
- Monitoring and health checks

**Verify deployment**:
```bash
# Check operator status
kubectl get pods -n kubeagentic-system

# List agents
kubectl get agents

# Check agent logs
kubectl logs -l kubeagentic.ai/agent=test-agent
```

## 🔧 Environment Configuration

### Create Environment File

```bash
cp local-testing/env.example .env
```

Edit `.env` with your settings:

```bash
# Required: At least one AI provider
OPENAI_API_KEY=sk-your-openai-key-here
CLAUDE_API_KEY=sk-ant-your-claude-key-here
GEMINI_API_KEY=your-gemini-key-here

# Optional: Self-hosted vLLM
VLLM_ENDPOINT=http://your-vllm-server:8000/v1
VLLM_API_KEY=your-vllm-token

# Optional: Custom settings
DEFAULT_MODEL=gpt-3.5-turbo
MAX_TOKENS=2048
TEMPERATURE=0.7

# Debug settings
LOG_LEVEL=INFO
DEBUG_MODE=false
```

### Kubernetes Secrets

For Kubernetes testing, secrets are created automatically:

```bash
# Manual secret creation (if needed)
kubectl create secret generic openai-secret \
  --from-literal=api-key="$OPENAI_API_KEY"

kubectl create secret generic claude-secret \
  --from-literal=api-key="$CLAUDE_API_KEY"
```

## 📁 Directory Structure

```
local-testing/
├── test-local.sh              # Main test runner
├── env.example                # Environment template
├── docs/
│   └── TESTING.md            # Detailed testing docs
├── scripts/
│   ├── local-deploy.sh       # Kubernetes deployment
│   ├── test-basic.sh         # Basic functionality tests
│   ├── cleanup.sh            # Resource cleanup
│   └── build-images.sh       # Docker image building
└── docker/
    ├── docker-compose.yml    # Multi-service setup
    ├── Dockerfile.mock-vllm  # Mock server image
    └── mock-vllm/           # Mock implementation
        ├── app.py
        └── requirements.txt
```

## 🐳 Docker Testing Details

### Manual Docker Compose

```bash
# Start all services
cd local-testing
docker-compose -f docker/docker-compose.yml up -d --build

# Check service status
docker-compose -f docker/docker-compose.yml ps

# View logs
docker-compose -f docker/docker-compose.yml logs openai-agent

# Stop services
docker-compose -f docker/docker-compose.yml down
```

### Mock vLLM Server

The included mock server provides:
- OpenAI-compatible API endpoints
- Realistic response formatting
- No actual model weights required
- Configurable response delays
- Perfect for CI/CD pipelines

**Mock server endpoints**:
```bash
# Health check
curl http://localhost:8084/health

# Model list
curl http://localhost:8084/v1/models

# Chat completion
curl -X POST http://localhost:8084/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "llama2-7b-chat",
    "messages": [{"role": "user", "content": "Hello!"}],
    "max_tokens": 100
  }'
```

## ☸️ Kubernetes Testing Details

### Cluster Setup

**Using kind**:
```bash
# Create cluster with specific configuration
cat <<EOF | kind create cluster --name kubeagentic --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
  - containerPort: 443
    hostPort: 443
EOF
```

**Using minikube**:
```bash
minikube start --driver=docker --cpus=4 --memory=8192
minikube addons enable ingress
```

**Using k3d**:
```bash
k3d cluster create kubeagentic --port "8080:80@loadbalancer"
```

### Deployment Steps

1. **Build and Load Images**:
   ```bash
   # Build operator image
   make docker-build IMG=kubeagentic:latest
   
   # Load into kind cluster
   kind load docker-image kubeagentic:latest --name kubeagentic
   ```

2. **Deploy CRDs and RBAC**:
   ```bash
   kubectl apply -f deploy/crds/
   kubectl apply -f deploy/rbac/
   ```

3. **Deploy Operator**:
   ```bash
   kubectl apply -f deploy/operator.yaml
   ```

4. **Create Test Agents**:
   ```bash
   # Apply test configurations
   kubectl apply -f local-testing/configs/test-agents.yaml
   ```

### Verification

```bash
# Check operator
kubectl get pods -n kubeagentic-system
kubectl logs -n kubeagentic-system deployment/kubeagentic-operator

# Check agents
kubectl get agents
kubectl describe agent test-openai-agent

# Test agent functionality
kubectl port-forward service/test-openai-agent-service 8080:80
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Test message"}'
```

## 🔍 Testing Scenarios

### Basic Functionality Tests

```bash
./local-testing/scripts/test-basic.sh
```

Tests include:
- Agent startup and health checks
- API endpoint responses
- Message processing
- Error handling
- Resource cleanup

### Load Testing

```bash
# Install hey (HTTP load testing tool)
go install github.com/rakyll/hey@latest

# Run load test
hey -n 100 -c 10 -m POST \
  -H "Content-Type: application/json" \
  -d '{"message": "Load test message"}' \
  http://localhost:8080/chat
```

### Multi-Provider Testing

```bash
# Test all providers simultaneously
./local-testing/scripts/test-providers.sh

# Compare response times
./local-testing/scripts/benchmark-providers.sh
```

### Security Testing

```bash
# Test with invalid API keys
export OPENAI_API_KEY="invalid-key"
./local-testing/test-local.sh standalone

# Test network policies (Kubernetes only)
kubectl apply -f local-testing/configs/network-policies.yaml
./local-testing/scripts/test-network-isolation.sh
```

## 🐛 Troubleshooting

### Common Issues

**Python virtual environment errors**:
```bash
# Clean and recreate
rm -rf agent/venv
python3 -m venv agent/venv
source agent/venv/bin/activate
pip install -r agent/requirements.txt
```

**Docker build failures**:
```bash
# Clean Docker cache
docker system prune -a
docker builder prune

# Rebuild with no cache
docker-compose build --no-cache
```

**Kubernetes deployment issues**:
```bash
# Check cluster connectivity
kubectl cluster-info
kubectl get nodes

# Verify images (kind)
docker exec -it kubeagentic-control-plane crictl images

# Check resource constraints
kubectl describe nodes
kubectl top nodes
```

**API connection failures**:
```bash
# Verify API keys
echo $OPENAI_API_KEY | wc -c  # Should be ~51 characters

# Test direct API access
curl -H "Authorization: Bearer $OPENAI_API_KEY" \
  https://api.openai.com/v1/models
```

### Debug Mode

Enable verbose logging:

```bash
# For standalone testing
export LOG_LEVEL=DEBUG
./local-testing/test-local.sh standalone

# For Docker testing
echo "LOG_LEVEL=DEBUG" >> .env
./local-testing/test-local.sh docker

# For Kubernetes testing
export OPERATOR_DEBUG=true
./local-testing/scripts/local-deploy.sh
```

### Log Analysis

```bash
# Agent logs
tail -f agent/logs/agent.log

# Docker logs
docker-compose logs -f openai-agent

# Kubernetes logs
kubectl logs -f deployment/kubeagentic-operator -n kubeagentic-system
kubectl logs -l kubeagentic.ai/agent=test-agent
```

## 🧹 Cleanup

### Complete Cleanup

```bash
# Clean all test resources
./local-testing/test-local.sh clean
```

### Selective Cleanup

```bash
# Docker only
docker-compose -f local-testing/docker/docker-compose.yml down
docker system prune -f

# Kubernetes only
kubectl delete -f local-testing/configs/
kind delete cluster --name kubeagentic

# Python environments only
rm -rf agent/venv
```

## 🚀 Performance Optimization

### Resource Tuning

**For development**:
```yaml
resources:
  requests:
    memory: "64Mi"
    cpu: "50m"
  limits:
    memory: "128Mi"
    cpu: "200m"
```

**For testing**:
```yaml
resources:
  requests:
    memory: "128Mi"
    cpu: "100m"
  limits:
    memory: "256Mi"
    cpu: "500m"
```

### Caching

Enable response caching:
```bash
export ENABLE_CACHE=true
export CACHE_TTL=300  # 5 minutes
```

## 📊 Monitoring

### Metrics Collection

```bash
# Enable metrics
export ENABLE_METRICS=true

# View metrics
curl http://localhost:9090/metrics
```

### Health Monitoring

```bash
# Health check endpoint
curl http://localhost:8080/health

# Kubernetes health
kubectl get pods -w
```

## 🤝 Contributing Tests

When adding new features:

1. **Add unit tests**:
   ```bash
   cd agent
   python -m pytest tests/
   ```

2. **Add integration tests**:
   ```bash
   # Update test-basic.sh with new test cases
   vim local-testing/scripts/test-basic.sh
   ```

3. **Update documentation**:
   ```bash
   # Update this guide and TESTING.md
   vim local-testing/docs/TESTING.md
   ```

4. **Test all modes**:
   ```bash
   ./local-testing/test-local.sh standalone
   ./local-testing/test-local.sh docker
   ./local-testing/test-local.sh kubernetes
   ```

## 📚 Additional Resources

- [Comprehensive Testing Documentation](https://github.com/sudeshmu/kubeagentic/blob/main/local-testing/docs/TESTING.md)
- [Main Project README](https://github.com/sudeshmu/kubeagentic/blob/main/README.md)
- [Examples and Use Cases](examples.html)
- [API Reference](api-reference.html)

## 💡 Testing Best Practices

1. **Start Small**: Begin with standalone testing
2. **Use Mocks**: Leverage mock services for consistent results
3. **Test Incrementally**: Validate each component separately
4. **Monitor Resources**: Keep an eye on CPU and memory usage
5. **Clean Between Tests**: Always clean up between test runs
6. **Version API Keys**: Use different keys for testing vs production
7. **Document Changes**: Update tests when adding features

Happy testing! 🧪✨