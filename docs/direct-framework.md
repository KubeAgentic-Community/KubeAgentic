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

## Use Cases & Examples

### 1. Customer Support Chat Bot

Simple customer support agent that can answer FAQs and look up basic information.

```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: support-chatbot
  namespace: customer-service
spec:
  framework: direct
  provider: openai
  model: gpt-4
  systemPrompt: |
    You are a friendly customer support agent for TechCorp.
    
    Guidelines:
    - Always greet customers warmly
    - Provide helpful and accurate information
    - Be concise but thorough
    - If you cannot help, suggest contacting human support
    - Use the order_lookup tool when customers ask about orders
    
    Company Info:
    - Business hours: Mon-Fri 9AM-6PM EST
    - Return policy: 30 days with receipt
    - Shipping: 3-5 business days standard
    
  apiSecretRef:
    name: openai-secret
    key: api-key
  
  tools:
  - name: order_lookup
    description: Look up customer order status and details
    inputSchema:
      type: object
      properties:
        order_id:
          type: string
          description: Customer order ID (format: ORD-XXXXXX)
        email:
          type: string
          format: email
          description: Customer email address for verification
      required: ["order_id"]
  
  - name: product_search
    description: Search for product information and availability
    inputSchema:
      type: object
      properties:
        product_name:
          type: string
          description: Product name or keyword to search
        category:
          type: string
          enum: ["electronics", "clothing", "home", "books"]
          description: Product category filter
      required: ["product_name"]
  
  replicas: 3
  resources:
    requests:
      cpu: 100m
      memory: 256Mi
    limits:
      cpu: 200m
      memory: 512Mi
```

### 2. Code Review Assistant

Technical agent that helps with code reviews and programming questions.

```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: code-reviewer
  namespace: engineering
spec:
  framework: direct
  provider: claude
  model: claude-3-sonnet
  systemPrompt: |
    You are a senior software engineer providing code review feedback.
    
    Focus Areas:
    - Code quality and best practices
    - Security vulnerabilities
    - Performance optimizations
    - Readability and maintainability
    - Testing recommendations
    
    Style:
    - Be constructive and educational
    - Provide specific examples
    - Suggest improvements with rationale
    - Acknowledge good practices when you see them
    
  apiSecretRef:
    name: claude-secret
    key: api-key
  
  tools:
  - name: code_analysis
    description: Analyze code for common issues and patterns
    inputSchema:
      type: object
      properties:
        code:
          type: string
          description: Code snippet to analyze
        language:
          type: string
          enum: ["python", "javascript", "go", "java", "typescript"]
          description: Programming language
        context:
          type: string
          description: Additional context about the code's purpose
      required: ["code", "language"]
  
  - name: security_scan
    description: Scan code for potential security vulnerabilities
    inputSchema:
      type: object
      properties:
        code:
          type: string
          description: Code to scan for security issues
        framework:
          type: string
          description: Framework or library being used
      required: ["code"]
  
  replicas: 2
  resources:
    requests:
      cpu: 200m
      memory: 512Mi
    limits:
      cpu: 500m
      memory: 1Gi
```

### 3. Content Generation Agent

Marketing content generator for social media, blogs, and campaigns.

```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: content-generator
  namespace: marketing
spec:
  framework: direct
  provider: openai
  model: gpt-4
  systemPrompt: |
    You are a creative content writer specializing in marketing content.
    
    Expertise:
    - Social media posts (Twitter, LinkedIn, Instagram)
    - Blog post outlines and content
    - Email marketing campaigns
    - Product descriptions
    - Ad copy and headlines
    
    Style Guidelines:
    - Match the requested tone and brand voice
    - Use engaging, action-oriented language
    - Include relevant hashtags for social media
    - Keep content concise and impactful
    - Always consider the target audience
    
  apiSecretRef:
    name: openai-secret
    key: api-key
  
  tools:
  - name: brand_guidelines
    description: Get brand guidelines and voice for content creation
    inputSchema:
      type: object
      properties:
        brand:
          type: string
          description: Brand name to get guidelines for
        content_type:
          type: string
          enum: ["social", "blog", "email", "ad", "product"]
          description: Type of content being created
      required: ["brand", "content_type"]
  
  - name: trend_analysis
    description: Get current trends and hashtags for content optimization
    inputSchema:
      type: object
      properties:
        platform:
          type: string
          enum: ["twitter", "linkedin", "instagram", "tiktok", "facebook"]
          description: Social media platform
        industry:
          type: string
          description: Industry or niche for trend analysis
      required: ["platform"]
  
  replicas: 2
  serviceType: LoadBalancer  # External access for marketing team
```

### 4. Data Analysis Assistant

Agent that helps with data interpretation and basic analytics.

```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: data-analyst
  namespace: analytics
spec:
  framework: direct
  provider: openai
  model: gpt-4
  systemPrompt: |
    You are a data analyst specializing in business intelligence and reporting.
    
    Capabilities:
    - Interpret charts, graphs, and statistical data
    - Explain trends and patterns in business metrics
    - Suggest data-driven recommendations
    - Help with basic statistical analysis
    - Create simple data visualizations descriptions
    
    Communication Style:
    - Use clear, non-technical language for business stakeholders
    - Provide actionable insights
    - Highlight key findings and recommendations
    - Ask clarifying questions when data context is unclear
    
  apiSecretRef:
    name: openai-secret
    key: api-key
  
  tools:
  - name: metric_lookup
    description: Look up specific business metrics and KPIs
    inputSchema:
      type: object
      properties:
        metric_name:
          type: string
          description: Name of the metric to look up
        time_period:
          type: string
          description: Time period for the metric (e.g., last 30 days, Q3 2024)
        department:
          type: string
          enum: ["sales", "marketing", "support", "product", "finance"]
          description: Department or area for the metric
      required: ["metric_name", "time_period"]
  
  - name: trend_calculation
    description: Calculate trends and growth rates from data
    inputSchema:
      type: object
      properties:
        data_points:
          type: array
          items:
            type: number
          description: Array of numerical data points
        labels:
          type: array
          items:
            type: string
          description: Labels for the data points (dates, periods, etc.)
      required: ["data_points"]
  
  replicas: 1
  resources:
    requests:
      cpu: 150m
      memory: 384Mi
    limits:
      cpu: 300m
      memory: 768Mi
```

### 5. Educational Tutor

Personalized learning assistant for students.

```yaml
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: math-tutor
  namespace: education
spec:
  framework: direct
  provider: openai
  model: gpt-4
  systemPrompt: |
    You are a patient and encouraging math tutor for students aged 10-18.
    
    Teaching Style:
    - Break down complex problems into simple steps
    - Use relatable examples and analogies
    - Encourage students when they struggle
    - Provide hints rather than direct answers
    - Adapt explanations to the student's level
    - Celebrate progress and correct answers
    
    Subjects:
    - Basic arithmetic through calculus
    - Algebra, geometry, trigonometry
    - Statistics and probability
    - Word problems and practical applications
    
  apiSecretRef:
    name: openai-secret
    key: api-key
  
  tools:
  - name: problem_generator
    description: Generate practice problems based on topic and difficulty
    inputSchema:
      type: object
      properties:
        topic:
          type: string
          description: Math topic (e.g., quadratic equations, derivatives)
        difficulty:
          type: string
          enum: ["beginner", "intermediate", "advanced"]
          description: Difficulty level for the student
        problem_count:
          type: integer
          minimum: 1
          maximum: 10
          description: Number of problems to generate
      required: ["topic", "difficulty"]
  
  - name: step_checker
    description: Check if a student's work step is correct
    inputSchema:
      type: object
      properties:
        original_problem:
          type: string
          description: The original math problem
        student_step:
          type: string
          description: The step the student performed
        previous_steps:
          type: array
          items:
            type: string
          description: Previous correct steps in the solution
      required: ["original_problem", "student_step"]
  
  replicas: 2
```

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

**Heavy workloads** (complex reasoning):
```yaml
resources:
  requests:
    cpu: 300m
    memory: 768Mi
  limits:
    cpu: 1
    memory: 2Gi
```

### Scaling Strategies

**High throughput** (many concurrent users):
```yaml
replicas: 5
resources:
  requests:
    cpu: 100m    # Lower per-replica resource usage
    memory: 256Mi
```

**Complex processing** (fewer concurrent users):
```yaml
replicas: 2
resources:
  requests:
    cpu: 500m    # Higher per-replica resource usage
    memory: 1Gi
```

### Tool Design Tips

1. **Keep tools focused**: Each tool should have a single, clear purpose
2. **Validate inputs**: Use proper JSON schema validation
3. **Provide good descriptions**: Help the AI understand when to use each tool
4. **Handle errors gracefully**: Tools should return meaningful error messages
5. **Optimize for speed**: Direct framework excels with fast tool responses

### Monitoring and Debugging

**Essential metrics to track:**
- Response latency (target: <500ms)
- Request throughput
- Error rates
- Resource utilization
- Tool usage patterns

**Common troubleshooting:**
- High latency → Check tool performance, consider caching
- High error rates → Validate tool schemas and API connections
- Resource pressure → Adjust limits or increase replicas
- Poor responses → Review system prompts and tool descriptions

## Migration from LangGraph

If you have a LangGraph workflow that's become too simple, consider migrating to Direct Framework:

**LangGraph** (overkill for simple case):
```yaml
langgraphConfig:
  nodes:
  - name: respond
    type: llm
    prompt: "{user_input}"
  entrypoint: respond
  endpoints: [respond]
```

**Direct** (simpler and faster):
```yaml
framework: direct
# Just use the system prompt - no workflow needed
```

The Direct Framework is perfect when you need speed, simplicity, and reliability for straightforward AI interactions.
