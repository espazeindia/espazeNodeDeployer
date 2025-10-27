package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Deployment represents a deployed application on a Kubernetes cluster
type Deployment struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	NodeID            primitive.ObjectID `bson:"node_id" json:"nodeId"`                   // Reference to the node
	UserID            primitive.ObjectID `bson:"user_id" json:"userId"`                   // Reference to the user
	Name              string             `bson:"name" json:"name"`
	ContextPath       string             `bson:"context_path" json:"contextPath"`         // URL path
	Namespace         string             `bson:"namespace" json:"namespace"`
	Status            DeploymentStatus   `bson:"status" json:"status"`
	GitHubRepo        GitHubRepository   `bson:"github_repo" json:"githubRepo"`
	Configuration     DeploymentConfig   `bson:"configuration" json:"configuration"`
	KubernetesInfo    K8sDeploymentInfo  `bson:"kubernetes_info" json:"kubernetesInfo"`
	Metrics           DeploymentMetrics  `bson:"metrics" json:"metrics"`
	CreatedAt         time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updatedAt"`
	DeployedAt        time.Time          `bson:"deployed_at" json:"deployedAt"`
	LastHealthCheckAt time.Time          `bson:"last_health_check_at" json:"lastHealthCheckAt"`
}

// GitHubRepository contains repository information
type GitHubRepository struct {
	Owner       string `bson:"owner" json:"owner"`
	Name        string `bson:"name" json:"name"`
	FullName    string `bson:"full_name" json:"fullName"`
	Branch      string `bson:"branch" json:"branch"`
	CommitSHA   string `bson:"commit_sha" json:"commitSha"`
	CloneURL    string `bson:"clone_url" json:"cloneUrl"`
	Private     bool   `bson:"private" json:"private"`
	Language    string `bson:"language" json:"language"`
	Description string `bson:"description" json:"description"`
}

// DeploymentConfig contains configuration for deployment
type DeploymentConfig struct {
	Replicas           int32                  `bson:"replicas" json:"replicas"`
	ContainerPort      int32                  `bson:"container_port" json:"containerPort"`
	ServicePort        int32                  `bson:"service_port" json:"servicePort"`
	MemoryRequest      string                 `bson:"memory_request" json:"memoryRequest"`
	MemoryLimit        string                 `bson:"memory_limit" json:"memoryLimit"`
	CPURequest         string                 `bson:"cpu_request" json:"cpuRequest"`
	CPULimit           string                 `bson:"cpu_limit" json:"cpuLimit"`
	EnvironmentVars    map[string]string      `bson:"environment_vars" json:"environmentVars"`
	AutoScaling        AutoScalingConfig      `bson:"auto_scaling" json:"autoScaling"`
	HealthCheck        HealthCheckConfig      `bson:"health_check" json:"healthCheck"`
	ImagePullPolicy    string                 `bson:"image_pull_policy" json:"imagePullPolicy"`
	RestartPolicy      string                 `bson:"restart_policy" json:"restartPolicy"`
	BuildConfig        BuildConfig            `bson:"build_config" json:"buildConfig"`
}

// AutoScalingConfig contains HPA configuration
type AutoScalingConfig struct {
	Enabled                    bool  `bson:"enabled" json:"enabled"`
	MinReplicas                int32 `bson:"min_replicas" json:"minReplicas"`
	MaxReplicas                int32 `bson:"max_replicas" json:"maxReplicas"`
	TargetCPUUtilization       int32 `bson:"target_cpu_utilization" json:"targetCPUUtilization"`
	TargetMemoryUtilization    int32 `bson:"target_memory_utilization" json:"targetMemoryUtilization"`
}

// HealthCheckConfig contains health check settings
type HealthCheckConfig struct {
	Enabled             bool   `bson:"enabled" json:"enabled"`
	Path                string `bson:"path" json:"path"`
	Port                int32  `bson:"port" json:"port"`
	InitialDelaySeconds int32  `bson:"initial_delay_seconds" json:"initialDelaySeconds"`
	PeriodSeconds       int32  `bson:"period_seconds" json:"periodSeconds"`
	TimeoutSeconds      int32  `bson:"timeout_seconds" json:"timeoutSeconds"`
	SuccessThreshold    int32  `bson:"success_threshold" json:"successThreshold"`
	FailureThreshold    int32  `bson:"failure_threshold" json:"failureThreshold"`
}

// BuildConfig contains Docker build configuration
type BuildConfig struct {
	Dockerfile     string            `bson:"dockerfile" json:"dockerfile"`
	BuildContext   string            `bson:"build_context" json:"buildContext"`
	BuildArgs      map[string]string `bson:"build_args" json:"buildArgs"`
	ImageName      string            `bson:"image_name" json:"imageName"`
	ImageTag       string            `bson:"image_tag" json:"imageTag"`
	RegistryURL    string            `bson:"registry_url" json:"registryUrl"`
}

// K8sDeploymentInfo contains actual Kubernetes deployment info
type K8sDeploymentInfo struct {
	DeploymentName  string    `bson:"deployment_name" json:"deploymentName"`
	ServiceName     string    `bson:"service_name" json:"serviceName"`
	IngressName     string    `bson:"ingress_name" json:"ingressName"`
	ConfigMapName   string    `bson:"configmap_name" json:"configMapName"`
	SecretName      string    `bson:"secret_name" json:"secretName"`
	URL             string    `bson:"url" json:"url"`
	InternalURL     string    `bson:"internal_url" json:"internalUrl"`
	PodSelector     string    `bson:"pod_selector" json:"podSelector"`
}

// DeploymentMetrics contains runtime metrics
type DeploymentMetrics struct {
	ActivePods       int       `bson:"active_pods" json:"activePods"`
	DesiredPods      int       `bson:"desired_pods" json:"desiredPods"`
	ReadyPods        int       `bson:"ready_pods" json:"readyPods"`
	CPUUsage         float64   `bson:"cpu_usage" json:"cpuUsage"`
	MemoryUsage      float64   `bson:"memory_usage" json:"memoryUsage"`
	NetworkIn        int64     `bson:"network_in" json:"networkIn"`
	NetworkOut       int64     `bson:"network_out" json:"networkOut"`
	RequestsPerMin   int64     `bson:"requests_per_min" json:"requestsPerMin"`
	ErrorRate        float64   `bson:"error_rate" json:"errorRate"`
	Uptime           int64     `bson:"uptime" json:"uptime"` // seconds
	LastRestartCount int       `bson:"last_restart_count" json:"lastRestartCount"`
	LastRestartTime  time.Time `bson:"last_restart_time" json:"lastRestartTime"`
}

// DeploymentStatus represents deployment status
type DeploymentStatus string

const (
	DeploymentStatusPending    DeploymentStatus = "pending"
	DeploymentStatusBuilding   DeploymentStatus = "building"
	DeploymentStatusDeploying  DeploymentStatus = "deploying"
	DeploymentStatusRunning    DeploymentStatus = "running"
	DeploymentStatusFailed     DeploymentStatus = "failed"
	DeploymentStatusStopped    DeploymentStatus = "stopped"
	DeploymentStatusUpdating   DeploymentStatus = "updating"
)

// DeploymentRequest is used to create a new deployment
type DeploymentRequest struct {
	Name          string           `json:"name" binding:"required"`
	ContextPath   string           `json:"contextPath" binding:"required"`
	GitHubRepo    GitHubRepository `json:"githubRepo" binding:"required"`
	Configuration DeploymentConfig `json:"configuration"`
	Namespace     string           `json:"namespace"`
}

// DeploymentUpdateRequest is used to update deployment
type DeploymentUpdateRequest struct {
	Replicas        *int32                 `json:"replicas,omitempty"`
	EnvironmentVars map[string]string      `json:"environmentVars,omitempty"`
	AutoScaling     *AutoScalingConfig     `json:"autoScaling,omitempty"`
}

