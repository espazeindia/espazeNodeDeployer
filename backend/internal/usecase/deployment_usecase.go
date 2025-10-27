package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/espazeindia/espazeNodeDeployer/internal/domain/entities"
	"github.com/espazeindia/espazeNodeDeployer/internal/github"
	"github.com/espazeindia/espazeNodeDeployer/internal/k8s"
	"github.com/espazeindia/espazeNodeDeployer/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeploymentUseCase interface {
	CreateDeployment(ctx context.Context, userID, nodeID primitive.ObjectID, req *entities.DeploymentRequest, githubToken string) (*entities.Deployment, error)
	GetDeployment(ctx context.Context, id primitive.ObjectID) (*entities.Deployment, error)
	GetDeploymentsByNode(ctx context.Context, nodeID primitive.ObjectID) ([]*entities.Deployment, error)
	GetDeploymentsByUser(ctx context.Context, userID primitive.ObjectID) ([]*entities.Deployment, error)
	GetAllDeployments(ctx context.Context, filters map[string]interface{}) ([]*entities.Deployment, error)
	UpdateDeployment(ctx context.Context, id primitive.ObjectID, req *entities.DeploymentUpdateRequest) error
	DeleteDeployment(ctx context.Context, id primitive.ObjectID, k8sClient *k8s.Client) error
	RestartDeployment(ctx context.Context, id primitive.ObjectID, k8sClient *k8s.Client) error
	ScaleDeployment(ctx context.Context, id primitive.ObjectID, replicas int32, k8sClient *k8s.Client) error
	UpdateDeploymentMetrics(ctx context.Context, id primitive.ObjectID, k8sClient *k8s.Client) error
	GetDeploymentStats(ctx context.Context, nodeID *primitive.ObjectID) (map[string]interface{}, error)
}

type deploymentUseCase struct {
	deploymentRepo  repository.DeploymentRepository
	k8sClient       *k8s.Client
	githubClient    *github.Client
	githubTokenRepo repository.GitHubTokenRepository
}

func NewDeploymentUseCase(
	deploymentRepo repository.DeploymentRepository,
	k8sClient *k8s.Client,
	githubClient *github.Client,
	githubTokenRepo repository.GitHubTokenRepository,
) DeploymentUseCase {
	return &deploymentUseCase{
		deploymentRepo:  deploymentRepo,
		k8sClient:       k8sClient,
		githubClient:    githubClient,
		githubTokenRepo: githubTokenRepo,
	}
}

func (uc *deploymentUseCase) CreateDeployment(
	ctx context.Context,
	userID, nodeID primitive.ObjectID,
	req *entities.DeploymentRequest,
	githubToken string,
) (*entities.Deployment, error) {
	// Validate request
	if err := uc.validateDeploymentRequest(req); err != nil {
		return nil, err
	}

	// Get GitHub repository info
	repo, err := uc.githubClient.GetRepository(ctx, githubToken, req.GitHubRepo.Owner, req.GitHubRepo.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}

	// Check for Dockerfile
	hasDockerfile, dockerfilePath, err := uc.githubClient.CheckDockerfile(
		ctx,
		githubToken,
		req.GitHubRepo.Owner,
		req.GitHubRepo.Name,
		req.GitHubRepo.Branch,
	)
	if err != nil || !hasDockerfile {
		return nil, errors.New("repository must contain a Dockerfile")
	}

	// Set default configuration if not provided
	if req.Configuration.Replicas == 0 {
		req.Configuration.Replicas = 2
	}
	if req.Configuration.MemoryRequest == "" {
		req.Configuration.MemoryRequest = "256Mi"
	}
	if req.Configuration.MemoryLimit == "" {
		req.Configuration.MemoryLimit = "512Mi"
	}
	if req.Configuration.CPURequest == "" {
		req.Configuration.CPURequest = "250m"
	}
	if req.Configuration.CPULimit == "" {
		req.Configuration.CPULimit = "500m"
	}
	if req.Configuration.ContainerPort == 0 {
		req.Configuration.ContainerPort = 8080
	}
	if req.Configuration.ServicePort == 0 {
		req.Configuration.ServicePort = 80
	}
	if req.Configuration.ImagePullPolicy == "" {
		req.Configuration.ImagePullPolicy = "IfNotPresent"
	}
	if req.Configuration.RestartPolicy == "" {
		req.Configuration.RestartPolicy = "Always"
	}

	// Set build configuration
	if req.Configuration.BuildConfig.Dockerfile == "" {
		req.Configuration.BuildConfig.Dockerfile = dockerfilePath
	}
	if req.Configuration.BuildConfig.ImageName == "" {
		req.Configuration.BuildConfig.ImageName = fmt.Sprintf("%s/%s", req.GitHubRepo.Owner, req.GitHubRepo.Name)
	}
	if req.Configuration.BuildConfig.ImageTag == "" {
		req.Configuration.BuildConfig.ImageTag = "latest"
	}

	// Create deployment entity
	deployment := &entities.Deployment{
		NodeID:      nodeID,
		UserID:      userID,
		Name:        req.Name,
		ContextPath: req.ContextPath,
		Namespace:   req.Namespace,
		Status:      entities.DeploymentStatusPending,
		GitHubRepo: entities.GitHubRepository{
			Owner:       repo.Owner,
			Name:        repo.Name,
			FullName:    repo.FullName,
			Branch:      req.GitHubRepo.Branch,
			CloneURL:    repo.CloneURL,
			Private:     repo.Private,
			Language:    repo.Language,
			Description: repo.Description,
		},
		Configuration: req.Configuration,
		Metrics: entities.DeploymentMetrics{
			DesiredPods: int(req.Configuration.Replicas),
		},
	}

	if deployment.Namespace == "" {
		deployment.Namespace = "espaze-node-deployer-apps"
	}

	// Save to database
	if err := uc.deploymentRepo.Create(ctx, deployment); err != nil {
		return nil, err
	}

	// Deploy to Kubernetes asynchronously
	go func() {
		deployCtx := context.Background()
		
		// Update status to building
		uc.deploymentRepo.UpdateStatus(deployCtx, deployment.ID, entities.DeploymentStatusBuilding)
		
		// In a real implementation, you would:
		// 1. Clone the repository
		// 2. Build Docker image
		// 3. Push to registry
		// For now, we'll simulate this with a delay
		time.Sleep(5 * time.Second)
		
		// Update status to deploying
		uc.deploymentRepo.UpdateStatus(deployCtx, deployment.ID, entities.DeploymentStatusDeploying)
		
		// Deploy to Kubernetes
		if err := uc.k8sClient.DeployApplication(deployCtx, deployment); err != nil {
			uc.deploymentRepo.UpdateStatus(deployCtx, deployment.ID, entities.DeploymentStatusFailed)
			return
		}
		
		// Update deployment with Kubernetes info
		update := map[string]interface{}{
			"kubernetes_info": deployment.KubernetesInfo,
			"deployed_at":     time.Now(),
		}
		uc.deploymentRepo.Update(deployCtx, deployment.ID, update)
		
		// Update status to running
		uc.deploymentRepo.UpdateStatus(deployCtx, deployment.ID, entities.DeploymentStatusRunning)
	}()

	return deployment, nil
}

func (uc *deploymentUseCase) GetDeployment(ctx context.Context, id primitive.ObjectID) (*entities.Deployment, error) {
	deployment, err := uc.deploymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if deployment == nil {
		return nil, errors.New("deployment not found")
	}
	return deployment, nil
}

func (uc *deploymentUseCase) GetDeploymentsByNode(ctx context.Context, nodeID primitive.ObjectID) ([]*entities.Deployment, error) {
	return uc.deploymentRepo.GetByNodeID(ctx, nodeID)
}

func (uc *deploymentUseCase) GetDeploymentsByUser(ctx context.Context, userID primitive.ObjectID) ([]*entities.Deployment, error) {
	return uc.deploymentRepo.GetByUserID(ctx, userID)
}

func (uc *deploymentUseCase) GetAllDeployments(ctx context.Context, filters map[string]interface{}) ([]*entities.Deployment, error) {
	return uc.deploymentRepo.GetAll(ctx, filters)
}

func (uc *deploymentUseCase) UpdateDeployment(ctx context.Context, id primitive.ObjectID, req *entities.DeploymentUpdateRequest) error {
	deployment, err := uc.deploymentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if deployment == nil {
		return errors.New("deployment not found")
	}

	update := make(map[string]interface{})

	if req.Replicas != nil {
		update["configuration.replicas"] = *req.Replicas
		
		// Update in Kubernetes
		if err := uc.k8sClient.ScaleDeployment(ctx, deployment.Namespace, deployment.KubernetesInfo.DeploymentName, *req.Replicas); err != nil {
			return fmt.Errorf("failed to scale deployment: %w", err)
		}
	}

	if req.EnvironmentVars != nil {
		update["configuration.environment_vars"] = req.EnvironmentVars
	}

	if req.AutoScaling != nil {
		update["configuration.auto_scaling"] = req.AutoScaling
	}

	return uc.deploymentRepo.Update(ctx, id, update)
}

func (uc *deploymentUseCase) DeleteDeployment(ctx context.Context, id primitive.ObjectID, k8sClient *k8s.Client) error {
	deployment, err := uc.deploymentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if deployment == nil {
		return errors.New("deployment not found")
	}

	// Delete from Kubernetes
	if err := k8sClient.DeleteApplication(ctx, deployment.Namespace, deployment.KubernetesInfo.DeploymentName); err != nil {
		return fmt.Errorf("failed to delete from Kubernetes: %w", err)
	}

	// Delete from database
	return uc.deploymentRepo.Delete(ctx, id)
}

func (uc *deploymentUseCase) RestartDeployment(ctx context.Context, id primitive.ObjectID, k8sClient *k8s.Client) error {
	deployment, err := uc.deploymentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if deployment == nil {
		return errors.New("deployment not found")
	}

	return k8sClient.RestartDeployment(ctx, deployment.Namespace, deployment.KubernetesInfo.DeploymentName)
}

func (uc *deploymentUseCase) ScaleDeployment(ctx context.Context, id primitive.ObjectID, replicas int32, k8sClient *k8s.Client) error {
	deployment, err := uc.deploymentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if deployment == nil {
		return errors.New("deployment not found")
	}

	if err := k8sClient.ScaleDeployment(ctx, deployment.Namespace, deployment.KubernetesInfo.DeploymentName, replicas); err != nil {
		return err
	}

	// Update database
	update := map[string]interface{}{
		"configuration.replicas": replicas,
	}
	return uc.deploymentRepo.Update(ctx, id, update)
}

func (uc *deploymentUseCase) UpdateDeploymentMetrics(ctx context.Context, id primitive.ObjectID, k8sClient *k8s.Client) error {
	deployment, err := uc.deploymentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if deployment == nil {
		return errors.New("deployment not found")
	}

	// Get pods
	pods, err := k8sClient.GetPods(ctx, deployment.Namespace)
	if err != nil {
		return err
	}

	// Count pods matching this deployment
	activePods := 0
	readyPods := 0
	for _, pod := range pods.Items {
		if strings.HasPrefix(pod.Name, deployment.KubernetesInfo.DeploymentName) {
			activePods++
			if pod.Status.Phase == "Running" {
				readyPods++
			}
		}
	}

	metrics := &entities.DeploymentMetrics{
		ActivePods:  activePods,
		DesiredPods: int(deployment.Configuration.Replicas),
		ReadyPods:   readyPods,
	}

	return uc.deploymentRepo.UpdateMetrics(ctx, id, metrics)
}

func (uc *deploymentUseCase) GetDeploymentStats(ctx context.Context, nodeID *primitive.ObjectID) (map[string]interface{}, error) {
	return uc.deploymentRepo.GetDeploymentStats(ctx, nodeID)
}

func (uc *deploymentUseCase) validateDeploymentRequest(req *entities.DeploymentRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.ContextPath == "" {
		return errors.New("context path is required")
	}
	if req.GitHubRepo.Owner == "" || req.GitHubRepo.Name == "" {
		return errors.New("GitHub repository owner and name are required")
	}
	if req.GitHubRepo.Branch == "" {
		return errors.New("GitHub branch is required")
	}
	return nil
}

