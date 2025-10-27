import axios from 'axios'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const authStorage = localStorage.getItem('auth-storage')
    if (authStorage) {
      const { state } = JSON.parse(authStorage)
      if (state.token) {
        config.headers.Authorization = `Bearer ${state.token}`
      }
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('auth-storage')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// Auth API
export const authAPI = {
  register: (data) => api.post('/auth/register', data),
  login: (data) => api.post('/auth/login', data),
  validate: () => api.get('/auth/validate'),
}

// Node API
export const nodeAPI = {
  getAll: (params) => api.get('/nodes', { params }),
  getById: (id) => api.get(`/nodes/${id}`),
  register: (data) => api.post('/nodes/register', data),
  update: (id, data) => api.put(`/nodes/${id}`, data),
  delete: (id) => api.delete(`/nodes/${id}`),
  heartbeat: (id) => api.post(`/nodes/${id}/heartbeat`),
  getStats: () => api.get('/nodes/stats'),
  getCurrent: () => api.get('/nodes/current'),
}

// Deployment API
export const deploymentAPI = {
  getAll: () => api.get('/deployments'),
  getById: (id) => api.get(`/deployments/${id}`),
  getByNode: (nodeId) => api.get(`/deployments/node/${nodeId}`),
  create: (data, nodeId, githubToken) => api.post(`/deployments?nodeId=${nodeId}`, data, {
    headers: {
      'X-GitHub-Token': githubToken,
    },
  }),
  update: (id, data) => api.put(`/deployments/${id}`, data),
  delete: (id) => api.delete(`/deployments/${id}`),
  restart: (id) => api.post(`/deployments/${id}/restart`),
  scale: (id, replicas) => api.post(`/deployments/${id}/scale`, { replicas }),
  getStats: (nodeId) => api.get('/deployments/stats', { params: { nodeId } }),
}

// GitHub API
export const githubAPI = {
  saveToken: (token) => api.post('/github/token', { token }),
  getUser: () => api.get('/github/user'),
  getRepositories: (page = 1, perPage = 30) => api.get('/github/repos', {
    params: { page, perPage },
  }),
  getRepository: (owner, repo) => api.get(`/github/repos/${owner}/${repo}`),
  getBranches: (owner, repo) => api.get(`/github/repos/${owner}/${repo}/branches`),
  search: (query, page = 1, perPage = 30) => api.get('/github/search', {
    params: { q: query, page, perPage },
  }),
}

// Kubernetes API
export const k8sAPI = {
  getClusterInfo: () => api.get('/k8s/cluster/info'),
  getNamespaces: () => api.get('/k8s/namespaces'),
  getPods: (namespace) => api.get('/k8s/pods', { params: { namespace } }),
  getPod: (namespace, name) => api.get(`/k8s/pods/${namespace}/${name}`),
  getPodLogs: (namespace, name, tail = 100) => api.get(`/k8s/pods/${namespace}/${name}/logs`, {
    params: { tail },
  }),
  getServices: (namespace) => api.get('/k8s/services', { params: { namespace } }),
  getNodes: () => api.get('/k8s/nodes'),
  getEvents: (namespace) => api.get('/k8s/events', { params: { namespace } }),
}

// Metrics API
export const metricsAPI = {
  getPodMetrics: (namespace) => api.get('/metrics/pods', { params: { namespace } }),
  getClusterMetrics: () => api.get('/metrics/cluster'),
  getDeploymentMetrics: (namespace, name) => api.get(`/metrics/deployments/${namespace}/${name}`),
}

export default api

