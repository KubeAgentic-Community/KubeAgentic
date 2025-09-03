# Local Testing Guide for KubeAgentic

This guide provides multiple approaches to test KubeAgentic locally, from quick standalone tests to full Kubernetes deployments.

## 🚀 Quick Testing Options

### Option 1: Standalone Agent Testing (Fastest)

Test just the agent application without Kubernetes:

```bash
# 1. Set up Python environment
cd agent/
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
pip install -r requirements.txt

# 2. Set environment variables
export AGENT_PROVIDER="openai"
export AGENT_MODEL="gpt-3.5-turbo" 
export AGENT_SYSTEM_PROMPT="You are a helpful AI assistant."
export AGENT_API_KEY="your-openai-api-key-here"
export PORT="8080"

# 3. Run the agent
python main.py
```

Then test in another terminal:
```bash
# Test health endpoint
curl http://localhost:8080/health

# Test chat endpoint
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello! Can you help me test this agent?"}'
```

### Option 2: Docker Compose (Recommended for Quick Testing)

```bash
# 1. Create .env file
cat > .env << EOF
OPENAI_API_KEY=your-openai-api-key-here
CLAUDE_API_KEY=your-claude-api-key-here
EOF

# 2. Start services
docker-compose up --build

# 3. Test different agents
curl -X POST http://localhost:8081/chat -H "Content-Type: application/json" -d '{"message": "Hello from OpenAI agent!"}'
curl -X POST http://localhost:8082/chat -H "Content-Type: application/json" -d '{"message": "Hello from Claude agent!"}'
```

### Option 3: Local Kubernetes Cluster

Choose your preferred local Kubernetes:

#### Using Kind (Recommended)
```bash
# 1. Install kind
# macOS: brew install kind
# Windows: choco install kind
# Linux: see https://kind.sigs.k8s.io/docs/user/quick-start/

# 2. Create cluster
kind create cluster --name kubeagentic-test

# 3. Deploy KubeAgentic
./scripts/local-deploy.sh
```

#### Using Minikube
```bash
# 1. Start minikube
minikube start

# 2. Deploy KubeAgentic
./scripts/local-deploy.sh
```

#### Using K3d
```bash
# 1. Create cluster
k3d cluster create kubeagentic-test

# 2. Deploy KubeAgentic
./scripts/local-deploy.sh
```

### Option 4: Operator Development Mode

Run the operator locally while connecting to a real cluster:

```bash
# 1. Make sure you have a kubeconfig pointing to your cluster
export KUBECONFIG=~/.kube/config

# 2. Install CRDs and RBAC (operator will run locally)
kubectl apply -f deploy/namespace.yaml
kubectl apply -f crd/agent-crd.yaml
kubectl apply -f deploy/rbac.yaml

# 3. Run operator locally
make run
```

## 🧪 Testing Scenarios

### Basic Functionality Tests

```bash
# Run all basic tests
./scripts/test-basic.sh
```

This includes:
- Health check endpoints
- Configuration validation
- Secret management
- Basic chat functionality

### Multi-Provider Tests

```bash
# Test OpenAI
./scripts/test-provider.sh openai

# Test Claude  
./scripts/test-provider.sh claude

# Test Gemini
./scripts/test-provider.sh gemini

# Test vLLM (requires local vLLM instance)
./scripts/test-provider.sh vllm
```

### Load Testing

```bash
# Simple load test
./scripts/load-test.sh

# Or use Apache Bench
ab -n 100 -c 10 -p test-data/chat-request.json -T application/json http://localhost:8080/chat
```

### Integration Tests

```bash
# Full integration test suite
./scripts/integration-test.sh
```

## 🐛 Debugging and Troubleshooting

### Common Issues

**1. Agent won't start**
```bash
# Check logs
kubectl logs -l kubeagentic.ai/agent=your-agent-name

# Check operator logs
kubectl logs -n kubeagentic-system deployment/kubeagentic-operator
```

**2. Secret issues**
```bash
# Verify secret exists and has correct key
kubectl get secret openai-secret -o yaml
kubectl describe agent your-agent-name
```

**3. Network issues**
```bash
# Test service connectivity
kubectl get svc
kubectl port-forward service/your-agent-service 8080:80

# Check endpoints
kubectl get endpoints
```

### Debug Mode

**For Agent (Python)**:
```bash
# Enable debug logging
export LOG_LEVEL=DEBUG
python agent/main.py
```

**For Operator (Go)**:
```bash
# Run with verbose logging
go run main.go --zap-devel --zap-log-level=debug
```

## 📊 Monitoring and Metrics

### Local Monitoring Setup

```bash
# Deploy monitoring stack (optional)
kubectl apply -f monitoring/

# Access Grafana
kubectl port-forward -n monitoring service/grafana 3000:3000
# Open http://localhost:3000 (admin/admin)
```

### Health Checks

```bash
# Check all components
./scripts/health-check.sh

# Manual health checks
curl http://localhost:8080/health
curl http://localhost:8080/ready
curl http://localhost:8080/metrics  # If metrics enabled
```

## 🔧 Development Workflow

### Making Changes

1. **Code Changes**:
   ```bash
   # Edit files
   vim controllers/agent_controller.go
   vim agent/main.py
   ```

2. **Test Changes**:
   ```bash
   # Test Go code
   make test
   
   # Test Python code
   cd agent && python -m pytest tests/
   ```

3. **Build and Deploy**:
   ```bash
   # For operator changes
   make docker-build-operator
   kubectl rollout restart deployment/kubeagentic-operator -n kubeagentic-system
   
   # For agent changes  
   make docker-build-agent
   kubectl rollout restart deployment/your-agent-name
   ```

### Hot Reloading

**Python Agent**: Use `uvicorn` with `--reload`:
```bash
cd agent/
uvicorn main:app --reload --host 0.0.0.0 --port 8080
```

**Go Operator**: Use `air` for hot reloading:
```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

## 🎯 Test Data and Examples

### Sample Requests

**Simple Chat**:
```json
{
  "message": "Hello, how are you today?",
  "conversation_id": "test-session-1"
}
```

**Complex Chat with Context**:
```json
{
  "message": "Based on our previous discussion about the order issue, what should I do next?",
  "conversation_id": "support-session-123",
  "context": {
    "order_id": "ORD-12345",
    "customer_tier": "premium",
    "issue_type": "shipping"
  }
}
```

### Test Agent Configurations

See `examples/` directory for various agent configurations:
- `openai-agent.yaml` - Customer support agent
- `claude-agent.yaml` - Code review agent  
- `vllm-agent.yaml` - Internal Q&A agent
- `gemini-agent.yaml` - Content creation agent

## 🚨 Performance Testing

### Resource Usage Monitoring

```bash
# Monitor resource usage
kubectl top pods -l kubeagentic.ai/agent
kubectl top nodes

# Detailed resource monitoring
kubectl describe pod your-agent-pod
```

### Scaling Tests

```bash
# Test horizontal scaling
kubectl scale agent your-agent --replicas=5

# Monitor scaling
kubectl get pods -w -l kubeagentic.ai/agent=your-agent
```

## 🔒 Security Testing

### Secret Management

```bash
# Verify secrets are not logged
kubectl logs -l kubeagentic.ai/agent | grep -i "api.*key" || echo "No API keys found in logs ✅"

# Check secret access
kubectl auth can-i get secrets --as=system:serviceaccount:default:your-agent
```

### Network Policies

```bash
# Apply network policies
kubectl apply -f security/network-policies.yaml

# Test connectivity
kubectl exec -it your-agent-pod -- curl http://external-service/
```

This comprehensive testing guide should help you validate KubeAgentic in various scenarios!
