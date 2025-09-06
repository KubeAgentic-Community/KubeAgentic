package controllers

import (
	"context"
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

// SimpleAgentReconciler is a simplified Kubernetes operator for AI agents
type SimpleAgentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=ai.example.com,resources=agents,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ai.example.com,resources=agents/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch

// Reconcile handles the main logic for managing AI agents
func (r *SimpleAgentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithValues("agent", req.NamespacedName)
	logger.Info("🔄 Starting simple reconciliation")

	// 1. Get the Agent resource
	var agent aiv1.Agent
	if err := r.Get(ctx, req.NamespacedName, &agent); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("❌ Agent not found - assuming deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "❌ Failed to get Agent")
		return ctrl.Result{}, err
	}

	// 2. Set finalizer for cleanup
	if !controllerutil.ContainsFinalizer(&agent, "kubeagentic.ai/finalizer") {
		controllerutil.AddFinalizer(&agent, "kubeagentic.ai/finalizer")
		if err := r.Update(ctx, &agent); err != nil {
			return ctrl.Result{}, err
		}
	}

	// 3. Handle deletion
	if agent.DeletionTimestamp != nil {
		logger.Info("🗑️ Agent being deleted")
		controllerutil.RemoveFinalizer(&agent, "kubeagentic.ai/finalizer")
		return ctrl.Result{}, r.Update(ctx, &agent)
	}

	// 4. Update status to pending
	if agent.Status.Phase == "" {
		agent.Status.Phase = aiv1.AgentPhasePending
		agent.Status.Message = "Creating agent resources"
		now := metav1.NewTime(time.Now())
		agent.Status.LastUpdated = &now
		if err := r.Status().Update(ctx, &agent); err != nil {
			logger.Error(err, "❌ Failed to update status")
			return ctrl.Result{}, err
		}
	}

	// 5. Create Deployment
	if err := r.createDeployment(ctx, &agent); err != nil {
		logger.Error(err, "❌ Failed to create deployment")
		return r.updateStatusFailed(ctx, &agent, err.Error()), nil
	}

	// 6. Create Service
	if err := r.createService(ctx, &agent); err != nil {
		logger.Error(err, "❌ Failed to create service")
		return r.updateStatusFailed(ctx, &agent, err.Error()), nil
	}

	// 7. Update status to running
	agent.Status.Phase = aiv1.AgentPhaseRunning
	agent.Status.Message = "Agent is running"
	now := metav1.NewTime(time.Now())
	agent.Status.LastUpdated = &now
	if err := r.Status().Update(ctx, &agent); err != nil {
		logger.Error(err, "❌ Failed to update status to running")
		return ctrl.Result{}, err
	}

	logger.Info("✅ Reconciliation completed successfully")
	return ctrl.Result{RequeueAfter: time.Minute * 5}, nil
}

// createDeployment creates a simple deployment for the agent
func (r *SimpleAgentReconciler) createDeployment(ctx context.Context, agent *aiv1.Agent) error {
	logger := log.FromContext(ctx)

	// Set defaults
	replicas := int32(1)
	if agent.Spec.Replicas != nil {
		replicas = *agent.Spec.Replicas
	}

	// Create deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      agent.Name,
			Namespace: agent.Namespace,
			Labels: map[string]string{
				"app": agent.Name,
				"kubeagentic.ai/agent": agent.Name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": agent.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": agent.Name,
						"kubeagentic.ai/agent": agent.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "agent",
							Image: "kubeagentic/agent:latest",
							Ports: []corev1.ContainerPort{
								{ContainerPort: 8080, Protocol: corev1.ProtocolTCP},
							},
							Env: []corev1.EnvVar{
								{Name: "AGENT_PROVIDER", Value: agent.Spec.Provider},
								{Name: "AGENT_MODEL", Value: agent.Spec.Model},
								{Name: "AGENT_SYSTEM_PROMPT", Value: agent.Spec.SystemPrompt},
								{Name: "AGENT_FRAMEWORK", Value: stringOrDefault(agent.Spec.Framework, "direct")},
								{
									Name: "AGENT_API_KEY",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &agent.Spec.ApiSecretRef,
									},
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse("256Mi"),
									corev1.ResourceCPU:    resource.MustParse("100m"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse("512Mi"),
									corev1.ResourceCPU:    resource.MustParse("200m"),
								},
							},
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

	// Set controller reference
	if err := controllerutil.SetControllerReference(agent, deployment, r.Scheme); err != nil {
		return err
	}

	// Create or update
	found := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("📦 Creating deployment", "name", deployment.Name)
		return r.Create(ctx, deployment)
	} else if err != nil {
		return err
	}

	logger.Info("🔄 Updating deployment", "name", deployment.Name)
	found.Spec = deployment.Spec
	return r.Update(ctx, found)
}

// createService creates a simple service for the agent
func (r *SimpleAgentReconciler) createService(ctx context.Context, agent *aiv1.Agent) error {
	logger := log.FromContext(ctx)

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      agent.Name + "-service",
			Namespace: agent.Namespace,
			Labels: map[string]string{
				"app": agent.Name,
				"kubeagentic.ai/agent": agent.Name,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": agent.Name,
			},
			Ports: []corev1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt(8080),
					Protocol:   corev1.ProtocolTCP,
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	// Set controller reference
	if err := controllerutil.SetControllerReference(agent, service, r.Scheme); err != nil {
		return err
	}

	// Create or update
	found := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("🌐 Creating service", "name", service.Name)
		return r.Create(ctx, service)
	} else if err != nil {
		return err
	}

	logger.Info("🔄 Updating service", "name", service.Name)
	found.Spec = service.Spec
	return r.Update(ctx, found)
}

// updateStatusFailed updates the agent status to failed
func (r *SimpleAgentReconciler) updateStatusFailed(ctx context.Context, agent *aiv1.Agent, message string) ctrl.Result {
	agent.Status.Phase = aiv1.AgentPhaseFailed
	agent.Status.Message = message
	now := metav1.NewTime(time.Now())
	agent.Status.LastUpdated = &now
	r.Status().Update(ctx, agent)
	return ctrl.Result{RequeueAfter: time.Minute * 2}
}

// SetupWithManager sets up the controller with the Manager
func (r *SimpleAgentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aiv1.Agent{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}

// Helper function for string defaults
func stringOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
