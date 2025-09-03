#!/bin/bash

# Basic functionality test script for KubeAgentic

set -e

echo "🧪 KubeAgentic Basic Tests"
echo "========================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
TESTS_RUN=0
TESTS_PASSED=0

run_test() {
    local test_name="$1"
    local test_command="$2"
    
    echo -n "Testing: $test_name... "
    TESTS_RUN=$((TESTS_RUN + 1))
    
    if eval "$test_command" &> /dev/null; then
        echo -e "${GREEN}✅ PASS${NC}"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        echo -e "${RED}❌ FAIL${NC}"
    fi
}

# Check if running in Kubernetes mode or standalone mode
if kubectl cluster-info &> /dev/null; then
    KUBERNETES_MODE=true
    echo "🔍 Detected Kubernetes environment"
else
    KUBERNETES_MODE=false
    echo "🔍 Running in standalone mode"
fi

echo ""

if [ "$KUBERNETES_MODE" = true ]; then
    echo "📦 Testing Kubernetes Components"
    echo "--------------------------------"
    
    # Test CRD installation
    run_test "CRD exists" "kubectl get crd agents.ai.example.com"
    
    # Test operator deployment
    run_test "Operator namespace exists" "kubectl get namespace kubeagentic-system"
    run_test "Operator deployment exists" "kubectl get deployment kubeagentic-operator -n kubeagentic-system"
    run_test "Operator is ready" "kubectl get deployment kubeagentic-operator -n kubeagentic-system -o jsonpath='{.status.readyReplicas}' | grep -q '1'"
    
    # Test RBAC
    run_test "ServiceAccount exists" "kubectl get serviceaccount kubeagentic-operator -n kubeagentic-system"
    run_test "ClusterRole exists" "kubectl get clusterrole kubeagentic-operator-role"
    
    echo ""
    echo "🤖 Testing Agent Functionality"
    echo "------------------------------"
    
    # Check if any agents are deployed
    if kubectl get agents &> /dev/null && [ "$(kubectl get agents --no-headers | wc -l)" -gt 0 ]; then
        AGENT_NAME=$(kubectl get agents -o jsonpath='{.items[0].metadata.name}')
        AGENT_NAMESPACE=$(kubectl get agents -o jsonpath='{.items[0].metadata.namespace}')
        
        run_test "Agent exists" "kubectl get agent $AGENT_NAME -n $AGENT_NAMESPACE"
        run_test "Agent service exists" "kubectl get service ${AGENT_NAME}-service -n $AGENT_NAMESPACE"
        
        # Test agent status
        AGENT_PHASE=$(kubectl get agent $AGENT_NAME -n $AGENT_NAMESPACE -o jsonpath='{.status.phase}' 2>/dev/null || echo "")
        if [ "$AGENT_PHASE" = "Running" ]; then
            echo -e "Testing: Agent is running... ${GREEN}✅ PASS${NC}"
            TESTS_PASSED=$((TESTS_PASSED + 1))
        else
            echo -e "Testing: Agent is running... ${YELLOW}⚠️  PENDING (Phase: $AGENT_PHASE)${NC}"
        fi
        TESTS_RUN=$((TESTS_RUN + 1))
        
        # Test agent endpoints (if port-forward is possible)
        echo ""
        echo "🌐 Testing Agent API (requires port-forward)"
        echo "-------------------------------------------"
        echo "💡 To test agent API manually, run:"
        echo "   kubectl port-forward service/${AGENT_NAME}-service 8080:80 -n $AGENT_NAMESPACE"
        echo "   curl http://localhost:8080/health"
        echo "   curl -X POST http://localhost:8080/chat -H 'Content-Type: application/json' -d '{\"message\": \"test\"}'"
        
    else
        echo "ℹ️  No agents deployed. To deploy a test agent:"
        echo "   kubectl apply -f examples/openai-agent.yaml"
        echo "   # (Make sure to create the secret first)"
    fi
    
else
    # Standalone mode tests
    echo "🐍 Testing Standalone Agent"
    echo "---------------------------"
    
    # Check if agent is running
    if pgrep -f "python.*main.py" > /dev/null || curl -s http://localhost:8080/health > /dev/null 2>&1; then
        run_test "Agent process running" "curl -s http://localhost:8080/health"
        run_test "Health endpoint responds" "curl -s http://localhost:8080/health | grep -q 'healthy\\|ready'"
        run_test "Config endpoint responds" "curl -s http://localhost:8080/config"
        run_test "Root endpoint responds" "curl -s http://localhost:8080/"
        
        # Test chat endpoint (basic)
        echo -n "Testing: Chat endpoint responds... "
        TESTS_RUN=$((TESTS_RUN + 1))
        
        RESPONSE=$(curl -s -X POST http://localhost:8080/chat \
                       -H "Content-Type: application/json" \
                       -d '{"message": "test"}' || echo "")
        
        if echo "$RESPONSE" | grep -q '"response"'; then
            echo -e "${GREEN}✅ PASS${NC}"
            TESTS_PASSED=$((TESTS_PASSED + 1))
        else
            echo -e "${YELLOW}⚠️  PARTIAL (Check API key and provider config)${NC}"
        fi
        
    else
        echo "❌ No agent found running on localhost:8080"
        echo "💡 To start standalone agent:"
        echo "   cd agent/"
        echo "   export AGENT_PROVIDER=openai"
        echo "   export AGENT_MODEL=gpt-3.5-turbo"
        echo "   export AGENT_API_KEY=your-api-key"
        echo "   export AGENT_SYSTEM_PROMPT='You are a test assistant'"
        echo "   python main.py"
    fi
fi

echo ""
echo "📊 Test Results"
echo "==============="
echo "Tests run: $TESTS_RUN"
echo "Tests passed: $TESTS_PASSED"
echo "Tests failed: $((TESTS_RUN - TESTS_PASSED))"

if [ $TESTS_PASSED -eq $TESTS_RUN ]; then
    echo -e "${GREEN}🎉 All tests passed!${NC}"
    exit 0
elif [ $TESTS_PASSED -gt 0 ]; then
    echo -e "${YELLOW}⚠️  Some tests failed or are pending${NC}"
    exit 1
else
    echo -e "${RED}❌ All tests failed${NC}"
    exit 1
fi
