# KubeAgentic - Multi-Architecture Kubernetes AI Agent Operator

![Docker Pulls](https://img.shields.io/docker/pulls/sudeshmu/kubeagentic)
![Image Size](https://img.shields.io/docker/image-size/sudeshmu/kubeagentic/operator-latest)
![GitHub Stars](https://img.shields.io/github/stars/KubeAgentic-Community/KubeAgentic)
![Architectures](https://img.shields.io/badge/architectures-AMD64%20%7C%20ARM64-blue)

Deploy and manage AI agents on Kubernetes with simple YAML configurations. KubeAgentic is a powerful Kubernetes operator that simplifies the deployment, management, and scaling of AI agents in your cluster.

## 🏗️ **Multi-Architecture Support**

**✅ Native support for both AMD64 and ARM64 architectures**
- **Intel/AMD x86_64**: Traditional servers, VMs, most cloud instances
- **ARM64**: Apple Silicon (M1/M2), AWS Graviton, GCP T2A, Azure Ampere
- **Automatic Selection**: Kubernetes automatically picks the right architecture
- **Single Manifest**: One image reference works on all platforms

## 🏷️ Available Images

### Operator (Recommended)
```bash
docker pull sudeshmu/kubeagentic:operator-latest
```
- **Size:** ~219MB (Multi-stage optimized)
- **Base:** Red Hat UBI Micro
- **Architectures:** `linux/amd64`, `linux/arm64`
- **Purpose:** Kubernetes operator for managing agents
- **Build:** Multi-architecture using Docker Buildx

### Agent Runtime
```bash
docker pull sudeshmu/kubeagentic:agent-fixed
```
- **Size:** ~1.25GB (Includes all AI frameworks)
- **Base:** Red Hat UBI Minimal  
- **Architectures:** `linux/amd64`, `linux/arm64`
- **Purpose:** Python 3.11 + FastAPI + LangGraph/LangChain
- **Features:** Direct + LangGraph framework support

### ⚡ **Architecture-Specific Pulls**
```bash
# Automatically selects your platform
docker pull sudeshmu/kubeagentic:agent-fixed

# Force specific architecture
docker pull --platform linux/amd64 sudeshmu/kubeagentic:agent-fixed
docker pull --platform linux/arm64 sudeshmu/kubeagentic:agent-fixed
```

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
  framework: direct  # Choose: direct or langgraph
  systemPrompt: "You are a helpful assistant."
  apiSecretRef:
    name: openai-secret
    key: api-key
```

### 3. Test Multi-Architecture Deployment
```bash
# View image manifests
docker buildx imagetools inspect sudeshmu/kubeagentic:operator-latest
docker buildx imagetools inspect sudeshmu/kubeagentic:agent-fixed

# Test on different architectures
docker run --rm --platform linux/amd64 sudeshmu/kubeagentic:agent-fixed python --version
docker run --rm --platform linux/arm64 sudeshmu/kubeagentic:agent-fixed python --version
```

## ✨ Multi-Architecture Build Features

### 🏗️ **Docker Buildx Integration**
- **Cross-platform builds**: Single command builds for AMD64 + ARM64
- **Manifest lists**: One image tag serves all architectures automatically
- **Parallel builds**: Simultaneous compilation for faster CI/CD
- **Registry optimization**: Efficient layer sharing between architectures

### 🔒 **Security Hardening**
- 🔒 **Non-root execution** with user ID 1001 (agent) and 65532 (operator)
- 🛡️ **Red Hat Universal Base Images** (UBI) for enterprise security
- 🔐 **Minimal attack surface** with only required packages
- 🚫 **No package managers** in runtime images
- 🏷️ **Attestation manifests** for supply chain security

### ⚡ **Performance Optimizations**
- 🖥️ **Native execution** on both x86_64 and ARM64
- ⚡ **Virtual environments** for isolated Python dependencies
- 🗜️ **Layer optimization** with combined commands
- 💾 **Efficient caching** with .dockerignore patterns
- 🚀 **Platform-specific optimizations** during build

## 🏗️ Supported AI Providers

| Provider | Models | Authentication |
|----------|---------|----------------|
| **OpenAI** | GPT-4, GPT-3.5-turbo | API Key |
| **Anthropic** | Claude-3 (Opus, Sonnet, Haiku) | API Key |
| **Google** | Gemini Pro, Gemini Pro Vision | API Key |
| **vLLM** | Self-hosted models | Optional API Key |

## 🔧 Configuration Examples

### 🚀 **High-Performance Setup (Direct Framework)**
```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: high-performance-agent
spec:
  provider: openai
  model: gpt-4
  framework: direct  # Low latency, simple workflows
  replicas: 5
  systemPrompt: "You are a high-performance assistant."
  resources:
    requests: {cpu: 200m, memory: 256Mi}
    limits: {cpu: 500m, memory: 512Mi}
  apiSecretRef:
    name: openai-secret
    key: api-key
```

### 🧠 **Complex Workflow Setup (LangGraph Framework)**
```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: workflow-agent
spec:
  provider: anthropic
  model: claude-3-sonnet-20240229
  framework: langgraph  # Complex multi-step workflows
  systemPrompt: "You are a workflow automation assistant."
  langgraphConfig:
    graphType: conditional
    nodes:
      - name: analyze
        type: llm
      - name: tools
        type: tool
    edges:
      - from: analyze
        to: tools
        condition: needs_tools
    entrypoint: analyze
  tools:
    - calculator
    - web_search
  apiSecretRef:
    name: anthropic-secret
    key: api-key
```

### 🌐 **Multi-Architecture Deployment**
```yaml
# Works automatically on both AMD64 and ARM64 nodes
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kubeagentic-operator
spec:
  replicas: 2
  template:
    spec:
      containers:
      - name: manager
        image: sudeshmu/kubeagentic:operator-latest  # ← Auto-selects architecture
        resources:
          requests: {cpu: 100m, memory: 128Mi}
          limits: {cpu: 200m, memory: 256Mi}
```

## 📊 Multi-Architecture Image Details

| Component | AMD64 Size | ARM64 Size | Architectures | Base Image |
|-----------|------------|------------|---------------|------------|
| **Operator** | ~219MB | ~219MB | ✅ Both | UBI Micro |
| **Agent** | ~1.25GB | ~1.25GB | ✅ Both | UBI Minimal |

### 🏷️ **Image Manifests**
```bash
# Each image tag contains multiple architecture-specific manifests
sudeshmu/kubeagentic:operator-latest
├── linux/amd64 → sha256:2335acc4...
├── linux/arm64 → sha256:08d4833d...
└── attestations (security metadata)

sudeshmu/kubeagentic:agent-fixed  
├── linux/amd64 → sha256:c33d00cb...
├── linux/arm64 → sha256:2cdf8f8e...
└── attestations (security metadata)
```

## 🛠️ Development & Customization

### 🏗️ **Build Multi-Architecture Images**
```bash
# Clone the repository
git clone https://github.com/KubeAgentic-Community/KubeAgentic.git
cd KubeAgentic

# Setup buildx for multi-architecture builds
make buildx-setup

# Build and push multi-architecture images
make docker-buildx-all

# Or build individually
make docker-buildx-operator  # Build operator for AMD64 + ARM64
make docker-buildx-agent     # Build agent for AMD64 + ARM64

# Build locally without pushing (for development)
make docker-buildx-local-all

# Legacy single-architecture builds (if needed)
docker build -f Dockerfile.operator -t my-kubeagentic:operator .
docker build -f Dockerfile.agent -t my-kubeagentic:agent .
```

### 🔍 **Inspect Multi-Architecture Images**
```bash
# View detailed manifest information
make inspect-images

# Or manually inspect
docker buildx imagetools inspect sudeshmu/kubeagentic:operator-latest
docker buildx imagetools inspect sudeshmu/kubeagentic:agent-fixed
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

### 🖥️ **Multi-Architecture Support**
- **AMD64 nodes**: Intel/AMD x86_64 processors
- **ARM64 nodes**: Apple Silicon, AWS Graviton, GCP T2A, Azure Ampere
- **Mixed clusters**: Automatic architecture selection

### Minimum Requirements
- **Kubernetes:** v1.19+
- **CPU:** 100m per agent (both architectures)
- **Memory:** 128Mi per agent (both architectures)
- **Storage:** 1.5Gi for images (both architectures cached)

### Recommended for Production
- **CPU:** 200-500m per agent
- **Memory:** 256Mi-1Gi per agent  
- **Replicas:** 2+ for high availability (can span architectures)
- **Monitoring:** Prometheus + Grafana
- **Node selection**: Mix of AMD64 and ARM64 for cost optimization

### ☁️ **Cloud Provider Compatibility**

| Provider | AMD64 Support | ARM64 Support | ARM Instance Types | Cost Savings |
|----------|---------------|---------------|-------------------|--------------|
| **AWS** | ✅ | ✅ | Graviton2/3 instances | Up to 40% |
| **GCP** | ✅ | ✅ | T2A instances | Up to 35% |
| **Azure** | ✅ | ✅ | Ampere Altra instances | Up to 50% |
| **Local** | ✅ | ✅ | M1/M2 Macs, ARM SBCs | Native performance |

## 🔗 Links & Resources

- **🌐 Website:** [kubeagentic.com](https://kubeagentic.com)
- **📚 Documentation:** [GitHub Docs](https://github.com/KubeAgentic-Community/KubeAgentic/tree/main/docs)
- **🐛 Issues:** [GitHub Issues](https://github.com/KubeAgentic-Community/KubeAgentic/issues)
- **💬 Discussions:** [GitHub Discussions](https://github.com/KubeAgentic-Community/KubeAgentic/discussions)
- **⭐ Source Code:** [GitHub Repository](https://github.com/KubeAgentic-Community/KubeAgentic)

## 🏷️ Tags & Versioning

| Tag | Description | Architectures | Update Frequency |
|-----|-------------|---------------|------------------|
| `operator-latest` | Latest stable operator | AMD64, ARM64 | On releases |
| `agent-fixed` | Production-ready agent with LangGraph | AMD64, ARM64 | On stable releases |
| `agent-latest` | Development agent builds | AMD64, ARM64 | On commits |

### 📦 **Multi-Architecture Tags**
```bash
# These tags automatically serve the correct architecture:
sudeshmu/kubeagentic:operator-latest  # ✅ AMD64 + ARM64
sudeshmu/kubeagentic:agent-fixed      # ✅ AMD64 + ARM64 (RECOMMENDED)
sudeshmu/kubeagentic:agent-latest     # ✅ AMD64 + ARM64 (Development)

# Architecture-specific tags (if needed):
sudeshmu/kubeagentic:operator-latest@sha256:2335acc4...  # AMD64 only
sudeshmu/kubeagentic:operator-latest@sha256:08d4833d...  # ARM64 only
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](https://github.com/KubeAgentic-Community/KubeAgentic/blob/main/CONTRIBUTING.md) for details.

## 📄 License

Apache License 2.0 - see [LICENSE](https://github.com/KubeAgentic-Community/KubeAgentic/blob/main/LICENSE)

---

**🚀 Built with ❤️ for the Multi-Architecture Kubernetes AI Community**

*These multi-architecture images are built using Docker Buildx with Red Hat UBI base images, ensuring security, performance, and native execution on both AMD64 and ARM64 architectures. Deploy once, run anywhere!*

### 🎯 **Why Multi-Architecture?**
- **💰 Cost savings**: Use ARM-based cloud instances (up to 50% cheaper)  
- **⚡ Performance**: Native execution on Apple Silicon and ARM servers
- **🌍 Flexibility**: Deploy on any Kubernetes cluster architecture
- **🔮 Future-proof**: Ready for the ARM64 adoption wave

**Experience the power of truly portable AI agents! 🤖✨**
