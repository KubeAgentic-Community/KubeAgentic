---
layout: home
title: Home
---

# KubeAgentic 🤖

**Deploy and manage AI agents on Kubernetes with simple YAML configurations.**

KubeAgentic is a powerful Kubernetes operator that simplifies the deployment, management, and scaling of AI agents in your cluster. Define your agent's configuration in a simple YAML file and let Kubernetes handle the rest.

## ✨ Key Features

- **🤖 Multi-Provider Support**: OpenAI, Anthropic (Claude), Google (Gemini), and self-hosted vLLM models
- **📝 Declarative Configuration**: Standard Kubernetes Custom Resources (CRDs) 
- **🔄 Autoscaling**: Automatic scaling based on demand
- **🔒 Secure by Default**: API keys managed with Kubernetes Secrets
- **📊 Built-in Monitoring**: Real-time health checks and status reporting
- **🛠️ Tool Integration**: Extend agents with custom tools and services
- **🔗 Framework Choice**: Direct API calls or LangGraph workflows for complex reasoning

## 🚀 Quick Start

```bash
# Install KubeAgentic
kubectl apply -f deploy/all.yaml

# Create API key secret
kubectl create secret generic openai-secret \
  --from-literal=api-key='your-openai-api-key'

# Deploy your first agent
kubectl apply -f examples/openai-agent.yaml

# Interact with your agent
kubectl port-forward service/my-assistant-service 8080:80
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello! How can you help me?"}'
```

## 📚 Documentation

- [📖 User Guide](docs/) - Complete documentation and tutorials
- [🔧 API Reference](api-reference) - Detailed API specification
- [💡 Examples](examples) - Real-world usage examples
- [🧪 Local Testing](local-testing) - Development and testing guide

## 🎯 Use Cases

- **Customer Support**: Deploy scalable support bots
- **Code Review**: Automated code analysis and feedback
- **Knowledge Management**: Internal Q&A assistants
- **Content Generation**: AI-powered content creation

## 🤝 Community

- [GitHub Repository](https://github.com/sudeshmu/KubeAgentic)
- [Issues & Bug Reports](https://github.com/sudeshmu/KubeAgentic/issues)
- [Discussions](https://github.com/sudeshmu/KubeAgentic/discussions)

## 📄 License

Licensed under the Apache License 2.0. See [LICENSE](https://github.com/sudeshmu/KubeAgentic/blob/main/LICENSE) for details.