package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

type Client struct {
	clientID     string
	clientSecret string
}

func NewClient(clientID, clientSecret string) *Client {
	return &Client{
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

// CreateAuthenticatedClient creates a GitHub client with user token
func (c *Client) CreateAuthenticatedClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

// GetUserRepositories returns all repositories for the authenticated user
func (c *Client) GetUserRepositories(ctx context.Context, token string, page, perPage int) ([]*Repository, *Pagination, error) {
	client := c.CreateAuthenticatedClient(ctx, token)

	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: perPage,
		},
		Sort:      "updated",
		Direction: "desc",
	}

	repos, resp, err := client.Repositories.List(ctx, "", opts)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list repositories: %w", err)
	}

	repositories := make([]*Repository, len(repos))
	for i, repo := range repos {
		repositories[i] = convertRepository(repo)
	}

	pagination := &Pagination{
		Page:       page,
		PerPage:    perPage,
		TotalPages: resp.LastPage,
		Total:      resp.NextPage,
	}

	return repositories, pagination, nil
}

// GetRepository returns a single repository
func (c *Client) GetRepository(ctx context.Context, token, owner, repo string) (*Repository, error) {
	client := c.CreateAuthenticatedClient(ctx, token)

	repository, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}

	return convertRepository(repository), nil
}

// GetBranches returns all branches for a repository
func (c *Client) GetBranches(ctx context.Context, token, owner, repo string) ([]*Branch, error) {
	client := c.CreateAuthenticatedClient(ctx, token)

	branches, _, err := client.Repositories.ListBranches(ctx, owner, repo, &github.BranchListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	result := make([]*Branch, len(branches))
	for i, branch := range branches {
		result[i] = &Branch{
			Name:      branch.GetName(),
			Protected: branch.GetProtected(),
			CommitSHA: branch.GetCommit().GetSHA(),
		}
	}

	return result, nil
}

// GetDefaultBranch returns the default branch for a repository
func (c *Client) GetDefaultBranch(ctx context.Context, token, owner, repo string) (string, error) {
	client := c.CreateAuthenticatedClient(ctx, token)

	repository, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return "", fmt.Errorf("failed to get repository: %w", err)
	}

	return repository.GetDefaultBranch(), nil
}

// GetCommit returns information about a specific commit
func (c *Client) GetCommit(ctx context.Context, token, owner, repo, sha string) (*Commit, error) {
	client := c.CreateAuthenticatedClient(ctx, token)

	commit, _, err := client.Repositories.GetCommit(ctx, owner, repo, sha, &github.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get commit: %w", err)
	}

	return &Commit{
		SHA:     commit.GetSHA(),
		Message: commit.GetCommit().GetMessage(),
		Author:  commit.GetCommit().GetAuthor().GetName(),
		Date:    commit.GetCommit().GetAuthor().GetDate().Time,
		URL:     commit.GetHTMLURL(),
	}, nil
}

// SearchRepositories searches for repositories
func (c *Client) SearchRepositories(ctx context.Context, token, query string, page, perPage int) ([]*Repository, error) {
	client := c.CreateAuthenticatedClient(ctx, token)

	opts := &github.SearchOptions{
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: perPage,
		},
	}

	result, _, err := client.Search.Repositories(ctx, query, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to search repositories: %w", err)
	}

	repositories := make([]*Repository, len(result.Repositories))
	for i, repo := range result.Repositories {
		repositories[i] = convertRepository(&repo)
	}

	return repositories, nil
}

// GetAuthenticatedUser returns information about the authenticated user
func (c *Client) GetAuthenticatedUser(ctx context.Context, token string) (*User, error) {
	client := c.CreateAuthenticatedClient(ctx, token)

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &User{
		Login:     user.GetLogin(),
		ID:        user.GetID(),
		AvatarURL: user.GetAvatarURL(),
		Name:      user.GetName(),
		Email:     user.GetEmail(),
		Bio:       user.GetBio(),
		Location:  user.GetLocation(),
	}, nil
}

// CheckDockerfile checks if a Dockerfile exists in the repository
func (c *Client) CheckDockerfile(ctx context.Context, token, owner, repo, branch string) (bool, string, error) {
	client := c.CreateAuthenticatedClient(ctx, token)

	// Check for Dockerfile in root
	_, _, resp, err := client.Repositories.GetContents(ctx, owner, repo, "Dockerfile", &github.RepositoryContentGetOptions{
		Ref: branch,
	})

	if resp.StatusCode == 200 {
		return true, "Dockerfile", nil
	}

	// Check for dockerfile (lowercase)
	_, _, resp, err = client.Repositories.GetContents(ctx, owner, repo, "dockerfile", &github.RepositoryContentGetOptions{
		Ref: branch,
	})

	if resp.StatusCode == 200 {
		return true, "dockerfile", nil
	}

	// Check for Dockerfile in .docker directory
	_, _, resp, err = client.Repositories.GetContents(ctx, owner, repo, ".docker/Dockerfile", &github.RepositoryContentGetOptions{
		Ref: branch,
	})

	if resp.StatusCode == 200 {
		return true, ".docker/Dockerfile", nil
	}

	return false, "", err
}

// GetFileContent retrieves the content of a file from the repository
func (c *Client) GetFileContent(ctx context.Context, token, owner, repo, path, branch string) (string, error) {
	client := c.CreateAuthenticatedClient(ctx, token)

	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{
		Ref: branch,
	})

	if err != nil {
		return "", fmt.Errorf("failed to get file content: %w", err)
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return "", fmt.Errorf("failed to decode file content: %w", err)
	}

	return content, nil
}

// Helper function to convert GitHub repository to our format
func convertRepository(repo *github.Repository) *Repository {
	return &Repository{
		ID:          repo.GetID(),
		Name:        repo.GetName(),
		FullName:    repo.GetFullName(),
		Owner:       repo.GetOwner().GetLogin(),
		Private:     repo.GetPrivate(),
		HTMLURL:     repo.GetHTMLURL(),
		Description: repo.GetDescription(),
		CloneURL:    repo.GetCloneURL(),
		GitURL:      repo.GetGitURL(),
		SSHURL:      repo.GetSSHURL(),
		Language:    repo.GetLanguage(),
		StarCount:   repo.GetStargazersCount(),
		ForkCount:   repo.GetForksCount(),
		WatchCount:  repo.GetWatchersCount(),
		OpenIssues:  repo.GetOpenIssuesCount(),
		CreatedAt:   repo.GetCreatedAt().Time,
		UpdatedAt:   repo.GetUpdatedAt().Time,
		PushedAt:    repo.GetPushedAt().Time,
		Size:        repo.GetSize(),
		DefaultBranch: repo.GetDefaultBranch(),
	}
}

