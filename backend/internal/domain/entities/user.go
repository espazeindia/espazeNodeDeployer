package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a system user
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email" json:"email"`
	Username  string             `bson:"username" json:"username"`
	Password  string             `bson:"password" json:"-"` // Never expose in JSON
	FullName  string             `bson:"full_name" json:"fullName"`
	Role      UserRole           `bson:"role" json:"role"`
	Status    UserStatus         `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
	LastLogin time.Time          `bson:"last_login" json:"lastLogin"`
}

// UserRole defines user roles
type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

// UserStatus defines user status
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

// LoginRequest is used for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest is used for user registration
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"fullName" binding:"required"`
}

// LoginResponse is returned after successful login
type LoginResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refreshToken"`
	User         *User     `json:"user"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

// GitHubToken stores GitHub personal access token
type GitHubToken struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"userId"`
	Token     string             `bson:"token" json:"-"` // Encrypted, never expose
	Scope     string             `bson:"scope" json:"scope"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
	ExpiresAt *time.Time         `bson:"expires_at" json:"expiresAt"`
}

// GitHubTokenRequest is used to save GitHub token
type GitHubTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

