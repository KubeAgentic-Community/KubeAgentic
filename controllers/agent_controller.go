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

	aiv1 "github.com/KubeAgentic-Community/kubeagentic/api/v1"
)

// AgentReconciler reconciles an Agent object.
// It's the core component of the operator, responsible for managing the lifecycle of Agent resources.
type AgentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// RBAC annotations setup the necessary permissions for the controller to manage resources.
// +kubebuilder:rbac:groups=ai.example.com,resources=agents,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ai.example.com,resources=agents/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ai.example.com,resources=agents/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch

// Reconcile is the main reconciliation loop for the Agent controller.
// It's triggered by changes to Agent resources or the resources it owns.
func (r *AgentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithValues("agent", req.NamespacedName)
	logger.Info("Starting reconciliation")

	// Fetch the Agent instance
	var agent aiv1.Agent
	if err := r.Get(ctx, req.NamespacedName, &agent); err != nil {
		if errors.IsNotFound(err) {
			// The Agent resource was not found, likely deleted.
			// There's nothing to do, so we can return without error.
			logger.Info("Agent resource not found, assuming it's been deleted")
			return ctrl.Result{}, nil
		}
		// An unexpected error occurred while fetching the Agent resource.
		logger.Error(err, "Failed to get Agent resource")
		return ctrl.Result{}, err
	}

	// Set the initial status of the Agent resource.
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

	// Validate the secret reference to ensure the API key is available.
	if err := r.validateSecretRef(ctx, &agent); err != nil {
		logger.Error(err, "Secret validation failed")
		return r.updateStatusFailed(ctx, &agent, fmt.Sprintf("Secret validation failed: %v", err))
	}

	// Reconcile the Deployment for the Agent.
	if err := r.reconcileDeployment(ctx, &agent); err != nil {
		logger.Error(err, "Failed to reconcile Deployment")
		return r.updateStatusFailed(ctx, &agent, fmt.Sprintf("Failed to reconcile Deployment: %v", err))
	}

	// Reconcile the Service for the Agent.
	if err := r.reconcileService(ctx, &agent); err != nil {
		logger.Error(err, "Failed to reconcile Service")
		return r.updateStatusFailed(ctx, &agent, fmt.Sprintf("Failed to reconcile Service: %v", err))
	}

	// Update the Agent's status based on the state of its owned resources.
	if err := r.updateAgentStatus(ctx, &agent); err != nil {
		logger.Error(err, "Failed to update Agent status")
		return ctrl.Result{}, err
	}

	logger.Info("Reconciliation completed successfully")
	return ctrl.Result{RequeueAfter: time.Minute * 5}, nil
}

// validateSecretRef ensures that the secret referenced by the Agent exists and contains the required key.
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

// reconcileDeployment manages the Deployment resource for the Agent.
func (r *AgentReconciler) reconcileDeployment(ctx context.Context, agent *aiv1.Agent) error {
	deployment := r.buildDeployment(agent)
	if err := controllerutil.SetControllerReference(agent, deployment, r.Scheme); err != nil {
		return err
	}

	found := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.FromContext(ctx).Info("Creating new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		return r.Create(ctx, deployment)
	} else if err != nil {
		return err
	}

	log.FromContext(ctx).Info("Updating existing Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
	found.Spec = deployment.Spec
	return r.Update(ctx, found)
}

// reconcileService manages the Service resource for the Agent.
func (r *AgentReconciler) reconcileService(ctx context.Context, agent *aiv1.Agent) error {
	service := r.buildService(agent)
	if err := controllerutil.SetControllerReference(agent, service, r.Scheme); err != nil {
		return err
	}

	foundService := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		log.FromContext(ctx).Info("Creating new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
		return r.Create(ctx, service)
	} else if err != nil {
		return err
	}

	log.FromContext(ctx).Info("Updating existing Service", "Service.Namespace", foundService.Namespace, "Service.Name", foundService.Name)
	foundService.Spec.Ports = service.Spec.Ports
	foundService.Spec.Selector = service.Spec.Selector
	foundService.Spec.Type = service.Spec.Type
	return r.Update(ctx, foundService)
}

// buildDeployment creates a new Deployment resource based on the Agent's specification.
func (r *AgentReconciler) buildDeployment(agent *aiv1.Agent) *appsv1.Deployment {
	replicas := int32(1)
	if agent.Spec.Replicas != nil {
		replicas = *agent.Spec.Replicas
	}

	// Default resource requirements, can be overridden by the user.
	resources := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse("256Mi"),
			corev1.ResourceCPU:    resource.MustParse("100m"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse("512Mi"),
			corev1.ResourceCPU:    resource.MustParse("200m"),
		},
	}

	if agent.Spec.Resources != nil {
		resources = *agent.Spec.Resources
	}

	// Construct environment variables for the agent container.
	env := []corev1.EnvVar{
		{Name: "AGENT_PROVIDER", Value: agent.Spec.Provider},
		{Name: "AGENT_MODEL", Value: agent.Spec.Model},
		{Name: "AGENT_SYSTEM_PROMPT", Value: agent.Spec.SystemPrompt},
		{
			Name: "AGENT_API_KEY",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &agent.Spec.ApiSecretRef,
			},
		},
	}

	if agent.Spec.Endpoint != "" {
		env = append(env, corev1.EnvVar{
			Name:  "AGENT_ENDPOINT",
			Value: agent.Spec.Endpoint,
		})
	}

	// Add framework configuration
	framework := "direct" // default
	if agent.Spec.Framework != "" {
		framework = agent.Spec.Framework
	}
	env = append(env, corev1.EnvVar{
		Name:  "AGENT_FRAMEWORK",
		Value: framework,
	})

	// Add LangGraph configuration if present
	if agent.Spec.LanggraphConfig != nil && framework == "langgraph" {
		configBytes, err := json.Marshal(agent.Spec.LanggraphConfig)
		if err == nil {
			env = append(env, corev1.EnvVar{
				Name:  "AGENT_LANGGRAPH_CONFIG",
				Value: string(configBytes),
			})
		}
	}

	// A simple way to pass tools to the agent. A more robust implementation might use a ConfigMap.
	if len(agent.Spec.Tools) > 0 {
		env = append(env, corev1.EnvVar{
			Name:  "AGENT_TOOLS_COUNT",
			Value: fmt.Sprintf("%d", len(agent.Spec.Tools)),
		})
	}

	labels := map[string]string{
		"app.kubernetes.io/name":     "kubeagentic-agent",
		"app.kubernetes.io/instance": agent.Name,
		"kubeagentic.ai/agent":       agent.Name,
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      agent.Name,
			Namespace: agent.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "agent",
							Image: "sudeshmu/kubeagentic:agent-fixed", // This should be configurable in a real-world scenario.
							Ports: []corev1.ContainerPort{
								{ContainerPort: 8080, Protocol: corev1.ProtocolTCP},
							},
							Env:       env,
							Resources: resources,
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/health",
										Port: intstr.FromInt(8080),
									},
								},
								InitialDelaySeconds: 30,
								PeriodSeconds:       10,
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/ready",
										Port: intstr.FromInt(8080),
									},
								},
								InitialDelaySeconds: 5,
								PeriodSeconds:       5,
							},
						},
					},
				},
			},
		},
	}
}

// buildService creates a new Service resource to expose the Agent's Deployment.
func (r *AgentReconciler) buildService(agent *aiv1.Agent) *corev1.Service {
	serviceType := corev1.ServiceTypeClusterIP
	if agent.Spec.ServiceType != "" {
		serviceType = agent.Spec.ServiceType
	}

	labels := map[string]string{
		"app.kubernetes.io/name":     "kubeagentic-agent",
		"app.kubernetes.io/instance": agent.Name,
		"kubeagentic.ai/agent":       agent.Name,
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      agent.Name + "-service",
			Namespace: agent.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type:     serviceType,
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt(8080),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}
}

// updateAgentStatus updates the status of the Agent resource based on the state of the Deployment.
func (r *AgentReconciler) updateAgentStatus(ctx context.Context, agent *aiv1.Agent) error {
	deployment := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: agent.Name, Namespace: agent.Namespace}, deployment)
	if err != nil {
		return fmt.Errorf("failed to get deployment for status update: %w", err)
	}

	// Update replica status from the deployment.
	agent.Status.ReplicaStatus.Desired = *deployment.Spec.Replicas
	agent.Status.ReplicaStatus.Ready = deployment.Status.ReadyReplicas
	agent.Status.ReplicaStatus.Available = deployment.Status.AvailableReplicas

	// Determine the phase of the Agent based on the deployment's status.
	if deployment.Status.ReadyReplicas == *deployment.Spec.Replicas && deployment.Status.ReadyReplicas > 0 {
		agent.Status.Phase = aiv1.AgentPhaseRunning
		agent.Status.Message = "Agent is running and ready"
	} else if deployment.Status.Replicas == 0 {
		agent.Status.Phase = aiv1.AgentPhasePending
		agent.Status.Message = "Agent deployment is scaling up"
	} else {
		agent.Status.Phase = aiv1.AgentPhasePending
		agent.Status.Message = fmt.Sprintf("Agent deployment in progress (%d/%d ready)", deployment.Status.ReadyReplicas, *deployment.Spec.Replicas)
	}

	now := metav1.NewTime(time.Now())
	agent.Status.LastUpdated = &now

	// Set the Ready condition based on the Agent's phase.
	readyCondition := aiv1.AgentCondition{
		Type:               aiv1.AgentConditionReady,
		LastTransitionTime: &now,
	}

	if agent.Status.Phase == aiv1.AgentPhaseRunning {
		readyCondition.Status = corev1.ConditionTrue
		readyCondition.Reason = "DeploymentReady"
		readyCondition.Message = "All replicas are ready"
	} else {
		readyCondition.Status = corev1.ConditionFalse
		readyCondition.Reason = "DeploymentNotReady"
		readyCondition.Message = "Deployment is not yet ready"
	}

	agent.Status.Conditions = r.updateCondition(agent.Status.Conditions, readyCondition)

	return r.Status().Update(ctx, agent)
}

// updateStatusFailed is a helper function to update the Agent's status to Failed.
func (r *AgentReconciler) updateStatusFailed(ctx context.Context, agent *aiv1.Agent, message string) (ctrl.Result, error) {
	agent.Status.Phase = aiv1.AgentPhaseFailed
	agent.Status.Message = message
	now := metav1.NewTime(time.Now())
	agent.Status.LastUpdated = &now

	degradedCondition := aiv1.AgentCondition{
		Type:               aiv1.AgentConditionDegraded,
		Status:             corev1.ConditionTrue,
		Reason:             "ReconciliationFailed",
		Message:            message,
		LastTransitionTime: &now,
	}
	agent.Status.Conditions = r.updateCondition(agent.Status.Conditions, degradedCondition)

	if err := r.Status().Update(ctx, agent); err != nil {
		// Log the error but return the original error to avoid masking the root cause.
		log.FromContext(ctx).Error(err, "Failed to update agent status to Failed")
	}

	// Requeue after a short period to allow for manual intervention or for the issue to be resolved.
	return ctrl.Result{RequeueAfter: time.Minute * 2}, nil
}

// updateCondition is a helper function to update a condition in the Agent's status.
func (r *AgentReconciler) updateCondition(conditions []aiv1.AgentCondition, newCondition aiv1.AgentCondition) []aiv1.AgentCondition {
	for i, condition := range conditions {
		if condition.Type == newCondition.Type {
			// If the status of the condition has not changed, we don't need to update it.
			if condition.Status == newCondition.Status {
				newCondition.LastTransitionTime = condition.LastTransitionTime
			}
			conditions[i] = newCondition
			return conditions
		}
	}
	return append(conditions, newCondition)
}

// SetupWithManager sets up the controller with the Manager.
// This is how the controller is registered with the controller-runtime.
func (r *AgentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aiv1.Agent{}).
		// Owns specifies the resources that are owned by the Agent resource.
		// This allows the controller to watch for changes to these resources.
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
