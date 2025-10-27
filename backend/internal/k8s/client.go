package k8s

import (
	"context"
	"fmt"
	"io"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

type Client struct {
	clientset        *kubernetes.Clientset
	metricsClientset *metricsv.Clientset
	config           *rest.Config
}

func NewClient(kubeconfig string) (*Client, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	metricsClientset, err := metricsv.NewForConfig(config)
	if err != nil {
		// Metrics server might not be available, log but don't fail
		fmt.Printf("Warning: metrics server not available: %v\n", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{Limit: 1})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to cluster: %w", err)
	}

	return &Client{
		clientset:        clientset,
		metricsClientset: metricsClientset,
		config:           config,
	}, nil
}

func (c *Client) GetClientset() *kubernetes.Clientset {
	return c.clientset
}

func (c *Client) GetMetricsClientset() *metricsv.Clientset {
	return c.metricsClientset
}

// Namespace operations
func (c *Client) CreateNamespace(ctx context.Context, name string, labels map[string]string) error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
	}

	_, err := c.clientset.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create namespace: %w", err)
	}

	return nil
}

func (c *Client) GetNamespaces(ctx context.Context) ([]string, error) {
	nsList, err := c.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	namespaces := make([]string, len(nsList.Items))
	for i, ns := range nsList.Items {
		namespaces[i] = ns.Name
	}

	return namespaces, nil
}

// Pod operations
func (c *Client) GetPods(ctx context.Context, namespace string) (*corev1.PodList, error) {
	return c.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
}

func (c *Client) GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	return c.clientset.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (c *Client) GetPodLogs(ctx context.Context, namespace, name string, tailLines int64) (string, error) {
	podLogOpts := corev1.PodLogOptions{
		TailLines: &tailLines,
	}

	req := c.clientset.CoreV1().Pods(namespace).GetLogs(name, &podLogOpts)
	podLogs, err := req.Stream(ctx)
	if err != nil {
		return "", err
	}
	defer podLogs.Close()

	buf := make([]byte, 2000)
	logs := ""
	for {
		n, err := podLogs.Read(buf)
		if n == 0 {
			break
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return logs, err
		}
		logs += string(buf[:n])
	}

	return logs, nil
}

func (c *Client) DeletePod(ctx context.Context, namespace, name string) error {
	return c.clientset.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// Service operations
func (c *Client) GetServices(ctx context.Context, namespace string) (*corev1.ServiceList, error) {
	return c.clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
}

func (c *Client) GetService(ctx context.Context, namespace, name string) (*corev1.Service, error) {
	return c.clientset.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
}

// Node operations
func (c *Client) GetNodes(ctx context.Context) (*corev1.NodeList, error) {
	return c.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
}

func (c *Client) GetNode(ctx context.Context, name string) (*corev1.Node, error) {
	return c.clientset.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
}

// Events operations
func (c *Client) GetEvents(ctx context.Context, namespace string) (*corev1.EventList, error) {
	return c.clientset.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{})
}

// Resource parsing helpers
func ParseMemory(memory string) (*resource.Quantity, error) {
	return resource.ParseQuantity(memory)
}

func ParseCPU(cpu string) (*resource.Quantity, error) {
	return resource.ParseQuantity(cpu)
}

// Cluster info
func (c *Client) GetClusterInfo(ctx context.Context) (map[string]interface{}, error) {
	version, err := c.clientset.Discovery().ServerVersion()
	if err != nil {
		return nil, err
	}

	nodes, err := c.GetNodes(ctx)
	if err != nil {
		return nil, err
	}

	namespaces, err := c.GetNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	pods, err := c.GetPods(ctx, "")
	if err != nil {
		return nil, err
	}

	info := map[string]interface{}{
		"version":         version.GitVersion,
		"nodesCount":      len(nodes.Items),
		"namespacesCount": len(namespaces),
		"podsCount":       len(pods.Items),
		"platform":        version.Platform,
	}

	return info, nil
}

