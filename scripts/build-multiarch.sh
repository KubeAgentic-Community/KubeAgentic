#!/bin/bash

# Multi-Architecture Build Script for KubeAgentic
# Builds and pushes both operator and agent images for AMD64 and ARM64

set -e

echo "🏗️ KubeAgentic Multi-Architecture Build"
echo "======================================="

# Check if buildx is available
if ! command -v docker &> /dev/null || ! docker buildx version &> /dev/null; then
    echo "❌ Docker buildx is not available"
    echo "Please install Docker with buildx support"
    exit 1
fi

# Create or use multiarch builder
echo "🔧 Setting up buildx builder..."
if ! docker buildx inspect multiarch &> /dev/null; then
    echo "Creating new multiarch builder..."
    docker buildx create --name multiarch --driver docker-container --use
    docker buildx inspect --bootstrap
else
    echo "Using existing multiarch builder..."
    docker buildx use multiarch
fi

# Configuration
PLATFORMS="linux/amd64,linux/arm64"
OPERATOR_IMAGE="sudeshmu/kubeagentic:operator-latest"
AGENT_IMAGE="sudeshmu/kubeagentic:agent-latest"

# Build and push operator image
echo ""
echo "🏗️ Building operator image for platforms: $PLATFORMS"
echo "Image: $OPERATOR_IMAGE"
docker buildx build \
    --platform $PLATFORMS \
    -f Dockerfile.operator \
    -t $OPERATOR_IMAGE \
    --push .

echo "✅ Operator image built and pushed successfully!"

# Build and push agent image
echo ""
echo "🏗️ Building agent image for platforms: $PLATFORMS"
echo "Image: $AGENT_IMAGE"
docker buildx build \
    --platform $PLATFORMS \
    -f Dockerfile.agent \
    -t $AGENT_IMAGE \
    --push .

echo "✅ Agent image built and pushed successfully!"

# Verify multi-architecture support
echo ""
echo "🔍 Verifying multi-architecture support..."
echo ""
echo "Operator image manifests:"
docker buildx imagetools inspect $OPERATOR_IMAGE

echo ""
echo "Agent image manifests:"
docker buildx imagetools inspect $AGENT_IMAGE

echo ""
echo "🎉 Multi-architecture build complete!"
echo ""
echo "📊 Image Information:"
echo "   Operator: $OPERATOR_IMAGE (supports AMD64 + ARM64)"
echo "   Agent:    $AGENT_IMAGE (supports AMD64 + ARM64)"
echo ""
echo "🚀 Your images are now available on both:"
echo "   - x86_64/AMD64 (Intel/AMD servers, VMs)"
echo "   - ARM64 (Apple Silicon, ARM servers)"
echo ""
echo "Test on different architectures:"
echo "   docker run --rm $OPERATOR_IMAGE --help"
echo "   docker run --rm $AGENT_IMAGE python --version"
