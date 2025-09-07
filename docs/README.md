# KubeAgentic Documentation

Welcome to the comprehensive documentation for KubeAgentic! This directory contains detailed guides, examples, and references to help you deploy and manage AI agents on Kubernetes.

## 📋 Documentation Structure

### Framework Guides

Choose the right framework for your use case:

**[Direct Framework Guide](direct-framework.md)** 
- Simple, fast API calls directly to LLM providers
- Perfect for chat bots, Q&A systems, and basic tool usage
- Low latency (~100-500ms) and minimal resource requirements
- High concurrency and excellent horizontal scaling
- **Best for**: Customer support bots, code review assistants, content generators

**[LangGraph Framework Guide](langgraph-framework.md)**
- Complex, stateful workflows with multi-step reasoning  
- Advanced tool orchestration and conditional logic
- Higher latency (~1-5s) but sophisticated capabilities
- State management and conversation persistence
- **Best for**: Complex customer service, research pipelines, decision engines

### Reference Documentation

**[API Reference](api.md)**
- Complete Custom Resource Definition (CRD) specification
- All configuration fields with types, defaults, and descriptions
- Validation rules and constraints
- Complete examples for every configuration option

### Framework Comparison

| Aspect | Direct Framework | LangGraph Framework |
|--------|------------------|-------------------|
| **Use Case** | Simple interactions | Complex workflows |
| **Response Time** | ~100-500ms | ~1-5 seconds |
| **State Management** | Stateless | Stateful with persistence |
| **Resource Usage** | Low CPU/memory | Higher CPU/memory |
| **Concurrency** | High | Moderate |
| **Debugging** | Simple logs | Visual workflow debugging |
| **Learning Curve** | Easy | Moderate |

## 🚀 Quick Start Guide

### 1. Choose Your Framework

**For simple use cases** (chat, Q&A, basic tools):
```yaml
spec:
  framework: direct
  provider: openai
  model: gpt-4
  systemPrompt: "You are a helpful assistant."
```

**For complex workflows** (multi-step reasoning, conditional logic):
```yaml
spec:
  framework: langgraph
  provider: openai
  model: gpt-4
  langgraphConfig:
    graphType: conditional
    nodes: [...]  # Define your workflow
    edges: [...]  # Define transitions
```

### 2. Configure Your Agent

Start with these essential fields:
- `provider` - Choose your LLM provider (openai, claude, gemini, vllm)
- `model` - Specify the model (gpt-4, claude-3-sonnet, etc.)
- `systemPrompt` - Define your agent's behavior and persona
- `apiSecretRef` - Reference to your API key secret

### 3. Add Tools (Optional)

Both frameworks support tools for extending agent capabilities:
```yaml
tools:
- name: weather_lookup
  description: "Get current weather for a location"
  inputSchema:
    type: object
    properties:
      location: {type: string}
    required: ["location"]
```

### 4. Deploy and Scale

Configure resources and replicas based on your needs:
```yaml
replicas: 3
resources:
  requests:
    cpu: 200m
    memory: 512Mi
  limits:
    cpu: 500m
    memory: 1Gi
```

## 📖 Learning Path

### Beginners
1. Start with [Direct Framework Guide](direct-framework.md)
2. Try the [simple examples](../examples/direct-agent.yaml)
3. Learn about [tools and configurations](api.md)
4. Set up [local testing](../local-testing/)

### Intermediate  
1. Explore [LangGraph Framework Guide](langgraph-framework.md)
2. Study [complex workflow examples](../examples/langgraph-workflow-agent.yaml)
3. Learn [state management](langgraph-framework.md#state-management)
4. Practice with [conditional logic](langgraph-framework.md#conditional-logic)

### Advanced
1. Design [custom workflows](langgraph-framework.md#workflow-components)
2. Implement [error handling](langgraph-framework.md#error-handling)
3. Optimize [performance and scaling](langgraph-framework.md#performance-optimization)
4. Set up [monitoring and debugging](langgraph-framework.md#monitoring-and-debugging)

## 🎯 Use Case Examples

### Customer Support
- **Simple FAQ Bot** → [Direct Framework](direct-framework.md#1-customer-support-chat-bot)
- **Complex Service Workflow** → [LangGraph Framework](langgraph-framework.md#1-advanced-customer-service-workflow)

### Content Creation  
- **Marketing Copy Generator** → [Direct Framework](direct-framework.md#3-content-generation-agent)
- **Content Production Pipeline** → [LangGraph Framework](langgraph-framework.md#4-content-production-pipeline)

### Analysis & Research
- **Data Analysis Helper** → [Direct Framework](direct-framework.md#4-data-analysis-assistant)  
- **Multi-Source Research** → [LangGraph Framework](langgraph-framework.md#2-research-and-analysis-workflow)

### Development
- **Code Review Assistant** → [Direct Framework](direct-framework.md#2-code-review-assistant)
- **Complex Decision Engine** → [LangGraph Framework](langgraph-framework.md#3-complex-decision-engine)

## 🔧 Configuration Patterns

### High Throughput (Direct)
```yaml
framework: direct
replicas: 10
resources:
  requests: {cpu: 100m, memory: 256Mi}
  limits: {cpu: 200m, memory: 512Mi}
```

### Complex Processing (LangGraph)
```yaml
framework: langgraph  
replicas: 2
resources:
  requests: {cpu: 500m, memory: 1Gi}
  limits: {cpu: 1, memory: 2Gi}
```

### Development/Testing
```yaml
replicas: 1
resources:
  requests: {cpu: 100m, memory: 128Mi}
  limits: {cpu: 500m, memory: 1Gi}
```

## 🛠️ Tools and Integrations

Both frameworks support tools for extending functionality:

**Common tool patterns:**
- Database lookups
- API integrations  
- File operations
- Calculations
- External service calls

**Framework-specific usage:**
- **Direct**: Tools called independently based on AI decision
- **LangGraph**: Tools orchestrated within workflow nodes

## 📊 Monitoring and Operations

### Essential Metrics
- Response latency
- Request throughput  
- Error rates
- Resource utilization
- Tool usage patterns

### Health Checks
```bash
kubectl get agents
kubectl describe agent my-agent
kubectl logs deployment/my-agent-deployment
```

### Scaling
```bash
kubectl scale agent my-agent --replicas=5
kubectl get hpa  # If autoscaling enabled
```

## 🤝 Contributing

Help improve the documentation:
1. Report issues or gaps in coverage
2. Submit examples for new use cases
3. Contribute performance optimization tips
4. Share best practices from production usage

## 🔗 Additional Resources

- **Website**: [https://KubeAgentic-Community.github.io/KubeAgentic/](https://KubeAgentic-Community.github.io/KubeAgentic/)
- **GitHub**: [https://github.com/KubeAgentic-Community/KubeAgentic](https://github.com/KubeAgentic-Community/KubeAgentic)
- **Examples**: [/examples directory](../examples/)
- **Local Testing**: [/local-testing directory](../local-testing/)

## 📞 Support

- **Documentation Issues**: [GitHub Issues](https://github.com/KubeAgentic-Community/KubeAgentic/issues)
- **Feature Requests**: [GitHub Discussions](https://github.com/KubeAgentic-Community/KubeAgentic/discussions)
- **Community**: Check the repository for community guidelines

Choose your framework, configure your agent, and start building intelligent systems on Kubernetes!
