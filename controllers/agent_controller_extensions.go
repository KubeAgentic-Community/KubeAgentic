package controllers

import (
	"context"
	"fmt"

	autoscalingv2 "k8s.io/api/autoscaling/v2"
	networkingv1 "k8s.io/api/networking/v1"
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

// reconcileHPA creates or updates HorizontalPodAutoscaler for the agent
func (r *AgentReconciler) reconcileHPA(ctx context.Context, agent *aiv1.Agent) error {
	// Only create HPA if replicas > 1 or if explicitly enabled
	if agent.Spec.Replicas != nil && *agent.Spec.Replicas == 1 {
		// Check if HPA exists and delete it
		hpa := &autoscalingv2.HorizontalPodAutoscaler{}
		err := r.Get(ctx, types.NamespacedName{Name: agent.Name + "-hpa", Namespace: agent.Namespace}, hpa)
		if err == nil {
			log.FromContext(ctx).Info("Deleting HPA for single replica agent", "HPA.Name", hpa.Name)
			return r.Delete(ctx, hpa)
		}
		return nil
	}

	hpa := r.buildHPA(agent)
	if err := controllerutil.SetControllerReference(agent, hpa, r.Scheme); err != nil {
		return err
	}

	found := &autoscalingv2.HorizontalPodAutoscaler{}
	err := r.Get(ctx, types.NamespacedName{Name: hpa.Name, Namespace: hpa.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.FromContext(ctx).Info("Creating new HPA", "HPA.Namespace", hpa.Namespace, "HPA.Name", hpa.Name)
		return r.Create(ctx, hpa)
	} else if err != nil {
		return err
	}

	log.FromContext(ctx).Info("Updating existing HPA", "HPA.Namespace", found.Namespace, "HPA.Name", found.Name)
	found.Spec = hpa.Spec
	return r.Update(ctx, found)
}

// buildHPA creates a HorizontalPodAutoscaler for the agent
func (r *AgentReconciler) buildHPA(agent *aiv1.Agent) *autoscalingv2.HorizontalPodAutoscaler {
	labels := map[string]string{
		"app.kubernetes.io/name":     "kubeagentic-agent",
		"app.kubernetes.io/instance": agent.Name,
		"kubeagentic.ai/agent":       agent.Name,
	}

	minReplicas := int32(1)
	maxReplicas := int32(10)
	if agent.Spec.Replicas != nil {
		minReplicas = *agent.Spec.Replicas
		maxReplicas = *agent.Spec.Replicas * 3 // Scale up to 3x
	}

	return &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      agent.Name + "-hpa",
			Namespace: agent.Namespace,
			Labels:    labels,
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       agent.Name,
			},
			MinReplicas: &minReplicas,
			MaxReplicas: maxReplicas,
			Metrics: []autoscalingv2.MetricSpec{
				{
					Type: autoscalingv2.ResourceMetricSourceType,
					Resource: &autoscalingv2.ResourceMetricSource{
						Name: "cpu",
						Target: autoscalingv2.MetricTarget{
							Type:               autoscalingv2.UtilizationMetricType,
							AverageUtilization: int32Ptr(70),
						},
					},
				},
				{
					Type: autoscalingv2.ResourceMetricSourceType,
					Resource: &autoscalingv2.ResourceMetricSource{
						Name: "memory",
						Target: autoscalingv2.MetricTarget{
							Type:               autoscalingv2.UtilizationMetricType,
							AverageUtilization: int32Ptr(80),
						},
					},
				},
			},
		},
	}
}

// reconcileIngress creates or updates Ingress for the agent
func (r *AgentReconciler) reconcileIngress(ctx context.Context, agent *aiv1.Agent) error {
	// Only create Ingress if service type is LoadBalancer or if explicitly configured
	if agent.Spec.ServiceType != "LoadBalancer" {
		// Check if Ingress exists and delete it
		ingress := &networkingv1.Ingress{}
		err := r.Get(ctx, types.NamespacedName{Name: agent.Name + "-ingress", Namespace: agent.Namespace}, ingress)
		if err == nil {
			log.FromContext(ctx).Info("Deleting Ingress for non-LoadBalancer service", "Ingress.Name", ingress.Name)
			return r.Delete(ctx, ingress)
		}
		return nil
	}

	ingress := r.buildIngress(agent)
	if err := controllerutil.SetControllerReference(agent, ingress, r.Scheme); err != nil {
		return err
	}

	found := &networkingv1.Ingress{}
	err := r.Get(ctx, types.NamespacedName{Name: ingress.Name, Namespace: ingress.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.FromContext(ctx).Info("Creating new Ingress", "Ingress.Namespace", ingress.Namespace, "Ingress.Name", ingress.Name)
		return r.Create(ctx, ingress)
	} else if err != nil {
		return err
	}

	log.FromContext(ctx).Info("Updating existing Ingress", "Ingress.Namespace", found.Namespace, "Ingress.Name", found.Name)
	found.Spec = ingress.Spec
	return r.Update(ctx, found)
}

// buildIngress creates an Ingress for the agent
func (r *AgentReconciler) buildIngress(agent *aiv1.Agent) *networkingv1.Ingress {
	labels := map[string]string{
		"app.kubernetes.io/name":     "kubeagentic-agent",
		"app.kubernetes.io/instance": agent.Name,
		"kubeagentic.ai/agent":       agent.Name,
	}

	hostname := fmt.Sprintf("%s.%s.local", agent.Name, agent.Namespace)
	pathType := networkingv1.PathTypePrefix

	return &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      agent.Name + "-ingress",
			Namespace: agent.Namespace,
			Labels:    labels,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target": "/",
				"nginx.ingress.kubernetes.io/ssl-redirect":   "false",
			},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: hostname,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: agent.Name + "-service",
											Port: networkingv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Helper function to create int32 pointer
func int32Ptr(i int32) *int32 { return &i }
