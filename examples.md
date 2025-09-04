---
layout: page
title: Examples
description: Real-world examples and templates for deploying AI agents with KubeAgenticpermalink: /examples/
---

# KubeAgentic Examples

This page showcases real-world examples of AI agents deployed with KubeAgentic. Each example includes complete YAML configurations and explains the use case, benefits, and customization options.

## 🎯 Quick Examples

### Basic Assistant

The simplest possible agent configuration:

```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: simple-assistant
spec:
  provider: openai
  model: gpt-3.5-turbo
  systemPrompt: "You are a helpful assistant."
  apiSecretRef:
    name: openai-secret
    key: api-key
```

Deploy it:
```bash
kubectl apply -f - <<EOF
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: simple-assistant
spec:
  provider: openai
  model: gpt-3.5-turbo
  systemPrompt: "You are a helpful assistant."
  apiSecretRef:
    name: openai-secret
    key: api-key
EOF
```

---

## 🛎️ Customer Support Agent

**Use Case**: Automated customer support with order lookup and inventory checking capabilities.

**Features**:
- Multi-replica deployment for high availability
- Tool integration for order and inventory systems
- Professional customer service training

```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: customer-support-agent
  namespace: default
spec:
  provider: "openai"
  model: "gpt-4"
  systemPrompt: |
    You are a helpful customer support agent for an e-commerce company.
    You should be friendly, professional, and always try to resolve customer issues.
    If you cannot resolve an issue, escalate it to a human agent.
    Always ask for order numbers when dealing with order-related issues.
  apiSecretRef:
    name: openai-secret
    key: api-key
  replicas: 2
  resources:
    requests:
      memory: "256Mi"
      cpu: "100m"
    limits:
      memory: "512Mi"
      cpu: "200m"
  serviceType: "ClusterIP"
  tools:
  - name: "order_lookup"
    description: "Look up order information by order ID"
    inputSchema:
      type: "object"
      properties:
        order_id:
          type: "string"
          description: "The order ID to look up"
      required: ["order_id"]
  - name: "inventory_check"
    description: "Check if a product is in stock"
    inputSchema:
      type: "object"
      properties:
        product_id:
          type: "string"
          description: "The product ID to check"
      required: ["product_id"]
```

**Deployment**:
```bash
# Create the API secret first
kubectl create secret generic openai-secret \
  --from-literal=api-key='your-openai-api-key'

# Deploy the agent
kubectl apply -f examples/openai-agent.yaml

# Access the agent
kubectl port-forward service/customer-support-agent-service 8080:80
```

**Testing**:
```bash
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "I need help with my order #12345"}'
```

---

## 🔍 Code Review Agent

**Use Case**: Automated code review for pull requests with security and performance analysis.

**Features**:
- Powered by Claude for advanced reasoning
- Static code analysis integration
- Focus on security, performance, and best practices

```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: code-review-agent
  namespace: default
spec:
  provider: "claude"
  model: "claude-3-sonnet-20240229"
  systemPrompt: |
    You are an expert code reviewer specializing in security, performance, and best practices.
    Review code submissions and provide detailed feedback including:
    - Security vulnerabilities
    - Performance improvements
    - Code style and maintainability
    - Best practice recommendations
    Always be constructive and educational in your feedback.
  apiSecretRef:
    name: claude-secret
    key: api-key
  replicas: 1
  resources:
    requests:
      memory: "512Mi"
      cpu: "200m"
    limits:
      memory: "1Gi"
      cpu: "500m"
  serviceType: "ClusterIP"
  tools:
  - name: "static_analysis"
    description: "Run static code analysis on submitted code"
    inputSchema:
      type: "object"
      properties:
        code:
          type: "string"
          description: "The code to analyze"
        language:
          type: "string"
          description: "Programming language"
      required: ["code", "language"]
```

**Deployment**:
```bash
# Create Claude API secret
kubectl create secret generic claude-secret \
  --from-literal=api-key='sk-ant-your-claude-api-key'

# Deploy the agent
kubectl apply -f examples/claude-agent.yaml
```

**Usage with GitHub Actions**:
```yaml
# .github/workflows/code-review.yml
name: AI Code Review
on: [pull_request]
jobs:
  review:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Review Code
      run: |
        curl -X POST http://code-review-agent.default.svc.cluster.local/review \
          -H "Content-Type: application/json" \
          -d '{"code": "'$(cat changed-files.txt)'", "language": "python"}'
```

---

## 🏢 Internal Knowledge Assistant

**Use Case**: Self-hosted AI assistant for internal company knowledge using vLLM.

**Features**:
- Self-hosted model (cost-effective)
- LoadBalancer service for employee access
- Integration with company systems
- Scalable for organization size

```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: internal-qa-agent
  namespace: default
spec:
  provider: "vllm"
  model: "llama2-7b-chat"
  endpoint: "http://vllm-service.ml-inference.svc.cluster.local:8000/v1"
  systemPrompt: |
    You are an internal Q&A assistant for company employees.
    You have access to company documentation and policies.
    Provide accurate, helpful answers based on internal knowledge.
    If you don't know something, be honest and direct users to appropriate resources.
  apiSecretRef:
    name: vllm-secret
    key: api-key
  replicas: 3
  resources:
    requests:
      memory: "128Mi"
      cpu: "50m"
    limits:
      memory: "256Mi"
      cpu: "100m"
  serviceType: "LoadBalancer"  # Expose externally for employee access
  tools:
  - name: "policy_search"
    description: "Search company policies and documentation"
    inputSchema:
      type: "object"
      properties:
        query:
          type: "string"
          description: "Search query for policies"
      required: ["query"]
  - name: "employee_directory"
    description: "Look up employee contact information"
    inputSchema:
      type: "object"
      properties:
        name:
          type: "string"
          description: "Employee name to search for"
      required: ["name"]
```

**Setup with vLLM**:
```bash
# Deploy vLLM server first
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: vllm-server
  namespace: ml-inference
spec:
  replicas: 1
  selector:
    matchLabels:
      app: vllm-server
  template:
    metadata:
      labels:
        app: vllm-server
    spec:
      containers:
      - name: vllm
        image: vllm/vllm-openai:latest
        args: ["--model", "meta-llama/Llama-2-7b-chat-hf", "--served-model-name", "llama2-7b-chat"]
        ports:
        - containerPort: 8000
        resources:
          requests:
            nvidia.com/gpu: 1
EOF

# Create service
kubectl expose deployment vllm-server --port=8000 --name=vllm-service -n ml-inference

# Deploy the agent
kubectl apply -f examples/vllm-agent.yaml
```

---

## 🔄 Auto-scaling Examples

### High Traffic Agent

For agents that need to handle variable load:

```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: high-traffic-agent
spec:
  provider: openai
  model: gpt-3.5-turbo
  systemPrompt: "You are a fast, efficient assistant."
  apiSecretRef:
    name: openai-secret
    key: api-key
  
  # Auto-scaling configuration
  replicas: 2  # Initial replicas
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
    behavior:
      scaleUp:
        stabilizationWindowSeconds: 60
        policies:
        - type: Percent
          value: 100
          periodSeconds: 15
      scaleDown:
        stabilizationWindowSeconds: 300
        policies:
        - type: Percent
          value: 10
          periodSeconds: 60
```

---

## 🛡️ Security-focused Examples

### Restricted Network Agent

Agent with network policies and security constraints:

```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: secure-agent
spec:
  provider: openai
  model: gpt-4
  systemPrompt: "You are a security-conscious assistant."
  apiSecretRef:
    name: openai-secret
    key: api-key
  
  # Security settings
  securityContext:
    runAsNonRoot: true
    runAsUser: 65534
    readOnlyRootFilesystem: true
    capabilities:
      drop:
      - ALL
  
  # Resource constraints
  resources:
    requests:
      memory: "64Mi"
      cpu: "50m"
    limits:
      memory: "128Mi"
      cpu: "100m"
  
  # Environment restrictions
  env:
  - name: LOG_LEVEL
    value: "WARN"  # Minimal logging
  - name: MAX_TOKENS
    value: "1000"  # Token limit
---
# Network Policy
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: secure-agent-netpol
spec:
  podSelector:
    matchLabels:
      kubeagentic.ai/agent: secure-agent
  policyTypes:
  - Egress
  egress:
  - to: []
    ports:
    - protocol: TCP
      port: 443  # HTTPS only
    - protocol: TCP
      port: 53   # DNS
    - protocol: UDP
      port: 53   # DNS
```

---

## 🎨 Multi-Model Examples

### Model Comparison Agent

Deploy multiple models for A/B testing:

```yaml
# GPT-4 version
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: assistant-gpt4
  labels:
    version: gpt4
    experiment: model-comparison
spec:
  provider: openai
  model: gpt-4
  systemPrompt: "You are a helpful assistant (GPT-4)."
  apiSecretRef:
    name: openai-secret
    key: api-key
  replicas: 1
---
# Claude version
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: assistant-claude
  labels:
    version: claude
    experiment: model-comparison
spec:
  provider: claude
  model: claude-3-sonnet-20240229
  systemPrompt: "You are a helpful assistant (Claude)."
  apiSecretRef:
    name: claude-secret
    key: api-key
  replicas: 1
---
# Load balancer service
apiVersion: v1
kind: Service
metadata:
  name: assistant-comparison
spec:
  selector:
    experiment: model-comparison
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

---

## 📊 Monitoring Examples

### Agent with Full Observability

```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: monitored-agent
spec:
  provider: openai
  model: gpt-3.5-turbo
  systemPrompt: "You are a monitored assistant."
  apiSecretRef:
    name: openai-secret
    key: api-key
  
  # Monitoring configuration
  metrics:
    enabled: true
    port: 9090
    path: /metrics
  
  # Health checks
  healthCheck:
    enabled: true
    path: "/health"
    initialDelaySeconds: 30
    periodSeconds: 10
    timeoutSeconds: 5
    failureThreshold: 3
  
  # Logging
  logging:
    level: INFO
    format: json
    outputs:
    - console
    - file
  
  # Tracing
  tracing:
    enabled: true
    endpoint: "http://jaeger-collector:14268/api/traces"
    samplingRate: 0.1
```

---

## 🔧 Development Examples

### Development Agent

Agent configured for local development:

```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: dev-agent
spec:
  provider: openai
  model: gpt-3.5-turbo  # Cheaper for development
  systemPrompt: "You are a development assistant."
  apiSecretRef:
    name: openai-secret
    key: api-key
  
  # Development settings
  replicas: 1
  resources:
    requests:
      memory: "64Mi"
      cpu: "25m"
    limits:
      memory: "128Mi"
      cpu: "100m"
  
  # Debug settings
  env:
  - name: LOG_LEVEL
    value: "DEBUG"
  - name: DEV_MODE
    value: "true"
  
  # Quick iteration
  serviceType: NodePort
  
  # Development tools
  tools:
  - name: "code_generator"
    description: "Generate development code"
    endpoint: "http://localhost:3000/generate"
```

---

## 🚀 Getting Started

To use any of these examples:

1. **Create API Key Secrets**:
   ```bash
   kubectl create secret generic openai-secret --from-literal=api-key='your-key'
   kubectl create secret generic claude-secret --from-literal=api-key='your-key'
   ```

2. **Deploy an Example**:
   ```bash
   kubectl apply -f examples/openai-agent.yaml
   ```

3. **Test the Agent**:
   ```bash
   kubectl port-forward service/customer-support-agent-service 8080:80
   curl -X POST http://localhost:8080/chat \
     -H "Content-Type: application/json" \
     -d '{"message": "Hello!"}'
   ```

4. **Monitor Status**:
   ```bash
   kubectl get agents
   kubectl describe agent customer-support-agent
   ```

## 📚 More Resources

- [Documentation](docs/) - Complete setup guide
- [API Reference](api-reference.html) - Detailed API specification  
- [Local Testing](local-testing.html) - Development environment setup
- [GitHub Repository](https://github.com/sudeshmu/kubeagentic) - Source code and issues

Ready to deploy your own AI agents? Start with the [Quick Start Guide](docs/#quick-start-guide)!