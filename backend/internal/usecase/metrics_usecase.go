package usecase

import (
	"context"
	"fmt"

	"github.com/espazeindia/espazeNodeDeployer/internal/k8s"
)

type MetricsUseCase interface {
	GetPodMetrics(ctx context.Context, namespace string) ([]PodMetrics, error)
	GetClusterMetrics(ctx context.Context) (*ClusterMetrics, error)
	GetDeploymentMetrics(ctx context.Context, namespace, deploymentName string) (*DeploymentMetricsInfo, error)
}

type metricsUseCase struct {
	k8sClient *k8s.Client
}

type PodMetrics struct {
	Name         string  `json:"name"`
	Namespace    string  `json:"namespace"`
	CPUUsage     string  `json:"cpuUsage"`
	MemoryUsage  string  `json:"memoryUsage"`
	Status       string  `json:"status"`
	RestartCount int32   `json:"restartCount"`
	Age          string  `json:"age"`
}

type ClusterMetrics struct {
	TotalNodes       int     `json:"totalNodes"`
	TotalPods        int     `json:"totalPods"`
	RunningPods      int     `json:"runningPods"`
	PendingPods      int     `json:"pendingPods"`
	FailedPods       int     `json:"failedPods"`
	TotalCPU         string  `json:"totalCpu"`
	TotalMemory      string  `json:"totalMemory"`
	CPUUsagePercent  float64 `json:"cpuUsagePercent"`
	MemUsagePercent  float64 `json:"memUsagePercent"`
	NamespacesCount  int     `json:"namespacesCount"`
}

type DeploymentMetricsInfo struct {
	Name              string       `json:"name"`
	Namespace         string       `json:"namespace"`
	DesiredReplicas   int32        `json:"desiredReplicas"`
	CurrentReplicas   int32        `json:"currentReplicas"`
	AvailableReplicas int32        `json:"availableReplicas"`
	ReadyReplicas     int32        `json:"readyReplicas"`
	Pods              []PodMetrics `json:"pods"`
	Status            string       `json:"status"`
}

func NewMetricsUseCase(k8sClient *k8s.Client) MetricsUseCase {
	return &metricsUseCase{
		k8sClient: k8sClient,
	}
}

func (uc *metricsUseCase) GetPodMetrics(ctx context.Context, namespace string) ([]PodMetrics, error) {
	if namespace == "" {
		namespace = "espaze-node-deployer-apps"
	}

	pods, err := uc.k8sClient.GetPods(ctx, namespace)
	if err != nil {
		return nil, err
	}

	metrics := make([]PodMetrics, 0, len(pods.Items))
	for _, pod := range pods.Items {
		var restartCount int32
		if len(pod.Status.ContainerStatuses) > 0 {
			restartCount = pod.Status.ContainerStatuses[0].RestartCount
		}

		age := fmt.Sprintf("%v", pod.CreationTimestamp.Time)

		metrics = append(metrics, PodMetrics{
			Name:         pod.Name,
			Namespace:    pod.Namespace,
			Status:       string(pod.Status.Phase),
			RestartCount: restartCount,
			Age:          age,
			CPUUsage:     "N/A",     // Would need metrics-server
			MemoryUsage:  "N/A",     // Would need metrics-server
		})
	}

	return metrics, nil
}

func (uc *metricsUseCase) GetClusterMetrics(ctx context.Context) (*ClusterMetrics, error) {
	// Get nodes
	nodes, err := uc.k8sClient.GetNodes(ctx)
	if err != nil {
		return nil, err
	}

	// Get all pods
	pods, err := uc.k8sClient.GetPods(ctx, "")
	if err != nil {
		return nil, err
	}

	// Get namespaces
	namespaces, err := uc.k8sClient.GetNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	// Count pod statuses
	runningPods := 0
	pendingPods := 0
	failedPods := 0

	for _, pod := range pods.Items {
		switch pod.Status.Phase {
		case "Running":
			runningPods++
		case "Pending":
			pendingPods++
		case "Failed":
			failedPods++
		}
	}

	// Calculate total CPU and Memory
	totalCPU := int64(0)
	totalMemory := int64(0)

	for _, node := range nodes.Items {
		cpu := node.Status.Capacity.Cpu().Value()
		memory := node.Status.Capacity.Memory().Value()
		totalCPU += cpu
		totalMemory += memory
	}

	metrics := &ClusterMetrics{
		TotalNodes:       len(nodes.Items),
		TotalPods:        len(pods.Items),
		RunningPods:      runningPods,
		PendingPods:      pendingPods,
		FailedPods:       failedPods,
		TotalCPU:         fmt.Sprintf("%d", totalCPU),
		TotalMemory:      fmt.Sprintf("%dGi", totalMemory/(1024*1024*1024)),
		CPUUsagePercent:  0.0, // Would need metrics-server
		MemUsagePercent:  0.0, // Would need metrics-server
		NamespacesCount:  len(namespaces),
	}

	return metrics, nil
}

func (uc *metricsUseCase) GetDeploymentMetrics(ctx context.Context, namespace, deploymentName string) (*DeploymentMetricsInfo, error) {
	if namespace == "" {
		namespace = "espaze-node-deployer-apps"
	}

	// Get deployment
	deployment, err := uc.k8sClient.GetClientset().AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// Get pods for this deployment
	pods, err := uc.k8sClient.GetPods(ctx, namespace)
	if err != nil {
		return nil, err
	}

	// Filter pods for this deployment
	deploymentPods := []PodMetrics{}
	for _, pod := range pods.Items {
		if pod.Labels["app"] == deploymentName {
			var restartCount int32
			if len(pod.Status.ContainerStatuses) > 0 {
				restartCount = pod.Status.ContainerStatuses[0].RestartCount
			}

			age := fmt.Sprintf("%v", pod.CreationTimestamp.Time)

			deploymentPods = append(deploymentPods, PodMetrics{
				Name:         pod.Name,
				Namespace:    pod.Namespace,
				Status:       string(pod.Status.Phase),
				RestartCount: restartCount,
				Age:          age,
				CPUUsage:     "N/A",
				MemoryUsage:  "N/A",
			})
		}
	}

	status := "Healthy"
	if deployment.Status.AvailableReplicas < *deployment.Spec.Replicas {
		status = "Degraded"
	}
	if deployment.Status.AvailableReplicas == 0 {
		status = "Unavailable"
	}

	metrics := &DeploymentMetricsInfo{
		Name:              deployment.Name,
		Namespace:         deployment.Namespace,
		DesiredReplicas:   *deployment.Spec.Replicas,
		CurrentReplicas:   deployment.Status.Replicas,
		AvailableReplicas: deployment.Status.AvailableReplicas,
		ReadyReplicas:     deployment.Status.ReadyReplicas,
		Pods:              deploymentPods,
		Status:            status,
	}

	return metrics, nil
}

// Import metav1
import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

