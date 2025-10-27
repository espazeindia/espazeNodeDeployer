import { useQuery } from '@tanstack/react-query'
import { FiPackage, FiServer, FiActivity, FiTrendingUp } from 'react-icons/fi'
import { deploymentAPI, nodeAPI, k8sAPI, metricsAPI } from '../services/api'
import { Link } from 'react-router-dom'

function StatCard({ title, value, icon: Icon, trend, color, link }) {
  return (
    <Link to={link} className="card hover:shadow-lg transition-shadow">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm font-medium text-gray-600 dark:text-gray-400">{title}</p>
          <p className="text-3xl font-bold text-gray-900 dark:text-white mt-2">{value}</p>
          {trend && (
            <div className="flex items-center gap-1 mt-2 text-sm text-green-600">
              <FiTrendingUp className="w-4 h-4" />
              <span>{trend}</span>
            </div>
          )}
        </div>
        <div className={`w-16 h-16 rounded-2xl ${color} bg-opacity-10 flex items-center justify-center`}>
          <Icon className={`w-8 h-8 ${color.replace('bg-', 'text-')}`} />
        </div>
      </div>
    </Link>
  )
}

function DeploymentCard({ deployment }) {
  const statusColors = {
    running: 'bg-green-100 text-green-800',
    pending: 'bg-yellow-100 text-yellow-800',
    failed: 'bg-red-100 text-red-800',
    deploying: 'bg-blue-100 text-blue-800',
  }

  return (
    <Link 
      to={`/deployments/${deployment.id}`}
      className="card hover:shadow-lg transition-all"
    >
      <div className="flex items-start justify-between mb-4">
        <div>
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
            {deployment.name}
          </h3>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
            {deployment.githubRepo?.fullName}
          </p>
        </div>
        <span className={`badge ${statusColors[deployment.status] || 'badge-info'}`}>
          {deployment.status}
        </span>
      </div>
      
      <div className="grid grid-cols-3 gap-4 text-center mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
        <div>
          <p className="text-2xl font-bold text-gray-900 dark:text-white">
            {deployment.metrics?.activePods || 0}
          </p>
          <p className="text-xs text-gray-500 dark:text-gray-400">Pods</p>
        </div>
        <div>
          <p className="text-2xl font-bold text-gray-900 dark:text-white">
            {deployment.configuration?.replicas || 0}
          </p>
          <p className="text-xs text-gray-500 dark:text-gray-400">Replicas</p>
        </div>
        <div>
          <p className="text-2xl font-bold text-gray-900 dark:text-white">
            {deployment.metrics?.cpuUsage?.toFixed(1) || 0}%
          </p>
          <p className="text-xs text-gray-500 dark:text-gray-400">CPU</p>
        </div>
      </div>
    </Link>
  )
}

export default function Dashboard() {
  const { data: deployments } = useQuery({
    queryKey: ['deployments'],
    queryFn: () => deploymentAPI.getAll().then(res => res.data),
  })

  const { data: nodes } = useQuery({
    queryKey: ['nodes'],
    queryFn: () => nodeAPI.getAll().then(res => res.data),
  })

  const { data: clusterInfo } = useQuery({
    queryKey: ['clusterInfo'],
    queryFn: () => k8sAPI.getClusterInfo().then(res => res.data),
  })

  const { data: clusterMetrics } = useQuery({
    queryKey: ['clusterMetrics'],
    queryFn: () => metricsAPI.getClusterMetrics().then(res => res.data),
  })

  const activeDeployments = deployments?.filter(d => d.status === 'running')?.length || 0
  const totalNodes = nodes?.length || 0
  const onlineNodes = nodes?.filter(n => n.status === 'online')?.length || 0
  const totalPods = clusterMetrics?.totalPods || 0

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Dashboard</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Overview of your Kubernetes deployments
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <StatCard
          title="Active Deployments"
          value={activeDeployments}
          icon={FiPackage}
          trend="+12% this week"
          color="bg-primary-600"
          link="/deployments"
        />
        <StatCard
          title="Total Nodes"
          value={totalNodes}
          icon={FiServer}
          trend={`${onlineNodes} online`}
          color="bg-green-600"
          link="/nodes"
        />
        <StatCard
          title="Running Pods"
          value={totalPods}
          icon={FiActivity}
          trend={`${clusterMetrics?.runningPods || 0} healthy`}
          color="bg-blue-600"
          link="/observability"
        />
        <StatCard
          title="Cluster Health"
          value="98%"
          icon={FiTrendingUp}
          trend="Excellent"
          color="bg-purple-600"
          link="/observability"
        />
      </div>

      {/* Cluster Info */}
      {clusterInfo && (
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
            Cluster Information
          </h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div>
              <p className="text-sm text-gray-600 dark:text-gray-400">Version</p>
              <p className="text-lg font-semibold text-gray-900 dark:text-white mt-1">
                {clusterInfo.version}
              </p>
            </div>
            <div>
              <p className="text-sm text-gray-600 dark:text-gray-400">Nodes</p>
              <p className="text-lg font-semibold text-gray-900 dark:text-white mt-1">
                {clusterInfo.nodesCount}
              </p>
            </div>
            <div>
              <p className="text-sm text-gray-600 dark:text-gray-400">Namespaces</p>
              <p className="text-lg font-semibold text-gray-900 dark:text-white mt-1">
                {clusterInfo.namespacesCount}
              </p>
            </div>
            <div>
              <p className="text-sm text-gray-600 dark:text-gray-400">Pods</p>
              <p className="text-lg font-semibold text-gray-900 dark:text-white mt-1">
                {clusterInfo.podsCount}
              </p>
            </div>
          </div>
        </div>
      )}

      {/* Recent Deployments */}
      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white">
            Recent Deployments
          </h2>
          <Link to="/deployments" className="text-primary-600 hover:text-primary-700 text-sm font-medium">
            View All →
          </Link>
        </div>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {deployments?.slice(0, 6).map((deployment) => (
            <DeploymentCard key={deployment.id} deployment={deployment} />
          ))}
          
          {!deployments?.length && (
            <div className="col-span-full card text-center py-12">
              <FiPackage className="w-16 h-16 text-gray-400 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">
                No deployments yet
              </h3>
              <p className="text-gray-600 dark:text-gray-400 mb-4">
                Get started by deploying your first application
              </p>
              <Link to="/repositories" className="btn btn-primary">
                Browse Repositories
              </Link>
            </div>
          )}
        </div>
      </div>

      {/* Nodes Overview */}
      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white">
            Registered Nodes
          </h2>
          <Link to="/nodes" className="text-primary-600 hover:text-primary-700 text-sm font-medium">
            View All →
          </Link>
        </div>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {nodes?.slice(0, 3).map((node) => (
            <Link key={node.id} to={`/nodes/${node.id}`} className="card hover:shadow-lg transition-all">
              <div className="flex items-start justify-between mb-3">
                <div>
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                    {node.nodeName}
                  </h3>
                  <p className="text-sm text-gray-500 dark:text-gray-400">
                    {node.publicIp}
                  </p>
                </div>
                <span className={`badge ${node.status === 'online' ? 'badge-success' : 'badge-danger'}`}>
                  {node.status}
                </span>
              </div>
              
              <div className="space-y-2">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600 dark:text-gray-400">CPU Usage</span>
                  <span className="font-medium text-gray-900 dark:text-white">
                    {node.resources?.cpuUsage?.toFixed(1) || 0}%
                  </span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600 dark:text-gray-400">Memory</span>
                  <span className="font-medium text-gray-900 dark:text-white">
                    {node.resources?.memoryUsage?.toFixed(1) || 0}%
                  </span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600 dark:text-gray-400">Pods Running</span>
                  <span className="font-medium text-gray-900 dark:text-white">
                    {node.resources?.podsRunning || 0}
                  </span>
                </div>
              </div>
            </Link>
          ))}
        </div>
      </div>
    </div>
  )
}

