import { useParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { deploymentAPI, k8sAPI } from '../services/api'
import { FiPackage, FiRefreshCw, FiTrash2, FiExternalLink, FiGithub } from 'react-icons/fi'

export default function DeploymentDetails() {
  const { id } = useParams()

  const { data: deployment } = useQuery({
    queryKey: ['deployment', id],
    queryFn: () => deploymentAPI.getById(id).then(res => res.data),
    refetchInterval: 10000,
  })

  if (!deployment) {
    return <div className="animate-pulse">Loading...</div>
  }

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">{deployment.name}</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            {deployment.githubRepo?.fullName}
          </p>
        </div>
        <div className="flex gap-2">
          <button className="btn btn-secondary">
            <FiRefreshCw className="w-4 h-4" />
            Restart
          </button>
          <button className="btn btn-danger">
            <FiTrash2 className="w-4 h-4" />
            Delete
          </button>
        </div>
      </div>

      {/* Status and Info */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="card">
          <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400 mb-2">Status</h3>
          <span className="badge badge-success">{deployment.status}</span>
        </div>
        <div className="card">
          <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400 mb-2">URL</h3>
          {deployment.kubernetesInfo?.url ? (
            <a
              href={deployment.kubernetesInfo.url}
              target="_blank"
              rel="noopener noreferrer"
              className="text-primary-600 hover:text-primary-700 flex items-center gap-1"
            >
              Open <FiExternalLink className="w-4 h-4" />
            </a>
          ) : (
            <span className="text-gray-500">Not available</span>
          )}
        </div>
        <div className="card">
          <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400 mb-2">Repository</h3>
          <a
            href={`https://github.com/${deployment.githubRepo?.fullName}`}
            target="_blank"
            rel="noopener noreferrer"
            className="text-primary-600 hover:text-primary-700 flex items-center gap-1"
          >
            <FiGithub className="w-4 h-4" /> View on GitHub
          </a>
        </div>
      </div>

      {/* Configuration */}
      <div className="card">
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">Configuration</h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div>
            <p className="text-sm text-gray-600 dark:text-gray-400">Replicas</p>
            <p className="text-lg font-semibold text-gray-900 dark:text-white">
              {deployment.configuration?.replicas}
            </p>
          </div>
          <div>
            <p className="text-sm text-gray-600 dark:text-gray-400">Memory Limit</p>
            <p className="text-lg font-semibold text-gray-900 dark:text-white">
              {deployment.configuration?.memoryLimit}
            </p>
          </div>
          <div>
            <p className="text-sm text-gray-600 dark:text-gray-400">CPU Limit</p>
            <p className="text-lg font-semibold text-gray-900 dark:text-white">
              {deployment.configuration?.cpuLimit}
            </p>
          </div>
          <div>
            <p className="text-sm text-gray-600 dark:text-gray-400">Namespace</p>
            <p className="text-lg font-semibold text-gray-900 dark:text-white">
              {deployment.namespace}
            </p>
          </div>
        </div>
      </div>

      {/* Metrics */}
      <div className="card">
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">Metrics</h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div>
            <p className="text-sm text-gray-600 dark:text-gray-400">Active Pods</p>
            <p className="text-lg font-semibold text-gray-900 dark:text-white">
              {deployment.metrics?.activePods || 0}
            </p>
          </div>
          <div>
            <p className="text-sm text-gray-600 dark:text-gray-400">Ready Pods</p>
            <p className="text-lg font-semibold text-gray-900 dark:text-white">
              {deployment.metrics?.readyPods || 0}
            </p>
          </div>
          <div>
            <p className="text-sm text-gray-600 dark:text-gray-400">CPU Usage</p>
            <p className="text-lg font-semibold text-gray-900 dark:text-white">
              {deployment.metrics?.cpuUsage?.toFixed(1) || 0}%
            </p>
          </div>
          <div>
            <p className="text-sm text-gray-600 dark:text-gray-400">Memory Usage</p>
            <p className="text-lg font-semibold text-gray-900 dark:text-white">
              {deployment.metrics?.memoryUsage?.toFixed(1) || 0}%
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}

