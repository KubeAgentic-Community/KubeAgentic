# KubeAgentic Operator

A comprehensive Kubernetes operator for deploying and managing AI agents with advanced features including autoscaling, monitoring, webhooks, and multi-framework support.

## 🚀 Features

### Core Features
- **Multi-Provider Support**: OpenAI, Anthropic Claude, Google Gemini, vLLM
- **Framework Flexibility**: Direct API calls or complex LangGraph workflows
- **Declarative Configuration**: Simple YAML-based agent definitions
- **Automatic Scaling**: Horizontal Pod Autoscaler (HPA) integration
- **Service Management**: ClusterIP, NodePort, and LoadBalancer support
- **Ingress Integration**: Automatic ingress creation for LoadBalancer services

### Advanced Features
- **Webhook Validation**: Admission webhooks for configuration validation
- **Monitoring Integration**: Prometheus metrics and Grafana dashboards
- **Resource Management**: Configurable CPU and memory limits
- **Tool Integration**: Support for custom tools and functions
- **LangGraph Workflows**: Complex multi-step reasoning workflows
- **Health Checks**: Liveness and readiness probes
- **Finalizers**: Proper cleanup on resource deletion

## 📋 Prerequisites

- Kubernetes cluster (v1.19+)
- kubectl configured to access your cluster
- Go 1.19+ (for building from source)
- Docker (for building images)

## 🐳 Docker Images

**Optimized Multi-Architecture Images on Docker Hub:**

| Component | Image | Size | Architecture | Base |
|-----------|-------|------|--------------|------|
| Operator | `sudeshmu/kubeagentic:operator-latest` | 108MB | linux/amd64, linux/arm64 | UBI Micro |
| Agent Runtime | `sudeshmu/kubeagentic:agent-latest` | 625MB | linux/amd64, linux/arm64 | UBI Minimal |

**Image Optimization Highlights:**
- **66% size reduction** for agent runtime (1.85GB → 625MB)
- Multi-stage builds for minimal final image size
- Red Hat Universal Base Images (UBI) for enterprise security
- Non-root execution with proper security contexts
- Optimized Python virtual environments

**Available Tags:**
- `operator-latest`: Latest stable operator release
- `agent-latest`: Latest optimized agent runtime
- `agent-optimized`: Explicitly tagged optimized version

```bash
# Quick verification of image sizes
docker images sudeshmu/kubeagentic
```

## 🛠️ Installation

### Option 1: Deploy with kubectl

```bash
# Deploy the operator
kubectl apply -f deploy/operator-enhanced.yaml

# Verify the installation
kubectl get pods -n kubeagentic-system
kubectl get crd agents.ai.example.com
```

### Option 2: Build and Deploy from Source

```bash
# Build the operator
make build

# Build and push the Docker image (or use pre-built images)
make docker-build docker-push

# Deploy with pre-built optimized images
kubectl apply -f deploy/operator-enhanced.yaml

# Or pull optimized images from Docker Hub
docker pull sudeshmu/kubeagentic:operator-latest  # 108MB optimized operator
docker pull sudeshmu/kubeagentic:agent-latest     # 625MB optimized agent runtime
```

## �� Usage

### 1. Create API Secret

First, create a secret with your API credentials:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: openai-secret
  namespace: default
type: Opaque
stringData:
  api-key: "your-openai-api-key-here"
```

### 2. Deploy an Agent

Create a simple agent:

```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: my-chatbot
  namespace: default
spec:
  provider: openai
  model: gpt-4
  systemPrompt: "You are a helpful AI assistant."
  apiSecretRef:
    name: openai-secret
    key: api-key
  framework: direct
  replicas: 2
  serviceType: LoadBalancer
```

### 3. Check Agent Status

```bash
# List all agents
kubectl get agents

# Get detailed information
kubectl describe agent my-chatbot

# Check agent logs
kubectl logs -l app.kubernetes.io/instance=my-chatbot
```

## 🔧 Configuration

### Agent Specification

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `provider` | string | Yes | LLM provider (openai, claude, gemini, vllm) |
| `model` | string | Yes | Specific model to use |
| `systemPrompt` | string | Yes | Agent's persona and instructions |
| `apiSecretRef` | object | Yes | Reference to API credentials secret |
| `endpoint` | string | No | Custom endpoint for self-hosted models |
| `framework` | string | No | Execution framework (direct, langgraph) |
| `langgraphConfig` | object | No | LangGraph workflow configuration |
| `tools` | array | No | Available tools for the agent |
| `replicas` | integer | No | Number of replicas (1-10) |
| `resources` | object | No | CPU and memory requirements |
| `serviceType` | string | No | Kubernetes service type |

### LangGraph Configuration

For complex workflows, use the LangGraph framework:

```yaml
spec:
  framework: langgraph
  langgraphConfig:
    graphType: sequential
    entrypoint: "start"
    endpoints: ["end"]
    nodes:
    - name: "start"
      type: "llm"
      prompt: "Analyze the user's request"
      inputs: ["user_input"]
      outputs: ["analysis"]
    - name: "search"
      type: "tool"
      tool: "web_search"
      inputs: ["analysis"]
      outputs: ["search_results"]
    edges:
    - from: "start"
      to: "search"
    - from: "search"
      to: "end"
```

### Tools Configuration

Define custom tools for your agents:

```yaml
spec:
  tools:
  - name: "calculator"
    description: "Perform mathematical calculations"
    inputSchema:
      type: "object"
      properties:
        expression:
          type: "string"
          description: "Mathematical expression to evaluate"
      required: ["expression"]
```

## 📊 Monitoring

The operator automatically creates monitoring resources:

### Prometheus Metrics

- `kubeagentic_requests_total`: Total number of requests
- `kubeagentic_response_duration_seconds`: Response time histogram
- `kubeagentic_errors_total`: Total number of errors
- `kubeagentic_active_connections`: Number of active connections

### Grafana Dashboards

Automatic Grafana dashboard creation with:
- Request rate graphs
- Response time percentiles
- Error rate monitoring
- Resource utilization

### Health Checks

- **Liveness Probe**: `/health` endpoint
- **Readiness Probe**: `/ready` endpoint
- **Metrics Endpoint**: `/metrics` for Prometheus

## 🔒 Security

### RBAC Permissions

The operator uses minimal required permissions:
- Agent CRD management
- Deployment and Service creation
- ConfigMap and Secret access
- HPA and Ingress management

### Security Context

- Non-root user execution
- Read-only root filesystem
- Dropped capabilities
- Security context constraints

### Webhook Validation

Admission webhooks validate:
- Provider and model compatibility
- Required field presence
- Resource limits validation
- Framework configuration

## 🚀 Advanced Features

### Autoscaling

Automatic HPA creation based on:
- CPU utilization (70% threshold)
- Memory utilization (80% threshold)
- Custom metrics support

### Ingress Integration

Automatic ingress creation for LoadBalancer services:
- Nginx ingress controller support
- SSL/TLS configuration
- Custom hostname support

### Multi-Framework Support

- **Direct Framework**: Simple API calls for basic interactions
- **LangGraph Framework**: Complex workflows with state management

### Tool Integration

- Custom tool definitions
- JSON schema validation
- Tool chaining and composition
- External service integration

## 🐛 Troubleshooting

### Common Issues

1. **Agent not starting**
   ```bash
   kubectl describe agent <agent-name>
   kubectl logs -l app.kubernetes.io/instance=<agent-name>
   ```

2. **Secret not found**
   ```bash
   kubectl get secrets
   kubectl describe secret <secret-name>
   ```

3. **Webhook issues**
   ```bash
   kubectl get validatingwebhookconfigurations
   kubectl describe validatingwebhookconfiguration <webhook-name>
   ```

### Debug Mode

Enable debug logging:

```bash
kubectl set env deployment/kubeagentic-operator -n kubeagentic-system LOG_LEVEL=debug
```

## 🔄 Updates and Maintenance

### Updating the Operator

```bash
# Update the operator image (using optimized Docker Hub image)
kubectl set image deployment/kubeagentic-operator -n kubeagentic-system operator=sudeshmu/kubeagentic:operator-latest

# Restart the operator
kubectl rollout restart deployment/kubeagentic-operator -n kubeagentic-system
```

### Backup and Recovery

```bash
# Backup agent configurations
kubectl get agents -o yaml > agents-backup.yaml

# Restore from backup
kubectl apply -f agents-backup.yaml
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

- **Documentation**: [https://kubeagentic.com](https://kubeagentic.com)
- **Issues**: [GitHub Issues](https://github.com/KubeAgentic-Community/KubeAgentic/issues)
- **Discussions**: [GitHub Discussions](https://github.com/KubeAgentic-Community/KubeAgentic/discussions)
- **Email**: contact@kubeagentic.com
