---
layout: home
title: Home
---

<div style="text-align: center; margin-bottom: 3rem;">
  <img src="/assets/logo.jpg" alt="KubeAgentic Logo" style="max-width: 200px; height: auto;">
</div>

<div style="text-align: center; margin-bottom: 4rem;">
  <h1 style="font-size: 3.5rem; font-weight: 600; margin-bottom: 1rem; color: #2c3e50;">KubeAgentic</h1>
  <p style="font-size: 1.25rem; color: #7f8c8d; margin-bottom: 2rem; max-width: 600px; margin-left: auto; margin-right: auto;">
    Deploy and manage AI agents on Kubernetes with simple YAML configurations
  </p>
  <div style="margin-bottom: 2rem;">
    <a href="#quick-start" style="background: #42b883; color: white; padding: 12px 24px; border-radius: 6px; text-decoration: none; font-weight: 500; margin-right: 1rem; display: inline-block;">
      Get Started
    </a>
    <a href="https://github.com/sudeshmu/KubeAgentic" style="background: transparent; color: #42b883; padding: 12px 24px; border: 2px solid #42b883; border-radius: 6px; text-decoration: none; font-weight: 500; display: inline-block;">
      View on GitHub
    </a>
  </div>
</div>

<div style="background: #f8f9fa; padding: 3rem 0; margin: 3rem 0; border-radius: 8px;">
  <div style="max-width: 1200px; margin: 0 auto; padding: 0 2rem;">
    <h2 style="text-align: center; margin-bottom: 2rem; color: #2c3e50;">✨ Key Features</h2>
    <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 2rem;">
      <div style="background: white; padding: 1.5rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
        <h3 style="color: #42b883; margin-bottom: 1rem;">🤖 Multi-Provider Support</h3>
        <p>OpenAI, Anthropic (Claude), Google (Gemini), and self-hosted vLLM models</p>
      </div>
      <div style="background: white; padding: 1.5rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
        <h3 style="color: #42b883; margin-bottom: 1rem;">📝 Declarative Configuration</h3>
        <p>Standard Kubernetes Custom Resources (CRDs) for easy management</p>
      </div>
      <div style="background: white; padding: 1.5rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
        <h3 style="color: #42b883; margin-bottom: 1rem;">🔄 Autoscaling</h3>
        <p>Automatic scaling based on demand and resource usage</p>
      </div>
      <div style="background: white; padding: 1.5rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
        <h3 style="color: #42b883; margin-bottom: 1rem;">🔒 Secure by Default</h3>
        <p>API keys managed with Kubernetes Secrets</p>
      </div>
      <div style="background: white; padding: 1.5rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
        <h3 style="color: #42b883; margin-bottom: 1rem;">📊 Built-in Monitoring</h3>
        <p>Real-time health checks and status reporting</p>
      </div>
      <div style="background: white; padding: 1.5rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
        <h3 style="color: #42b883; margin-bottom: 1rem;">🛠️ Tool Integration</h3>
        <p>Extend agents with custom tools and services</p>
      </div>
    </div>
  </div>
</div>

<div id="quick-start" style="max-width: 1200px; margin: 0 auto; padding: 0 2rem;">
  <h2 style="text-align: center; margin-bottom: 2rem; color: #2c3e50;">🚀 Quick Start</h2>
  
  <div style="background: #2c3e50; color: #ecf0f1; padding: 1.5rem; border-radius: 8px; margin-bottom: 2rem; overflow-x: auto;">
    <pre style="margin: 0; font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;"><code># Install KubeAgentic
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
  -d '{"message": "Hello! How can you help me?"}'</code></pre>
  </div>
</div>

<div style="background: #f8f9fa; padding: 3rem 0; margin: 3rem 0; border-radius: 8px;">
  <div style="max-width: 1200px; margin: 0 auto; padding: 0 2rem;">
    <h2 style="text-align: center; margin-bottom: 2rem; color: #2c3e50;">📚 Documentation</h2>
    
    <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 2rem;">
      <div style="background: white; padding: 1.5rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); text-align: center;">
        <h3 style="color: #42b883; margin-bottom: 1rem;">⚡ Direct Framework</h3>
        <p style="margin-bottom: 1rem;">Simple, fast API calls for basic interactions</p>
        <a href="direct-framework" style="color: #42b883; text-decoration: none; font-weight: 500;">Learn More →</a>
      </div>
      
      <div style="background: white; padding: 1.5rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); text-align: center;">
        <h3 style="color: #42b883; margin-bottom: 1rem;">🔗 LangGraph Framework</h3>
        <p style="margin-bottom: 1rem;">Complex workflows with multi-step reasoning</p>
        <a href="langgraph-framework" style="color: #42b883; text-decoration: none; font-weight: 500;">Learn More →</a>
      </div>
      
      <div style="background: white; padding: 1.5rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); text-align: center;">
        <h3 style="color: #42b883; margin-bottom: 1rem;">🔧 API Reference</h3>
        <p style="margin-bottom: 1rem;">Detailed API specification</p>
        <a href="api-reference" style="color: #42b883; text-decoration: none; font-weight: 500;">View Docs →</a>
      </div>
      
      <div style="background: white; padding: 1.5rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); text-align: center;">
        <h3 style="color: #42b883; margin-bottom: 1rem;">💡 Examples</h3>
        <p style="margin-bottom: 1rem;">Real-world usage examples</p>
        <a href="examples" style="color: #42b883; text-decoration: none; font-weight: 500;">View Examples →</a>
      </div>
    </div>
  </div>
</div>

<div style="max-width: 1200px; margin: 0 auto; padding: 0 2rem;">
  <h2 style="text-align: center; margin-bottom: 2rem; color: #2c3e50;">🎯 Use Cases</h2>
  
  <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 2rem; margin-bottom: 3rem;">
    <div style="text-align: center; padding: 1.5rem;">
      <h3 style="color: #42b883; margin-bottom: 1rem;">Customer Support</h3>
      <p>Deploy scalable support bots that can handle multiple conversations simultaneously</p>
    </div>
    <div style="text-align: center; padding: 1.5rem;">
      <h3 style="color: #42b883; margin-bottom: 1rem;">Code Review</h3>
      <p>Automated code analysis and feedback for improved code quality</p>
    </div>
    <div style="text-align: center; padding: 1.5rem;">
      <h3 style="color: #42b883; margin-bottom: 1rem;">Knowledge Management</h3>
      <p>Internal Q&A assistants for company documentation and procedures</p>
    </div>
    <div style="text-align: center; padding: 1.5rem;">
      <h3 style="color: #42b883; margin-bottom: 1rem;">Content Generation</h3>
      <p>AI-powered content creation for marketing and documentation</p>
    </div>
  </div>
</div>

<div style="background: #2c3e50; color: #ecf0f1; padding: 3rem 0; margin-top: 4rem;">
  <div style="max-width: 1200px; margin: 0 auto; padding: 0 2rem; text-align: center;">
    <h2 style="margin-bottom: 2rem; color: #ecf0f1;">🤝 Community & Support</h2>
    
    <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 2rem; margin-bottom: 2rem;">
      <div>
        <h3 style="color: #42b883; margin-bottom: 1rem;">GitHub</h3>
        <p style="margin-bottom: 1rem;">Star us on GitHub</p>
        <a href="https://github.com/sudeshmu/KubeAgentic" style="color: #42b883; text-decoration: none; font-weight: 500;">View Repository →</a>
      </div>
      <div>
        <h3 style="color: #42b883; margin-bottom: 1rem;">Issues</h3>
        <p style="margin-bottom: 1rem;">Report bugs and request features</p>
        <a href="https://github.com/sudeshmu/KubeAgentic/issues" style="color: #42b883; text-decoration: none; font-weight: 500;">Report Issue →</a>
      </div>
      <div>
        <h3 style="color: #42b883; margin-bottom: 1rem;">Discussions</h3>
        <p style="margin-bottom: 1rem;">Join the community</p>
        <a href="https://github.com/sudeshmu/KubeAgentic/discussions" style="color: #42b883; text-decoration: none; font-weight: 500;">Join Discussion →</a>
      </div>
      <div>
        <h3 style="color: #42b883; margin-bottom: 1rem;">Contact</h3>
        <p style="margin-bottom: 1rem;">Get in touch with us</p>
        <a href="mailto:contact@kubeagentic.com" style="color: #42b883; text-decoration: none; font-weight: 500;">contact@kubeagentic.com</a>
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
