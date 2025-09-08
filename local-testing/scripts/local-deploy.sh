#!/bin/bash

# Local Deployment Script for KubeAgentic
# Deploys KubeAgentic to a local Kubernetes cluster (kind, minikube, k3d)

set -e

echo "🚀 KubeAgentic Local Deployment"
echo "==============================="

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl is not installed or not in PATH"
    exit 1
fi

# Check if we can connect to cluster
if ! kubectl cluster-info &> /dev/null; then
    echo "❌ Cannot connect to Kubernetes cluster"
    echo "💡 Make sure you have a local cluster running:"
    echo "   kind create cluster --name kubeagentic-test"
    echo "   # OR"
    echo "   minikube start"
    echo "   # OR" 
    echo "   k3d cluster create kubeagentic-test"
    exit 1
fi

echo "✅ Kubernetes cluster connection verified"

# Build images for local registry
echo "🔨 Building container images..."

# For kind, we need to load images
if kubectl config current-context | grep -q "kind"; then
    echo "  → Detected kind cluster, building and loading images..."
    
    # Build operator image
    docker build -f Dockerfile.operator -t kubeagentic/operator:local .
    kind load docker-image kubeagentic/operator:local
    
    # Build agent image  
    docker build -f Dockerfile.agent -t kubeagentic/agent:local .
    kind load docker-image kubeagentic/agent:local
    
    # Update image references for local testing
    sed -i.bak 's|sudeshmu/kubeagentic:operator-latest|kubeagentic/operator:local|g' deploy/operator.yaml
    sed -i.bak 's|sudeshmu/kubeagentic:agent-latest|kubeagentic/agent:local|g' controllers/agent_controller.go

# For minikube, point docker to minikube's docker daemon
elif kubectl config current-context | grep -q "minikube"; then
    echo "  → Detected minikube cluster, configuring docker daemon..."
    eval $(minikube docker-env)
    
    # Build images in minikube's docker
    docker build -f Dockerfile.operator -t kubeagentic/operator:local .
    docker build -f Dockerfile.agent -t kubeagentic/agent:local .
    
    # Update image references
    sed -i.bak 's|sudeshmu/kubeagentic:operator-latest|kubeagentic/operator:local|g' deploy/operator.yaml
    sed -i.bak 's|sudeshmu/kubeagentic:agent-latest|kubeagentic/agent:local|g' controllers/agent_controller.go

# For other clusters (like k3d), try direct build
else
    echo "  → Building images for local registry..."
    docker build -f Dockerfile.operator -t kubeagentic/operator:local .
    docker build -f Dockerfile.agent -t kubeagentic/agent:local .
    
    # Update image references
    sed -i.bak 's|sudeshmu/kubeagentic:operator-latest|kubeagentic/operator:local|g' deploy/operator.yaml
    sed -i.bak 's|sudeshmu/kubeagentic:agent-latest|kubeagentic/agent:local|g' controllers/agent_controller.go
fi

# Deploy the system
echo "📦 Deploying KubeAgentic..."

echo "  → Creating namespace..."
kubectl apply -f deploy/namespace.yaml

echo "  → Installing CRD..."
kubectl apply -f crd/agent-crd.yaml

echo "  → Setting up RBAC..."
kubectl apply -f deploy/rbac.yaml

echo "  → Deploying operator..."
# Use local image tag
sed 's|sudeshmu/kubeagentic:operator-latest|kubeagentic/operator:local|g' deploy/operator.yaml | kubectl apply -f -

echo "⏳ Waiting for operator to be ready..."
kubectl wait --for=condition=Available -n kubeagentic-system deployment/kubeagentic-operator --timeout=300s

echo "✅ KubeAgentic operator is ready!"

# Restore original files
if [ -f deploy/operator.yaml.bak ]; then
    mv deploy/operator.yaml.bak deploy/operator.yaml
fi

# Prompt for test deployment
echo ""
echo "🧪 Would you like to deploy a test agent? (y/n)"
read -r deploy_test

if [[ $deploy_test =~ ^[Yy]$ ]]; then
    echo "Choose a test agent:"
    echo "1) OpenAI (requires API key)"
    echo "2) Mock vLLM (no API key needed)"
    echo "3) Both"
    read -r choice
    
    case $choice in
        1)
            echo "Enter your OpenAI API key:"
            read -s api_key
            if [[ -n $api_key ]]; then
                kubectl create secret generic openai-secret \
                    --from-literal=api-key="$api_key" \
                    --dry-run=client -o yaml | kubectl apply -f -
                
                # Use local agent image
                sed 's|sudeshmu/kubeagentic:agent-latest|kubeagentic/agent:local|g' examples/openai-agent.yaml | kubectl apply -f -
                echo "✅ OpenAI agent deployed!"
            fi
            ;;
        2)
            # Deploy mock vLLM server first
            kubectl create deployment mock-vllm \
                --image=kubeagentic/mock-vllm:local \
                --port=8000
            kubectl expose deployment mock-vllm --port=8000 --target-port=8000
            
            # Create mock secret
            kubectl create secret generic vllm-secret \
                --from-literal=api-key="mock-api-key" \
                --dry-run=client -o yaml | kubectl apply -f -
            
            # Deploy vLLM agent
            sed 's|sudeshmu/kubeagentic:agent-latest|kubeagentic/agent:local|g' examples/vllm-agent.yaml | kubectl apply -f -
            echo "✅ Mock vLLM agent deployed!"
            ;;
        3)
            echo "Enter your OpenAI API key:"
            read -s api_key
            if [[ -n $api_key ]]; then
                kubectl create secret generic openai-secret \
                    --from-literal=api-key="$api_key" \
                    --dry-run=client -o yaml | kubectl apply -f -
                
                sed 's|sudeshmu/kubeagentic:agent-latest|kubeagentic/agent:local|g' examples/openai-agent.yaml | kubectl apply -f -
            fi
            
            # Deploy mock vLLM
            kubectl create deployment mock-vllm \
                --image=kubeagentic/mock-vllm:local \
                --port=8000
            kubectl expose deployment mock-vllm --port=8000 --target-port=8000
            
            kubectl create secret generic vllm-secret \
                --from-literal=api-key="mock-api-key" \
                --dry-run=client -o yaml | kubectl apply -f -
            
            sed 's|sudeshmu/kubeagentic:agent-latest|kubeagentic/agent:local|g' examples/vllm-agent.yaml | kubectl apply -f -
            echo "✅ Both agents deployed!"
            ;;
    esac
    
    echo ""
    echo "⏳ Waiting for agents to be ready..."
    kubectl wait --for=condition=Ready agents --all --timeout=300s || true
    
    echo ""
    echo "🎉 Deployment complete!"
    echo ""
    echo "To test your agents:"
    echo "1. Check status: kubectl get agents"
    echo "2. Port forward: kubectl port-forward service/AGENT-NAME-service 8080:80"
    echo "3. Test: curl -X POST http://localhost:8080/chat -H 'Content-Type: application/json' -d '{\"message\": \"Hello!\"}'"
fi

echo ""
echo "📚 Useful commands:"
echo "- kubectl get agents                     # List all agents"
echo "- kubectl describe agent AGENT-NAME     # Get agent details"
echo "- kubectl logs -n kubeagentic-system deployment/kubeagentic-operator  # Operator logs"
echo "- kubectl delete agents --all           # Clean up agents"
echo ""
echo "🎉 KubeAgentic local deployment complete!"
