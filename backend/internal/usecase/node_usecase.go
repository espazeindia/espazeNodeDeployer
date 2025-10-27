package usecase

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"runtime"
	"time"

	"github.com/espazeindia/espazeNodeDeployer/internal/domain/entities"
	"github.com/espazeindia/espazeNodeDeployer/internal/k8s"
	"github.com/espazeindia/espazeNodeDeployer/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NodeUseCase interface {
	RegisterNode(ctx context.Context, req *entities.NodeRegistrationRequest) (*entities.Node, error)
	GetNode(ctx context.Context, id primitive.ObjectID) (*entities.Node, error)
	GetNodeByMac(ctx context.Context, macAddress string) (*entities.Node, error)
	GetAllNodes(ctx context.Context, filters map[string]interface{}) ([]*entities.Node, error)
	UpdateNode(ctx context.Context, id primitive.ObjectID, req *entities.NodeUpdateRequest) error
	DeleteNode(ctx context.Context, id primitive.ObjectID) error
	UpdateNodeResources(ctx context.Context, id primitive.ObjectID, k8sClient *k8s.Client) error
	GetNodeStats(ctx context.Context) (map[string]interface{}, error)
	GetNodesByLocation(ctx context.Context, latitude, longitude, radiusKm float64) ([]*entities.Node, error)
	GetCurrentNodeInfo() (*entities.NodeRegistrationRequest, error)
	Heartbeat(ctx context.Context, nodeID primitive.ObjectID) error
}

type nodeUseCase struct {
	nodeRepo repository.NodeRepository
}

func NewNodeUseCase(nodeRepo repository.NodeRepository) NodeUseCase {
	return &nodeUseCase{
		nodeRepo: nodeRepo,
	}
}

func (uc *nodeUseCase) RegisterNode(ctx context.Context, req *entities.NodeRegistrationRequest) (*entities.Node, error) {
	// Check if node with same MAC address already exists
	existingNode, err := uc.nodeRepo.GetByMacAddress(ctx, req.MacAddress)
	if err != nil {
		return nil, err
	}

	if existingNode != nil {
		// Update existing node
		update := &entities.NodeUpdateRequest{
			Status:      entities.NodeStatusOnline,
			Location:    &req.Location,
			Resources:   &req.Resources,
			ClusterInfo: &req.ClusterInfo,
		}
		if err := uc.nodeRepo.Update(ctx, existingNode.ID, update); err != nil {
			return nil, err
		}
		return uc.nodeRepo.GetByID(ctx, existingNode.ID)
	}

	// Create new node
	node := &entities.Node{
		NodeName:    req.NodeName,
		MacAddress:  req.MacAddress,
		PublicIP:    req.PublicIP,
		PrivateIP:   req.PrivateIP,
		Location:    req.Location,
		Status:      entities.NodeStatusOnline,
		ClusterInfo: req.ClusterInfo,
		Resources:   req.Resources,
		Metadata:    req.Metadata,
	}

	if err := uc.nodeRepo.Create(ctx, node); err != nil {
		return nil, err
	}

	return node, nil
}

func (uc *nodeUseCase) GetNode(ctx context.Context, id primitive.ObjectID) (*entities.Node, error) {
	node, err := uc.nodeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, errors.New("node not found")
	}
	return node, nil
}

func (uc *nodeUseCase) GetNodeByMac(ctx context.Context, macAddress string) (*entities.Node, error) {
	return uc.nodeRepo.GetByMacAddress(ctx, macAddress)
}

func (uc *nodeUseCase) GetAllNodes(ctx context.Context, filters map[string]interface{}) ([]*entities.Node, error) {
	return uc.nodeRepo.GetAll(ctx, filters)
}

func (uc *nodeUseCase) UpdateNode(ctx context.Context, id primitive.ObjectID, req *entities.NodeUpdateRequest) error {
	return uc.nodeRepo.Update(ctx, id, req)
}

func (uc *nodeUseCase) DeleteNode(ctx context.Context, id primitive.ObjectID) error {
	return uc.nodeRepo.Delete(ctx, id)
}

func (uc *nodeUseCase) UpdateNodeResources(ctx context.Context, id primitive.ObjectID, k8sClient *k8s.Client) error {
	// Get cluster info and resource usage
	clusterInfo, err := k8sClient.GetClusterInfo(ctx)
	if err != nil {
		return fmt.Errorf("failed to get cluster info: %w", err)
	}

	// Get nodes to calculate resources
	nodes, err := k8sClient.GetNodes(ctx)
	if err != nil {
		return fmt.Errorf("failed to get nodes: %w", err)
	}

	totalCPU := int64(0)
	totalMemory := int64(0)
	
	for _, node := range nodes.Items {
		cpu := node.Status.Capacity.Cpu().Value()
		memory := node.Status.Capacity.Memory().Value()
		totalCPU += cpu
		totalMemory += memory
	}

	pods, err := k8sClient.GetPods(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to get pods: %w", err)
	}

	resources := &entities.NodeResources{
		CPUCores:     int(totalCPU),
		MemoryTotal:  totalMemory,
		PodsRunning:  len(pods.Items),
		PodsCapacity: 110 * len(nodes.Items), // Default pod capacity per node
	}

	// Try to get metrics if available
	if k8sClient.GetMetricsClientset() != nil {
		// TODO: Implement metrics collection
	}

	return uc.nodeRepo.UpdateResources(ctx, id, resources)
}

func (uc *nodeUseCase) GetNodeStats(ctx context.Context) (map[string]interface{}, error) {
	return uc.nodeRepo.GetNodeStats(ctx)
}

func (uc *nodeUseCase) GetNodesByLocation(ctx context.Context, latitude, longitude, radiusKm float64) ([]*entities.Node, error) {
	return uc.nodeRepo.GetNodesByLocation(ctx, latitude, longitude, radiusKm)
}

func (uc *nodeUseCase) GetCurrentNodeInfo() (*entities.NodeRegistrationRequest, error) {
	// Get MAC address
	macAddress, err := getMacAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get MAC address: %w", err)
	}

	// Get public IP
	publicIP, err := getPublicIP()
	if err != nil {
		publicIP = "unknown"
	}

	// Get private IP
	privateIP, err := getPrivateIP()
	if err != nil {
		privateIP = "unknown"
	}

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// Get location (this would need a geolocation service)
	location := entities.Location{
		Latitude:  0,
		Longitude: 0,
		City:      "Unknown",
		Country:   "Unknown",
		Timezone:  time.Local.String(),
	}

	metadata := entities.NodeMetadata{
		OSType:       runtime.GOOS,
		Architecture: runtime.GOARCH,
		Hostname:     hostname,
		Labels:       make(map[string]string),
		Tags:         []string{},
	}

	req := &entities.NodeRegistrationRequest{
		NodeName:   hostname,
		MacAddress: macAddress,
		PublicIP:   publicIP,
		PrivateIP:  privateIP,
		Location:   location,
		Metadata:   metadata,
	}

	return req, nil
}

func (uc *nodeUseCase) Heartbeat(ctx context.Context, nodeID primitive.ObjectID) error {
	return uc.nodeRepo.UpdateLastSeen(ctx, nodeID)
}

// Helper functions

func getMacAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		// Skip loopback and inactive interfaces
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		mac := iface.HardwareAddr.String()
		if mac != "" {
			return mac, nil
		}
	}

	return "", errors.New("no MAC address found")
}

func getPublicIP() (string, error) {
	// This is a simplified version - in production, use a proper service
	// For now, return empty string
	return "", errors.New("public IP detection not implemented")
}

func getPrivateIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", errors.New("no private IP found")
}

