package usecase

import (
	"context"

	"github.com/espaze/espazeNodeDeployer/internal/domain/entities"
	"github.com/espaze/espazeNodeDeployer/internal/github"
	"github.com/espaze/espazeNodeDeployer/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GitHubUseCase interface {
	SaveToken(ctx context.Context, userID primitive.ObjectID, token string) error
	GetRepositories(ctx context.Context, userID primitive.ObjectID, page, perPage int) ([]*github.Repository, *github.Pagination, error)
	GetRepository(ctx context.Context, userID primitive.ObjectID, owner, repo string) (*github.Repository, error)
	GetBranches(ctx context.Context, userID primitive.ObjectID, owner, repo string) ([]*github.Branch, error)
	SearchRepositories(ctx context.Context, userID primitive.ObjectID, query string, page, perPage int) ([]*github.Repository, error)
	GetUser(ctx context.Context, userID primitive.ObjectID) (*github.User, error)
}

type githubUseCase struct {
	githubClient    *github.Client
	githubTokenRepo repository.GitHubTokenRepository
}

func NewGitHubUseCase(
	githubClient *github.Client,
	githubTokenRepo repository.GitHubTokenRepository,
) GitHubUseCase {
	return &githubUseCase{
		githubClient:    githubClient,
		githubTokenRepo: githubTokenRepo,
	}
}

func (uc *githubUseCase) SaveToken(ctx context.Context, userID primitive.ObjectID, token string) error {
	// In production, encrypt the token before storing
	return uc.githubTokenRepo.Update(ctx, userID, token)
}

func (uc *githubUseCase) GetRepositories(ctx context.Context, userID primitive.ObjectID, page, perPage int) ([]*github.Repository, *github.Pagination, error) {
	token, err := uc.getToken(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	return uc.githubClient.GetUserRepositories(ctx, token, page, perPage)
}

func (uc *githubUseCase) GetRepository(ctx context.Context, userID primitive.ObjectID, owner, repo string) (*github.Repository, error) {
	token, err := uc.getToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	return uc.githubClient.GetRepository(ctx, token, owner, repo)
}

func (uc *githubUseCase) GetBranches(ctx context.Context, userID primitive.ObjectID, owner, repo string) ([]*github.Branch, error) {
	token, err := uc.getToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	return uc.githubClient.GetBranches(ctx, token, owner, repo)
}

func (uc *githubUseCase) SearchRepositories(ctx context.Context, userID primitive.ObjectID, query string, page, perPage int) ([]*github.Repository, error) {
	token, err := uc.getToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	return uc.githubClient.SearchRepositories(ctx, token, query, page, perPage)
}

func (uc *githubUseCase) GetUser(ctx context.Context, userID primitive.ObjectID) (*github.User, error) {
	token, err := uc.getToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	return uc.githubClient.GetAuthenticatedUser(ctx, token)
}

func (uc *githubUseCase) getToken(ctx context.Context, userID primitive.ObjectID) (string, error) {
	tokenEntity, err := uc.githubTokenRepo.GetByUserID(ctx, userID)
	if err != nil {
		return "", err
	}
	if tokenEntity == nil {
		return "", entities.ErrGitHubTokenNotFound
	}

	// In production, decrypt the token here
	return tokenEntity.Token, nil
}

// Error definitions
var (
	ErrGitHubTokenNotFound = entities.ErrGitHubTokenNotFound
)

