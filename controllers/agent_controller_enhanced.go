package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	aiv1 "github.com/sudeshmu/kubeagentic/api/v1"
)

// AgentReconciler reconciles an Agent object with enhanced features
type AgentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=ai.example.com,resources=agents,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ai.example.com,resources=agents/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ai.example.com,resources=agents/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=autoscaling,resources=horizontalpodautoscalers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=policy,resources=poddisruptionbudgets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is the main reconciliation loop with enhanced features
func (r *AgentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithValues("agent", req.NamespacedName)
	logger.Info("Starting enhanced reconciliation")

	// Fetch the Agent instance
	var agent aiv1.Agent
	if err := r.Get(ctx, req.NamespacedName, &agent); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Agent resource not found, assuming it's been deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get Agent resource")
		return ctrl.Result{}, err
	}

	// Add finalizer for cleanup
	if agent.DeletionTimestamp == nil {
		if !controllerutil.ContainsFinalizer(&agent, "kubeagentic.ai/finalizer") {
			controllerutil.AddFinalizer(&agent, "kubeagentic.ai/finalizer")
			if err := r.Update(ctx, &agent); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// Handle deletion
		if controllerutil.ContainsFinalizer(&agent, "kubeagentic.ai/finalizer") {
			if err := r.cleanupResources(ctx, &agent); err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(&agent, "kubeagentic.ai/finalizer")
			if err := r.Update(ctx, &agent); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Initialize status
	if agent.Status.Phase == "" {
		logger.Info("Initializing Agent status")
		agent.Status.Phase = aiv1.AgentPhasePending
		agent.Status.Message = "Initializing agent deployment"
		now := metav1.NewTime(time.Now())
		agent.Status.LastUpdated = &now
		if err := r.Status().Update(ctx, &agent); err != nil {
			logger.Error(err, "Failed to update Agent status to Pending")
			return ctrl.Result{}, err
		}
	}

	// Validate configuration
	if err := r.validateConfiguration(ctx, &agent); err != nil {
		logger.Error(err, "Configuration validation failed")
		return r.updateStatusFailed(ctx, &agent, fmt.Sprintf("Configuration validation failed: %v", err))
	}

	// Validate secret reference
	if err := r.validateSecretRef(ctx, &agent); err != nil {
		logger.Error(err, "Secret validation failed")
		return r.updateStatusFailed(ctx, &agent, fmt.Sprintf("Secret validation failed: %v", err))
	}

	// Reconcile ConfigMap for tools and configuration
	if err := r.reconcileConfigMap(ctx, &agent); err != nil {
		logger.Error(err, "Failed to reconcile ConfigMap")
		return r.updateStatusFailed(ctx, &agent, fmt.Sprintf("Failed to reconcile ConfigMap: %v", err))
	}

	// Reconcile Deployment
	if err := r.reconcileDeployment(ctx, &agent); err != nil {
		logger.Error(err, "Failed to reconcile Deployment")
		return r.updateStatusFailed(ctx, &agent, fmt.Sprintf("Failed to reconcile Deployment: %v", err))
	}

	// Reconcile Service
	if err := r.reconcileService(ctx, &agent); err != nil {
		logger.Error(err, "Failed to reconcile Service")
		return r.updateStatusFailed(ctx, &agent, fmt.Sprintf("Failed to reconcile Service: %v", err))
	}

	// Reconcile HPA if enabled
	if err := r.reconcileHPA(ctx, &agent); err != nil {
		logger.Error(err, "Failed to reconcile HPA")
		return r.updateStatusFailed(ctx, &agent, fmt.Sprintf("Failed to reconcile HPA: %v", err))
	}

	// Reconcile Ingress if configured
	if err := r.reconcileIngress(ctx, &agent); err != nil {
		logger.Error(err, "Failed to reconcile Ingress")
		return r.updateStatusFailed(ctx, &agent, fmt.Sprintf("Failed to reconcile Ingress: %v", err))
	}

	// Update status
	if err := r.updateAgentStatus(ctx, &agent); err != nil {
		logger.Error(err, "Failed to update Agent status")
		return ctrl.Result{}, err
	}

	logger.Info("Enhanced reconciliation completed successfully")
	return ctrl.Result{RequeueAfter: time.Minute * 5}, nil
}

// validateConfiguration validates the agent configuration
func (r *AgentReconciler) validateConfiguration(ctx context.Context, agent *aiv1.Agent) error {
	// Validate provider
	validProviders := []string{"openai", "gemini", "claude", "vllm"}
	valid := false
	for _, provider := range validProviders {
		if agent.Spec.Provider == provider {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid provider: %s, must be one of %v", agent.Spec.Provider, validProviders)
	}

	// Validate framework
	if agent.Spec.Framework != "" && agent.Spec.Framework != "direct" && agent.Spec.Framework != "langgraph" {
		return fmt.Errorf("invalid framework: %s, must be 'direct' or 'langgraph'", agent.Spec.Framework)
	}

	// Validate LangGraph configuration if framework is langgraph
	if agent.Spec.Framework == "langgraph" && agent.Spec.LanggraphConfig == nil {
		return fmt.Errorf("langgraphConfig is required when framework is 'langgraph'")
	}

	// Validate replicas
	if agent.Spec.Replicas != nil && (*agent.Spec.Replicas < 1 || *agent.Spec.Replicas > 10) {
		return fmt.Errorf("replicas must be between 1 and 10, got %d", *agent.Spec.Replicas)
	}

	return nil
}

// validateSecretRef ensures that the secret referenced by the Agent exists
func (r *AgentReconciler) validateSecretRef(ctx context.Context, agent *aiv1.Agent) error {
	secret := &corev1.Secret{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      agent.Spec.ApiSecretRef.Name,
		Namespace: agent.Namespace,
	}, secret)
	if err != nil {
		return fmt.Errorf("failed to get secret %s: %w", agent.Spec.ApiSecretRef.Name, err)
	}

	if _, exists := secret.Data[agent.Spec.ApiSecretRef.Key]; !exists {
		return fmt.Errorf("key %s not found in secret %s", agent.Spec.ApiSecretRef.Key, agent.Spec.ApiSecretRef.Name)
	}

	return nil
}

// reconcileConfigMap creates a ConfigMap for tools and configuration
func (r *AgentReconciler) reconcileConfigMap(ctx context.Context, agent *aiv1.Agent) error {
	configMap := r.buildConfigMap(agent)
	if err := controllerutil.SetControllerReference(agent, configMap, r.Scheme); err != nil {
		return err
	}

	found := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: configMap.Name, Namespace: configMap.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.FromContext(ctx).Info("Creating new ConfigMap", "ConfigMap.Namespace", configMap.Namespace, "ConfigMap.Name", configMap.Name)
		return r.Create(ctx, configMap)
	} else if err != nil {
		return err
	}

	log.FromContext(ctx).Info("Updating existing ConfigMap", "ConfigMap.Namespace", found.Namespace, "ConfigMap.Name", found.Name)
	found.Data = configMap.Data
	return r.Update(ctx, found)
}

// buildConfigMap creates a ConfigMap with tools and configuration
func (r *AgentReconciler) buildConfigMap(agent *aiv1.Agent) *corev1.ConfigMap {
	labels := map[string]string{
		"app.kubernetes.io/name":     "kubeagentic-agent",
		"app.kubernetes.io/instance": agent.Name,
		"kubeagentic.ai/agent":       agent.Name,
	}

	data := make(map[string]string)
	
	// Add tools configuration
	if len(agent.Spec.Tools) > 0 {
		toolsJSON, _ := json.Marshal(agent.Spec.Tools)
		data["tools.json"] = string(toolsJSON)
	}

	// Add LangGraph configuration
	if agent.Spec.LanggraphConfig != nil {
		configJSON, _ := json.Marshal(agent.Spec.LanggraphConfig)
		data["langgraph-config.json"] = string(configJSON)
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      agent.Name + "-config",
			Namespace: agent.Namespace,
			Labels:    labels,
		},
		Data: data,
	}
}

// cleanupResources handles cleanup when agent is deleted
func (r *AgentReconciler) cleanupResources(ctx context.Context, agent *aiv1.Agent) error {
	logger := log.FromContext(ctx)
	logger.Info("Cleaning up resources for agent", "agent", agent.Name)

	// Update status to indicate cleanup
	agent.Status.Phase = aiv1.AgentPhaseFailed
	agent.Status.Message = "Agent is being deleted"
	now := metav1.NewTime(time.Now())
	agent.Status.LastUpdated = &now
	r.Status().Update(ctx, agent)

	return nil
}

// SetupWithManager sets up the controller with the Manager
func (r *AgentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aiv1.Agent{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
