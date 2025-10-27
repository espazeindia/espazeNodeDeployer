import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { FiActivity, FiCpu, FiHardDrive, FiPackage, FiServer } from 'react-icons/fi'
import { metricsAPI, k8sAPI } from '../services/api'
import { LineChart, Line, BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts'

function PodMetricCard({ pod }) {
  const statusColors = {
    Running: 'badge-success',
    Pending: 'badge-warning',
    Failed: 'badge-danger',
    Succeeded: 'badge-info',
  }

  return (
    <div className="card">
      <div className="flex items-start justify-between mb-3">
        <div className="flex-1">
          <h4 className="font-semibold text-gray-900 dark:text-white">{pod.name}</h4>
          <p className="text-sm text-gray-500 dark:text-gray-400">{pod.namespace}</p>
        </div>
        <span className={`badge ${statusColors[pod.status] || 'badge-info'}`}>
          {pod.status}
        </span>
      </div>

      <div className="grid grid-cols-3 gap-3 text-center">
        <div>
          <FiCpu className="w-5 h-5 text-primary-600 mx-auto mb-1" />
          <p className="text-xs text-gray-500 dark:text-gray-400">CPU</p>
          <p className="text-sm font-semibold text-gray-900 dark:text-white">{pod.cpuUsage}</p>
        </div>
        <div>
          <FiHardDrive className="w-5 h-5 text-green-600 mx-auto mb-1" />
          <p className="text-xs text-gray-500 dark:text-gray-400">Memory</p>
          <p className="text-sm font-semibold text-gray-900 dark:text-white">{pod.memoryUsage}</p>
        </div>
        <div>
          <FiActivity className="w-5 h-5 text-blue-600 mx-auto mb-1" />
          <p className="text-xs text-gray-500 dark:text-gray-400">Restarts</p>
          <p className="text-sm font-semibold text-gray-900 dark:text-white">{pod.restartCount}</p>
        </div>
      </div>
    </div>
  )
}

export default function Observability() {
  const [selectedNamespace, setSelectedNamespace] = useState('espaze-node-deployer-apps')

  const { data: clusterMetrics } = useQuery({
    queryKey: ['clusterMetrics'],
    queryFn: () => metricsAPI.getClusterMetrics().then(res => res.data),
    refetchInterval: 10000,
  })

  const { data: podMetrics } = useQuery({
    queryKey: ['podMetrics', selectedNamespace],
    queryFn: () => metricsAPI.getPodMetrics(selectedNamespace).then(res => res.data),
    refetchInterval: 10000,
  })

  const { data: namespaces } = useQuery({
    queryKey: ['namespaces'],
    queryFn: () => k8sAPI.getNamespaces().then(res => res.data),
  })

  const { data: events } = useQuery({
    queryKey: ['events', selectedNamespace],
    queryFn: () => k8sAPI.getEvents(selectedNamespace).then(res => res.data),
  })

  // Sample data for charts (in production, this would come from time-series metrics)
  const cpuData = [
    { time: '12:00', usage: 45 },
    { time: '12:05', usage: 52 },
    { time: '12:10', usage: 48 },
    { time: '12:15', usage: 65 },
    { time: '12:20', usage: 58 },
    { time: '12:25', usage: 72 },
  ]

  const memoryData = [
    { time: '12:00', usage: 62 },
    { time: '12:05', usage: 65 },
    { time: '12:10', usage: 68 },
    { time: '12:15', usage: 70 },
    { time: '12:20', usage: 67 },
    { time: '12:25', usage: 71 },
  ]

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Observability</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Monitor cluster and pod metrics
          </p>
        </div>
        <select
          value={selectedNamespace}
          onChange={(e) => setSelectedNamespace(e.target.value)}
          className="input w-64"
        >
          <option value="">All Namespaces</option>
          {namespaces?.map((ns) => (
            <option key={ns} value={ns}>{ns}</option>
          ))}
        </select>
      </div>

      {/* Cluster Metrics */}
      {clusterMetrics && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <div className="card">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Total Nodes</p>
                <p className="text-3xl font-bold text-gray-900 dark:text-white mt-2">
                  {clusterMetrics.totalNodes}
                </p>
              </div>
              <FiServer className="w-12 h-12 text-primary-600" />
            </div>
          </div>

          <div className="card">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Running Pods</p>
                <p className="text-3xl font-bold text-green-600 mt-2">
                  {clusterMetrics.runningPods}
                </p>
                <p className="text-xs text-gray-500 mt-1">of {clusterMetrics.totalPods} total</p>
              </div>
              <FiPackage className="w-12 h-12 text-green-600" />
            </div>
          </div>

          <div className="card">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600 dark:text-gray-400">CPU Usage</p>
                <p className="text-3xl font-bold text-blue-600 mt-2">
                  {clusterMetrics.cpuUsagePercent?.toFixed(1) || 0}%
                </p>
                <p className="text-xs text-gray-500 mt-1">{clusterMetrics.totalCpu}</p>
              </div>
              <FiCpu className="w-12 h-12 text-blue-600" />
            </div>
          </div>

          <div className="card">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600 dark:text-gray-400">Memory Usage</p>
                <p className="text-3xl font-bold text-purple-600 mt-2">
                  {clusterMetrics.memUsagePercent?.toFixed(1) || 0}%
                </p>
                <p className="text-xs text-gray-500 mt-1">{clusterMetrics.totalMemory}</p>
              </div>
              <FiHardDrive className="w-12 h-12 text-purple-600" />
            </div>
          </div>
        </div>
      )}

      {/* Charts */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="card">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
            CPU Usage Over Time
          </h3>
          <ResponsiveContainer width="100%" height={250}>
            <LineChart data={cpuData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="time" />
              <YAxis />
              <Tooltip />
              <Legend />
              <Line type="monotone" dataKey="usage" stroke="#8b5cf6" strokeWidth={2} name="CPU %" />
            </LineChart>
          </ResponsiveContainer>
        </div>

        <div className="card">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
            Memory Usage Over Time
          </h3>
          <ResponsiveContainer width="100%" height={250}>
            <LineChart data={memoryData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="time" />
              <YAxis />
              <Tooltip />
              <Legend />
              <Line type="monotone" dataKey="usage" stroke="#10b981" strokeWidth={2} name="Memory %" />
            </LineChart>
          </ResponsiveContainer>
        </div>
      </div>

      {/* Pod Metrics */}
      {podMetrics && podMetrics.length > 0 && (
        <div>
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
            Pod Metrics
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {podMetrics.map((pod, index) => (
              <PodMetricCard key={index} pod={pod} />
            ))}
          </div>
        </div>
      )}

      {/* Recent Events */}
      {events && events.items && events.items.length > 0 && (
        <div className="card">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
            Recent Events
          </h3>
          <div className="space-y-3">
            {events.items.slice(0, 10).map((event, index) => (
              <div
                key={index}
                className="flex items-start gap-3 p-3 bg-gray-50 dark:bg-gray-700 rounded-lg"
              >
                <div className={`w-2 h-2 rounded-full mt-2 ${
                  event.type === 'Normal' ? 'bg-green-500' : 'bg-yellow-500'
                }`} />
                <div className="flex-1">
                  <p className="text-sm font-medium text-gray-900 dark:text-white">
                    {event.reason}
                  </p>
                  <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">
                    {event.message}
                  </p>
                  <p className="text-xs text-gray-500 dark:text-gray-500 mt-1">
                    {event.metadata?.name} • {new Date(event.lastTimestamp).toLocaleString()}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Empty State */}
      {(!podMetrics || podMetrics.length === 0) && (
        <div className="card text-center py-12">
          <FiActivity className="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">
            No metrics available
          </h3>
          <p className="text-gray-600 dark:text-gray-400">
            Deploy an application to see metrics
          </p>
        </div>
      )}
    </div>
  )
}

