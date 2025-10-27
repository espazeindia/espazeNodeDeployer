#!/bin/bash

# Quick Start Script for Espaze Node Deployer
# This script automates the entire setup process

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m'

echo -e "${PURPLE}"
echo "╔═══════════════════════════════════════════════════════════╗"
echo "║                                                           ║"
echo "║        Espaze Node Deployer - Quick Start                ║"
echo "║                                                           ║"
echo "║  Automated setup for Kubernetes Deployment Platform      ║"
echo "║                                                           ║"
echo "╚═══════════════════════════════════════════════════════════╝"
echo -e "${NC}"

# Function to print steps
print_step() {
    echo -e "\n${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Check if running on macOS
if [[ "$OSTYPE" != "darwin"* ]]; then
    print_error "This script is designed for macOS only"
    exit 1
fi

# Step 1: Install Kubernetes
print_step "Step 1/6: Installing Kubernetes"
if command -v kubectl &> /dev/null && command -v kind &> /dev/null; then
    print_success "Kubernetes tools already installed"
else
    print_warning "Running Kubernetes installation script..."
    chmod +x scripts/install-k8s-macos.sh
    ./scripts/install-k8s-macos.sh
fi

# Step 2: Setup MongoDB
print_step "Step 2/6: Setting up MongoDB"
if docker ps | grep -q espaze-node-deployer-mongodb; then
    print_success "MongoDB container already running"
else
    print_warning "Starting MongoDB container..."
    docker run -d \
        --name espaze-node-deployer-mongodb \
        -p 27017:27017 \
        -v espaze-node-deployer-mongo-data:/data/db \
        mongo:7.0
    sleep 3
    print_success "MongoDB started successfully"
fi

# Step 3: Setup Backend
print_step "Step 3/6: Setting up Backend"
cd backend

if [ ! -f ".env" ]; then
    cp .env.example .env
    print_warning "Created .env file - please update with your settings"
fi

print_warning "Installing Go dependencies..."
go mod download

print_success "Backend setup complete"
cd ..

# Step 4: Setup Frontend
print_step "Step 4/6: Setting up Frontend"
cd frontend

if [ ! -f ".env" ]; then
    cp .env.example .env
    print_warning "Created frontend .env file"
fi

if [ ! -d "node_modules" ]; then
    print_warning "Installing npm dependencies (this may take a few minutes)..."
    npm install
fi

print_success "Frontend setup complete"
cd ..

# Step 5: Generate startup scripts
print_step "Step 5/6: Creating startup scripts"

# Create backend start script
cat > start-backend.sh << 'EOF'
#!/bin/bash
cd backend
echo "🚀 Starting Backend Server..."
echo "📍 API will be available at http://localhost:8080"
go run cmd/server/main.go
EOF
chmod +x start-backend.sh

# Create frontend start script
cat > start-frontend.sh << 'EOF'
#!/bin/bash
cd frontend
echo "🚀 Starting Frontend Server..."
echo "📍 App will be available at http://localhost:5173"
npm run dev
EOF
chmod +x start-frontend.sh

print_success "Startup scripts created"

# Step 6: Display instructions
print_step "Step 6/6: Setup Complete! 🎉"

echo -e "${GREEN}"
echo "┌─────────────────────────────────────────────────────────┐"
echo "│                                                         │"
echo "│  ✅ Espaze Node Deployer is ready to use!             │"
echo "│                                                         │"
echo "└─────────────────────────────────────────────────────────┘"
echo -e "${NC}"

echo -e "\n${YELLOW}📋 Next Steps:${NC}\n"
echo "1. Open TWO terminal windows"
echo ""
echo -e "   ${BLUE}Terminal 1 - Start Backend:${NC}"
echo "   $ ./start-backend.sh"
echo ""
echo -e "   ${BLUE}Terminal 2 - Start Frontend:${NC}"
echo "   $ ./start-frontend.sh"
echo ""
echo -e "2. Open your browser to ${GREEN}http://localhost:5173${NC}"
echo ""
echo "3. Register a new account"
echo ""
echo "4. Configure GitHub token in Settings"
echo ""
echo "5. Start deploying! 🚀"
echo ""

echo -e "${PURPLE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${PURPLE}  Useful Commands:${NC}"
echo -e "${PURPLE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "  kubectl get pods -n espaze-node-deployer-apps  # View deployed pods"
echo "  kubectl get services                            # View services"
echo "  k9s                                             # Interactive K8s UI"
echo "  docker logs espaze-node-deployer-mongodb       # View MongoDB logs"
echo ""

echo -e "${GREEN}Happy deploying! 🎊${NC}\n"

