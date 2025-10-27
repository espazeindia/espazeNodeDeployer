import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from './store/authStore'
import Layout from './components/Layout'
import Login from './pages/Auth/Login'
import Register from './pages/Auth/Register'
import Dashboard from './pages/Dashboard'
import Deployments from './pages/Deployments'
import DeploymentDetails from './pages/DeploymentDetails'
import CreateDeployment from './pages/CreateDeployment'
import Repositories from './pages/Repositories'
import Nodes from './pages/Nodes'
import NodeDetails from './pages/NodeDetails'
import Observability from './pages/Observability'
import Settings from './pages/Settings'

function PrivateRoute({ children }) {
  const { token } = useAuthStore()
  return token ? children : <Navigate to="/login" />
}

function PublicRoute({ children }) {
  const { token } = useAuthStore()
  return token ? <Navigate to="/" /> : children
}

function App() {
  return (
    <Routes>
      <Route path="/login" element={<PublicRoute><Login /></PublicRoute>} />
      <Route path="/register" element={<PublicRoute><Register /></PublicRoute>} />
      
      <Route path="/" element={<PrivateRoute><Layout /></PrivateRoute>}>
        <Route index element={<Dashboard />} />
        <Route path="deployments" element={<Deployments />} />
        <Route path="deployments/:id" element={<DeploymentDetails />} />
        <Route path="deployments/new" element={<CreateDeployment />} />
        <Route path="repositories" element={<Repositories />} />
        <Route path="nodes" element={<Nodes />} />
        <Route path="nodes/:id" element={<NodeDetails />} />
        <Route path="observability" element={<Observability />} />
        <Route path="settings" element={<Settings />} />
      </Route>
    </Routes>
  )
}

export default App

