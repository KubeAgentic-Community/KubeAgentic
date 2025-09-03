---
layout: page
title: Documentation
permalink: /docs/
---

# KubeAgentic Documentation

Welcome to the comprehensive documentation for KubeAgentic, the Kubernetes operator for managing AI agents.

## 📚 Table of Contents

### Getting Started
- [Quick Start Guide](#quick-start-guide)
- [Installation](#installation)
- [Your First Agent](#your-first-agent)

### Configuration
- [Agent Configuration](#agent-configuration)
- [Providers](#providers)
- [Security](#security)
- [Scaling](#scaling)

### Advanced Topics
- [Monitoring & Observability](#monitoring--observability)
- [Troubleshooting](#troubleshooting)
- [Best Practices](#best-practices)

---

## Quick Start Guide

### Prerequisites

Before you begin, ensure you have:
- A running Kubernetes cluster (v1.19+)
- `kubectl` installed and configured
- Basic understanding of Kubernetes concepts

### Installation

Install KubeAgentic using one of these methods:

**Option 1: Direct Installation**
```bash
kubectl apply -f https://raw.githubusercontent.com/sudeshmu/kubeagentic/main/deploy/all.yaml
```

**Option 2: Local Installation**
```bash
git clone https://github.com/sudeshmu/kubeagentic.git
cd kubeagentic
kubectl apply -f deploy/all.yaml
```

**Option 3: Helm Chart** (Coming Soon)
```bash
helm repo add kubeagentic https://your-username.github.io/kubeagentic
helm install kubeagentic kubeagentic/kubeagentic
```

### Verify Installation

```bash
# Check if the operator is running
kubectl get pods -n kubeagentic-system

# Verify CRDs are installed
kubectl get crd agents.ai.example.com
```

---

## Your First Agent

### 1. Create API Key Secret

```bash
# For OpenAI
kubectl create secret generic openai-secret \
  --from-literal=api-key='sk-your-openai-api-key'

# For Anthropic Claude
kubectl create secret generic claude-secret \
  --from-literal=api-key='sk-ant-your-claude-api-key'

# For Google Gemini
kubectl create secret generic gemini-secret \
  --from-literal=api-key='your-gemini-api-key'
```

### 2. Deploy a Simple Agent

Create `simple-agent.yaml`:

```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: helpful-assistant
spec:
  provider: openai
  model: gpt-4
  systemPrompt: "You are a helpful assistant that answers questions concisely."
  apiSecretRef:
    name: openai-secret
    key: api-key
  replicas: 1
  resources:
    requests:
      memory: "128Mi"
      cpu: "100m"
    limits:
      memory: "256Mi"
      cpu: "200m"
```

Deploy it:
```bash
kubectl apply -f simple-agent.yaml
```

### 3. Interact with Your Agent

```bash
# Port forward to access the agent
kubectl port-forward service/helpful-assistant-service 8080:80

# Send a message (in a new terminal)
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "What is Kubernetes?"}'
```

---

## Agent Configuration

### Basic Configuration

| Field | Description | Required | Default |
|-------|-------------|----------|---------|
| `provider` | AI provider (openai, claude, gemini, vllm) | Yes | - |
| `model` | Model name | Yes | - |
| `systemPrompt` | System prompt for the agent | No | "" |
| `apiSecretRef` | Reference to API key secret | Yes | - |
| `replicas` | Number of agent instances | No | 1 |
| `resources` | CPU/Memory requests and limits | No | See defaults |

### Advanced Configuration

```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: advanced-agent
spec:
  provider: openai
  model: gpt-4
  systemPrompt: |
    You are an expert software engineer.
    Follow these guidelines:
    - Write clean, well-documented code
    - Consider security implications
    - Optimize for performance
  
  # API Configuration
  apiSecretRef:
    name: openai-secret
    key: api-key
  
  # Scaling Configuration
  replicas: 3
  autoscaling:
    enabled: true
    minReplicas: 1
    maxReplicas: 10
    targetCPUUtilizationPercentage: 70
  
  # Resource Management
  resources:
    requests:
      memory: "256Mi"
      cpu: "200m"
    limits:
      memory: "512Mi"
      cpu: "500m"
  
  # Service Configuration
  serviceType: ClusterIP
  servicePort: 80
  
  # Health Checks
  healthCheck:
    enabled: true
    path: "/health"
    intervalSeconds: 30
    timeoutSeconds: 5
    
  # Environment Variables
  env:
    - name: LOG_LEVEL
      value: "INFO"
    - name: MAX_TOKENS
      value: "2048"
      
  # Tool Integration
  tools:
    - name: web-search
      endpoint: "http://search-service:8080"
    - name: database
      endpoint: "postgresql://db:5432/agents"
      secretRef:
        name: db-credentials
        key: connection-string
```

---

## Providers

### OpenAI

```yaml
spec:
  provider: openai
  model: gpt-4  # or gpt-3.5-turbo, gpt-4-turbo
  apiSecretRef:
    name: openai-secret
    key: api-key
```

**Supported Models:**
- `gpt-4`
- `gpt-4-turbo`
- `gpt-3.5-turbo`

### Anthropic Claude

```yaml
spec:
  provider: claude
  model: claude-3-sonnet-20240229
  apiSecretRef:
    name: claude-secret
    key: api-key
```

**Supported Models:**
- `claude-3-opus-20240229`
- `claude-3-sonnet-20240229`
- `claude-3-haiku-20240307`

### Google Gemini

```yaml
spec:
  provider: gemini
  model: gemini-pro
  apiSecretRef:
    name: gemini-secret
    key: api-key
```

**Supported Models:**
- `gemini-pro`
- `gemini-pro-vision`

### Self-hosted vLLM

```yaml
spec:
  provider: vllm
  model: llama2-7b-chat
  endpoint: http://vllm-server:8000/v1
  apiSecretRef:
    name: vllm-secret
    key: api-key  # Optional for public endpoints
```

---

## Security

### API Key Management

Always store API keys in Kubernetes Secrets:

```bash
# Create secret
kubectl create secret generic my-ai-secret \
  --from-literal=api-key='your-secret-key'

# Reference in agent spec
spec:
  apiSecretRef:
    name: my-ai-secret
    key: api-key
```

### Network Policies

Restrict agent network access:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: agent-network-policy
spec:
  podSelector:
    matchLabels:
      kubeagentic.ai/agent: my-agent
  policyTypes:
  - Egress
  egress:
  - to:
    - namespaceSelector:
        matchLabels:
          name: kube-system
    ports:
    - protocol: TCP
      port: 443  # HTTPS only
```

### RBAC

Configure fine-grained permissions:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: agent-reader
rules:
- apiGroups: [""]
  resources: ["configmaps", "secrets"]
  verbs: ["get", "list"]
- apiGroups: ["ai.example.com"]
  resources: ["agents"]
  verbs: ["get", "list", "watch"]
```

---

## Scaling

### Manual Scaling

```bash
# Scale to 5 replicas
kubectl patch agent my-agent --type='merge' -p='{"spec":{"replicas":5}}'
```

### Automatic Scaling

```yaml
spec:
  autoscaling:
    enabled: true
    minReplicas: 2
    maxReplicas: 20
    metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80
```

---

## Monitoring & Observability

### Health Checks

```bash
# Check agent status
kubectl get agents

# Detailed agent information
kubectl describe agent my-agent

# Agent logs
kubectl logs -l kubeagentic.ai/agent=my-agent

# Resource usage
kubectl top pods -l kubeagentic.ai/agent=my-agent
```

### Metrics

KubeAgentic exposes Prometheus metrics:

```yaml
# Enable metrics in agent spec
spec:
  metrics:
    enabled: true
    port: 9090
    path: /metrics
```

**Key Metrics:**
- `kubeagentic_requests_total` - Total API requests
- `kubeagentic_request_duration_seconds` - Request latency
- `kubeagentic_errors_total` - Error count
- `kubeagentic_tokens_used_total` - Token usage

---

## Troubleshooting

### Common Issues

**Agent Stuck in Pending**
```bash
# Check events
kubectl describe agent my-agent

# Check pod status
kubectl get pods -l kubeagentic.ai/agent=my-agent

# Check operator logs
kubectl logs -n kubeagentic-system deployment/kubeagentic-operator
```

**API Connection Errors**
- Verify API key is correct
- Check network connectivity
- Ensure sufficient API quota

**Resource Issues**
```bash
# Check node resources
kubectl describe nodes

# Adjust resource requests
kubectl patch agent my-agent --type='merge' -p='{"spec":{"resources":{"requests":{"memory":"64Mi","cpu":"50m"}}}}'
```

### Debug Mode

Enable debug logging:

```yaml
spec:
  env:
  - name: LOG_LEVEL
    value: "DEBUG"
```

---

## Best Practices

### Performance
- Set appropriate resource limits
- Use autoscaling for variable workloads
- Monitor token usage and costs
- Implement request caching where possible

### Security
- Always use Secrets for API keys
- Apply network policies
- Run containers as non-root users
- Regularly rotate API keys

### Operations
- Use GitOps for agent configurations
- Implement proper monitoring and alerting
- Test agents in staging before production
- Keep the operator updated

### Cost Optimization
- Monitor API usage and costs
- Set token limits per request
- Use smaller models for simple tasks
- Implement request throttling

---

For more detailed information, see our [API Reference](../api-reference.html) and [Examples](../examples.html).