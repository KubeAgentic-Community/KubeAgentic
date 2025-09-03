# Agent API Reference

This document describes the complete API specification for the `Agent` Custom Resource in KubeAgentic.

## Agent Resource

The `Agent` resource is the core component of KubeAgentic. It defines an AI agent's configuration and desired state.

### API Version

- **API Version**: `ai.example.com/v1`
- **Kind**: `Agent`

### Basic Structure

```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: agent-name
  namespace: default
spec:
  # Agent specification
status:
  # Agent status (managed by operator)
```

## Spec Fields

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `provider` | string | LLM provider to use |
| `model` | string | Specific model name |
| `systemPrompt` | string | Agent's system prompt |
| `apiSecretRef` | object | Reference to API key secret |

#### provider

Specifies which LLM provider to use for this agent.

**Type**: `string`  
**Required**: Yes  
**Allowed Values**: `openai`, `claude`, `gemini`, `vllm`

```yaml
spec:
  provider: openai
```

#### model

The specific model to use from the provider.

**Type**: `string`  
**Required**: Yes  

**Examples by Provider**:
- **OpenAI**: `gpt-4`, `gpt-3.5-turbo`, `gpt-4-turbo`
- **Claude**: `claude-3-sonnet-20240229`, `claude-3-opus-20240229`, `claude-3-haiku-20240307`
- **Gemini**: `gemini-pro`, `gemini-pro-vision`
- **vLLM**: Any model supported by your vLLM deployment

```yaml
spec:
  model: gpt-4
```

#### systemPrompt

The system prompt that defines the agent's behavior and personality.

**Type**: `string`  
**Required**: Yes  

```yaml
spec:
  systemPrompt: |
    You are a helpful customer service agent.
    Be friendly, professional, and always try to solve customer problems.
```

#### apiSecretRef

Reference to a Kubernetes Secret containing the API key for the LLM provider.

**Type**: `object`  
**Required**: Yes  

**Properties**:
- `name` (string, required): Name of the Secret
- `key` (string, required): Key within the Secret containing the API key

```yaml
spec:
  apiSecretRef:
    name: openai-secret
    key: api-key
```

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `endpoint` | string | - | Custom endpoint URL |
| `replicas` | integer | 1 | Number of replicas |
| `resources` | object | See below | Resource requirements |
| `serviceType` | string | `ClusterIP` | Kubernetes service type |
| `tools` | array | `[]` | Available tools |

#### endpoint

Custom endpoint URL for self-hosted models or alternative API endpoints.

**Type**: `string`  
**Required**: No  
**Use Cases**: vLLM deployments, OpenAI-compatible APIs, custom endpoints

```yaml
spec:
  endpoint: http://my-vllm-server:8000/v1
```

#### framework

Specifies which framework to use for agent execution.

**Type**: `string`  
**Required**: No  
**Default**: `direct`  
**Allowed Values**: `direct`, `langgraph`

- **`direct`**: Simple, fast API calls directly to the LLM
- **`langgraph`**: Complex, stateful workflows with multi-step reasoning

```yaml
spec:
  framework: langgraph
```

#### langgraphConfig

Configuration for LangGraph workflows. Only used when `framework` is set to `langgraph`.

**Type**: `object`  
**Required**: No (required when framework is `langgraph`)

**Properties**:
- `graphType` (string, required): Type of workflow (`sequential`, `parallel`, `conditional`, `hierarchical`)
- `nodes` (array, required): Workflow nodes definition
- `edges` (array, required): Workflow edges definition  
- `state` (object, optional): State schema for the workflow
- `entrypoint` (string, required): Entry node name
- `endpoints` (array, optional): Possible end nodes

```yaml
spec:
  endpoint: http://my-vllm-server:8000/v1
```

#### replicas

Number of agent pod replicas to run.

**Type**: `integer`  
**Required**: No  
**Default**: `1`  
**Minimum**: `1`  
**Maximum**: `10`

```yaml
spec:
  replicas: 3
```

#### resources

Resource requests and limits for agent pods.

**Type**: `object`  
**Required**: No  

**Default Values**:
```yaml
spec:
  resources:
    requests:
      cpu: "100m"
      memory: "256Mi"
    limits:
      cpu: "200m"
      memory: "512Mi"
```

**Custom Example**:
```yaml
spec:
  resources:
    requests:
      cpu: "200m"
      memory: "512Mi"
    limits:
      cpu: "1000m"
      memory: "2Gi"
```

#### serviceType

Kubernetes Service type for exposing the agent.

**Type**: `string`  
**Required**: No  
**Default**: `ClusterIP`  
**Allowed Values**: `ClusterIP`, `NodePort`, `LoadBalancer`

```yaml
spec:
  serviceType: LoadBalancer
```

#### tools

Array of tools available to the agent.

**Type**: `array`  
**Required**: No  
**Default**: `[]`

**Tool Object Properties**:
- `name` (string, required): Tool identifier
- `description` (string, required): Human-readable description
- `inputSchema` (object, optional): JSON schema for input validation

```yaml
spec:
  tools:
  - name: weather_lookup
    description: Get current weather for a location
    inputSchema:
      type: object
      properties:
        location:
          type: string
          description: City name or coordinates
      required: ["location"]
  - name: calculator
    description: Perform mathematical calculations
    inputSchema:
      type: object
      properties:
        expression:
          type: string
          description: Mathematical expression to evaluate
      required: ["expression"]
```

## Status Fields

The status section is managed by the KubeAgentic operator and reflects the current state of the agent.

### Status Properties

| Field | Type | Description |
|-------|------|-------------|
| `phase` | string | Current deployment phase |
| `message` | string | Human-readable status message |
| `replicaStatus` | object | Replica status information |
| `lastUpdated` | string | Last update timestamp |
| `conditions` | array | Detailed status conditions |

#### phase

Current phase of the agent deployment.

**Type**: `string`  
**Possible Values**:
- `Pending`: Agent is being created or updated
- `Running`: Agent is running and ready
- `Failed`: Agent deployment failed
- `Succeeded`: Agent completed successfully (rare)

#### replicaStatus

Information about agent replicas.

**Type**: `object`  
**Properties**:
- `desired` (integer): Number of desired replicas
- `ready` (integer): Number of ready replicas
- `available` (integer): Number of available replicas

#### conditions

Array of status conditions providing detailed state information.

**Type**: `array`  
**Condition Properties**:
- `type` (string): Condition type (`Ready`, `Progressing`, `Degraded`)
- `status` (string): Condition status (`True`, `False`, `Unknown`)
- `reason` (string): Brief reason for the condition
- `message` (string): Human-readable message
- `lastTransitionTime` (string): When the condition last changed

## Complete Examples

### Direct Framework Example

```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: simple-support
  namespace: production
spec:
  # Framework choice
  framework: direct
  
  # Required fields
  provider: openai
  model: gpt-4
  systemPrompt: "You are a helpful customer support agent."
  apiSecretRef:
    name: openai-secret
    key: api-key
  
  # Simple tools
  tools:
  - name: order_lookup
    description: Look up customer order information
    inputSchema:
      type: object
      properties:
        order_id: {type: string}
      required: ["order_id"]
  
  replicas: 2
  serviceType: ClusterIP
```

### LangGraph Workflow Example

```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: workflow-support
  namespace: production
spec:
  # Framework choice
  framework: langgraph
  
  # Required fields
  provider: openai
  model: gpt-4
  systemPrompt: "You are an advanced customer service agent with systematic problem-solving capabilities."
  apiSecretRef:
    name: openai-secret
    key: api-key
  
  # LangGraph workflow configuration
  langgraphConfig:
    graphType: conditional
    nodes:
    - name: classify_issue
      type: llm
      prompt: "Classify this customer issue: {user_input}"
      outputs: ["issue_type"]
    - name: lookup_order
      type: tool
      tool: order_lookup
      condition: "issue_type == 'order'"
      inputs: ["order_id"]
      outputs: ["order_data"]
    - name: resolve_issue
      type: llm
      prompt: "Resolve the issue based on: {order_data}"
      outputs: ["resolution"]
    edges:
    - from: classify_issue
      to: lookup_order
      condition: "issue_type == 'order'"
    - from: lookup_order
      to: resolve_issue
    entrypoint: classify_issue
    endpoints: [resolve_issue]
  
  # Tools for workflow
  tools:
  - name: order_lookup
    description: Look up order details
    inputSchema:
      type: object
      properties:
        order_id: {type: string}
      required: ["order_id"]
  
  # Higher resources for workflow processing
  replicas: 1
  resources:
    requests:
      cpu: "200m"
      memory: "512Mi"
    limits:
      cpu: "500m"
      memory: "1Gi"
```

## Validation Rules

The following validation rules are enforced by the CRD:

1. **Provider Enum**: Must be one of `openai`, `claude`, `gemini`, `vllm`
2. **Replica Limits**: Must be between 1 and 10 inclusive
3. **Service Type Enum**: Must be `ClusterIP`, `NodePort`, or `LoadBalancer`
4. **Required Fields**: `provider`, `model`, `systemPrompt`, and `apiSecretRef` are mandatory
5. **Secret Reference**: `apiSecretRef` must have both `name` and `key` fields
6. **Tool Schema**: Each tool must have `name` and `description` fields

## Error Conditions

Common error conditions and their meanings:

| Condition | Cause | Resolution |
|-----------|-------|------------|
| `SecretNotFound` | Referenced secret doesn't exist | Create the required secret with API key |
| `InvalidProvider` | Unsupported provider specified | Use a supported provider |
| `ModelNotFound` | Model not available for provider | Check model name and provider compatibility |
| `ResourceConstraints` | Insufficient cluster resources | Adjust resource requests or add capacity |
| `EndpointUnreachable` | Custom endpoint not accessible | Verify endpoint URL and network connectivity |

For more troubleshooting information, see the [main documentation](../README.md).
