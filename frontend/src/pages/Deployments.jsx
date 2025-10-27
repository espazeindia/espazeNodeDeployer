import { useQuery } from '@tanstack/react-query'
import { FiPackage, FiPlus, FiActivity, FiClock, FiCheckCircle, FiXCircle } from 'react-icons/fi'
import { deploymentAPI } from '../services/api'
import { Link } from 'react-router-dom'
import { formatDistanceToNow } from 'date-fns'

export default function Deployments() {
  const { data: deployments, isLoading } = useQuery({
    queryKey: ['deployments'],
    queryFn: () => deploymentAPI.getAll().then(res => res.data),
    refetchInterval: 15000,
  })

  const { data: stats } = useQuery({
    queryKey: ['deploymentStats'],
    queryFn: () => deploymentAPI.getStats().then(res => res.data),
  })

  const statusConfig = {
    running: { icon: FiCheckCircle, color: 'text-green-600', bg: 'bg-green-100', label: 'Running' },
    pending: { icon: FiClock, color: 'text-yellow-600', bg: 'bg-yellow-100', label: 'Pending' },
    deploying: { icon: FiActivity, color: 'text-blue-600', bg: 'bg-blue-100', label: 'Deploying' },
    failed: { icon: FiXCircle, color: 'text-red-600', bg: 'bg-red-100', label: 'Failed' },
    stopped: { icon: FiXCircle, color: 'text-gray-600', bg: 'bg-gray-100', label: 'Stopped' },
  }

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Deployments</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Manage your application deployments
          </p>
        </div>
        <Link to="/repositories" className="btn btn-primary">
          <FiPlus className="w-4 h-4" />
          New Deployment
        </Link>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="card">
          <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Total</p>
          <p className="text-3xl font-bold text-gray-900 dark:text-white mt-2">
            {stats?.total || 0}
          </p>
        </div>
        <div className="card">
          <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Running</p>
          <p className="text-3xl font-bold text-green-600 mt-2">
            {stats?.running || 0}
          </p>
        </div>
        <div className="card">
          <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Deploying</p>
          <p className="text-3xl font-bold text-blue-600 mt-2">
            {stats?.deploying || 0}
          </p>
        </div>
        <div className="card">
          <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Failed</p>
          <p className="text-3xl font-bold text-red-600 mt-2">
            {stats?.failed || 0}
          </p>
        </div>
      </div>

      {/* Loading */}
      {isLoading && (
        <div className="space-y-4">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="card animate-pulse">
              <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-3/4 mb-4"></div>
              <div className="h-3 bg-gray-200 dark:bg-gray-700 rounded w-full"></div>
            </div>
          ))}
        </div>
      )}

      {/* Deployments List */}
      {!isLoading && deployments && deployments.length > 0 && (
        <div className="space-y-4">
          {deployments.map((deployment) => {
            const status = statusConfig[deployment.status] || statusConfig.pending
            const StatusIcon = status.icon

            return (
              <Link
                key={deployment.id}
                to={`/deployments/${deployment.id}`}
                className="card hover:shadow-lg transition-all"
              >
                <div className="flex items-start gap-4">
                  <div className={`w-12 h-12 rounded-lg ${status.bg} flex items-center justify-center flex-shrink-0`}>
                    <StatusIcon className={`w-6 h-6 ${status.color}`} />
                  </div>

                  <div className="flex-1 min-w-0">
                    <div className="flex items-start justify-between mb-2">
                      <div>
                        <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                          {deployment.name}
                        </h3>
                        <p className="text-sm text-gray-500 dark:text-gray-400">
                          {deployment.githubRepo?.fullName}
                        </p>
                      </div>
                      <span className={`badge ${status.bg} ${status.color}`}>
                        {status.label}
                      </span>
                    </div>

                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mt-4">
                      <div>
                        <p className="text-xs text-gray-500 dark:text-gray-400">Namespace</p>
                        <p className="text-sm font-medium text-gray-900 dark:text-white">
                          {deployment.namespace}
                        </p>
                      </div>
                      <div>
                        <p className="text-xs text-gray-500 dark:text-gray-400">Replicas</p>
                        <p className="text-sm font-medium text-gray-900 dark:text-white">
                          {deployment.metrics?.activePods || 0} / {deployment.configuration?.replicas || 0}
                        </p>
                      </div>
                      <div>
                        <p className="text-xs text-gray-500 dark:text-gray-400">Context Path</p>
                        <p className="text-sm font-medium text-gray-900 dark:text-white font-mono">
                          {deployment.contextPath}
                        </p>
                      </div>
                      <div>
                        <p className="text-xs text-gray-500 dark:text-gray-400">Created</p>
                        <p className="text-sm font-medium text-gray-900 dark:text-white">
                          {deployment.createdAt ? formatDistanceToNow(new Date(deployment.createdAt), { addSuffix: true }) : 'N/A'}
                        </p>
                      </div>
                    </div>

                    {deployment.kubernetesInfo?.url && (
                      <div className="mt-3 pt-3 border-t border-gray-200 dark:border-gray-700">
                        <p className="text-xs text-gray-500 dark:text-gray-400 mb-1">URL</p>
                        <a
                          href={deployment.kubernetesInfo.url}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-sm text-primary-600 hover:text-primary-700 font-mono"
                          onClick={(e) => e.stopPropagation()}
                        >
                          {deployment.kubernetesInfo.url}
                        </a>
                      </div>
                    )}
                  </div>
                </div>
              </Link>
            )
          })}
        </div>
      )}

      {/* Empty State */}
      {!isLoading && (!deployments || deployments.length === 0) && (
        <div className="card text-center py-12">
          <FiPackage className="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">
            No deployments yet
          </h3>
          <p className="text-gray-600 dark:text-gray-400 mb-6">
            Get started by deploying your first application from GitHub
          </p>
          <Link to="/repositories" className="btn btn-primary">
            <FiPlus className="w-4 h-4" />
            Browse Repositories
          </Link>
        </div>
      )}
    </div>
  )
}

