---
layout: page
title: API Reference
description: Complete API specification for Agent Custom Resources and componentspermalink: /api-reference/
---

# KubeAgentic API Reference

Complete API specification for the `Agent` Custom Resource and related components.

## Table of Contents

- [Agent Resource](#agent-resource)
- [Spec Fields](#spec-fields)
- [Status Fields](#status-fields)
- [Examples](#examples)
- [Field Validation](#field-validation)

---

## Agent Resource

The `Agent` resource is the core component of KubeAgentic, defining an AI agent's configuration and desired state.

### API Version

- **API Version**: `ai.example.com/v1`
- **Kind**: `Agent`
- **Scope**: Namespaced

### Basic Structure

```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: agent-name
  namespace: default
  labels:
    app: my-agent
    version: v1.0.0
spec:
  # Agent specification (see below)
status:
  # Agent status (managed by operator)
```

---

## Spec Fields

### Required Fields

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| `provider` | string | LLM provider | `openai`, `claude`, `gemini`, `vllm` |
| `model` | string | Model name | `gpt-4`, `claude-3-sonnet` |
| `systemPrompt` | string | Agent instructions | `"You are a helpful assistant"` |
| `apiSecretRef` | object | API key secret reference | See [apiSecretRef](#apisecretref) |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `replicas` | integer | `1` | Number of agent instances |
| `resources` | object | See defaults | CPU/memory requests and limits |
| `serviceType` | string | `ClusterIP` | Kubernetes service type |
| `servicePort` | integer | `80` | Service port |
| `env` | array | `[]` | Environment variables |
| `tools` | array | `[]` | Agent tools configuration |
| `autoscaling` | object | `nil` | Horizontal Pod Autoscaler config |
| `healthCheck` | object | See defaults | Health check configuration |
| `metrics` | object | `nil` | Metrics configuration |
| `securityContext` | object | See defaults | Pod security context |

---

### Detailed Field Specifications

#### provider

Specifies the LLM provider for the agent.

**Type**: `string`  
**Required**: Yes  
**Allowed Values**: 
- `openai` - OpenAI GPT models
- `claude` - Anthropic Claude models  
- `gemini` - Google Gemini models
- `vllm` - Self-hosted vLLM models

```yaml
spec:
  provider: openai
```

#### model

The specific model to use from the provider.

**Type**: `string`  
**Required**: Yes

**OpenAI Models**:
- `gpt-4`
- `gpt-4-turbo`
- `gpt-3.5-turbo`

**Claude Models**:
- `claude-3-opus-20240229`
- `claude-3-sonnet-20240229`
- `claude-3-haiku-20240307`

**Gemini Models**:
- `gemini-pro`
- `gemini-pro-vision`

**vLLM Models**: Any model name supported by your vLLM server

```yaml
spec:
  model: gpt-4
```

#### systemPrompt

Instructions that define the agent's behavior and personality.

**Type**: `string`  
**Required**: Yes

```yaml
spec:
  systemPrompt: |
    You are a helpful customer service agent.
    Always be polite and professional.
    Ask for order numbers when helping with orders.
```

#### apiSecretRef

Reference to a Kubernetes Secret containing the API key.

**Type**: `object`  
**Required**: Yes

**Fields**:
- `name` (string, required): Secret name
- `key` (string, required): Key within the secret

```yaml
spec:
  apiSecretRef:
    name: openai-secret
    key: api-key
```

#### endpoint

Custom API endpoint (primarily for vLLM).

**Type**: `string`  
**Required**: Only for `vllm` provider

```yaml
spec:
  provider: vllm
  endpoint: http://vllm-server:8000/v1
```

#### replicas

Number of agent instances to run.

**Type**: `integer`  
**Required**: No  
**Default**: `1`  
**Minimum**: `0`  
**Maximum**: `100`

```yaml
spec:
  replicas: 3
```

#### resources

CPU and memory resource requests and limits.

**Type**: `object`  
**Required**: No

**Default values**:
```yaml
spec:
  resources:
    requests:
      memory: "128Mi"
      cpu: "100m"
    limits:
      memory: "256Mi" 
      cpu: "200m"
```

**Custom configuration**:
```yaml
spec:
  resources:
    requests:
      memory: "512Mi"
      cpu: "250m"
    limits:
      memory: "1Gi"
      cpu: "500m"
```

#### serviceType

Type of Kubernetes service to create.

**Type**: `string`  
**Required**: No  
**Default**: `ClusterIP`  
**Allowed Values**: `ClusterIP`, `NodePort`, `LoadBalancer`

```yaml
spec:
  serviceType: LoadBalancer
```

#### servicePort

Port for the agent service.

**Type**: `integer`  
**Required**: No  
**Default**: `80`  
**Range**: `1-65535`

```yaml
spec:
  servicePort: 8080
```

#### env

Environment variables for the agent container.

**Type**: `array`  
**Required**: No

```yaml
spec:
  env:
  - name: LOG_LEVEL
    value: "INFO"
  - name: MAX_TOKENS
    value: "2048"
  - name: TEMPERATURE
    value: "0.7"
  - name: DB_PASSWORD
    valueFrom:
      secretKeyRef:
        name: db-secret
        key: password
```

#### tools

External tools the agent can use.

**Type**: `array`  
**Required**: No

```yaml
spec:
  tools:
  - name: "web_search"
    description: "Search the web for information"
    endpoint: "http://search-service:8080/search"
    inputSchema:
      type: "object"
      properties:
        query:
          type: "string"
          description: "Search query"
      required: ["query"]
  - name: "database_query"
    description: "Query the company database"
    endpoint: "http://db-service:5432/query"
    secretRef:
      name: db-credentials
      key: connection-string
```

#### autoscaling

Horizontal Pod Autoscaler configuration.

**Type**: `object`  
**Required**: No

```yaml
spec:
  autoscaling:
    enabled: true
    minReplicas: 2
    maxReplicas: 10
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

#### healthCheck

Health check configuration for the agent.

**Type**: `object`  
**Required**: No

**Default values**:
```yaml
spec:
  healthCheck:
    enabled: true
    path: "/health"
    initialDelaySeconds: 30
    periodSeconds: 10
    timeoutSeconds: 5
    failureThreshold: 3
    successThreshold: 1
```

#### metrics

Prometheus metrics configuration.

**Type**: `object`  
**Required**: No

```yaml
spec:
  metrics:
    enabled: true
    port: 9090
    path: "/metrics"
    interval: "30s"
```

#### securityContext

Pod security context.

**Type**: `object`  
**Required**: No

**Default values**:
```yaml
spec:
  securityContext:
    runAsNonRoot: true
    runAsUser: 65534
    readOnlyRootFilesystem: true
    capabilities:
      drop:
      - ALL
```

---

## Status Fields

The `status` section is managed by the KubeAgentic operator and reflects the current state of the agent.

### Status Structure

```yaml
status:
  phase: "Running"
  conditions:
  - type: "Ready"
    status: "True"
    lastTransitionTime: "2024-01-15T10:30:00Z"
    reason: "AgentReady"
    message: "Agent is running and healthy"
  replicas:
    desired: 3
    ready: 3
    available: 3
  lastUpdated: "2024-01-15T10:30:00Z"
  observedGeneration: 1
```

### Status Fields

| Field | Type | Description |
|-------|------|-------------|
| `phase` | string | Current agent phase |
| `conditions` | array | Detailed condition information |
| `replicas` | object | Replica status counts |
| `lastUpdated` | string | Last status update time |
| `observedGeneration` | integer | Last observed spec generation |

#### phase

Current lifecycle phase of the agent.

**Type**: `string`  
**Possible Values**:
- `Pending` - Agent is being created
- `Running` - Agent is running normally
- `Scaling` - Agent is scaling up/down
- `Failed` - Agent has failed
- `Terminating` - Agent is being deleted

#### conditions

Detailed status conditions.

**Type**: `array`

**Condition Types**:
- `Ready` - Agent is ready to serve requests
- `Available` - Agent has available replicas
- `Progressing` - Agent is progressing towards desired state
- `ReplicaFailure` - Replica creation has failed

**Condition Fields**:
- `type` (string): Condition type
- `status` (string): `True`, `False`, or `Unknown`
- `lastTransitionTime` (string): When condition last changed
- `reason` (string): Brief reason for condition
- `message` (string): Human-readable message

---

## Examples

### Basic Agent

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

### Production-Ready Agent

```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: customer-support
  labels:
    app: customer-support
    tier: production
spec:
  provider: openai
  model: gpt-4
  systemPrompt: |
    You are a professional customer support agent.
    Be helpful, polite, and always ask for order numbers.
  apiSecretRef:
    name: openai-secret
    key: api-key
  
  # Scaling
  replicas: 3
  autoscaling:
    enabled: true
    minReplicas: 2
    maxReplicas: 10
    metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
  
  # Resources
  resources:
    requests:
      memory: "256Mi"
      cpu: "200m"
    limits:
      memory: "512Mi"
      cpu: "500m"
  
  # Service
  serviceType: LoadBalancer
  servicePort: 80
  
  # Monitoring
  metrics:
    enabled: true
    port: 9090
  
  # Security
  securityContext:
    runAsNonRoot: true
    runAsUser: 65534
    readOnlyRootFilesystem: true
  
  # Environment
  env:
  - name: LOG_LEVEL
    value: "INFO"
  - name: MAX_TOKENS
    value: "2048"
  
  # Tools
  tools:
  - name: order_lookup
    description: "Look up customer orders"
    endpoint: "http://order-service:8080/lookup"
    inputSchema:
      type: object
      properties:
        order_id:
          type: string
      required: ["order_id"]
```

### Self-hosted vLLM Agent

```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: internal-assistant
spec:
  provider: vllm
  model: llama2-7b-chat
  endpoint: http://vllm-server.ml.svc.cluster.local:8000/v1
  systemPrompt: "You are an internal company assistant."
  apiSecretRef:
    name: vllm-secret
    key: api-key
  
  replicas: 2
  serviceType: ClusterIP
  
  resources:
    requests:
      memory: "128Mi"
      cpu: "100m"
    limits:
      memory: "256Mi"
      cpu: "200m"
```

---

## Field Validation

### Validation Rules

#### provider + model combinations

```yaml
# Valid combinations
openai + gpt-4
openai + gpt-3.5-turbo
claude + claude-3-sonnet-20240229
gemini + gemini-pro
vllm + any-model-name
```

#### Resource constraints

```yaml
resources:
  requests:
    memory: "64Mi"    # Minimum
    cpu: "50m"        # Minimum
  limits:
    memory: "8Gi"     # Maximum
    cpu: "4000m"      # Maximum
```

#### Replica limits

- `replicas`: 0-100
- `autoscaling.minReplicas`: 1-50
- `autoscaling.maxReplicas`: 1-100
- `maxReplicas` must be ≥ `minReplicas`

#### Port ranges

- `servicePort`: 1-65535 (excluding system ports 1-1023 for non-root)

### Common Validation Errors

**Invalid provider/model combination**:
```yaml
# ❌ Error: Invalid model for provider
spec:
  provider: openai
  model: claude-3-sonnet  # Wrong model for OpenAI
```

**Missing required endpoint for vLLM**:
```yaml
# ❌ Error: endpoint required for vLLM
spec:
  provider: vllm
  model: llama2
  # endpoint: missing!
```

**Invalid resource format**:
```yaml
# ❌ Error: Invalid resource format
resources:
  requests:
    memory: "invalid"  # Should be like "128Mi"
    cpu: "invalid"     # Should be like "100m"
```

---

## API Evolution

### Version Compatibility

- **v1**: Current stable version
- **v1beta1**: Previous beta version (deprecated)

### Deprecated Fields

| Field | Deprecated In | Removed In | Replacement |
|-------|---------------|------------|-------------|
| `image` | v1.0.0 | v1.2.0 | Managed by operator |
| `port` | v1.1.0 | v1.3.0 | `servicePort` |

### Migration Guide

**From v1beta1 to v1**:
```bash
# Update apiVersion
sed -i 's/ai.example.com\/v1beta1/ai.example.com\/v1/g' agent.yaml
```

---

For more examples and use cases, see the [Examples page](examples.html).