# 🏗️ Docker Multi-Architecture Build Guide

This guide explains how to build KubeAgentic Docker images for multiple architectures (AMD64 and ARM64).

## 🚀 Quick Start

### Build and Push Multi-Architecture Images

```bash
# Build and push both operator and agent images for AMD64 and ARM64
make docker-buildx-all

# Or individually:
make docker-buildx-operator  # Build operator image
make docker-buildx-agent     # Build agent image
```

### Local Development (Build Only, No Push)

```bash
# Build locally for testing (no push to registry)
make docker-buildx-local-all

# Or individually:
make docker-buildx-local-operator
make docker-buildx-local-agent
```

## 🔧 Configuration

### Image Tags

Configure image names and tags:

```bash
# Set custom image names
export OPERATOR_IMG=your-registry/kubeagentic:operator-v1.0.0
export AGENT_IMG=your-registry/kubeagentic:agent-v1.0.0

# Build with custom images
make docker-buildx-all
```

### Supported Platforms

Default platforms: `linux/amd64,linux/arm64`

```bash
# Build for specific platforms
export PLATFORMS=linux/amd64,linux/arm64,linux/arm/v7
make docker-buildx-all
```

### Builder Configuration

```bash
# Use custom builder name
export BUILDX_BUILDER=my-custom-builder
make buildx-setup
```

## 📋 Available Commands

### Setup and Initialization

- `make buildx-setup` - Initialize buildx builder for multi-arch builds

### Multi-Architecture Build Commands

- `make docker-buildx-all` - Build and push all images (multi-arch)
- `make docker-buildx-operator` - Build and push operator image (multi-arch)
- `make docker-buildx-agent` - Build and push agent image (multi-arch)
- `make docker-buildx-local-all` - Build all images locally (no push)
- `make docker-buildx-local-operator` - Build operator image locally (no push)
- `make docker-buildx-local-agent` - Build agent image locally (no push)

### Legacy Single-Architecture Commands

- `make docker-build-operator` - Build operator image (single arch)
- `make docker-build-agent` - Build agent image (single arch)
- `make docker-push-operator` - Push operator image
- `make docker-push-agent` - Push agent image

### Image Management

- `make inspect-images` - Inspect multi-architecture image manifests
- `make buildx-cleanup` - Clean up buildx builder

### Complete Workflows

- `make complete-deploy` - Build, push, and deploy (multi-arch by default)
- `make dev-deploy` - Build and deploy for development (multi-arch, no push)
- `make complete-deploy-single-arch` - Build, push, and deploy (single arch)
- `make dev-deploy-single-arch` - Build and deploy (single arch, no push)

## 🛠️ Alternative: Using the Build Script

You can also use the standalone script:

```bash
# Build using script (reads environment variables)
./scripts/build-multiarch.sh

# With custom settings
export OPERATOR_IMG=my-registry/kubeagentic:operator-v2.0.0
export AGENT_IMG=my-registry/kubeagentic:agent-v2.0.0
./scripts/build-multiarch.sh
```

## 🔍 Verification

### Check Multi-Architecture Support

```bash
# Inspect image manifests
make inspect-images

# Or manually:
docker buildx imagetools inspect sudeshmu/kubeagentic:operator-latest
docker buildx imagetools inspect sudeshmu/kubeagentic:agent-fixed
```

### Test Images on Different Architectures

```bash
# Test operator
docker run --rm --platform linux/amd64 sudeshmu/kubeagentic:operator-latest --help
docker run --rm --platform linux/arm64 sudeshmu/kubeagentic:operator-latest --help

# Test agent
docker run --rm --platform linux/amd64 sudeshmu/kubeagentic:agent-fixed python --version
docker run --rm --platform linux/arm64 sudeshmu/kubeagentic:agent-fixed python --version
```

## 🏗️ Architecture Support

### Supported Architectures

- **AMD64** (`linux/amd64`) - Intel/AMD x86_64
- **ARM64** (`linux/arm64`) - Apple Silicon, AWS Graviton, ARM servers

### Cloud Provider Compatibility

| Provider | AMD64 | ARM64 | Notes |
|----------|-------|-------|-------|
| AWS | ✅ | ✅ | Graviton instances |
| GCP | ✅ | ✅ | T2A instances |
| Azure | ✅ | ✅ | Ampere Altra |
| Docker Desktop | ✅ | ✅ | M1/M2 Macs |
| Local K8s | ✅ | ✅ | kind, minikube |

## ⚠️ Troubleshooting

### Common Issues

1. **Buildx not available**: Install Docker Desktop or enable buildx
2. **Builder not found**: Run `make buildx-setup`
3. **Permission denied**: Ensure Docker Hub login: `docker login`
4. **Platform not supported**: Check available platforms with:
   ```bash
   docker buildx inspect kubeagentic-builder
   ```

### Reset Builder

```bash
# Clean up and recreate builder
make buildx-cleanup
make buildx-setup
```

## 🚀 CI/CD Integration

For automated builds in CI/CD:

```yaml
# GitHub Actions example
- name: Set up Docker Buildx
  uses: docker/setup-buildx-action@v3

- name: Build and push
  run: make docker-buildx-all
  env:
    OPERATOR_IMG: ${{ secrets.REGISTRY }}/kubeagentic:operator-${{ github.sha }}
    AGENT_IMG: ${{ secrets.REGISTRY }}/kubeagentic:agent-${{ github.sha }}
```

## 📊 Image Sizes

Multi-architecture builds create a single manifest that points to architecture-specific images:

- **Operator**: ~50MB (UBI micro + Go binary)
- **Agent**: ~200MB (UBI minimal + Python + dependencies)

When pulled, Kubernetes automatically selects the correct architecture.
