#!/bin/bash

# Kubernetes Installation Script for macOS Apple Silicon
# This script installs and configures Kubernetes using kind (Kubernetes in Docker)

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_header() {
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
}

# Check if running on macOS
if [[ "$OSTYPE" != "darwin"* ]]; then
    print_error "This script is designed for macOS only."
    exit 1
fi

# Check if Apple Silicon
ARCH=$(uname -m)
if [[ "$ARCH" != "arm64" ]]; then
    print_warning "This script is optimized for Apple Silicon (M1/M2/M3). Detected: $ARCH"
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

print_header "Kubernetes Installation for macOS Apple Silicon"

# Step 1: Install Homebrew
print_info "Step 1: Checking Homebrew installation..."
if ! command -v brew &> /dev/null; then
    print_info "Installing Homebrew..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    
    # Add Homebrew to PATH for Apple Silicon
    if [[ "$ARCH" == "arm64" ]]; then
        echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zprofile
        eval "$(/opt/homebrew/bin/brew shellenv)"
    fi
    
    print_success "Homebrew installed successfully"
else
    print_success "Homebrew is already installed"
    brew update
fi

# Step 2: Install Docker Desktop
print_info "Step 2: Checking Docker installation..."
if ! command -v docker &> /dev/null; then
    print_warning "Docker is not installed."
    print_info "Please install Docker Desktop for Mac from:"
    print_info "https://www.docker.com/products/docker-desktop"
    print_info ""
    print_info "After installation, start Docker Desktop and run this script again."
    exit 1
else
    # Check if Docker is running
    if ! docker info &> /dev/null; then
        print_warning "Docker is installed but not running."
        print_info "Please start Docker Desktop and run this script again."
        exit 1
    fi
    print_success "Docker is installed and running"
fi

# Step 3: Install kubectl
print_info "Step 3: Installing kubectl..."
if ! command -v kubectl &> /dev/null; then
    brew install kubectl
    print_success "kubectl installed successfully"
else
    print_success "kubectl is already installed"
    KUBECTL_VERSION=$(kubectl version --client --short 2>/dev/null | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1)
    print_info "kubectl version: $KUBECTL_VERSION"
fi

# Step 4: Install kind (Kubernetes in Docker)
print_info "Step 4: Installing kind..."
if ! command -v kind &> /dev/null; then
    brew install kind
    print_success "kind installed successfully"
else
    print_success "kind is already installed"
    KIND_VERSION=$(kind version | grep -oE '[0-9]+\.[0-9]+\.[0-9]+')
    print_info "kind version: $KIND_VERSION"
fi

# Step 5: Install Helm
print_info "Step 5: Installing Helm..."
if ! command -v helm &> /dev/null; then
    brew install helm
    print_success "Helm installed successfully"
else
    print_success "Helm is already installed"
    HELM_VERSION=$(helm version --short | grep -oE '[0-9]+\.[0-9]+\.[0-9]+')
    print_info "Helm version: $HELM_VERSION"
fi

# Step 6: Create kind cluster configuration
print_info "Step 6: Creating kind cluster configuration..."
CLUSTER_NAME="espaze-node-deployer"

# Check if cluster already exists
if kind get clusters 2>/dev/null | grep -q "^${CLUSTER_NAME}$"; then
    print_warning "Cluster '${CLUSTER_NAME}' already exists."
    read -p "Do you want to delete and recreate it? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_info "Deleting existing cluster..."
        kind delete cluster --name "${CLUSTER_NAME}"
        print_success "Existing cluster deleted"
    else
        print_info "Using existing cluster"
        kubectl cluster-info --context "kind-${CLUSTER_NAME}"
        exit 0
    fi
fi

# Create kind cluster configuration file
cat > /tmp/kind-config.yaml <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: ${CLUSTER_NAME}
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
  - containerPort: 30000
    hostPort: 30000
    protocol: TCP
  - containerPort: 30001
    hostPort: 30001
    protocol: TCP
  - containerPort: 30002
    hostPort: 30002
    protocol: TCP
- role: worker
- role: worker
EOF

print_info "Creating kind cluster '${CLUSTER_NAME}'..."
kind create cluster --config /tmp/kind-config.yaml

print_success "Kind cluster created successfully"

# Step 7: Configure kubectl context
print_info "Step 7: Configuring kubectl context..."
kubectl cluster-info --context "kind-${CLUSTER_NAME}"
kubectl config use-context "kind-${CLUSTER_NAME}"
print_success "kubectl context configured"

# Step 8: Install NGINX Ingress Controller
print_info "Step 8: Installing NGINX Ingress Controller..."
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml

print_info "Waiting for ingress controller to be ready..."
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s

print_success "NGINX Ingress Controller installed"

# Step 9: Install metrics-server
print_info "Step 9: Installing metrics-server..."
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# Patch metrics-server for kind
kubectl patch deployment metrics-server -n kube-system --type='json' -p='[
  {
    "op": "add",
    "path": "/spec/template/spec/containers/0/args/-",
    "value": "--kubelet-insecure-tls"
  },
  {
    "op": "add",
    "path": "/spec/template/spec/containers/0/args/-",
    "value": "--kubelet-preferred-address-types=InternalIP"
  }
]'

print_info "Waiting for metrics-server to be ready..."
kubectl wait --namespace kube-system \
  --for=condition=ready pod \
  --selector=k8s-app=metrics-server \
  --timeout=90s

print_success "metrics-server installed"

# Step 10: Create default namespace for deployments
print_info "Step 10: Creating application namespace..."
kubectl create namespace espaze-node-deployer-apps 2>/dev/null || true
kubectl label namespace espaze-node-deployer-apps name=espaze-node-deployer-apps --overwrite
print_success "Namespace 'espaze-node-deployer-apps' created"

# Step 11: Install k9s (optional Kubernetes CLI manager)
print_info "Step 11: Installing k9s (optional)..."
if ! command -v k9s &> /dev/null; then
    brew install k9s
    print_success "k9s installed successfully"
else
    print_success "k9s is already installed"
fi

# Step 12: Create cluster info file
print_info "Step 12: Saving cluster information..."
cat > ~/espaze-node-deployer-cluster-info.txt <<EOF
Kubernetes Cluster Information
==============================

Cluster Name: ${CLUSTER_NAME}
Context: kind-${CLUSTER_NAME}

Installation Date: $(date)
Architecture: ${ARCH}

Installed Components:
- kubectl: $(kubectl version --client --short 2>/dev/null | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1)
- kind: $(kind version | grep -oE '[0-9]+\.[0-9]+\.[0-9]+')
- helm: $(helm version --short | grep -oE '[0-9]+\.[0-9]+\.[0-9]+')
- Docker: $(docker --version | grep -oE '[0-9]+\.[0-9]+\.[0-9]+')
- Ingress Controller: NGINX
- Metrics Server: Installed
- k9s: $(k9s version --short 2>/dev/null | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' || echo "N/A")

Useful Commands:
================

# View cluster info
kubectl cluster-info

# View all namespaces
kubectl get namespaces

# View nodes
kubectl get nodes

# View all pods
kubectl get pods --all-namespaces

# View services
kubectl get services --all-namespaces

# Open k9s (interactive terminal UI)
k9s

# Delete cluster
kind delete cluster --name ${CLUSTER_NAME}

# Restart cluster
kind delete cluster --name ${CLUSTER_NAME}
./install-k8s-macos.sh

Namespaces:
===========
- default: Default namespace
- kube-system: System components
- kube-public: Public resources
- kube-node-lease: Node heartbeats
- ingress-nginx: Ingress controller
- espaze-node-deployer-apps: Application deployments

Access:
=======
- HTTP: http://localhost
- HTTPS: https://localhost
- NodePort Range: 30000-30002

Configuration File:
===================
~/.kube/config

EOF

print_success "Cluster information saved to ~/espaze-node-deployer-cluster-info.txt"

# Final summary
print_header "Installation Complete! 🎉"

echo ""
print_success "Kubernetes cluster is ready!"
echo ""
print_info "Cluster Name: ${CLUSTER_NAME}"
print_info "Context: kind-${CLUSTER_NAME}"
echo ""
print_info "Next steps:"
echo "  1. View cluster info: kubectl cluster-info"
echo "  2. View nodes: kubectl get nodes"
echo "  3. Open k9s: k9s"
echo "  4. Start the espazeNodeDeployer application"
echo ""
print_info "Cluster information saved to: ~/espaze-node-deployer-cluster-info.txt"
echo ""
print_success "Happy deploying! 🚀"
echo ""

