---
layout: home
title: Home
---

<style>
/* Custom CSS for better styling */
body { margin: 0; padding: 0; }
.wrapper { max-width: none !important; }
.page-content { max-width: none !important; }
.site-header { max-width: none !important; }
.site-footer { max-width: none !important; }

/* Hero section with gradient background */
.hero-section {
  background: linear-gradient(135deg, #42b883 0%, #2c3e50 100%);
  color: white;
  padding: 4rem 0;
  text-align: center;
  width: 100%;
  margin: 0;
}

/* Button hover effects */
.cta-button {
  transition: all 0.3s ease;
  display: inline-block;
  text-decoration: none;
  font-weight: 500;
  border-radius: 6px;
  padding: 12px 24px;
  margin: 0 0.5rem;
}

.cta-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.2);
}

.cta-primary {
  background: #42b883;
  color: white;
}

.cta-secondary {
  background: transparent;
  color: #42b883;
  border: 2px solid #42b883;
}

.cta-primary:hover {
  background: #369870;
  color: white;
}

.cta-secondary:hover {
  background: #42b883;
  color: white;
}

/* Code block styling */
.code-block {
  background: #1e1e1e;
  color: #d4d4d4;
  padding: 2rem;
  border-radius: 8px;
  overflow-x: auto;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.5;
  margin: 2rem 0;
}

/* Feature icons */
.feature-icon {
  font-size: 2.5rem;
  margin-bottom: 1rem;
  display: block;
}

/* Full width sections */
.full-width {
  width: 100%;
  margin: 0;
  padding: 0;
}

/* External link indicator */
.external-link::after {
  content: " ↗";
  font-size: 0.8em;
  opacity: 0.7;
}
</style>

<!-- Hero Section -->
<div class="hero-section">
  <div style="max-width: 1200px; margin: 0 auto; padding: 0 1rem;">
    <img src="{{ '/assets/logo.png' | relative_url }}" alt="KubeAgentic Logo - AI agents on Kubernetes" style="max-width: 300px; height: auto; margin-bottom: 2rem;">
    <h1 style="font-size: 4rem; font-weight: 700; margin-bottom: 1rem; text-shadow: 0 2px 4px rgba(0,0,0,0.3);">
      Deploy AI Agents on Kubernetes
    </h1>
    <p style="font-size: 1.5rem; margin-bottom: 2rem; opacity: 0.9; max-width: 800px; margin-left: auto; margin-right: auto;">
      The simplest way to deploy, manage, and scale AI agents in your Kubernetes cluster with declarative YAML configurations
    </p>
    <div style="margin-bottom: 2rem;">
      <a href="#quick-start" class="cta-button cta-primary">
        🚀 Get Started Now
      </a>
      <a href="https://github.com/sudeshmu/KubeAgentic" class="cta-button cta-secondary external-link">
        View on GitHub
      </a>
    </div>
  </div>
</div>

<!-- Key Features Section -->
<div style="background: #f8f9fa; padding: 4rem 0;">
  <div style="max-width: 1400px; margin: 0 auto; padding: 0 1rem;">
    <h2 style="text-align: center; margin-bottom: 3rem; color: #2c3e50; font-size: 2.5rem;">✨ Why Choose KubeAgentic?</h2>
    <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(350px, 1fr)); gap: 2rem;">
      <div style="background: white; padding: 2rem; border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center; transition: transform 0.3s ease;">
        <span class="feature-icon">🤖</span>
        <h3 style="color: #42b883; margin-bottom: 1rem; font-size: 1.3rem;">Multi-Provider Support</h3>
        <p>OpenAI, Anthropic (Claude), Google (Gemini), and self-hosted vLLM models - all in one platform</p>
      </div>
      <div style="background: white; padding: 2rem; border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center; transition: transform 0.3s ease;">
        <span class="feature-icon">📝</span>
        <h3 style="color: #42b883; margin-bottom: 1rem; font-size: 1.3rem;">Declarative Configuration</h3>
        <p>Standard Kubernetes Custom Resources (CRDs) for easy management and GitOps workflows</p>
      </div>
      <div style="background: white; padding: 2rem; border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center; transition: transform 0.3s ease;">
        <span class="feature-icon">🔄</span>
        <h3 style="color: #42b883; margin-bottom: 1rem; font-size: 1.3rem;">Auto-Scaling</h3>
        <p>Automatic scaling based on demand and resource usage with Kubernetes HPA integration</p>
      </div>
      <div style="background: white; padding: 2rem; border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center; transition: transform 0.3s ease;">
        <span class="feature-icon">🔒</span>
        <h3 style="color: #42b883; margin-bottom: 1rem; font-size: 1.3rem;">Secure by Default</h3>
        <p>API keys managed with Kubernetes Secrets and RBAC for enterprise-grade security</p>
      </div>
      <div style="background: white; padding: 2rem; border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center; transition: transform 0.3s ease;">
        <span class="feature-icon">📊</span>
        <h3 style="color: #42b883; margin-bottom: 1rem; font-size: 1.3rem;">Built-in Monitoring</h3>
        <p>Real-time health checks, metrics, and status reporting with Prometheus integration</p>
      </div>
      <div style="background: white; padding: 2rem; border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center; transition: transform 0.3s ease;">
        <span class="feature-icon">🛠️</span>
        <h3 style="color: #42b883; margin-bottom: 1rem; font-size: 1.3rem;">Tool Integration</h3>
        <p>Extend agents with custom tools, APIs, and services for complex workflows</p>
      </div>
    </div>
  </div>
</div>

<!-- Quick Start Section -->
<div id="quick-start" style="padding: 4rem 0; background: white;">
  <div style="max-width: 1400px; margin: 0 auto; padding: 0 1rem;">
    <h2 style="text-align: center; margin-bottom: 3rem; color: #2c3e50; font-size: 2.5rem;">🚀 Quick Start</h2>
    
    <div class="code-block">
# Install KubeAgentic
kubectl apply -f https://raw.githubusercontent.com/sudeshmu/KubeAgentic/main/deploy/all.yaml

# Create API key secret
kubectl create secret generic openai-secret \
  --from-literal=api-key='your-openai-api-key'

# Deploy your first agent
kubectl apply -f https://raw.githubusercontent.com/sudeshmu/KubeAgentic/main/examples/openai-agent.yaml

# Interact with your agent
kubectl port-forward service/my-assistant-service 8080:80
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello! How can you help me?"}'
    </div>
    
    <div style="text-align: center; margin-top: 2rem;">
      <a href="https://github.com/sudeshmu/KubeAgentic/tree/main/examples" class="cta-button cta-primary external-link">
        View All Examples
      </a>
    </div>
  </div>
</div>

<!-- Architecture Diagram -->
<div style="background: #f8f9fa; padding: 4rem 0;">
  <div style="max-width: 1400px; margin: 0 auto; padding: 0 1rem;">
    <h2 style="text-align: center; margin-bottom: 3rem; color: #2c3e50; font-size: 2.5rem;">🏗️ How It Works</h2>
    <div style="background: white; padding: 3rem; border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center;">
      <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 2rem; margin-bottom: 2rem;">
        <div>
          <div style="background: #42b883; color: white; padding: 1rem; border-radius: 50%; width: 80px; height: 80px; margin: 0 auto 1rem; display: flex; align-items: center; justify-content: center; font-size: 2rem;">☸️</div>
          <h4 style="color: #2c3e50; margin-bottom: 0.5rem;">Kubernetes Cluster</h4>
          <p style="color: #7f8c8d; font-size: 0.9rem;">Your existing K8s infrastructure</p>
        </div>
        <div style="display: flex; align-items: center; justify-content: center; font-size: 2rem; color: #42b883;">→</div>
        <div>
          <div style="background: #42b883; color: white; padding: 1rem; border-radius: 50%; width: 80px; height: 80px; margin: 0 auto 1rem; display: flex; align-items: center; justify-content: center; font-size: 2rem;">🤖</div>
          <h4 style="color: #2c3e50; margin-bottom: 0.5rem;">KubeAgentic Operator</h4>
          <p style="color: #7f8c8d; font-size: 0.9rem;">Deploy & manage AI agents</p>
        </div>
        <div style="display: flex; align-items: center; justify-content: center; font-size: 2rem; color: #42b883;">→</div>
        <div>
          <div style="background: #42b883; color: white; padding: 1rem; border-radius: 50%; width: 80px; height: 80px; margin: 0 auto 1rem; display: flex; align-items: center; justify-content: center; font-size: 2rem;">🧠</div>
          <h4 style="color: #2c3e50; margin-bottom: 0.5rem;">AI Models</h4>
          <p style="color: #7f8c8d; font-size: 0.9rem;">OpenAI, Claude, Gemini, vLLM</p>
        </div>
        <div style="display: flex; align-items: center; justify-content: center; font-size: 2rem; color: #42b883;">→</div>
        <div>
          <div style="background: #42b883; color: white; padding: 1rem; border-radius: 50%; width: 80px; height: 80px; margin: 0 auto 1rem; display: flex; align-items: center; justify-content: center; font-size: 2rem;">📈</div>
          <h4 style="color: #2c3e50; margin-bottom: 0.5rem;">Auto-Scaling</h4>
          <p style="color: #7f8c8d; font-size: 0.9rem;">Monitor & scale automatically</p>
        </div>
      </div>
    </div>
  </div>
</div>

<!-- Documentation Section -->
<div style="background: white; padding: 4rem 0;">
  <div style="max-width: 1400px; margin: 0 auto; padding: 0 1rem;">
    <h2 style="text-align: center; margin-bottom: 3rem; color: #2c3e50; font-size: 2.5rem;">📚 Documentation</h2>
    
    <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); gap: 2rem;">
      <div style="background: #f8f9fa; padding: 2rem; border-radius: 12px; text-align: center; border: 2px solid transparent; transition: all 0.3s ease;">
        <div style="font-size: 3rem; margin-bottom: 1rem;">⚡</div>
        <h3 style="color: #42b883; margin-bottom: 1rem; font-size: 1.3rem;">Direct Framework</h3>
        <p style="margin-bottom: 1.5rem; color: #7f8c8d;">Simple, fast API calls for basic interactions</p>
        <a href="direct-framework" class="cta-button cta-secondary">Learn More →</a>
      </div>
      
      <div style="background: #f8f9fa; padding: 2rem; border-radius: 12px; text-align: center; border: 2px solid transparent; transition: all 0.3s ease;">
        <div style="font-size: 3rem; margin-bottom: 1rem;">🔗</div>
        <h3 style="color: #42b883; margin-bottom: 1rem; font-size: 1.3rem;">LangGraph Framework</h3>
        <p style="margin-bottom: 1.5rem; color: #7f8c8d;">Complex workflows with multi-step reasoning</p>
        <a href="langgraph-framework" class="cta-button cta-secondary">Learn More →</a>
      </div>
      
      <div style="background: #f8f9fa; padding: 2rem; border-radius: 12px; text-align: center; border: 2px solid transparent; transition: all 0.3s ease;">
        <div style="font-size: 3rem; margin-bottom: 1rem;">🔧</div>
        <h3 style="color: #42b883; margin-bottom: 1rem; font-size: 1.3rem;">API Reference</h3>
        <p style="margin-bottom: 1.5rem; color: #7f8c8d;">Detailed API specification and examples</p>
        <a href="api-reference" class="cta-button cta-secondary">View Docs →</a>
      </div>
      
      <div style="background: #f8f9fa; padding: 2rem; border-radius: 12px; text-align: center; border: 2px solid transparent; transition: all 0.3s ease;">
        <div style="font-size: 3rem; margin-bottom: 1rem;">💡</div>
        <h3 style="color: #42b883; margin-bottom: 1rem; font-size: 1.3rem;">Examples</h3>
        <p style="margin-bottom: 1.5rem; color: #7f8c8d;">Real-world usage examples and templates</p>
        <a href="examples" class="cta-button cta-secondary">View Examples →</a>
      </div>
    </div>
  </div>
</div>

<!-- Use Cases Section -->
<div style="background: #f8f9fa; padding: 4rem 0;">
  <div style="max-width: 1400px; margin: 0 auto; padding: 0 1rem;">
    <h2 style="text-align: center; margin-bottom: 3rem; color: #2c3e50; font-size: 2.5rem;">🎯 Use Cases</h2>
    
    <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 2rem;">
      <div style="background: white; padding: 2rem; border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center;">
        <div style="font-size: 3rem; margin-bottom: 1rem;">🎧</div>
        <h3 style="color: #42b883; margin-bottom: 1rem; font-size: 1.3rem;">Customer Support</h3>
        <p>Deploy scalable support bots that can handle multiple conversations simultaneously with context awareness</p>
      </div>
      <div style="background: white; padding: 2rem; border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center;">
        <div style="font-size: 3rem; margin-bottom: 1rem;">🔍</div>
        <h3 style="color: #42b883; margin-bottom: 1rem; font-size: 1.3rem;">Code Review</h3>
        <p>Automated code analysis and feedback for improved code quality, security, and best practices</p>
      </div>
      <div style="background: white; padding: 2rem; border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center;">
        <div style="font-size: 3rem; margin-bottom: 1rem;">📚</div>
        <h3 style="color: #42b883; margin-bottom: 1rem; font-size: 1.3rem;">Knowledge Management</h3>
        <p>Internal Q&A assistants for company documentation, procedures, and knowledge base queries</p>
      </div>
      <div style="background: white; padding: 2rem; border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center;">
        <div style="font-size: 3rem; margin-bottom: 1rem;">✍️</div>
        <h3 style="color: #42b883; margin-bottom: 1rem; font-size: 1.3rem;">Content Generation</h3>
        <p>AI-powered content creation for marketing, documentation, and automated report generation</p>
      </div>
    </div>
  </div>
</div>

<!-- Footer -->
<div style="background: #2c3e50; color: #ecf0f1; padding: 4rem 0 2rem; width: 100%; margin: 0;">
  <div style="max-width: 1400px; margin: 0 auto; padding: 0 1rem; text-align: center;">
    <h2 style="margin-bottom: 3rem; color: #ecf0f1; font-size: 2rem;">🤝 Community & Support</h2>
    
    <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 2rem; margin-bottom: 3rem;">
      <div>
        <h3 style="color: #42b883; margin-bottom: 1rem; font-size: 1.2rem;">GitHub</h3>
        <p style="margin-bottom: 1rem; color: #bdc3c7;">Star us on GitHub and contribute</p>
        <a href="https://github.com/sudeshmu/KubeAgentic" class="cta-button cta-secondary external-link" style="color: #42b883; border-color: #42b883;">View Repository →</a>
      </div>
      <div>
        <h3 style="color: #42b883; margin-bottom: 1rem; font-size: 1.2rem;">Issues</h3>
        <p style="margin-bottom: 1rem; color: #bdc3c7;">Report bugs and request features</p>
        <a href="https://github.com/sudeshmu/KubeAgentic/issues" class="cta-button cta-secondary external-link" style="color: #42b883; border-color: #42b883;">Report Issue →</a>
      </div>
      <div>
        <h3 style="color: #42b883; margin-bottom: 1rem; font-size: 1.2rem;">Discussions</h3>
        <p style="margin-bottom: 1rem; color: #bdc3c7;">Join the community discussions</p>
        <a href="https://github.com/sudeshmu/KubeAgentic/discussions" class="cta-button cta-secondary external-link" style="color: #42b883; border-color: #42b883;">Join Discussion →</a>
      </div>
      <div>
        <h3 style="color: #42b883; margin-bottom: 1rem; font-size: 1.2rem;">Contact</h3>
        <p style="margin-bottom: 1rem; color: #bdc3c7;">Get in touch with us directly</p>
        <a href="mailto:contact@kubeagentic.com" class="cta-button cta-secondary" style="color: #42b883; border-color: #42b883;">contact@kubeagentic.com</a>
      </div>
    </div>
    
    <div style="border-top: 1px solid #34495e; padding-top: 2rem; margin-top: 2rem;">
      <p style="color: #95a5a6; margin-bottom: 1rem;">
        Licensed under the Apache License 2.0. See <a href="https://github.com/sudeshmu/KubeAgentic/blob/main/LICENSE" style="color: #42b883;">LICENSE</a> for details.
      </p>
      <p style="color: #95a5a6;">
        © 2025 KubeAgentic. Built with ❤️ for the Kubernetes community.
      </p>
    </div>
  </div>
</div>
