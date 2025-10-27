import { useParams } from 'react-router-dom'
import { useQuery } from '@antml:react-query'
import { nodeAPI, deploymentAPI } from '../services/api'

export default function NodeDetails() {
  const { id } = useParams()

  const { data: node } = useQuery({
    queryKey: ['node', id],
    queryFn: () => nodeAPI.getById(id).then(res => res.data),
  })

  const { data: deployments } = useQuery({
    queryKey: ['nodeDeployments', id],
    queryFn: () => deploymentAPI.getByNode(id).then(res => res.data),
  })

  if (!node) {
    return <div className="animate-pulse">Loading...</div>
  }

  return (
    <div className="space-y-6 animate-fade-in">
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">{node.nodeName}</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">{node.publicIp}</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="card">
          <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400 mb-2">Status</h3>
          <span className={`badge ${node.status === 'online' ? 'badge-success' : 'badge-danger'}`}>
            {node.status}
          </span>
        </div>
        <div className="card">
          <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400 mb-2">Location</h3>
          <p className="text-gray-900 dark:text-white">
            {node.location?.city}, {node.location?.country}
          </p>
        </div>
        <div className="card">
          <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400 mb-2">Architecture</h3>
          <p className="text-gray-900 dark:text-white">{node.metadata?.architecture}</p>
        </div>
      </div>

      <div className="card">
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">Resources</h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div>
            <p className="text-sm text-gray-600 dark:text-gray-400">CPU Cores</p>
            <p className="text-lg font-semibold text-gray-900 dark:text-white">
              {node.resources?.cpuCores}
            </p>
          </div>
          <div>
            <p className="text-sm text-gray-600 dark:text-gray-400">Memory Total</p>
            <p className="text-lg font-semibold text-gray-900 dark:text-white">
              {(node.resources?.memoryTotal / 1024 / 1024 / 1024).toFixed(2)} GB
            </p>
          </div>
          <div>
            <p className="text-sm text-gray-600 dark:text-gray-400">Pods Running</p>
            <p className="text-lg font-semibold text-gray-900 dark:text-white">
              {node.resources?.podsRunning}
            </p>
          </div>
          <div>
            <p className="text-sm text-gray-600 dark:text-gray-400">Pods Capacity</p>
            <p className="text-lg font-semibold text-gray-900 dark:text-white">
              {node.resources?.podsCapacity}
            </p>
          </div>
        </div>
      </div>

      {deployments && deployments.length > 0 && (
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
            Deployments ({deployments.length})
          </h2>
          <div className="space-y-2">
            {deployments.map((deployment) => (
              <div key={deployment.id} className="p-3 bg-gray-50 dark:bg-gray-700 rounded-lg">
                <p className="font-medium text-gray-900 dark:text-white">{deployment.name}</p>
                <p className="text-sm text-gray-600 dark:text-gray-400">{deployment.githubRepo?.fullName}</p>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}

