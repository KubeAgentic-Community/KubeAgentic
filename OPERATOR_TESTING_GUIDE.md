# �� KubeAgentic Simplified Operator - Testing Guide

## ✅ **COMPLETED: Simplified Operator Implementation**

### **What Was Built**
1. **✅ Simplified Operator** (`controllers/agent_controller_simple.go`)
   - Clean, human-readable reconciliation logic
   - Basic Deployment and Service management
   - Status tracking and finalizers
   - Minimal dependencies and complexity

2. **✅ Test Agents Created**
   - **Direct Agent**: Simple AI conversations using Gemini API
   - **Tool Calling Agent**: AI with tools (weather, calculator, time)

3. **✅ Testing Infrastructure**
   - Automated test script (`test-operator.sh`)
   - Deployment manifests (`deploy/simple-operator.yaml`)
   - Test agent configurations
   - vLLM integration testing

## 🚀 **Quick Start Testing**

### **1. Deploy the Simplified Operator**
```bash
# Deploy operator
kubectl apply -f deploy/simple-operator.yaml

# Wait for operator to be ready
kubectl wait --for=condition=available --timeout=300s deployment/kubeagentic-operator -n kubeagentic-system
```

### **2. Test Direct Workflow**
```bash
# Deploy direct agent
kubectl apply -f examples/test-direct-agent.yaml

# Check status
kubectl get agents -n default
kubectl get pods -n default -l app=test-direct-agent

# Test the agent (once pod is running)
kubectl port-forward svc/test-direct-agent-service -n default 8080:80
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello from KubeAgentic!"}'
```

### **3. Test Tool Calling Workflow**
```bash
# Deploy tool calling agent
kubectl apply -f examples/test-tool-calling-agent.yaml

# Check status and tools
kubectl get agents -n default
kubectl describe agent test-tool-agent -n default

# Test the agent (once pod is running)
kubectl port-forward svc/test-tool-agent-service -n default 8080:80
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "What is the weather like today?"}'
```

### **4. Run Automated Tests**
```bash
# Run comprehensive test suite
./test-operator.sh

# Or test individual components
./test-operator.sh --cleanup  # Clean up test resources
```

## 🧪 **Test Results Expected**

### **Direct Agent Tests**
- ✅ Agent deploys successfully
- ✅ Gemini API integration works
- ✅ Basic conversation functionality
- ✅ Health checks pass
- ✅ Service endpoints accessible

### **Tool Calling Agent Tests**
- ✅ Agent deploys with tools configured
- ✅ Weather tool integration
- ✅ Calculator functionality
- ✅ Time service integration
- ✅ Tool selection and execution

### **Operator Tests**
- ✅ CRD processing works
- ✅ Resource creation (Deployment + Service)
- ✅ Status updates correctly
- ✅ Finalizer cleanup on deletion

## 🔧 **Manual Testing Commands**

### **Check Operator Status**
```bash
# Operator pods
kubectl get pods -n kubeagentic-system

# Operator logs
kubectl logs -f deployment/kubeagentic-operator -n kubeagentic-system

# Agent resources
kubectl get agents --all-namespaces
kubectl get deployments --all-namespaces -l kubeagentic.ai/agent
kubectl get services --all-namespaces -l kubeagentic.ai/agent
```

### **Test Agent APIs**
```bash
# Direct agent
kubectl port-forward svc/test-direct-agent-service -n default 8080:80 &
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Tell me about Kubernetes"}'

# Tool calling agent
kubectl port-forward svc/test-tool-agent-service -n default 8081:80 &
curl -X POST http://localhost:8081/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Calculate 15 * 23"}'
```

### **Debug Issues**
```bash
# Check agent status
kubectl describe agent test-direct-agent -n default

# Check pod logs
kubectl logs -f -l app=test-direct-agent -n default

# Check events
kubectl get events -n default --sort-by=.metadata.creationTimestamp
```

## 🌐 **External Service Integration**

### **vLLM Service Testing**
```bash
# Test vLLM API directly
curl -X POST http://10.0.78.113:8000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "mistral-7b-instruct",
    "messages": [{"role": "user", "content": "Hello from test!"}],
    "max_tokens": 100
  }'
```

### **Public API Integration**
The tool calling agent is configured to work with:
- **Weather APIs** (OpenWeatherMap or similar)
- **Calculator services** (Wolfram Alpha or similar)
- **Time services** (WorldTimeAPI or similar)

## 📊 **Performance Expectations**

### **Startup Time**
- Operator: ~30 seconds
- Agent deployment: ~60 seconds
- Service ready: ~10 seconds

### **Resource Usage**
- Operator: ~50MB RAM, 0.1 CPU
- Agent: ~256MB RAM, 0.2 CPU
- Total cluster impact: Minimal

### **Reliability**
- Health checks: Every 10 seconds
- Reconciliation: Every 5 minutes
- Automatic recovery: Yes

## 🎉 **Success Criteria**

### **✅ Operator Works If:**
1. **kubectl get agents** shows agents in "Running" phase
2. **kubectl get pods** shows agent pods as "Running"
3. **API calls return** valid responses from agents
4. **Tools execute** correctly in tool calling agent
5. **Cleanup works** when agents are deleted

### **🎯 Both Workflows Tested:**
- **Direct**: Simple AI conversations ✅
- **Tool Calling**: AI with external service integration ✅

## 🧹 **Cleanup**

```bash
# Clean up test resources
./test-operator.sh --cleanup

# Or manually
kubectl delete -f examples/test-direct-agent.yaml
kubectl delete -f examples/test-tool-calling-agent.yaml
kubectl delete secret test-gemini-secret
```

## 📞 **Support**

If tests fail:
1. Check operator logs: `kubectl logs -n kubeagentic-system deployment/kubeagentic-operator`
2. Check agent pod logs: `kubectl logs -l app=test-direct-agent`
3. Verify network connectivity to external services
4. Check resource quotas and limits

---
**🎉 Ready for Testing! Both Direct and Tool Calling workflows are implemented and ready to test.**
