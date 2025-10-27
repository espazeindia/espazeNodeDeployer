package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Node represents a physical/virtual machine running Kubernetes
type Node struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	NodeName    string             `bson:"node_name" json:"nodeName"`
	MacAddress  string             `bson:"mac_address" json:"macAddress"` // Hardware address
	PublicIP    string             `bson:"public_ip" json:"publicIp"`
	PrivateIP   string             `bson:"private_ip" json:"privateIp"`
	Location    Location           `bson:"location" json:"location"`
	Status      NodeStatus         `bson:"status" json:"status"`
	ClusterInfo ClusterInfo        `bson:"cluster_info" json:"clusterInfo"`
	Resources   NodeResources      `bson:"resources" json:"resources"`
	Metadata    NodeMetadata       `bson:"metadata" json:"metadata"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
	LastSeenAt  time.Time          `bson:"last_seen_at" json:"lastSeenAt"`
}

// Location represents GPS coordinates of the node
type Location struct {
	Latitude  float64 `bson:"latitude" json:"latitude"`
	Longitude float64 `bson:"longitude" json:"longitude"`
	City      string  `bson:"city" json:"city"`
	Country   string  `bson:"country" json:"country"`
	Region    string  `bson:"region" json:"region"`
	Timezone  string  `bson:"timezone" json:"timezone"`
}

// ClusterInfo contains Kubernetes cluster information
type ClusterInfo struct {
	ClusterName    string `bson:"cluster_name" json:"clusterName"`
	KubeVersion    string `bson:"kube_version" json:"kubeVersion"`
	Provider       string `bson:"provider" json:"provider"` // kind, minikube, eks, gke, etc.
	NodesCount     int    `bson:"nodes_count" json:"nodesCount"`
	NamespacesCount int   `bson:"namespaces_count" json:"namespacesCount"`
}

// NodeResources tracks available resources
type NodeResources struct {
	CPUCores       int     `bson:"cpu_cores" json:"cpuCores"`
	CPUUsage       float64 `bson:"cpu_usage" json:"cpuUsage"`           // Percentage
	MemoryTotal    int64   `bson:"memory_total" json:"memoryTotal"`     // Bytes
	MemoryUsed     int64   `bson:"memory_used" json:"memoryUsed"`       // Bytes
	MemoryUsage    float64 `bson:"memory_usage" json:"memoryUsage"`     // Percentage
	DiskTotal      int64   `bson:"disk_total" json:"diskTotal"`         // Bytes
	DiskUsed       int64   `bson:"disk_used" json:"diskUsed"`           // Bytes
	DiskUsage      float64 `bson:"disk_usage" json:"diskUsage"`         // Percentage
	PodsRunning    int     `bson:"pods_running" json:"podsRunning"`
	PodsCapacity   int     `bson:"pods_capacity" json:"podsCapacity"`
}

// NodeMetadata contains additional information
type NodeMetadata struct {
	OSType        string            `bson:"os_type" json:"osType"`           // darwin, linux, windows
	Architecture  string            `bson:"architecture" json:"architecture"` // arm64, amd64
	Hostname      string            `bson:"hostname" json:"hostname"`
	KernelVersion string            `bson:"kernel_version" json:"kernelVersion"`
	Labels        map[string]string `bson:"labels" json:"labels"`
	Tags          []string          `bson:"tags" json:"tags"`
}

// NodeStatus represents the current status of a node
type NodeStatus string

const (
	NodeStatusOnline      NodeStatus = "online"
	NodeStatusOffline     NodeStatus = "offline"
	NodeStatusMaintenance NodeStatus = "maintenance"
	NodeStatusError       NodeStatus = "error"
)

// NodeRegistrationRequest is used when a node registers itself
type NodeRegistrationRequest struct {
	NodeName     string            `json:"nodeName"`
	MacAddress   string            `json:"macAddress"`
	PublicIP     string            `json:"publicIp"`
	PrivateIP    string            `json:"privateIp"`
	Location     Location          `json:"location"`
	ClusterInfo  ClusterInfo       `json:"clusterInfo"`
	Metadata     NodeMetadata      `json:"metadata"`
	Resources    NodeResources     `json:"resources"`
}

// NodeUpdateRequest is used to update node information
type NodeUpdateRequest struct {
	Status      NodeStatus    `json:"status,omitempty"`
	Location    *Location     `json:"location,omitempty"`
	Resources   *NodeResources `json:"resources,omitempty"`
	ClusterInfo *ClusterInfo  `json:"clusterInfo,omitempty"`
}

