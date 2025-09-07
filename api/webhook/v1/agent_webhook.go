package v1

import (
	"context"
	"fmt"
	"net/http"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	aiv1 "github.com/KubeAgentic-Community/kubeagentic/api/v1"
)

// +kubebuilder:webhook:path=/mutate-ai-example-com-v1-agent,mutating=true,failurePolicy=fail,sideEffects=None,groups=ai.example.com,resources=agents,verbs=create;update,versions=v1,name=magent.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Agent{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Agent) Default() {
	log := logf.Log.WithName("agent-resource")

	log.Info("default", "name", r.Name)

	// Set default framework if not specified
	if r.Spec.Framework == "" {
		r.Spec.Framework = "direct"
	}

	// Set default replicas if not specified
	if r.Spec.Replicas == nil {
		defaultReplicas := int32(1)
		r.Spec.Replicas = &defaultReplicas
	}

	// Set default service type if not specified
	if r.Spec.ServiceType == "" {
		r.Spec.ServiceType = "ClusterIP"
	}

	// Set default resources if not specified
	if r.Spec.Resources == nil {
		r.Spec.Resources = &corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse("256Mi"),
				corev1.ResourceCPU:    resource.MustParse("100m"),
			},
			Limits: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse("512Mi"),
				corev1.ResourceCPU:    resource.MustParse("200m"),
			},
		}
	}
}

// +kubebuilder:webhook:path=/validate-ai-example-com-v1-agent,mutating=false,failurePolicy=fail,sideEffects=None,groups=ai.example.com,resources=agents,verbs=create;update,versions=v1,name=vagent.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Agent{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Agent) ValidateCreate() (admission.Warnings, error) {
	log := logf.Log.WithName("agent-resource")
	log.Info("validate create", "name", r.Name)

	return nil, r.validateAgent()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Agent) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	log := logf.Log.WithName("agent-resource")
	log.Info("validate update", "name", r.Name)

	return nil, r.validateAgent()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Agent) ValidateDelete() (admission.Warnings, error) {
	log := logf.Log.WithName("agent-resource")
	log.Info("validate delete", "name", r.Name)

	// Add any deletion validation logic here
	return nil, nil
}

// validateAgent validates the Agent resource
func (r *Agent) validateAgent() error {
	var allErrs field.ErrorList

	// Validate provider
	validProviders := []string{"openai", "gemini", "claude", "vllm"}
	valid := false
	for _, provider := range validProviders {
		if r.Spec.Provider == provider {
			valid = true
			break
		}
	}
	if !valid {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec").Child("provider"),
			r.Spec.Provider,
			fmt.Sprintf("must be one of %v", validProviders),
		))
	}

	// Validate model
	if r.Spec.Model == "" {
		allErrs = append(allErrs, field.Required(
			field.NewPath("spec").Child("model"),
			"model is required",
		))
	}

	// Validate system prompt
	if r.Spec.SystemPrompt == "" {
		allErrs = append(allErrs, field.Required(
			field.NewPath("spec").Child("systemPrompt"),
			"systemPrompt is required",
		))
	}

	// Validate API secret reference
	if r.Spec.ApiSecretRef.Name == "" {
		allErrs = append(allErrs, field.Required(
			field.NewPath("spec").Child("apiSecretRef").Child("name"),
			"apiSecretRef.name is required",
		))
	}
	if r.Spec.ApiSecretRef.Key == "" {
		allErrs = append(allErrs, field.Required(
			field.NewPath("spec").Child("apiSecretRef").Child("key"),
			"apiSecretRef.key is required",
		))
	}

	// Validate framework
	if r.Spec.Framework != "" && r.Spec.Framework != "direct" && r.Spec.Framework != "langgraph" {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec").Child("framework"),
			r.Spec.Framework,
			"must be 'direct' or 'langgraph'",
		))
	}

	// Validate LangGraph configuration
	if r.Spec.Framework == "langgraph" && r.Spec.LanggraphConfig == nil {
		allErrs = append(allErrs, field.Required(
			field.NewPath("spec").Child("langgraphConfig"),
			"langgraphConfig is required when framework is 'langgraph'",
		))
	}

	// Validate replicas
	if r.Spec.Replicas != nil && (*r.Spec.Replicas < 1 || *r.Spec.Replicas > 10) {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec").Child("replicas"),
			*r.Spec.Replicas,
			"must be between 1 and 10",
		))
	}

	// Validate service type
	validServiceTypes := []string{"ClusterIP", "NodePort", "LoadBalancer"}
	validServiceType := false
	for _, serviceType := range validServiceTypes {
		if r.Spec.ServiceType == serviceType {
			validServiceType = true
			break
		}
	}
	if r.Spec.ServiceType != "" && !validServiceType {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec").Child("serviceType"),
			r.Spec.ServiceType,
			fmt.Sprintf("must be one of %v", validServiceTypes),
		))
	}

	if len(allErrs) == 0 {
		return nil
	}

	return fmt.Errorf("validation failed: %v", allErrs)
}

// SetupWebhookWithManager sets up the webhook with the Manager
func (r *Agent) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}
