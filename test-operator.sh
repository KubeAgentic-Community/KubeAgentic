#!/bin/bash

# Simple Operator Test Script
# Tests both Direct and Tool Calling workflows

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Test Direct Agent
test_direct_agent() {
    log_info "🧪 Testing Direct Agent Workflow"

    # Deploy direct agent
    log_info "Deploying test direct agent..."
    kubectl apply -f examples/test-direct-agent.yaml

    # Wait for deployment
    log_info "Waiting for direct agent to be ready..."
    kubectl wait --for=condition=available --timeout=300s deployment/test-direct-agent -n default

    # Check status
    local phase=$(kubectl get agent test-direct-agent -n default -o jsonpath='{.status.phase}')
    if [ "$phase" = "Running" ]; then
        log_success "✅ Direct agent is running"
    else
        log_error "❌ Direct agent failed to start. Phase: $phase"
        kubectl describe agent test-direct-agent -n default
        return 1
    fi

    # Test the agent API
    log_info "Testing direct agent API..."
    local service_ip=$(kubectl get svc test-direct-agent-service -n default -o jsonpath='{.spec.clusterIP}')
    if [ -n "$service_ip" ]; then
        log_success "✅ Direct agent service created: $service_ip"
    else
        log_error "❌ Direct agent service not found"
        return 1
    fi

    log_success "🎉 Direct agent test completed successfully"
}

# Test Tool Calling Agent
test_tool_calling_agent() {
    log_info "🔧 Testing Tool Calling Agent Workflow"

    # Deploy tool calling agent
    log_info "Deploying test tool calling agent..."
    kubectl apply -f examples/test-tool-calling-agent.yaml

    # Wait for deployment
    log_info "Waiting for tool calling agent to be ready..."
    kubectl wait --for=condition=available --timeout=300s deployment/test-tool-agent -n default

    # Check status
    local phase=$(kubectl get agent test-tool-agent -n default -o jsonpath='{.status.phase}')
    if [ "$phase" = "Running" ]; then
        log_success "✅ Tool calling agent is running"
    else
        log_error "❌ Tool calling agent failed to start. Phase: $phase"
        kubectl describe agent test-tool-agent -n default
        return 1
    fi

    # Check tools configuration
    local tools_count=$(kubectl get agent test-tool-agent -n default -o jsonpath='{.spec.tools[*].name}' | wc -w)
    if [ "$tools_count" -gt 0 ]; then
        log_success "✅ Tool calling agent has $tools_count tools configured"
        kubectl get agent test-tool-agent -n default -o jsonpath='{.spec.tools[*].name}' | tr ' ' '\n' | while read tool; do
            log_info "  🔧 Tool: $tool"
        done
    else
        log_warning "⚠️ No tools found in tool calling agent"
    fi

    # Test the agent API
    log_info "Testing tool calling agent API..."
    local service_ip=$(kubectl get svc test-tool-agent-service -n default -o jsonpath='{.spec.clusterIP}')
    if [ -n "$service_ip" ]; then
        log_success "✅ Tool calling agent service created: $service_ip"
    else
        log_error "❌ Tool calling agent service not found"
        return 1
    fi

    log_success "🎉 Tool calling agent test completed successfully"
}

# Test with vLLM service
test_vllm_integration() {
    log_info "🤖 Testing vLLM Integration"

    # Check if vLLM service is accessible
    log_info "Checking vLLM service availability..."
    if curl -s http://10.0.78.113:8000/health > /dev/null; then
        log_success "✅ vLLM service is accessible"

        # Test vLLM API
        local response=$(curl -s -X POST http://10.0.78.113:8000/v1/chat/completions \
            -H "Content-Type: application/json" \
            -d '{
                "model": "mistral-7b-instruct",
                "messages": [{"role": "user", "content": "Hello from KubeAgentic test!"}],
                "max_tokens": 50
            }')

        if echo "$response" | grep -q "choices"; then
            log_success "✅ vLLM API test successful"
        else
            log_warning "⚠️ vLLM API test failed or returned unexpected response"
        fi
    else
        log_warning "⚠️ vLLM service not accessible at 10.0.78.113:8000"
    fi
}

# Main test execution
main() {
    log_info "🚀 Starting KubeAgentic Operator Tests"
    log_info "Testing both Direct and Tool Calling workflows"

    # Check prerequisites
    if ! kubectl cluster-info &> /dev/null; then
        log_error "❌ Kubernetes cluster not accessible"
        exit 1
    fi

    if ! kubectl get crd agents.ai.example.com &> /dev/null; then
        log_error "❌ Agent CRD not found. Deploy operator first."
        exit 1
    fi

    log_success "✅ Prerequisites check passed"

    # Run tests
    local test_results=()

    # Test Direct Agent
    if test_direct_agent; then
        test_results+=("Direct Agent: ✅ PASSED")
    else
        test_results+=("Direct Agent: ❌ FAILED")
    fi

    # Test Tool Calling Agent
    if test_tool_calling_agent; then
        test_results+=("Tool Calling Agent: ✅ PASSED")
    else
        test_results+=("Tool Calling Agent: ❌ FAILED")
    fi

    # Test vLLM Integration
    test_vllm_integration

    # Print results
    echo
    log_info "📊 Test Results Summary:"
    for result in "${test_results[@]}"; do
        echo "  $result"
    done

    echo
    log_info "🔍 Operator Status:"
    kubectl get agents -A
    echo
    kubectl get pods -n kubeagentic-system

    echo
    log_success "🎉 Operator testing completed!"
}

# Cleanup function
cleanup() {
    log_info "🧹 Cleaning up test resources..."
    kubectl delete -f examples/test-direct-agent.yaml --ignore-not-found=true
    kubectl delete -f examples/test-tool-calling-agent.yaml --ignore-not-found=true
    kubectl delete secret test-gemini-secret --ignore-not-found=true
    log_success "✅ Cleanup completed"
}

# Handle script arguments
case "${1:-}" in
    --cleanup)
        cleanup
        exit 0
        ;;
    --help)
        echo "Usage: $0 [OPTIONS]"
        echo ""
        echo "Test the KubeAgentic operator with both Direct and Tool Calling workflows"
        echo ""
        echo "Options:"
        echo "  --cleanup    Clean up test resources"
        echo "  --help       Show this help"
        echo ""
        echo "Examples:"
        echo "  $0           Run all tests"
        echo "  $0 --cleanup Clean up test resources"
        exit 0
        ;;
    *)
        main
        ;;
esac
