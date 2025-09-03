package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// AgentSpec defines the desired state of an Agent.
// It contains all the configurable parameters for an agent, such as the provider, model, and resources.
type AgentSpec struct {
	// Provider specifies the LLM provider to use for the agent.
	// This is a mandatory field and must be one of the supported providers.
	// +kubebuilder:validation:Enum=openai;gemini;claude;vllm
	Provider string `json:"provider"`

	// Model specifies the specific model to use from the selected provider.
	// For example, "gpt-4" for OpenAI or "claude-2" for Anthropic.
	Model string `json:"model"`

	// SystemPrompt defines the agent's persona, behavior, and instructions.
	// It's a crucial part of the agent's configuration that guides its responses.
	SystemPrompt string `json:"systemPrompt"`

	// ApiSecretRef references a Kubernetes Secret that holds the API credentials for the provider.
	// The secret must contain a key with the API key.
	ApiSecretRef corev1.SecretKeySelector `json:"apiSecretRef"`

	// Endpoint is an optional field to specify a custom endpoint URL.
	// This is particularly useful for self-hosted models like vLLM.
	// +optional
	Endpoint string `json:"endpoint,omitempty"`

	// Framework specifies which framework to use for agent execution.
	// "direct" uses simple API calls, "langgraph" enables complex workflows.
	// +kubebuilder:validation:Enum=direct;langgraph
	// +kubebuilder:default=direct
	// +optional
	Framework string `json:"framework,omitempty"`

	// LanggraphConfig contains configuration for LangGraph workflows.
	// Only used when Framework is set to "langgraph".
	// +optional
	LanggraphConfig *LanggraphConfig `json:"langgraphConfig,omitempty"`

	// Tools is a list of tools that the agent can use to perform actions.
	// Each tool has a name, description, and an optional input schema.
	// +optional
	Tools []Tool `json:"tools,omitempty"`

	// Replicas is the number of agent pod replicas to run.
	// Defaults to 1 if not specified.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10
	// +kubebuilder:default=1
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Resources defines the CPU and memory requests and limits for the agent pods.
	// If not specified, default resources will be allocated.
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// ServiceType specifies the type of Kubernetes service to create for the agent endpoint.
	// It can be ClusterIP, NodePort, or LoadBalancer. Defaults to ClusterIP.
	// +kubebuilder:validation:Enum=ClusterIP;NodePort;LoadBalancer
	// +kubebuilder:default=ClusterIP
	// +optional
	ServiceType corev1.ServiceType `json:"serviceType,omitempty"`
}

// Tool defines a tool that is available to the agent.
// Tools allow agents to interact with external systems and perform actions.
type Tool struct {
	// Name is the unique identifier for the tool.
	Name string `json:"name"`

	// Description is a human-readable explanation of what the tool does.
	// This is used by the agent to decide when to use the tool.
	Description string `json:"description"`

	// InputSchema is a JSON schema that describes the input parameters for the tool.
	// This helps the agent to correctly format the input for the tool.
	// +optional
	InputSchema *runtime.RawExtension `json:"inputSchema,omitempty"`
}

// LanggraphConfig defines the configuration for LangGraph workflows
type LanggraphConfig struct {
	// GraphType specifies the type of LangGraph workflow
	// +kubebuilder:validation:Enum=sequential;parallel;conditional;hierarchical
	GraphType string `json:"graphType"`

	// Nodes defines the workflow nodes
	Nodes []WorkflowNode `json:"nodes"`

	// Edges defines the workflow edges
	Edges []WorkflowEdge `json:"edges"`

	// State defines the state schema for the workflow
	State *runtime.RawExtension `json:"state,omitempty"`

	// Entrypoint specifies the entry node for the workflow
	Entrypoint string `json:"entrypoint"`

	// Endpoints specifies possible end nodes for the workflow
	Endpoints []string `json:"endpoints,omitempty"`
}

// WorkflowNode defines a node in the LangGraph workflow
type WorkflowNode struct {
	// Name is the unique identifier for the node
	Name string `json:"name"`

	// Type specifies the type of node
	// +kubebuilder:validation:Enum=llm;tool;action
	Type string `json:"type"`

	// Prompt is the template for LLM nodes
	Prompt string `json:"prompt,omitempty"`

	// Tool is the tool name for tool nodes
	Tool string `json:"tool,omitempty"`

	// Action is the action to execute for action nodes
	Action string `json:"action,omitempty"`

	// Condition is the conditional logic for conditional nodes
	Condition string `json:"condition,omitempty"`

	// Inputs are the input fields from state
	Inputs []string `json:"inputs,omitempty"`

	// Outputs are the output fields to state
	Outputs []string `json:"outputs,omitempty"`
}

// WorkflowEdge defines an edge in the LangGraph workflow
type WorkflowEdge struct {
	// From is the source node name
	From string `json:"from"`

	// To is the target node name
	To string `json:"to"`

	// Condition is the conditional logic for the edge
	Condition string `json:"condition,omitempty"`
}

// AgentConditionType represents the type of an Agent's condition.
type AgentConditionType string

const (
	// AgentConditionReady indicates that the agent is ready to serve requests.
	AgentConditionReady AgentConditionType = "Ready"
	// AgentConditionProgressing indicates that the agent's deployment is in progress.
	AgentConditionProgressing AgentConditionType = "Progressing"
	// AgentConditionDegraded indicates that the agent is in a degraded state.
	AgentConditionDegraded AgentConditionType = "Degraded"
)

// AgentCondition represents the condition of an Agent.
// It provides more detailed information about the agent's state.
type AgentCondition struct {
	// Type of the condition.
	Type AgentConditionType `json:"type"`

	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`

	// Reason is a brief, machine-readable reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`

	// Message is a human-readable message indicating details about the last transition.
	// +optional
	Message string `json:"message,omitempty"`

	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime *metav1.Time `json:"lastTransitionTime,omitempty"`
}

// AgentPhase represents the lifecycle phase of an Agent.
type AgentPhase string

const (
	// AgentPhasePending means the agent is being created and is not yet ready.
	AgentPhasePending AgentPhase = "Pending"
	// AgentPhaseRunning means the agent is running and ready to serve requests.
	AgentPhaseRunning AgentPhase = "Running"
	// AgentPhaseFailed means the agent has encountered an error and is not running.
	AgentPhaseFailed AgentPhase = "Failed"
	// AgentPhaseSucceeded is not currently used but is reserved for future use.
	AgentPhaseSucceeded AgentPhase = "Succeeded"
)

// ReplicaStatus represents the status of the agent's replicas.
type ReplicaStatus struct {
	// Ready is the number of replicas that are ready to serve requests.
	Ready int32 `json:"ready"`

	// Desired is the desired number of replicas.
	Desired int32 `json:"desired"`

	// Available is the number of replicas that are available.
	Available int32 `json:"available"`
}

// AgentStatus defines the observed state of an Agent.
// It provides a summary of the agent's current state.
type AgentStatus struct {
	// Phase represents the current lifecycle phase of the agent.
	// +optional
	Phase AgentPhase `json:"phase,omitempty"`

	// Message is a human-readable message about the agent's current state.
	// +optional
	Message string `json:"message,omitempty"`

	// ReplicaStatus shows the current status of the agent's replicas.
	// +optional
	ReplicaStatus ReplicaStatus `json:"replicaStatus,omitempty"`

	// LastUpdated is the timestamp of the last status update.
	// +optional
	LastUpdated *metav1.Time `json:"lastUpdated,omitempty"`

	// Conditions is a list of the latest available observations of the agent's state.
	// +optional
	Conditions []AgentCondition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=ag
// +kubebuilder:printcolumn:name="Provider",type="string",JSONPath=".spec.provider"
// +kubebuilder:printcolumn:name="Model",type="string",JSONPath=".spec.model"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.replicaStatus.ready"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Agent is the Schema for the agents API. It represents a single AI agent.
type Agent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AgentSpec   `json:"spec,omitempty"`
	Status AgentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AgentList contains a list of Agent resources.
type AgentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Agent `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Agent{}, &AgentList{})
}
