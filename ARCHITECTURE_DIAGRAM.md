# Espaze Node Deployer - System Architecture

## High-Level Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              USER LAYER                                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Browser    │  │   Mobile     │  │    Tablet    │  │     CLI      │      │
│  │  (Chrome)    │  │   (Safari)   │  │    (iPad)    │  │   (curl)     │      │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘      │
└─────────┼──────────────────┼──────────────────┼──────────────────┼─────────────┘
          │                  │                  │                  │
          └──────────────────┴──────────────────┴──────────────────┘
                                     │
                                     │ HTTPS
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         FRONTEND LAYER (React)                                  │
│  ┌───────────────────────────────────────────────────────────────────────────┐ │
│  │                          React Application                                 │ │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────┐         │ │
│  │  │   Pages    │  │ Components │  │   Hooks    │  │   Store    │         │ │
│  │  │            │  │            │  │            │  │            │         │ │
│  │  │ Dashboard  │  │   Layout   │  │  useQuery  │  │ authStore  │         │ │
│  │  │ Deployments│  │   Cards    │  │  useMutation│  │            │         │ │
│  │  │   Nodes    │  │   Tables   │  │            │  │            │         │ │
│  │  │Observatory │  │   Charts   │  │            │  │            │         │ │
│  │  └────────────┘  └────────────┘  └────────────┘  └────────────┘         │ │
│  │                                                                            │ │
│  │  ┌─────────────────────────────────────────────────────────────┐         │ │
│  │  │                    API Service Layer                         │         │ │
│  │  │  • axios with interceptors                                   │         │ │
│  │  │  • Automatic token injection                                 │         │ │
│  │  │  • Error handling & retry logic                              │         │ │
│  │  │  • Request/Response transformation                           │         │ │
│  │  └─────────────────────────────────────────────────────────────┘         │ │
│  └───────────────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────┬───────────────────────────────────────────┘
                                      │ REST API (JSON)
                                      │ JWT Authentication
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        BACKEND LAYER (Go + Fiber)                               │
│  ┌───────────────────────────────────────────────────────────────────────────┐ │
│  │                         API Handler Layer                                  │ │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │ │
│  │  │    Auth     │  │    Nodes    │  │ Deployments │  │   GitHub    │     │ │
│  │  │   Handler   │  │   Handler   │  │   Handler   │  │   Handler   │     │ │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘     │ │
│  │  ┌─────────────┐  ┌─────────────┐                                        │ │
│  │  │   K8s       │  │  Metrics    │                                        │ │
│  │  │  Handler    │  │   Handler   │                                        │ │
│  │  └─────────────┘  └─────────────┘                                        │ │
│  │                                                                            │ │
│  │  ┌─────────────────────────────────────────────────────────────┐         │ │
│  │  │                    Middleware Layer                          │         │ │
│  │  │  • JWT Validation                                            │         │ │
│  │  │  • CORS Headers                                              │         │ │
│  │  │  • Request Logging                                           │         │ │
│  │  │  • Error Recovery                                            │         │ │
│  │  └─────────────────────────────────────────────────────────────┘         │ │
│  └───────────────────────────────────────────────────────────────────────────┘ │
│                                       │                                         │
│  ┌───────────────────────────────────┼───────────────────────────────────────┐ │
│  │                       Use Case Layer (Business Logic)                      │ │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │ │
│  │  │    Auth      │  │    Node      │  │  Deployment  │  │   GitHub     │ │ │
│  │  │  Use Case    │  │  Use Case    │  │   Use Case   │  │  Use Case    │ │ │
│  │  │              │  │              │  │              │  │              │ │ │
│  │  │ • Login      │  │ • Register   │  │ • Create     │  │ • List Repos │ │ │
│  │  │ • Register   │  │ • Heartbeat  │  │ • Scale      │  │ • Get Branch │ │ │
│  │  │ • Validate   │  │ • Update     │  │ • Restart    │  │ • Search     │ │ │
│  │  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘ │ │
│  │  ┌──────────────┐  ┌──────────────┐                                      │ │
│  │  │     K8s      │  │   Metrics    │                                      │ │
│  │  │  Use Case    │  │  Use Case    │                                      │ │
│  │  │              │  │              │                                      │ │
│  │  │ • Deploy     │  │ • Collect    │                                      │ │
│  │  │ • Get Pods   │  │ • Aggregate  │                                      │ │
│  │  │ • Get Logs   │  │ • Monitor    │                                      │ │
│  │  └──────────────┘  └──────────────┘                                      │ │
│  └───────────────────────────────────────────────────────────────────────────┘ │
│                                       │                                         │
│  ┌───────────────────────────────────┼───────────────────────────────────────┐ │
│  │                  Repository Layer (Data Access)                            │ │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │ │
│  │  │    User     │  │    Node     │  │ Deployment  │  │GitHub Token │     │ │
│  │  │ Repository  │  │ Repository  │  │ Repository  │  │ Repository  │     │ │
│  │  │             │  │             │  │             │  │             │     │ │
│  │  │ • CRUD      │  │ • CRUD      │  │ • CRUD      │  │ • CRUD      │     │ │
│  │  │ • Find      │  │ • Stats     │  │ • Filter    │  │ • Encrypt   │     │ │
│  │  │ • Update    │  │ • Heartbeat │  │ • Status    │  │             │     │ │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘     │ │
│  └───────────────────────────────────────────────────────────────────────────┘ │
└────────────┬──────────────────────────────┬──────────────────┬─────────────────┘
             │                              │                  │
             ▼                              ▼                  ▼
┌────────────────────────┐    ┌────────────────────────┐    ┌─────────────────────┐
│   MONGODB DATABASE     │    │   KUBERNETES CLUSTER   │    │   GITHUB API        │
│                        │    │                        │    │                     │
│  Collections:          │    │  Resources:            │    │  Endpoints:         │
│  • users               │    │  • Deployments         │    │  • /repos           │
│  • nodes               │    │  • Services            │    │  • /branches        │
│  • deployments         │    │  • Ingress             │    │  • /contents        │
│  • github_tokens       │    │  • ConfigMaps          │    │  • /commits         │
│                        │    │  • Pods                │    │                     │
│  Features:             │    │  • Namespaces          │    │  Authentication:    │
│  • Indexed queries     │    │                        │    │  • Personal Token   │
│  • Aggregation         │    │  Monitoring:           │    │  • OAuth            │
│  • Replication         │    │  • Metrics Server      │    │                     │
│  • Sharding Ready      │    │  • Events              │    │                     │
└────────────────────────┘    └────────────────────────┘    └─────────────────────┘
```

## Data Flow Diagrams

### 1. User Registration & Login Flow

```
User Browser                Frontend                Backend                MongoDB
     │                         │                       │                      │
     │─── Register Form ──────>│                       │                      │
     │                         │─── POST /auth/register >│                      │
     │                         │                       │── Hash Password ──>│
     │                         │                       │                      │
     │                         │                       │<── Save User ───────│
     │                         │<── JWT Token ─────────│                      │
     │<── Success Message ─────│                       │                      │
     │                         │                       │                      │
     │─── Login Form ─────────>│                       │                      │
     │                         │─── POST /auth/login ──>│                      │
     │                         │                       │── Verify Password ─>│
     │                         │                       │<── Get User ────────│
     │                         │                       │── Generate JWT ──>  │
     │                         │<── JWT + User Info ───│                      │
     │<── Redirect Dashboard ──│                       │                      │
```

### 2. Node Registration Flow

```
Node Machine            Frontend                Backend                MongoDB
     │                     │                       │                      │
     │── Get System Info ──>│                       │                      │
     │   (MAC, IP, GPS)     │                       │                      │
     │                     │─── POST /nodes/register >│                      │
     │                     │   {                    │                      │
     │                     │     macAddress,        │                      │
     │                     │     publicIp,          │                      │
     │                     │     location           │                      │
     │                     │   }                    │                      │
     │                     │                       │── Check Existing ───>│
     │                     │                       │<── Return or Create ─│
     │                     │<── Node Object ────────│                      │
     │<── Display Success ──│                       │                      │
     │                     │                       │                      │
     │── Heartbeat Loop ───>│                       │                      │
     │   (every 30s)        │─── POST /nodes/:id/heartbeat >│              │
     │                     │                       │── Update lastSeenAt ─>│
     │                     │<── 200 OK ─────────────│                      │
```

### 3. Deployment Creation Flow

```
User                Frontend            Backend            GitHub API        Kubernetes        MongoDB
 │                     │                   │                    │                 │               │
 │─ Select Repo ──────>│                   │                    │                 │               │
 │                     │── GET /github/repos >│                 │                 │               │
 │                     │                   │── List Repos ─────>│                 │               │
 │                     │                   │<── Repositories ───│                 │               │
 │                     │<── Repo List ─────│                    │                 │               │
 │<─ Display Repos ────│                   │                    │                 │               │
 │                     │                   │                    │                 │               │
 │─ Configure & Deploy >│                   │                    │                 │               │
 │                     │── POST /deployments >│                 │                 │               │
 │                     │   {                │                    │                 │               │
 │                     │     name,          │                    │                 │               │
 │                     │     githubRepo,    │                    │                 │               │
 │                     │     configuration  │                    │                 │               │
 │                     │   }                │                    │                 │               │
 │                     │                   │── Validate Repo ───>│                 │               │
 │                     │                   │── Check Dockerfile ─>│                 │               │
 │                     │                   │<── Dockerfile Exists│                 │               │
 │                     │                   │                    │                 │               │
 │                     │                   │── Save Deployment ────────────────────────────────────>│
 │                     │                   │                    │                 │               │
 │                     │<── Deployment ID ──│                    │                 │               │
 │<─ "Deploying..." ───│                   │                    │                 │               │
 │                     │                   │                    │                 │               │
 │                     │                   │[Async Background]  │                 │               │
 │                     │                   │── Create K8s Resources ──────────────>│               │
 │                     │                   │                    │                 │               │
 │                     │                   │                    │      ┌──────────▼──────────┐    │
 │                     │                   │                    │      │  • Deployment       │    │
 │                     │                   │                    │      │  • Service          │    │
 │                     │                   │                    │      │  • Ingress          │    │
 │                     │                   │                    │      │  • ConfigMap        │    │
 │                     │                   │                    │      └──────────┬──────────┘    │
 │                     │                   │                    │                 │               │
 │                     │                   │                    │      <── Pod Running            │
 │                     │                   │                    │                 │               │
 │                     │                   │── Update Status ───────────────────────────────────>│
 │                     │                   │   (status: "running")│                 │               │
 │                     │                   │                    │                 │               │
 │[Polling Updates]    │                   │                    │                 │               │
 │<── Status Updates ──│<── GET /deployments/:id ──────────────────────────────────────────────>│
 │    "Building..."    │                   │                    │                 │               │
 │    "Deploying..."   │                   │                    │                 │               │
 │    "Running!"       │                   │                    │                 │               │
```

### 4. Observability Metrics Flow

```
Frontend            Backend            Kubernetes        Metrics Server        MongoDB
   │                   │                    │                   │                 │
   │── View Dashboard ─>│                   │                   │                 │
   │                   │── GET /metrics/cluster >│              │                 │
   │                   │                    │                   │                 │
   │                   │── Get Nodes ──────>│                   │                 │
   │                   │<── Node List ──────│                   │                 │
   │                   │                    │                   │                 │
   │                   │── Get Pods ───────>│                   │                 │
   │                   │<── Pod List ───────│                   │                 │
   │                   │                    │                   │                 │
   │                   │── Get Metrics ─────────────────────────>│                 │
   │                   │<── CPU/Memory Data ────────────────────│                 │
   │                   │                    │                   │                 │
   │                   │── Calculate Stats ─┤                   │                 │
   │                   │   • Total CPU      │                   │                 │
   │                   │   • Total Memory   │                   │                 │
   │                   │   • Pod Counts     │                   │                 │
   │                   │                    │                   │                 │
   │<── Metrics JSON ───│                    │                   │                 │
   │                   │                    │                   │                 │
   │── Render Charts ──┤                    │                   │                 │
   │   • Line Charts   │                    │                   │                 │
   │   • Bar Charts    │                    │                   │                 │
   │   • Stats Cards   │                    │                   │                 │
   │                   │                    │                   │                 │
   │[Auto-refresh 10s] │                    │                   │                 │
   │── Refetch ────────>│── GET /metrics/pods >│                │                 │
   │<── Updated Data ───│                    │                   │                 │
```

## Component Interaction Matrix

| Component      | Interacts With                           | Purpose                    |
|----------------|------------------------------------------|----------------------------|
| Frontend       | Backend API, Browser Storage             | UI & User Interaction      |
| Backend API    | MongoDB, Kubernetes, GitHub              | Business Logic             |
| MongoDB        | Backend API                              | Data Persistence           |
| Kubernetes     | Backend API, Containers                  | Container Orchestration    |
| GitHub API     | Backend API                              | Repository Management      |
| Metrics Server | Kubernetes, Backend API                  | Resource Monitoring        |

## Security Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Security Layers                          │
├─────────────────────────────────────────────────────────────┤
│ Layer 1: HTTPS/TLS                                          │
│   • Encrypted transport                                     │
│   • Certificate validation                                  │
├─────────────────────────────────────────────────────────────┤
│ Layer 2: Authentication                                     │
│   • JWT tokens with expiration                              │
│   • BCrypt password hashing (cost 10)                       │
│   • Token refresh mechanism                                 │
├─────────────────────────────────────────────────────────────┤
│ Layer 3: Authorization                                      │
│   • Role-based access control (RBAC)                        │
│   • Resource ownership validation                           │
│   • Middleware protection                                   │
├─────────────────────────────────────────────────────────────┤
│ Layer 4: Data Protection                                    │
│   • Encrypted credentials at rest                           │
│   • Secure token storage                                    │
│   • Input sanitization                                      │
├─────────────────────────────────────────────────────────────┤
│ Layer 5: Network Security                                   │
│   • CORS configuration                                      │
│   • API rate limiting (ready)                               │
│   • Request validation                                      │
└─────────────────────────────────────────────────────────────┘
```

## Scalability Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                 Horizontal Scaling                          │
│                                                             │
│  Frontend (Stateless)                                       │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                 │
│  │ Instance │  │ Instance │  │ Instance │ ...              │
│  │    1     │  │    2     │  │    3     │                 │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘                 │
│       └─────────────┴─────────────┘                        │
│                     │                                       │
│              Load Balancer                                  │
│                     │                                       │
│  Backend (Stateless)                                        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                 │
│  │ Instance │  │ Instance │  │ Instance │ ...              │
│  │    1     │  │    2     │  │    3     │                 │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘                 │
│       └─────────────┴─────────────┘                        │
│                     │                                       │
│  MongoDB (Replica Set)                                      │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐                    │
│  │ Primary │  │Secondary│  │Secondary│                     │
│  └─────────┘  └─────────┘  └─────────┘                    │
│                                                             │
│  Kubernetes Clusters                                        │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐                    │
│  │  Node 1 │  │  Node 2 │  │  Node 3 │ ...                │
│  └─────────┘  └─────────┘  └─────────┘                    │
└─────────────────────────────────────────────────────────────┘
```

This architecture supports:
- **Thousands of nodes** across multiple locations
- **Tens of thousands of deployments** simultaneously
- **Millions of API requests** per day
- **Real-time metrics** for all resources
- **High availability** with no single point of failure

