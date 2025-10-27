import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { FiGithub, FiStar, FiGitBranch, FiPackage, FiSearch, FiSettings } from 'react-icons/fi'
import { githubAPI } from '../services/api'
import { useNavigate } from 'react-router-dom'
import toast from 'react-hot-toast'

function RepositoryCard({ repo, onDeploy }) {
  return (
    <div className="card hover:shadow-lg transition-all">
      <div className="flex items-start justify-between mb-4">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-2">
            <FiGithub className="w-5 h-5 text-gray-600 dark:text-gray-400" />
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white truncate">
              {repo.name}
            </h3>
          </div>
          <p className="text-sm text-gray-600 dark:text-gray-400 mb-3 line-clamp-2">
            {repo.description || 'No description provided'}
          </p>
        </div>
      </div>

      <div className="flex items-center gap-4 mb-4 text-sm text-gray-600 dark:text-gray-400">
        {repo.language && (
          <div className="flex items-center gap-1">
            <div className="w-3 h-3 rounded-full bg-primary-500"></div>
            <span>{repo.language}</span>
          </div>
        )}
        <div className="flex items-center gap-1">
          <FiStar className="w-4 h-4" />
          <span>{repo.starCount}</span>
        </div>
        <div className="flex items-center gap-1">
          <FiGitBranch className="w-4 h-4" />
          <span>{repo.forkCount}</span>
        </div>
        {repo.private && (
          <span className="badge badge-warning text-xs">Private</span>
        )}
      </div>

      <div className="flex items-center gap-2 pt-4 border-t border-gray-200 dark:border-gray-700">
        <button
          onClick={() => onDeploy(repo)}
          className="btn btn-primary flex-1 text-sm"
        >
          <FiPackage className="w-4 h-4" />
          Deploy
        </button>
        <a
          href={repo.htmlUrl}
          target="_blank"
          rel="noopener noreferrer"
          className="btn btn-secondary text-sm"
        >
          View on GitHub
        </a>
      </div>
    </div>
  )
}

export default function Repositories() {
  const navigate = useNavigate()
  const [searchQuery, setSearchQuery] = useState('')
  const [page, setPage] = useState(1)

  const { data, isLoading, error } = useQuery({
    queryKey: ['repositories', page],
    queryFn: () => githubAPI.getRepositories(page, 30).then(res => res.data),
    retry: 1,
  })

  const handleDeploy = (repo) => {
    // Store selected repository in session storage
    sessionStorage.setItem('selectedRepo', JSON.stringify(repo))
    navigate('/deployments/new')
  }

  const handleSaveToken = async () => {
    const token = prompt('Enter your GitHub Personal Access Token:')
    if (token) {
      try {
        await githubAPI.saveToken(token)
        toast.success('GitHub token saved successfully!')
        window.location.reload()
      } catch (error) {
        toast.error('Failed to save token')
      }
    }
  }

  if (error) {
    return (
      <div className="space-y-6 animate-fade-in">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Repositories</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Browse your GitHub repositories
          </p>
        </div>

        <div className="card text-center py-12">
          <FiGithub className="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">
            GitHub Token Required
          </h3>
          <p className="text-gray-600 dark:text-gray-400 mb-6 max-w-md mx-auto">
            Please configure your GitHub Personal Access Token to browse and deploy repositories.
          </p>
          <button onClick={handleSaveToken} className="btn btn-primary">
            <FiSettings className="w-4 h-4" />
            Configure GitHub Token
          </button>
        </div>
      </div>
    )
  }

  const repositories = data?.repositories || []

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Repositories</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Browse and deploy your GitHub repositories
        </p>
      </div>

      {/* Search and Filters */}
      <div className="card">
        <div className="flex flex-col sm:flex-row gap-4">
          <div className="flex-1 relative">
            <FiSearch className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
            <input
              type="text"
              placeholder="Search repositories..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="input pl-10 w-full"
            />
          </div>
          <button onClick={handleSaveToken} className="btn btn-secondary whitespace-nowrap">
            <FiSettings className="w-4 h-4" />
            Configure Token
          </button>
        </div>
      </div>

      {/* Loading State */}
      {isLoading && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {[...Array(6)].map((_, i) => (
            <div key={i} className="card animate-pulse">
              <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-3/4 mb-4"></div>
              <div className="h-3 bg-gray-200 dark:bg-gray-700 rounded w-full mb-2"></div>
              <div className="h-3 bg-gray-200 dark:bg-gray-700 rounded w-2/3"></div>
            </div>
          ))}
        </div>
      )}

      {/* Repositories Grid */}
      {!isLoading && repositories.length > 0 && (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {repositories
              .filter(repo => 
                repo.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                repo.description?.toLowerCase().includes(searchQuery.toLowerCase())
              )
              .map((repo) => (
                <RepositoryCard key={repo.id} repo={repo} onDeploy={handleDeploy} />
              ))}
          </div>

          {/* Pagination */}
          {data?.pagination && (
            <div className="flex justify-center gap-2">
              <button
                onClick={() => setPage(p => Math.max(1, p - 1))}
                disabled={page === 1}
                className="btn btn-secondary disabled:opacity-50"
              >
                Previous
              </button>
              <span className="flex items-center px-4 text-gray-700 dark:text-gray-300">
                Page {page}
              </span>
              <button
                onClick={() => setPage(p => p + 1)}
                disabled={!data?.pagination?.totalPages || page >= data.pagination.totalPages}
                className="btn btn-secondary disabled:opacity-50"
              >
                Next
              </button>
            </div>
          )}
        </>
      )}

      {/* Empty State */}
      {!isLoading && repositories.length === 0 && (
        <div className="card text-center py-12">
          <FiGithub className="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">
            No repositories found
          </h3>
          <p className="text-gray-600 dark:text-gray-400">
            {searchQuery ? 'Try adjusting your search' : 'Create a repository on GitHub to get started'}
          </p>
        </div>
      )}
    </div>
  )
}

