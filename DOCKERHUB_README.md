# KubeAgentic - Kubernetes AI Agent Operator

![Docker Pulls](https://img.shields.io/docker/pulls/sudeshmu/kubeagentic)
![Image Size](https://img.shields.io/docker/image-size/sudeshmu/kubeagentic/operator-latest)
![GitHub Stars](https://img.shields.io/github/stars/KubeAgentic-Community/KubeAgentic)

Deploy and manage AI agents on Kubernetes with simple YAML configurations. KubeAgentic is a powerful Kubernetes operator that simplifies the deployment, management, and scaling of AI agents in your cluster.

## 🏷️ Available Images

### Operator (Recommended)
```bash
docker pull sudeshmu/kubeagentic:operator-latest
```
- **Size:** 108MB (Highly Optimized)
- **Base:** Red Hat UBI Micro
- **Architecture:** linux/amd64, linux/arm64
- **Purpose:** Kubernetes operator for managing agents

### Agent Runtime
```bash
docker pull sudeshmu/kubeagentic:agent-latest
```
- **Size:** 625MB (66% smaller than original!)
- **Base:** Red Hat UBI Minimal  
- **Architecture:** linux/amd64, linux/arm64
- **Purpose:** Python-based agent runtime with AI frameworks

## 🚀 Quick Start

### 1. Deploy the Operator
```bash
# Pull the optimized operator image
docker pull sudeshmu/kubeagentic:operator-latest

# Deploy to Kubernetes
kubectl apply -f https://raw.githubusercontent.com/KubeAgentic-Community/kubeagentic/main/deploy/all.yaml
```

### 2. Create Your First Agent
```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: my-assistant
spec:
  provider: openai
  model: gpt-4
  systemPrompt: "You are a helpful assistant."
  apiSecretRef:
    name: openai-secret
    key: api-key
```

## ✨ Image Optimization Features

### Multi-Stage Builds
- **Build Stage**: Full development environment for dependency installation
- **Runtime Stage**: Minimal production environment with only essential components
- **Result**: 66% size reduction for agent runtime (1.85GB → 625MB)

### Security Hardening
- 🔒 **Non-root execution** with user ID 1001
- 🛡️ **Red Hat Universal Base Images** (UBI) for enterprise security
- 🔐 **Minimal attack surface** with only required packages
- 🚫 **No package managers** in runtime images

### Performance Optimizations
- ⚡ **Virtual environments** for isolated Python dependencies
- 🗜️ **Layer optimization** with combined commands
- 🚀 **Multi-architecture support** (AMD64 + ARM64)
- 💾 **Efficient caching** with .dockerignore patterns

## 🏗️ Supported AI Providers

| Provider | Models | Authentication |
|----------|---------|----------------|
| **OpenAI** | GPT-4, GPT-3.5-turbo | API Key |
| **Anthropic** | Claude-3 (Opus, Sonnet, Haiku) | API Key |
| **Google** | Gemini Pro, Gemini Pro Vision | API Key |
| **vLLM** | Self-hosted models | Optional API Key |

## 🔧 Configuration Examples

### High-Performance Setup
```yaml
spec:
  provider: openai
  model: gpt-4
  replicas: 5
  resources:
    requests: {cpu: 200m, memory: 256Mi}
    limits: {cpu: 500m, memory: 512Mi}
  framework: direct  # Low latency
```

### Complex Workflow Setup
```yaml
spec:
  provider: claude
  model: claude-3-sonnet-20240229
  framework: langgraph
  langgraphConfig:
    graphType: conditional
    nodes: [...]  # Multi-step reasoning
```

## 📊 Image Size Comparison

| Version | Operator | Agent Runtime | Total |
|---------|----------|---------------|--------|
| **Optimized** | 108MB | 625MB | 733MB |
| Original | ~150MB | 1.85GB | ~2GB |
| **Savings** | 28% | 66% | 63% |

## 🛠️ Development & Customization

### Build Your Own Images
```bash
# Clone the repository
git clone https://github.com/KubeAgentic-Community/KubeAgentic.git
cd KubeAgentic

# Build operator
docker build -f Dockerfile.operator -t my-kubeagentic:operator .

# Build agent runtime  
docker build -f Dockerfile.agent -t my-kubeagentic:agent .
```

### Environment Variables
```bash
# Operator
- LOG_LEVEL=info
- METRICS_ADDR=:8080
- HEALTH_PROBE_ADDR=:8081

# Agent Runtime
- PORT=8080
- LOG_LEVEL=info
- PYTHONUNBUFFERED=1
```

## 📋 System Requirements

### Minimum Requirements
- **Kubernetes:** v1.19+
- **CPU:** 100m per agent
- **Memory:** 128Mi per agent
- **Storage:** 1Gi for images

### Recommended for Production
- **CPU:** 200-500m per agent
- **Memory:** 256Mi-1Gi per agent  
- **Replicas:** 2+ for high availability
- **Monitoring:** Prometheus + Grafana

## 🔗 Links & Resources

- **🌐 Website:** [kubeagentic.com](https://kubeagentic.com)
- **📚 Documentation:** [GitHub Docs](https://github.com/KubeAgentic-Community/KubeAgentic/tree/main/docs)
- **🐛 Issues:** [GitHub Issues](https://github.com/KubeAgentic-Community/KubeAgentic/issues)
- **💬 Discussions:** [GitHub Discussions](https://github.com/KubeAgentic-Community/KubeAgentic/discussions)
- **⭐ Source Code:** [GitHub Repository](https://github.com/KubeAgentic-Community/KubeAgentic)

## 🏷️ Tags & Versioning

| Tag | Description | Update Frequency |
|-----|-------------|------------------|
| `operator-latest` | Latest stable operator | On releases |
| `agent-latest` | Latest optimized agent | On releases |
| `agent-optimized` | Explicitly optimized version | On optimization updates |

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](https://github.com/KubeAgentic-Community/KubeAgentic/blob/main/CONTRIBUTING.md) for details.

## 📄 License

Apache License 2.0 - see [LICENSE](https://github.com/KubeAgentic-Community/KubeAgentic/blob/main/LICENSE)

---

**Built with ❤️ for the Kubernetes AI community**

*These optimized images are built using multi-stage Docker builds with Red Hat UBI base images, ensuring security, performance, and minimal size for production deployments.*
