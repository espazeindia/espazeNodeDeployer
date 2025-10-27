package usecase

import (
	"context"

	"github.com/espazeindia/espazeNodeDeployer/internal/k8s"
	corev1 "k8s.io/api/core/v1"
)

type K8sUseCase interface {
	GetNamespaces(ctx context.Context) ([]string, error)
	GetPods(ctx context.Context, namespace string) (*corev1.PodList, error)
	GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error)
	GetPodLogs(ctx context.Context, namespace, name string, tailLines int64) (string, error)
	GetServices(ctx context.Context, namespace string) (*corev1.ServiceList, error)
	GetNodes(ctx context.Context) (*corev1.NodeList, error)
	GetEvents(ctx context.Context, namespace string) (*corev1.EventList, error)
	GetClusterInfo(ctx context.Context) (map[string]interface{}, error)
}

type k8sUseCase struct {
	k8sClient *k8s.Client
}

func NewK8sUseCase(k8sClient *k8s.Client) K8sUseCase {
	return &k8sUseCase{
		k8sClient: k8sClient,
	}
}

func (uc *k8sUseCase) GetNamespaces(ctx context.Context) ([]string, error) {
	return uc.k8sClient.GetNamespaces(ctx)
}

func (uc *k8sUseCase) GetPods(ctx context.Context, namespace string) (*corev1.PodList, error) {
	if namespace == "" {
		namespace = "espaze-node-deployer-apps"
	}
	return uc.k8sClient.GetPods(ctx, namespace)
}

func (uc *k8sUseCase) GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	return uc.k8sClient.GetPod(ctx, namespace, name)
}

func (uc *k8sUseCase) GetPodLogs(ctx context.Context, namespace, name string, tailLines int64) (string, error) {
	if tailLines == 0 {
		tailLines = 100
	}
	return uc.k8sClient.GetPodLogs(ctx, namespace, name, tailLines)
}

func (uc *k8sUseCase) GetServices(ctx context.Context, namespace string) (*corev1.ServiceList, error) {
	if namespace == "" {
		namespace = "espaze-node-deployer-apps"
	}
	return uc.k8sClient.GetServices(ctx, namespace)
}

func (uc *k8sUseCase) GetNodes(ctx context.Context) (*corev1.NodeList, error) {
	return uc.k8sClient.GetNodes(ctx)
}

func (uc *k8sUseCase) GetEvents(ctx context.Context, namespace string) (*corev1.EventList, error) {
	if namespace == "" {
		namespace = "espaze-node-deployer-apps"
	}
	return uc.k8sClient.GetEvents(ctx, namespace)
}

func (uc *k8sUseCase) GetClusterInfo(ctx context.Context) (map[string]interface{}, error) {
	return uc.k8sClient.GetClusterInfo(ctx)
}

