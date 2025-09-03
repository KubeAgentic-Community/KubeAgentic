---
layout: default
title: Direct Framework Guide
nav_order: 4
---

# Direct Framework Guide

The **Direct Framework** is KubeAgentic's default execution mode, designed for fast, straightforward interactions with AI agents. It provides simple API calls directly to the LLM provider without complex workflow orchestration.

## When to Use Direct Framework

✅ **Perfect for:**
- Chat bots and conversational agents
- Simple Q&A systems
- Basic tool usage scenarios
- High-throughput applications
- Lightweight agents with minimal resource requirements
- Straightforward request-response patterns

❌ **Not ideal for:**
- Complex multi-step reasoning
- Stateful conversation workflows
- Advanced tool orchestration
- Conditional logic between tools
- Long-running task workflows

## Performance Characteristics

- **Response Time**: ~100-500ms
- **Resource Usage**: Low CPU and memory footprint
- **Concurrency**: High - supports many simultaneous requests
- **Scalability**: Excellent horizontal scaling
- **Debugging**: Simple request/response flow

## Quick Example

```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: support-chatbot
  namespace: customer-service
spec:
  framework: direct        # Simple, fast interactions
  provider: openai
  model: gpt-4
  systemPrompt: "You are a helpful customer support agent."
  apiSecretRef:
    name: openai-secret
    key: api-key
  
  tools:
  - name: order_lookup
    description: Look up customer order information
    inputSchema:
      type: object
      properties:
        order_id: {type: string}
      required: ["order_id"]
  
  replicas: 3
  resources:
    requests:
      cpu: 100m
      memory: 256Mi
    limits:
      cpu: 200m
      memory: 512Mi
```

## Common Use Cases

### 1. Customer Support Bot
- Handle FAQs and basic inquiries
- Look up order status and customer information
- Escalate complex issues to humans
- Provide quick, accurate responses

### 2. Code Review Assistant
- Analyze code for best practices
- Suggest improvements and optimizations
- Check for security vulnerabilities
- Provide educational feedback

### 3. Content Generator
- Create marketing copy and social media posts
- Generate product descriptions
- Write blog post outlines
- Adapt content for different channels

### 4. Data Analysis Helper
- Interpret charts and business metrics
- Explain trends and patterns
- Provide data-driven recommendations
- Create visualization descriptions

### 5. Educational Tutor
- Answer student questions
- Generate practice problems
- Provide step-by-step explanations
- Adapt to different learning levels

## Configuration Best Practices

### Resource Allocation

**Light workloads** (simple chat):
```yaml
resources:
  requests:
    cpu: 100m
    memory: 256Mi
  limits:
    cpu: 200m
    memory: 512Mi
```

**Medium workloads** (with tools):
```yaml
resources:
  requests:
    cpu: 200m
    memory: 512Mi
  limits:
    cpu: 500m
    memory: 1Gi
```

### Scaling for High Throughput

```yaml
replicas: 5
resources:
  requests:
    cpu: 100m    # Lower per-replica usage
    memory: 256Mi
serviceType: LoadBalancer  # External access
```

### Tool Design Tips

1. **Keep tools focused**: Single, clear purpose per tool
2. **Validate inputs**: Use proper JSON schema validation
3. **Provide good descriptions**: Help the AI understand when to use each tool
4. **Handle errors gracefully**: Return meaningful error messages
5. **Optimize for speed**: Direct framework excels with fast tool responses

## Monitoring and Troubleshooting

**Essential metrics:**
- Response latency (target: <500ms)
- Request throughput
- Error rates
- Resource utilization
- Tool usage patterns

**Common issues:**
- **High latency** → Check tool performance, consider caching
- **High error rates** → Validate tool schemas and API connections
- **Resource pressure** → Adjust limits or increase replicas
- **Poor responses** → Review system prompts and tool descriptions

## When to Consider LangGraph

Consider upgrading to [LangGraph Framework](langgraph-framework) when you need:
- Multi-step conditional logic
- State persistence across interactions
- Complex tool orchestration
- Decision trees and branching workflows
- Advanced reasoning capabilities

## Complete Examples

Visit our [examples directory](https://github.com/sudeshmu/KubeAgentic/tree/main/examples) for full configuration files:

- [`direct-agent.yaml`](https://github.com/sudeshmu/KubeAgentic/blob/main/examples/direct-agent.yaml) - Simple customer support agent
- [`claude-agent.yaml`](https://github.com/sudeshmu/KubeAgentic/blob/main/examples/claude-agent.yaml) - Code review assistant
- [`openai-agent.yaml`](https://github.com/sudeshmu/KubeAgentic/blob/main/examples/openai-agent.yaml) - General purpose agent

## Next Steps

- [View API Reference](api-reference) for complete configuration options
- [Try Local Testing](local-testing) to experiment with different configurations  
- [Compare with LangGraph](langgraph-framework) for complex workflows
- [Check out Examples](examples) for more use cases

The Direct Framework provides the perfect balance of simplicity, speed, and reliability for most AI agent use cases.
