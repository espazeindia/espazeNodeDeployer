# Espaze Node Deployer - Kubernetes Deployment Management Platform

A comprehensive Kubernetes deployment management application with an intuitive UI for deploying GitHub repositories to Kubernetes clusters with observability dashboards.

## Features

- 🚀 **Kubernetes Installation** - Automated installation script for macOS Apple Silicon
- 📦 **GitHub Integration** - Browse and select GitHub repositories
- ⚙️ **Smart Configuration** - Configure deployments with memory, CPU, and replica settings
- 🎯 **One-Click Deploy** - Deploy applications to Kubernetes with a single click
- 📊 **Observability Dashboards** - Real-time pod monitoring and metrics
- 🎨 **Modern UI** - Beautiful, responsive interface inspired by modern design principles
- 🌓 **Dark Mode** - Full dark mode support
- 🔐 **Secure** - JWT-based authentication and secure GitHub token management

## Tech Stack

### Backend
- **Go** - High-performance backend
- **Fiber** - Express-inspired web framework
- **Kubernetes Client-Go** - Native Kubernetes integration
- **GitHub API** - Repository management
- **MongoDB** - Configuration storage

### Frontend
- **React 18** - Modern UI library
- **Vite** - Lightning-fast build tool
- **TailwindCSS** - Utility-first styling
- **Recharts** - Beautiful charts for observability
- **React Query** - Data fetching and caching
- **Axios** - HTTP client

## Prerequisites

- macOS with Apple Silicon (M1/M2/M3)
- Node.js v18 or higher
- Go 1.21 or higher
- Docker Desktop for Mac
- Git
- GitHub Personal Access Token

## Quick Start

### 1. Install Kubernetes

Run the automated installation script:

```bash
cd espazeNodeDeployer
chmod +x scripts/install-k8s-macos.sh
./scripts/install-k8s-macos.sh
```

This will install:
- Homebrew (if not present)
- kubectl
- kind (Kubernetes in Docker)
- helm
- Creates a local Kubernetes cluster

### 2. Setup Backend

```bash
cd backend
go mod download
cp .env.example .env
# Edit .env with your configurations
go run main.go
```

Backend will start on `http://localhost:8080`

### 3. Setup Frontend

```bash
cd frontend
npm install
cp .env.example .env
# Edit .env with backend URL
npm run dev
```

Frontend will start on `http://localhost:5173`

### 4. Configure GitHub Token

1. Generate a GitHub Personal Access Token with `repo` scope
2. Add it in the application settings or via environment variable

## Project Structure

```
espazeNodeDeployer/
├── backend/                    # Go backend
│   ├── cmd/
│   │   └── server/            # Main application
│   ├── internal/
│   │   ├── api/               # HTTP handlers
│   │   ├── domain/            # Business entities
│   │   ├── k8s/               # Kubernetes client
│   │   ├── github/            # GitHub integration
│   │   ├── repository/        # Data layer
│   │   └── usecase/           # Business logic
│   ├── pkg/
│   │   ├── auth/              # Authentication
│   │   └── utils/             # Utilities
│   ├── go.mod
│   └── .env.example
├── frontend/                   # React frontend
│   ├── src/
│   │   ├── components/        # React components
│   │   ├── pages/             # Page components
│   │   ├── services/          # API services
│   │   ├── hooks/             # Custom hooks
│   │   ├── contexts/          # React contexts
│   │   └── utils/             # Utilities
│   ├── public/
│   ├── package.json
│   └── vite.config.js
├── scripts/
│   └── install-k8s-macos.sh   # K8s installation script
├── docker-compose.yml          # Local development
└── README.md
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/register` - Register
- `POST /api/v1/auth/github/token` - Save GitHub token

### GitHub
- `GET /api/v1/github/repos` - List user repositories
- `GET /api/v1/github/repos/:owner/:repo` - Get repository details
- `GET /api/v1/github/repos/:owner/:repo/branches` - List branches

### Deployments
- `POST /api/v1/deployments` - Create deployment
- `GET /api/v1/deployments` - List deployments
- `GET /api/v1/deployments/:id` - Get deployment details
- `PUT /api/v1/deployments/:id` - Update deployment
- `DELETE /api/v1/deployments/:id` - Delete deployment
- `POST /api/v1/deployments/:id/restart` - Restart deployment

### Kubernetes
- `GET /api/v1/k8s/namespaces` - List namespaces
- `GET /api/v1/k8s/pods` - List pods
- `GET /api/v1/k8s/pods/:namespace/:name` - Get pod details
- `GET /api/v1/k8s/pods/:namespace/:name/logs` - Get pod logs
- `GET /api/v1/k8s/pods/:namespace/:name/metrics` - Get pod metrics

### Observability
- `GET /api/v1/metrics/pods` - Pod metrics
- `GET /api/v1/metrics/cluster` - Cluster metrics
- `GET /api/v1/metrics/deployments/:id` - Deployment metrics

## Features in Detail

### 1. Kubernetes Installation
The automated script handles:
- Detecting existing installations
- Installing required tools
- Creating a local Kind cluster
- Configuring kubectl context
- Installing metrics-server for observability
- Setting up ingress controller

### 2. GitHub Integration
- OAuth GitHub login or Personal Access Token
- Browse all accessible repositories
- View repository details and branches
- Select branch for deployment
- Automatic Dockerfile detection

### 3. Deployment Configuration
- **Context Name**: URL path for accessing the app
- **Memory Limits**: Min and max memory allocation
- **CPU Limits**: CPU resource allocation
- **Replicas**: Number of pod replicas
- **Environment Variables**: Configure app environment
- **Port Configuration**: Container and service ports
- **Auto-scaling**: HPA configuration

### 4. Default Configurations
Pre-configured templates for common stacks:
- Node.js applications
- Go applications
- Python applications
- React/Vue/Angular frontends
- Java/Spring Boot applications
- Custom configurations

### 5. Observability
Real-time dashboards showing:
- Pod status and health
- CPU and memory usage
- Network I/O
- Request rates
- Error rates
- Logs viewer
- Events timeline

## Environment Variables

### Backend (.env)
```env
PORT=8080
MONGODB_URI=mongodb://localhost:27017
DATABASE_NAME=espaze_node_deployer
JWT_SECRET=your-secret-key
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
KUBECONFIG=/Users/yourusername/.kube/config
```

### Frontend (.env)
```env
VITE_API_URL=http://localhost:8080/api/v1
VITE_GITHUB_CLIENT_ID=your-github-client-id
```

## Deployment Flow

1. **Select Repository**: Browse and select GitHub repository
2. **Choose Branch**: Select the branch to deploy
3. **Configure**: Set deployment parameters or use defaults
4. **Deploy**: Application automatically builds and deploys
5. **Monitor**: View real-time metrics and logs
6. **Manage**: Scale, restart, or update deployments

## Design Principles

- **Clean Architecture**: Separation of concerns with clear boundaries
- **SOLID Principles**: Maintainable and extensible code
- **Repository Pattern**: Abstract data access
- **Dependency Injection**: Loose coupling
- **Error Handling**: Comprehensive error management
- **Testing**: Unit and integration tests
- **Documentation**: Well-documented code

## Security

- JWT-based authentication
- GitHub token encryption at rest
- RBAC for Kubernetes access
- Rate limiting on API endpoints
- Input validation and sanitization
- Secure environment variable handling

## Performance Optimization

- Connection pooling
- Caching with Redis (optional)
- Lazy loading in frontend
- Code splitting
- Image optimization
- Efficient Kubernetes queries

## Troubleshooting

### Kubernetes not starting
```bash
kind delete cluster --name espaze-node-deployer
./scripts/install-k8s-macos.sh
```

### Backend connection issues
```bash
# Check if MongoDB is running
brew services list

# Restart MongoDB
brew services restart mongodb-community
```

### Frontend build issues
```bash
# Clear cache and reinstall
rm -rf node_modules package-lock.json
npm install
```

## Repository

**GitHub**: https://github.com/espazeindia/espazeNodeDeployer  
**Organization**: espazeindia  
**Visibility**: Private

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see LICENSE file for details

## Support

For issues and questions:
- GitHub Issues: [Create an issue]
- Documentation: See `/docs` folder
- Email: support@espaze.com

## Roadmap

- [ ] Multi-cluster support
- [ ] GitOps integration
- [ ] Helm chart deployment
- [ ] CI/CD pipeline integration
- [ ] Cost optimization recommendations
- [ ] Advanced monitoring with Prometheus/Grafana
- [ ] Slack/Discord notifications
- [ ] Team collaboration features
- [ ] Role-based access control
- [ ] Audit logs

## Acknowledgments

Built with modern technologies and best practices for Kubernetes deployment management.

