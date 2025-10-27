import { useState } from 'react'
import { FiGithub, FiSave, FiCheck } from 'react-icons/fi'
import { githubAPI } from '../services/api'
import toast from 'react-hot-toast'

export default function Settings() {
  const [githubToken, setGithubToken] = useState('')
  const [saving, setSaving] = useState(false)

  const handleSaveGitHubToken = async (e) => {
    e.preventDefault()
    setSaving(true)

    try {
      await githubAPI.saveToken(githubToken)
      toast.success('GitHub token saved successfully!')
      setGithubToken('')
    } catch (error) {
      toast.error('Failed to save GitHub token')
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="space-y-6 animate-fade-in max-w-4xl">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Settings</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Configure your application settings
        </p>
      </div>

      {/* GitHub Integration */}
      <div className="card">
        <div className="flex items-center gap-3 mb-4">
          <div className="w-12 h-12 rounded-lg bg-gray-900 flex items-center justify-center">
            <FiGithub className="w-6 h-6 text-white" />
          </div>
          <div>
            <h2 className="text-xl font-semibold text-gray-900 dark:text-white">
              GitHub Integration
            </h2>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              Configure your GitHub Personal Access Token
            </p>
          </div>
        </div>

        <form onSubmit={handleSaveGitHubToken} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              GitHub Personal Access Token
            </label>
            <input
              type="password"
              value={githubToken}
              onChange={(e) => setGithubToken(e.target.value)}
              placeholder="ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
              className="input"
              required
            />
            <p className="mt-2 text-sm text-gray-600 dark:text-gray-400">
              Create a token at{' '}
              <a
                href="https://github.com/settings/tokens"
                target="_blank"
                rel="noopener noreferrer"
                className="text-primary-600 hover:text-primary-700"
              >
                GitHub Settings → Developer settings → Personal access tokens
              </a>
            </p>
            <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
              Required scope: <code className="px-1 py-0.5 bg-gray-100 dark:bg-gray-800 rounded">repo</code>
            </p>
          </div>

          <button
            type="submit"
            disabled={saving || !githubToken}
            className="btn btn-primary"
          >
            {saving ? (
              <>
                <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                Saving...
              </>
            ) : (
              <>
                <FiSave className="w-4 h-4" />
                Save Token
              </>
            )}
          </button>
        </form>
      </div>

      {/* Default Deployment Configuration */}
      <div className="card">
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
          Default Deployment Configuration
        </h2>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Default Replicas
            </label>
            <input type="number" defaultValue="2" className="input" />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Default Memory Limit
            </label>
            <input type="text" defaultValue="512Mi" className="input" />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Default CPU Limit
            </label>
            <input type="text" defaultValue="500m" className="input" />
          </div>

          <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Default Namespace
              </label>
              <input type="text" defaultValue="espaze-node-deployer-apps" className="input" />
          </div>
        </div>

        <button className="btn btn-primary mt-6">
          <FiSave className="w-4 h-4" />
          Save Configuration
        </button>
      </div>

      {/* About */}
      <div className="card">
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
          About Espaze Node Deployer
        </h2>
        
        <div className="space-y-3">
          <div className="flex justify-between">
            <span className="text-gray-600 dark:text-gray-400">Version</span>
            <span className="font-semibold text-gray-900 dark:text-white">1.0.0</span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-600 dark:text-gray-400">Backend API</span>
            <span className="font-semibold text-gray-900 dark:text-white">v1</span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-600 dark:text-gray-400">Status</span>
            <span className="flex items-center gap-2 text-green-600">
              <FiCheck className="w-4 h-4" />
              Connected
            </span>
          </div>
        </div>
      </div>
    </div>
  )
}

