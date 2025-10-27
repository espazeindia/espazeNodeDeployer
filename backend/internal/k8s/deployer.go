package k8s

import (
	"context"
	"fmt"
	"strings"

	"github.com/espaze/espazeNodeDeployer/internal/domain/entities"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// DeployApplication creates Kubernetes resources for a deployment
func (c *Client) DeployApplication(ctx context.Context, deployment *entities.Deployment) error {
	namespace := deployment.Namespace
	if namespace == "" {
		namespace = "espaze-node-deployer-apps"
	}

	// Ensure namespace exists
	_, err := c.clientset.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if err != nil {
		if err := c.CreateNamespace(ctx, namespace, map[string]string{
			"managed-by": "espaze-node-deployer",
		}); err != nil {
			return fmt.Errorf("failed to create namespace: %w", err)
		}
	}

	// Create deployment name (sanitized)
	deploymentName := sanitizeName(deployment.Name)
	
	// 1. Create ConfigMap for environment variables (if any)
	if len(deployment.Configuration.EnvironmentVars) > 0 {
		if err := c.createConfigMap(ctx, namespace, deploymentName, deployment.Configuration.EnvironmentVars); err != nil {
			return fmt.Errorf("failed to create configmap: %w", err)
		}
		deployment.KubernetesInfo.ConfigMapName = fmt.Sprintf("%s-config", deploymentName)
	}

	// 2. Create Deployment
	if err := c.createDeployment(ctx, namespace, deployment); err != nil {
		return fmt.Errorf("failed to create deployment: %w", err)
	}
	deployment.KubernetesInfo.DeploymentName = deploymentName

	// 3. Create Service
	if err := c.createService(ctx, namespace, deploymentName, deployment.Configuration.ServicePort, deployment.Configuration.ContainerPort); err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}
	deployment.KubernetesInfo.ServiceName = fmt.Sprintf("%s-service", deploymentName)

	// 4. Create Ingress
	if deployment.ContextPath != "" {
		if err := c.createIngress(ctx, namespace, deploymentName, deployment.ContextPath, deployment.Configuration.ServicePort); err != nil {
			return fmt.Errorf("failed to create ingress: %w", err)
		}
		deployment.KubernetesInfo.IngressName = fmt.Sprintf("%s-ingress", deploymentName)
		deployment.KubernetesInfo.URL = fmt.Sprintf("http://localhost%s", deployment.ContextPath)
	}

	deployment.KubernetesInfo.InternalURL = fmt.Sprintf("http://%s-service.%s.svc.cluster.local:%d", deploymentName, namespace, deployment.Configuration.ServicePort)
	deployment.KubernetesInfo.PodSelector = fmt.Sprintf("app=%s", deploymentName)

	return nil
}

func (c *Client) createConfigMap(ctx context.Context, namespace, name string, data map[string]string) error {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-config", name),
			Namespace: namespace,
			Labels: map[string]string{
				"app":        name,
				"managed-by": "espaze-node-deployer",
			},
		},
		Data: data,
	}

	_, err := c.clientset.CoreV1().ConfigMaps(namespace).Create(ctx, configMap, metav1.CreateOptions{})
	return err
}

func (c *Client) createDeployment(ctx context.Context, namespace string, deployment *entities.Deployment) error {
	deploymentName := sanitizeName(deployment.Name)
	config := deployment.Configuration

	// Parse resource quantities
	memoryRequest, _ := resource.ParseQuantity(config.MemoryRequest)
	memoryLimit, _ := resource.ParseQuantity(config.MemoryLimit)
	cpuRequest, _ := resource.ParseQuantity(config.CPURequest)
	cpuLimit, _ := resource.ParseQuantity(config.CPULimit)

	// Build environment variables
	envVars := []corev1.EnvVar{}
	for key, value := range config.EnvironmentVars {
		envVars = append(envVars, corev1.EnvVar{
			Name:  key,
			Value: value,
		})
	}

	// Create deployment spec
	k8sDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: namespace,
			Labels: map[string]string{
				"app":        deploymentName,
				"managed-by": "espaze-node-deployer",
				"repo":       deployment.GitHubRepo.FullName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &config.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deploymentName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": deploymentName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            deploymentName,
							Image:           config.BuildConfig.ImageName + ":" + config.BuildConfig.ImageTag,
							ImagePullPolicy: corev1.PullPolicy(config.ImagePullPolicy),
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: config.ContainerPort,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							Env: envVars,
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceMemory: memoryRequest,
									corev1.ResourceCPU:    cpuRequest,
								},
								Limits: corev1.ResourceList{
									corev1.ResourceMemory: memoryLimit,
									corev1.ResourceCPU:    cpuLimit,
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicy(config.RestartPolicy),
				},
			},
		},
	}

	// Add health checks if enabled
	if config.HealthCheck.Enabled {
		k8sDeployment.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: config.HealthCheck.Path,
					Port: intstr.FromInt(int(config.HealthCheck.Port)),
				},
			},
			InitialDelaySeconds: config.HealthCheck.InitialDelaySeconds,
			PeriodSeconds:       config.HealthCheck.PeriodSeconds,
			TimeoutSeconds:      config.HealthCheck.TimeoutSeconds,
			SuccessThreshold:    config.HealthCheck.SuccessThreshold,
			FailureThreshold:    config.HealthCheck.FailureThreshold,
		}

		k8sDeployment.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: config.HealthCheck.Path,
					Port: intstr.FromInt(int(config.HealthCheck.Port)),
				},
			},
			InitialDelaySeconds: 5,
			PeriodSeconds:       config.HealthCheck.PeriodSeconds,
			TimeoutSeconds:      config.HealthCheck.TimeoutSeconds,
			SuccessThreshold:    config.HealthCheck.SuccessThreshold,
			FailureThreshold:    config.HealthCheck.FailureThreshold,
		}
	}

	_, err := c.clientset.AppsV1().Deployments(namespace).Create(ctx, k8sDeployment, metav1.CreateOptions{})
	return err
}

func (c *Client) createService(ctx context.Context, namespace, name string, servicePort, targetPort int32) error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-service", name),
			Namespace: namespace,
			Labels: map[string]string{
				"app":        name,
				"managed-by": "espaze-node-deployer",
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": name,
			},
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       servicePort,
					TargetPort: intstr.FromInt(int(targetPort)),
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	_, err := c.clientset.CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
	return err
}

func (c *Client) createIngress(ctx context.Context, namespace, name, path string, servicePort int32) error {
	pathType := networkingv1.PathTypePrefix

	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-ingress", name),
			Namespace: namespace,
			Labels: map[string]string{
				"app":        name,
				"managed-by": "espaze-node-deployer",
			},
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target": "/",
			},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     path,
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: fmt.Sprintf("%s-service", name),
											Port: networkingv1.ServiceBackendPort{
												Number: servicePort,
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

	_, err := c.clientset.NetworkingV1().Ingresses(namespace).Create(ctx, ingress, metav1.CreateOptions{})
	return err
}

// DeleteApplication removes all Kubernetes resources for a deployment
func (c *Client) DeleteApplication(ctx context.Context, namespace, deploymentName string) error {
	name := sanitizeName(deploymentName)

	// Delete in reverse order
	propagationPolicy := metav1.DeletePropagationForeground

	// Delete Ingress
	c.clientset.NetworkingV1().Ingresses(namespace).Delete(ctx, fmt.Sprintf("%s-ingress", name), metav1.DeleteOptions{
		PropagationPolicy: &propagationPolicy,
	})

	// Delete Service
	c.clientset.CoreV1().Services(namespace).Delete(ctx, fmt.Sprintf("%s-service", name), metav1.DeleteOptions{})

	// Delete Deployment
	c.clientset.AppsV1().Deployments(namespace).Delete(ctx, name, metav1.DeleteOptions{
		PropagationPolicy: &propagationPolicy,
	})

	// Delete ConfigMap
	c.clientset.CoreV1().ConfigMaps(namespace).Delete(ctx, fmt.Sprintf("%s-config", name), metav1.DeleteOptions{})

	return nil
}

// ScaleDeployment scales a deployment to the specified number of replicas
func (c *Client) ScaleDeployment(ctx context.Context, namespace, name string, replicas int32) error {
	deployment, err := c.clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	deployment.Spec.Replicas = &replicas

	_, err = c.clientset.AppsV1().Deployments(namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	return err
}

// RestartDeployment restarts a deployment by updating an annotation
func (c *Client) RestartDeployment(ctx context.Context, namespace, name string) error {
	deployment, err := c.clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if deployment.Spec.Template.ObjectMeta.Annotations == nil {
		deployment.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
	}

	deployment.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = metav1.Now().Format("2006-01-02T15:04:05Z07:00")

	_, err = c.clientset.AppsV1().Deployments(namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	return err
}

// Helper function to sanitize names for Kubernetes
func sanitizeName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "_", "-")
	name = strings.ReplaceAll(name, " ", "-")
	// Remove any characters that aren't alphanumeric or hyphens
	result := ""
	for _, char := range name {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			result += string(char)
		}
	}
	// Ensure it doesn't start or end with a hyphen
	result = strings.Trim(result, "-")
	// Limit length to 63 characters (Kubernetes limit)
	if len(result) > 63 {
		result = result[:63]
	}
	return result
}

