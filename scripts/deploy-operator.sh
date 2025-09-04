#!/bin/bash

# KubeAgentic Operator Deployment Script
# This script deploys the enhanced KubeAgentic operator with all features

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
OPERATOR_NAMESPACE="kubeagentic-system"
OPERATOR_IMAGE="kubeagentic/operator:latest"
OPERATOR_MANIFEST="deploy/operator-enhanced.yaml"

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check if kubectl is installed
    if ! command -v kubectl &> /dev/null; then
        log_error "kubectl is not installed. Please install kubectl first."
        exit 1
    fi
    
    # Check if kubectl can connect to cluster
    if ! kubectl cluster-info &> /dev/null; then
        log_error "Cannot connect to Kubernetes cluster. Please check your kubeconfig."
        exit 1
    fi
    
    # Check if the manifest file exists
    if [ ! -f "$OPERATOR_MANIFEST" ]; then
        log_error "Operator manifest file not found: $OPERATOR_MANIFEST"
        exit 1
    fi
    
    log_success "Prerequisites check passed"
}

deploy_operator() {
    log_info "Deploying KubeAgentic operator..."
    
    # Apply the operator manifest
    kubectl apply -f "$OPERATOR_MANIFEST"
    
    if [ $? -eq 0 ]; then
        log_success "Operator manifest applied successfully"
    else
        log_error "Failed to apply operator manifest"
        exit 1
    fi
}

wait_for_operator() {
    log_info "Waiting for operator to be ready..."
    
    # Wait for the operator deployment to be ready
    kubectl wait --for=condition=available --timeout=300s deployment/kubeagentic-operator -n "$OPERATOR_NAMESPACE"
    
    if [ $? -eq 0 ]; then
        log_success "Operator is ready"
    else
        log_error "Operator failed to become ready within 5 minutes"
        exit 1
    fi
}

verify_deployment() {
    log_info "Verifying deployment..."
    
    # Check if the operator pod is running
    if kubectl get pods -n "$OPERATOR_NAMESPACE" -l app.kubernetes.io/component=operator | grep -q "Running"; then
        log_success "Operator pod is running"
    else
        log_error "Operator pod is not running"
        exit 1
    fi
    
    # Check if the CRD is installed
    if kubectl get crd agents.ai.example.com &> /dev/null; then
        log_success "Agent CRD is installed"
    else
        log_error "Agent CRD is not installed"
        exit 1
    fi
    
    # Check if the webhook is configured
    if kubectl get validatingwebhookconfigurations | grep -q "kubeagentic"; then
        log_success "Webhook is configured"
    else
        log_warning "Webhook is not configured (this is optional)"
    fi
}

show_status() {
    log_info "Deployment status:"
    echo ""
    
    echo "Operator Pods:"
    kubectl get pods -n "$OPERATOR_NAMESPACE" -l app.kubernetes.io/component=operator
    echo ""
    
    echo "Custom Resource Definitions:"
    kubectl get crd | grep "ai.example.com"
    echo ""
    
    echo "Operator Logs (last 10 lines):"
    kubectl logs -n "$OPERATOR_NAMESPACE" -l app.kubernetes.io/component=operator --tail=10
}

deploy_test_agent() {
    log_info "Deploying test agent..."
    
    # Create a test secret
    kubectl create secret generic test-secret --from-literal=api-key=test-key --dry-run=client -o yaml | kubectl apply -f -
    
    # Deploy test agent
    kubectl apply -f examples/enhanced-agent-example.yaml
    
    if [ $? -eq 0 ]; then
        log_success "Test agent deployed successfully"
    else
        log_warning "Failed to deploy test agent (this is optional)"
    fi
}

cleanup() {
    log_info "Cleaning up..."
    
    # Delete test agent
    kubectl delete -f examples/enhanced-agent-example.yaml --ignore-not-found=true
    
    # Delete test secret
    kubectl delete secret test-secret --ignore-not-found=true
    
    log_success "Cleanup completed"
}

show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help          Show this help message"
    echo "  -d, --deploy        Deploy the operator"
    echo "  -s, --status        Show operator status"
    echo "  -t, --test          Deploy test agent"
    echo "  -c, --cleanup       Clean up test resources"
    echo "  -u, --undeploy      Undeploy the operator"
    echo "  -f, --full          Full deployment with test agent"
    echo ""
    echo "Examples:"
    echo "  $0 --deploy         Deploy the operator"
    echo "  $0 --status         Show current status"
    echo "  $0 --full           Deploy operator and test agent"
    echo "  $0 --cleanup        Clean up test resources"
}

# Main script logic
case "${1:-}" in
    -h|--help)
        show_usage
        exit 0
        ;;
    -d|--deploy)
        check_prerequisites
        deploy_operator
        wait_for_operator
        verify_deployment
        show_status
        ;;
    -s|--status)
        show_status
        ;;
    -t|--test)
        deploy_test_agent
        ;;
    -c|--cleanup)
        cleanup
        ;;
    -u|--undeploy)
        log_info "Undeploying operator..."
        kubectl delete -f "$OPERATOR_MANIFEST" --ignore-not-found=true
        log_success "Operator undeployed"
        ;;
    -f|--full)
        check_prerequisites
        deploy_operator
        wait_for_operator
        verify_deployment
        deploy_test_agent
        show_status
        ;;
    *)
        log_error "Invalid option: ${1:-}"
        echo ""
        show_usage
        exit 1
        ;;
esac
