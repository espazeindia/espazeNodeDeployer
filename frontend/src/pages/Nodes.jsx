import { useQuery } from '@tanstack/react-query'
import { FiServer, FiMapPin, FiCpu, FiHardDrive, FiActivity, FiPlus } from 'react-icons/fi'
import { nodeAPI } from '../services/api'
import { Link } from 'react-router-dom'
import { formatDistanceToNow } from 'date-fns'

function NodeCard({ node }) {
  const statusColors = {
    online: 'badge-success',
    offline: 'badge-danger',
    maintenance: 'badge-warning',
    error: 'badge-danger',
  }

  return (
    <Link to={`/nodes/${node.id}`} className="card hover:shadow-lg transition-all">
      <div className="flex items-start justify-between mb-4">
        <div>
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-1">
            {node.nodeName}
          </h3>
          <div className="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
            <FiMapPin className="w-4 h-4" />
            <span>{node.location?.city || 'Unknown'}, {node.location?.country || 'Unknown'}</span>
          </div>
        </div>
        <span className={`badge ${statusColors[node.status]}`}>
          {node.status}
        </span>
      </div>

      <div className="space-y-3">
        <div className="flex items-center justify-between text-sm">
          <span className="text-gray-600 dark:text-gray-400">Public IP</span>
          <span className="font-mono text-gray-900 dark:text-white">{node.publicIp}</span>
        </div>
        
        <div className="flex items-center justify-between text-sm">
          <span className="text-gray-600 dark:text-gray-400">MAC Address</span>
          <span className="font-mono text-xs text-gray-900 dark:text-white">
            {node.macAddress}
          </span>
        </div>

        <div className="pt-3 border-t border-gray-200 dark:border-gray-700">
          <div className="grid grid-cols-3 gap-4 text-center">
            <div>
              <FiCpu className="w-5 h-5 text-primary-600 mx-auto mb-1" />
              <p className="text-xs text-gray-600 dark:text-gray-400">CPU</p>
              <p className="text-sm font-semibold text-gray-900 dark:text-white">
                {node.resources?.cpuUsage?.toFixed(1) || 0}%
              </p>
            </div>
            <div>
              <FiHardDrive className="w-5 h-5 text-green-600 mx-auto mb-1" />
              <p className="text-xs text-gray-600 dark:text-gray-400">Memory</p>
              <p className="text-sm font-semibold text-gray-900 dark:text-white">
                {node.resources?.memoryUsage?.toFixed(1) || 0}%
              </p>
            </div>
            <div>
              <FiActivity className="w-5 h-5 text-blue-600 mx-auto mb-1" />
              <p className="text-xs text-gray-600 dark:text-gray-400">Pods</p>
              <p className="text-sm font-semibold text-gray-900 dark:text-white">
                {node.resources?.podsRunning || 0}
              </p>
            </div>
          </div>
        </div>

        <div className="pt-3 border-t border-gray-200 dark:border-gray-700">
          <div className="flex justify-between text-xs text-gray-500 dark:text-gray-400">
            <span>Cluster: {node.clusterInfo?.clusterName || 'N/A'}</span>
            <span>
              Last seen: {node.lastSeenAt ? formatDistanceToNow(new Date(node.lastSeenAt), { addSuffix: true }) : 'Never'}
            </span>
          </div>
        </div>
      </div>
    </Link>
  )
}

export default function Nodes() {
  const { data: nodes, isLoading } = useQuery({
    queryKey: ['nodes'],
    queryFn: () => nodeAPI.getAll().then(res => res.data),
    refetchInterval: 30000, // Refresh every 30 seconds
  })

  const { data: stats } = useQuery({
    queryKey: ['nodeStats'],
    queryFn: () => nodeAPI.getStats().then(res => res.data),
  })

  const onlineNodes = nodes?.filter(n => n.status === 'online')?.length || 0
  const totalNodes = nodes?.length || 0

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Nodes</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Manage your Kubernetes nodes
          </p>
        </div>
        <Link to="/nodes/register" className="btn btn-primary">
          <FiPlus className="w-4 h-4" />
          Register Node
        </Link>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Total Nodes</p>
              <p className="text-3xl font-bold text-gray-900 dark:text-white mt-2">
                {totalNodes}
              </p>
            </div>
            <div className="w-12 h-12 rounded-lg bg-primary-100 dark:bg-primary-900 flex items-center justify-center">
              <FiServer className="w-6 h-6 text-primary-600" />
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Online</p>
              <p className="text-3xl font-bold text-green-600 mt-2">{onlineNodes}</p>
            </div>
            <div className="w-12 h-12 rounded-lg bg-green-100 dark:bg-green-900 flex items-center justify-center">
              <FiActivity className="w-6 h-6 text-green-600" />
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Offline</p>
              <p className="text-3xl font-bold text-red-600 mt-2">
                {stats?.offline || 0}
              </p>
            </div>
            <div className="w-12 h-12 rounded-lg bg-red-100 dark:bg-red-900 flex items-center justify-center">
              <FiServer className="w-6 h-6 text-red-600" />
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Maintenance</p>
              <p className="text-3xl font-bold text-yellow-600 mt-2">
                {stats?.maintenance || 0}
              </p>
            </div>
            <div className="w-12 h-12 rounded-lg bg-yellow-100 dark:bg-yellow-900 flex items-center justify-center">
              <FiServer className="w-6 h-6 text-yellow-600" />
            </div>
          </div>
        </div>
      </div>

      {/* Loading State */}
      {isLoading && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="card animate-pulse">
              <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-3/4 mb-4"></div>
              <div className="h-3 bg-gray-200 dark:bg-gray-700 rounded w-full mb-2"></div>
              <div className="h-3 bg-gray-200 dark:bg-gray-700 rounded w-2/3"></div>
            </div>
          ))}
        </div>
      )}

      {/* Nodes Grid */}
      {!isLoading && nodes && nodes.length > 0 && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {nodes.map((node) => (
            <NodeCard key={node.id} node={node} />
          ))}
        </div>
      )}

      {/* Empty State */}
      {!isLoading && (!nodes || nodes.length === 0) && (
        <div className="card text-center py-12">
          <FiServer className="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">
            No nodes registered
          </h3>
          <p className="text-gray-600 dark:text-gray-400 mb-6">
            Register your first node to start deploying applications
          </p>
          <Link to="/nodes/register" className="btn btn-primary">
            <FiPlus className="w-4 h-4" />
            Register This Node
          </Link>
        </div>
      )}
    </div>
  )
}

