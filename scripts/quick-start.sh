#!/bin/bash

# KubeAgentic Quick Start Script
# This script provides a streamlined way to deploy KubeAgentic for evaluation and testing.
# It checks for prerequisites, deploys the operator and its components, and optionally
# creates a sample agent to get you started.

# --- Script Configuration ---
# Exit immediately if a command exits with a non-zero status.
set -e
# Treat unset variables as an error when substituting.
set -u

echo "🚀 KubeAgentic Quick Start"
echo "========================="

# --- Prerequisite Checks ---

# Check if kubectl is installed and available in the system's PATH.
if ! command -v kubectl &> /dev/null; then
    echo "❌ Error: kubectl is not installed or not in your PATH."
    echo "Please install kubectl and ensure it's accessible to continue."
    exit 1
fi

# Verify that we can connect to a Kubernetes cluster.
if ! kubectl cluster-info &> /dev/null; then
    echo "❌ Error: Cannot connect to a Kubernetes cluster."
    echo "Please ensure your kubeconfig is set up correctly and you can access your cluster."
    exit 1
fi

echo "✅ Kubernetes cluster connection verified."

# --- KubeAgentic Deployment ---

echo "📦 Deploying KubeAgentic components..."

echo "  → Creating the 'kubeagentic-system' namespace..."
kubectl apply -f deploy/namespace.yaml

echo "  → Installing the Agent Custom Resource Definition (CRD)..."
kubectl apply -f crd/agent-crd.yaml

echo "  → Setting up Role-Based Access Control (RBAC)..."
kubectl apply -f deploy/rbac.yaml

echo "  → Deploying the KubeAgentic operator..."
kubectl apply -f deploy/operator.yaml

echo "⏳ Waiting for the KubeAgentic operator to become ready..."
# This command will wait until the 'Available' condition of the deployment is true.
kubectl wait --for=condition=Available -n kubeagentic-system deployment/kubeagentic-operator --timeout=300s

echo "✅ KubeAgentic operator is ready!"

# --- Optional: Create a Sample Agent ---

echo ""
echo "🔑 To create your first agent, you'll need an API key."
read -p "Would you like to create a sample OpenAI agent? (y/n) " -r create_agent

if [[ $create_agent =~ ^[Yy]$ ]]; then
    echo "Please enter your OpenAI API key (it will be stored in a Kubernetes secret):"
    read -s api_key
    
    if [[ -n $api_key ]]; then
        # Create a Kubernetes secret to store the API key.
        # Using --dry-run and piping to kubectl apply is a good practice to avoid errors.
        kubectl create secret generic openai-secret \
            --from-literal=api-key="$api_key" \
            --dry-run=client -o yaml | kubectl apply -f -
        
        echo "✅ Kubernetes secret 'openai-secret' created successfully."
        
        # Deploy the example OpenAI agent.
        kubectl apply -f examples/openai-agent.yaml
        
        echo "⏳ Waiting for the sample agent to become ready..."
        kubectl wait --for=condition=Ready agent/customer-support-agent --timeout=300s
        
        echo "🎉 Your sample agent is ready!"
        echo ""
        echo "To test your agent, follow these steps:"
        echo "1. Port-forward the agent's service: kubectl port-forward service/customer-support-agent-service 8080:80"
        echo "2. In a new terminal, send a request: curl -X POST http://localhost:8080/chat -H 'Content-Type: application/json' -d '{\"message\": \"Hello!\"}'"
        echo ""
        echo "To check the status of your agents, run: kubectl get agents"
    else
        echo "⚠️ No API key provided. You can create agents manually later by following the documentation."
    fi
fi

# --- Next Steps ---

echo ""
echo "📚 Next Steps:"
echo "- Explore the 'examples/' directory for more agent configurations."
echo "- Read the README.md for in-depth documentation."
echo "- Check out the 'local-testing/' directory for more advanced local testing options."
echo "- Run 'kubectl get agents' to see your deployed agents."
echo "- Use 'make status' to check the overall system status."

echo ""
echo "🎉 KubeAgentic setup is complete!"
