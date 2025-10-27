package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/espaze/espazeNodeDeployer/internal/domain/entities"
	"github.com/espaze/espazeNodeDeployer/internal/repository"
	"github.com/espaze/espazeNodeDeployer/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Register(ctx context.Context, req *entities.RegisterRequest) (*entities.User, error)
	Login(ctx context.Context, req *entities.LoginRequest) (*entities.LoginResponse, error)
	ValidateToken(token string) (*auth.Claims, error)
}

type authUseCase struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthUseCase(userRepo repository.UserRepository, jwtSecret string) AuthUseCase {
	return &authUseCase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (uc *authUseCase) Register(ctx context.Context, req *entities.RegisterRequest) (*entities.User, error) {
	// Check if user already exists
	existingUser, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Check if username is taken
	existingUser, err = uc.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already taken")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &entities.User{
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashedPassword),
		FullName: req.FullName,
		Role:     entities.UserRoleUser,
		Status:   entities.UserStatusActive,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *authUseCase) Login(ctx context.Context, req *entities.LoginRequest) (*entities.LoginResponse, error) {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if user.Status != entities.UserStatusActive {
		return nil, errors.New("user account is not active")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Update last login
	if err := uc.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail login
	}

	// Generate JWT token
	expiresAt := time.Now().Add(24 * time.Hour)
	token, err := auth.GenerateToken(user.ID.Hex(), user.Email, string(user.Role), uc.jwtSecret, expiresAt)
	if err != nil {
		return nil, err
	}

	// Generate refresh token (valid for 7 days)
	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)
	refreshToken, err := auth.GenerateToken(user.ID.Hex(), user.Email, string(user.Role), uc.jwtSecret, refreshExpiresAt)
	if err != nil {
		return nil, err
	}

	return &entities.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
		ExpiresAt:    expiresAt,
	}, nil
}

func (uc *authUseCase) ValidateToken(token string) (*auth.Claims, error) {
	return auth.ValidateToken(token, uc.jwtSecret)
}

