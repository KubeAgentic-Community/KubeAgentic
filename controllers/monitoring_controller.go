package controllers

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	aiv1 "github.com/sudeshmu/kubeagentic/api/v1"
)

// MonitoringReconciler handles monitoring and observability for agents
type MonitoringReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=ai.example.com,resources=agents,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile handles monitoring setup for agents
func (r *MonitoringReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithValues("monitoring", req.NamespacedName)
	logger.Info("Starting monitoring reconciliation")

	// Fetch all agents
	var agents aiv1.AgentList
	if err := r.List(ctx, &agents); err != nil {
		logger.Error(err, "Failed to list agents")
		return ctrl.Result{}, err
	}

	// Create or update monitoring resources for each agent
	for _, agent := range agents.Items {
		if err := r.setupMonitoringForAgent(ctx, &agent); err != nil {
			logger.Error(err, "Failed to setup monitoring for agent", "agent", agent.Name)
			continue
		}
	}

	logger.Info("Monitoring reconciliation completed")
	return ctrl.Result{RequeueAfter: time.Minute * 10}, nil
}

// setupMonitoringForAgent sets up monitoring resources for a specific agent
func (r *MonitoringReconciler) setupMonitoringForAgent(ctx context.Context, agent *aiv1.Agent) error {
	logger := log.FromContext(ctx).WithValues("agent", agent.Name)

	// Create ServiceMonitor for Prometheus
	if err := r.createServiceMonitor(ctx, agent); err != nil {
		logger.Error(err, "Failed to create ServiceMonitor")
		return err
	}

	// Create Grafana dashboard ConfigMap
	if err := r.createGrafanaDashboard(ctx, agent); err != nil {
		logger.Error(err, "Failed to create Grafana dashboard")
		return err
	}

	return nil
}

// createServiceMonitor creates a ServiceMonitor for Prometheus scraping
func (r *MonitoringReconciler) createServiceMonitor(ctx context.Context, agent *aiv1.Agent) error {
	// This would typically create a ServiceMonitor CRD for Prometheus Operator
	// For now, we'll create a ConfigMap with monitoring configuration
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      agent.Name + "-monitoring",
			Namespace: agent.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":     "kubeagentic-agent",
				"app.kubernetes.io/instance": agent.Name,
				"kubeagentic.ai/agent":       agent.Name,
				"kubeagentic.ai/monitoring":  "true",
			},
		},
		Data: map[string]string{
			"prometheus.yml": fmt.Sprintf(`
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'kubeagentic-agent-%s'
    static_configs:
      - targets: ['%s-service:80']
    metrics_path: '/metrics'
    scrape_interval: 30s
`, agent.Name, agent.Name),
		},
	}

	if err := controllerutil.SetControllerReference(agent, configMap, r.Scheme); err != nil {
		return err
	}

	found := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: configMap.Name, Namespace: configMap.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.FromContext(ctx).Info("Creating monitoring ConfigMap", "ConfigMap.Name", configMap.Name)
		return r.Create(ctx, configMap)
	} else if err != nil {
		return err
	}

	log.FromContext(ctx).Info("Updating monitoring ConfigMap", "ConfigMap.Name", found.Name)
	found.Data = configMap.Data
	return r.Update(ctx, found)
}

// createGrafanaDashboard creates a Grafana dashboard ConfigMap
func (r *MonitoringReconciler) createGrafanaDashboard(ctx context.Context, agent *aiv1.Agent) error {
	dashboard := fmt.Sprintf(`{
  "dashboard": {
    "id": null,
    "title": "KubeAgentic Agent - %s",
    "tags": ["kubeagentic", "ai", "agent"],
    "timezone": "browser",
    "panels": [
      {
        "id": 1,
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(kubeagentic_requests_total{agent=\"%s\"}[5m])",
            "legendFormat": "Requests/sec"
          }
        ],
        "yAxes": [
          {
            "label": "Requests/sec"
          }
        ]
      },
      {
        "id": 2,
        "title": "Response Time",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(kubeagentic_response_duration_seconds_bucket{agent=\"%s\"}[5m]))",
            "legendFormat": "95th percentile"
          }
        ],
        "yAxes": [
          {
            "label": "Seconds"
          }
        ]
      },
      {
        "id": 3,
        "title": "Error Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(kubeagentic_errors_total{agent=\"%s\"}[5m])",
            "legendFormat": "Errors/sec"
          }
        ],
        "yAxes": [
          {
            "label": "Errors/sec"
          }
        ]
      }
    ],
    "time": {
      "from": "now-1h",
      "to": "now"
    },
    "refresh": "30s"
  }
}`, agent.Name, agent.Name, agent.Name, agent.Name)

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      agent.Name + "-grafana-dashboard",
			Namespace: agent.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":     "kubeagentic-agent",
				"app.kubernetes.io/instance": agent.Name,
				"kubeagentic.ai/agent":       agent.Name,
				"grafana_dashboard":          "1",
			},
		},
		Data: map[string]string{
			"dashboard.json": dashboard,
		},
	}

	if err := controllerutil.SetControllerReference(agent, configMap, r.Scheme); err != nil {
		return err
	}

	found := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: configMap.Name, Namespace: configMap.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.FromContext(ctx).Info("Creating Grafana dashboard ConfigMap", "ConfigMap.Name", configMap.Name)
		return r.Create(ctx, configMap)
	} else if err != nil {
		return err
	}

	log.FromContext(ctx).Info("Updating Grafana dashboard ConfigMap", "ConfigMap.Name", found.Name)
	found.Data = configMap.Data
	return r.Update(ctx, found)
}

// SetupWithManager sets up the controller with the Manager
func (r *MonitoringReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aiv1.Agent{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
