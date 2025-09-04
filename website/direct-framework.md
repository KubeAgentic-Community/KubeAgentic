---
layout: page
title: Direct Framework Guide
description: Simple, fast API calls for basic AI agent interactions
---

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
- **Scalability**: Excellent horizontal scaling
- **Complexity**: Minimal configuration required

## Configuration Example

Here's a basic Direct Framework agent configuration:

```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: simple-chatbot
spec:
  framework: direct
  provider: openai
  model: gpt-4
  systemPrompt: "You are a helpful customer service assistant."
  apiSecretRef:
    name: openai-secret
    key: api-key
  replicas: 2
  resources:
    requests:
      memory: "128Mi"
      cpu: "100m"
    limits:
      memory: "256Mi"
      cpu: "200m"
```

## Key Features

### Simple Configuration
The Direct Framework requires minimal configuration. Just specify your provider, model, and system prompt.

### Fast Response Times
Direct API calls to LLM providers ensure the fastest possible response times for your applications.

### Resource Efficient
Minimal overhead means your agents use fewer resources, allowing for better cost optimization.

### Easy Debugging
Simple request-response patterns make it easy to debug and monitor your agents.

## Tool Integration

Even with the Direct Framework, you can still integrate tools:

```yaml
spec:
  framework: direct
  provider: openai
  model: gpt-4
  tools:
  - name: calculator
    description: "Basic math operations"
    endpoint: "http://calculator-service:8080/calculate"
  - name: weather
    description: "Get current weather information"
    endpoint: "http://weather-service:8080/weather"
```

## Monitoring and Observability

The Direct Framework provides built-in monitoring capabilities:

- **Health Checks**: Automatic health monitoring
- **Metrics**: Request/response metrics via Prometheus
- **Logging**: Structured logging for debugging
- **Tracing**: Request tracing for performance analysis

## Best Practices

1. **Keep it Simple**: Use Direct Framework for straightforward use cases
2. **Optimize Prompts**: Well-crafted system prompts improve response quality
3. **Monitor Performance**: Track response times and resource usage
4. **Scale Horizontally**: Add more replicas for high-throughput scenarios
5. **Use Appropriate Models**: Choose models based on your performance and cost requirements

## Migration from Other Frameworks

If you're currently using a more complex framework and want to simplify:

1. **Evaluate Complexity**: Determine if you really need complex workflows
2. **Simplify Logic**: Move complex logic to your application layer
3. **Test Performance**: Ensure Direct Framework meets your performance needs
4. **Update Configuration**: Modify your agent specs to use `framework: direct`

## Troubleshooting

### Common Issues

**Slow Response Times**
- Check your network connectivity to the LLM provider
- Verify your API key has sufficient quota
- Consider using a faster model or region

**High Resource Usage**
- Review your resource limits and requests
- Check for memory leaks in your application
- Monitor CPU usage patterns

**Tool Integration Issues**
- Verify tool endpoints are accessible
- Check tool response formats
- Ensure proper error handling

## Next Steps

- [View Examples](examples) - See real-world Direct Framework implementations
- [API Reference](api-reference) - Detailed configuration options
- [LangGraph Framework](langgraph-framework) - For complex workflows
- [Local Testing](local-testing) - Test your agents locally
