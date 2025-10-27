import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { FiGithub, FiPackage } from 'react-icons/fi'
import { deploymentAPI, nodeAPI } from '../services/api'
import toast from 'react-hot-toast'

export default function CreateDeployment() {
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const [selectedRepo, setSelectedRepo] = useState(null)
  const [nodes, setNodes] = useState([])
  
  const [formData, setFormData] = useState({
    name: '',
    contextPath: '',
    namespace: 'espaze-node-deployer-apps',
    nodeId: '',
    githubToken: '',
    branch: 'main',
    replicas: 2,
    memoryRequest: '256Mi',
    memoryLimit: '512Mi',
    cpuRequest: '250m',
    cpuLimit: '500m',
    containerPort: 8080,
    servicePort: 80,
  })

  useEffect(() => {
    // Get selected repository from session storage
    const repoData = sessionStorage.getItem('selectedRepo')
    if (repoData) {
      const repo = JSON.parse(repoData)
      setSelectedRepo(repo)
      setFormData(prev => ({
        ...prev,
        name: repo.name.toLowerCase().replace(/[^a-z0-9-]/g, '-'),
        contextPath: `/${repo.name.toLowerCase()}`,
      }))
    }

    // Load nodes
    nodeAPI.getAll().then(res => {
      setNodes(res.data || [])
      if (res.data?.length > 0) {
        setFormData(prev => ({ ...prev, nodeId: res.data[0].id }))
      }
    })
  }, [])

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)

    try {
      const deploymentData = {
        name: formData.name,
        contextPath: formData.contextPath,
        namespace: formData.namespace,
        githubRepo: {
          owner: selectedRepo.owner,
          name: selectedRepo.name,
          fullName: selectedRepo.fullName,
          branch: formData.branch,
        },
        configuration: {
          replicas: parseInt(formData.replicas),
          containerPort: parseInt(formData.containerPort),
          servicePort: parseInt(formData.servicePort),
          memoryRequest: formData.memoryRequest,
          memoryLimit: formData.memoryLimit,
          cpuRequest: formData.cpuRequest,
          cpuLimit: formData.cpuLimit,
          imagePullPolicy: 'IfNotPresent',
          restartPolicy: 'Always',
          environmentVars: {},
          healthCheck: {
            enabled: true,
            path: '/health',
            port: parseInt(formData.containerPort),
            initialDelaySeconds: 30,
            periodSeconds: 10,
            timeoutSeconds: 5,
            successThreshold: 1,
            failureThreshold: 3,
          },
          buildConfig: {
            dockerfile: 'Dockerfile',
            buildContext: '.',
            imageName: `${selectedRepo.owner}/${selectedRepo.name}`,
            imageTag: 'latest',
          },
        },
      }

      await deploymentAPI.create(deploymentData, formData.nodeId, formData.githubToken)
      toast.success('Deployment created successfully!')
      sessionStorage.removeItem('selectedRepo')
      navigate('/deployments')
    } catch (error) {
      toast.error(error.response?.data?.error || 'Failed to create deployment')
    } finally {
      setLoading(false)
    }
  }

  if (!selectedRepo) {
    return (
      <div className="card text-center py-12">
        <FiGithub className="w-16 h-16 text-gray-400 mx-auto mb-4" />
        <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">
          No repository selected
        </h3>
        <p className="text-gray-600 dark:text-gray-400 mb-6">
          Please select a repository from the repositories page
        </p>
        <button onClick={() => navigate('/repositories')} className="btn btn-primary">
          Browse Repositories
        </button>
      </div>
    )
  }

  return (
    <div className="space-y-6 animate-fade-in max-w-4xl">
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Create Deployment</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Deploy {selectedRepo.fullName} to Kubernetes
        </p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Basic Settings */}
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
            Basic Settings
          </h2>
          
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Deployment Name
              </label>
              <input
                type="text"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                className="input"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Context Path (URL)
              </label>
              <input
                type="text"
                value={formData.contextPath}
                onChange={(e) => setFormData({ ...formData, contextPath: e.target.value })}
                className="input"
                placeholder="/my-app"
                required
              />
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Branch
                </label>
                <input
                  type="text"
                  value={formData.branch}
                  onChange={(e) => setFormData({ ...formData, branch: e.target.value })}
                  className="input"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Target Node
                </label>
                <select
                  value={formData.nodeId}
                  onChange={(e) => setFormData({ ...formData, nodeId: e.target.value })}
                  className="input"
                  required
                >
                  {nodes.map(node => (
                    <option key={node.id} value={node.id}>{node.nodeName}</option>
                  ))}
                </select>
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                GitHub Token
              </label>
              <input
                type="password"
                value={formData.githubToken}
                onChange={(e) => setFormData({ ...formData, githubToken: e.target.value })}
                className="input"
                placeholder="ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
                required
              />
            </div>
          </div>
        </div>

        {/* Resource Configuration */}
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
            Resources
          </h2>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Replicas
              </label>
              <input
                type="number"
                value={formData.replicas}
                onChange={(e) => setFormData({ ...formData, replicas: e.target.value })}
                className="input"
                min="1"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Container Port
              </label>
              <input
                type="number"
                value={formData.containerPort}
                onChange={(e) => setFormData({ ...formData, containerPort: e.target.value })}
                className="input"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Memory Limit
              </label>
              <input
                type="text"
                value={formData.memoryLimit}
                onChange={(e) => setFormData({ ...formData, memoryLimit: e.target.value })}
                className="input"
                placeholder="512Mi"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                CPU Limit
              </label>
              <input
                type="text"
                value={formData.cpuLimit}
                onChange={(e) => setFormData({ ...formData, cpuLimit: e.target.value })}
                className="input"
                placeholder="500m"
                required
              />
            </div>
          </div>
        </div>

        {/* Submit */}
        <div className="flex gap-4">
          <button
            type="submit"
            disabled={loading}
            className="btn btn-primary flex-1"
          >
            {loading ? (
              <>
                <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin" />
                Deploying...
              </>
            ) : (
              <>
                <FiPackage className="w-4 h-4" />
                Deploy Application
              </>
            )}
          </button>
          <button
            type="button"
            onClick={() => navigate('/repositories')}
            className="btn btn-secondary"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  )
}

