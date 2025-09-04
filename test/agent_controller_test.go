package test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	aiv1 "github.com/sudeshmu/kubeagentic/api/v1"
)

var _ = Describe("Agent Controller", func() {
	const (
		AgentName      = "test-agent"
		AgentNamespace = "default"
		timeout        = time.Second * 10
		interval       = time.Millisecond * 250
	)

	Context("When creating an Agent", func() {
		It("Should create a Deployment", func() {
			By("Creating a new Agent")
			ctx := context.Background()
			agent := &aiv1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AgentName,
					Namespace: AgentNamespace,
				},
				Spec: aiv1.AgentSpec{
					Provider: "openai",
					Model:    "gpt-4",
					SystemPrompt: "You are a helpful AI assistant.",
					ApiSecretRef: corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "test-secret",
						},
						Key: "api-key",
					},
					Framework: "direct",
					Replicas:  int32Ptr(1),
				},
			}

			Expect(k8sClient.Create(ctx, agent)).Should(Succeed())

			agentLookupKey := types.NamespacedName{Name: AgentName, Namespace: AgentNamespace}
			createdAgent := &aiv1.Agent{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, agentLookupKey, createdAgent)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			By("Checking that the Agent status is updated")
			Eventually(func() string {
				err := k8sClient.Get(ctx, agentLookupKey, createdAgent)
				if err != nil {
					return ""
				}
				return string(createdAgent.Status.Phase)
			}, timeout, interval).Should(Equal(string(aiv1.AgentPhasePending)))

			By("Checking that a Deployment is created")
			deploymentLookupKey := types.NamespacedName{Name: AgentName, Namespace: AgentNamespace}
			createdDeployment := &appsv1.Deployment{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, deploymentLookupKey, createdDeployment)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			By("Checking that a Service is created")
			serviceLookupKey := types.NamespacedName{Name: AgentName + "-service", Namespace: AgentNamespace}
			createdService := &corev1.Service{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, serviceLookupKey, createdService)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			By("Checking that a ConfigMap is created")
			configMapLookupKey := types.NamespacedName{Name: AgentName + "-config", Namespace: AgentNamespace}
			createdConfigMap := &corev1.ConfigMap{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, configMapLookupKey, createdConfigMap)
				return err == nil
			}, timeout, interval).Should(BeTrue())
		})

		It("Should create an HPA for multiple replicas", func() {
			By("Creating an Agent with multiple replicas")
			ctx := context.Background()
			agent := &aiv1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AgentName + "-hpa",
					Namespace: AgentNamespace,
				},
				Spec: aiv1.AgentSpec{
					Provider: "openai",
					Model:    "gpt-4",
					SystemPrompt: "You are a helpful AI assistant.",
					ApiSecretRef: corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "test-secret",
						},
						Key: "api-key",
					},
					Framework: "direct",
					Replicas:  int32Ptr(3),
				},
			}

			Expect(k8sClient.Create(ctx, agent)).Should(Succeed())

			By("Checking that an HPA is created")
			hpaLookupKey := types.NamespacedName{Name: AgentName + "-hpa-hpa", Namespace: AgentNamespace}
			createdHPA := &autoscalingv2.HorizontalPodAutoscaler{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, hpaLookupKey, createdHPA)
				return err == nil
			}, timeout, interval).Should(BeTrue())
		})

		It("Should create an Ingress for LoadBalancer service", func() {
			By("Creating an Agent with LoadBalancer service")
			ctx := context.Background()
			agent := &aiv1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AgentName + "-ingress",
					Namespace: AgentNamespace,
				},
				Spec: aiv1.AgentSpec{
					Provider: "openai",
					Model:    "gpt-4",
					SystemPrompt: "You are a helpful AI assistant.",
					ApiSecretRef: corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "test-secret",
						},
						Key: "api-key",
					},
					Framework:   "direct",
					Replicas:    int32Ptr(1),
					ServiceType: corev1.ServiceTypeLoadBalancer,
				},
			}

			Expect(k8sClient.Create(ctx, agent)).Should(Succeed())

			By("Checking that an Ingress is created")
			ingressLookupKey := types.NamespacedName{Name: AgentName + "-ingress-ingress", Namespace: AgentNamespace}
			createdIngress := &networkingv1.Ingress{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, ingressLookupKey, createdIngress)
				return err == nil
			}, timeout, interval).Should(BeTrue())
		})
	})

	Context("When updating an Agent", func() {
		It("Should update the Deployment when spec changes", func() {
			By("Creating an Agent")
			ctx := context.Background()
			agent := &aiv1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AgentName + "-update",
					Namespace: AgentNamespace,
				},
				Spec: aiv1.AgentSpec{
					Provider: "openai",
					Model:    "gpt-4",
					SystemPrompt: "You are a helpful AI assistant.",
					ApiSecretRef: corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "test-secret",
						},
						Key: "api-key",
					},
					Framework: "direct",
					Replicas:  int32Ptr(1),
				},
			}

			Expect(k8sClient.Create(ctx, agent)).Should(Succeed())

			By("Updating the Agent spec")
			agentLookupKey := types.NamespacedName{Name: AgentName + "-update", Namespace: AgentNamespace}
			createdAgent := &aiv1.Agent{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, agentLookupKey, createdAgent)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			createdAgent.Spec.Replicas = int32Ptr(2)
			Expect(k8sClient.Update(ctx, createdAgent)).Should(Succeed())

			By("Checking that the Deployment is updated")
			deploymentLookupKey := types.NamespacedName{Name: AgentName + "-update", Namespace: AgentNamespace}
			updatedDeployment := &appsv1.Deployment{}

			Eventually(func() int32 {
				err := k8sClient.Get(ctx, deploymentLookupKey, updatedDeployment)
				if err != nil {
					return 0
				}
				return *updatedDeployment.Spec.Replicas
			}, timeout, interval).Should(Equal(int32(2)))
		})
	})

	Context("When deleting an Agent", func() {
		It("Should clean up all resources", func() {
			By("Creating an Agent")
			ctx := context.Background()
			agent := &aiv1.Agent{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AgentName + "-delete",
					Namespace: AgentNamespace,
				},
				Spec: aiv1.AgentSpec{
					Provider: "openai",
					Model:    "gpt-4",
					SystemPrompt: "You are a helpful AI assistant.",
					ApiSecretRef: corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "test-secret",
						},
						Key: "api-key",
					},
					Framework: "direct",
					Replicas:  int32Ptr(1),
				},
			}

			Expect(k8sClient.Create(ctx, agent)).Should(Succeed())

			By("Deleting the Agent")
			Expect(k8sClient.Delete(ctx, agent)).Should(Succeed())

			By("Checking that the Agent is deleted")
			agentLookupKey := types.NamespacedName{Name: AgentName + "-delete", Namespace: AgentNamespace}
			deletedAgent := &aiv1.Agent{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, agentLookupKey, deletedAgent)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())

			By("Checking that the Deployment is deleted")
			deploymentLookupKey := types.NamespacedName{Name: AgentName + "-delete", Namespace: AgentNamespace}
			deletedDeployment := &appsv1.Deployment{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, deploymentLookupKey, deletedDeployment)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())

			By("Checking that the Service is deleted")
			serviceLookupKey := types.NamespacedName{Name: AgentName + "-delete-service", Namespace: AgentNamespace}
			deletedService := &corev1.Service{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, serviceLookupKey, deletedService)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())
		})
	})
})

// Helper function to create int32 pointer
func int32Ptr(i int32) *int32 { return &i }
