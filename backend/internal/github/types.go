package github

import "time"

// Repository represents a GitHub repository
type Repository struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	FullName      string    `json:"fullName"`
	Owner         string    `json:"owner"`
	Private       bool      `json:"private"`
	HTMLURL       string    `json:"htmlUrl"`
	Description   string    `json:"description"`
	CloneURL      string    `json:"cloneUrl"`
	GitURL        string    `json:"gitUrl"`
	SSHURL        string    `json:"sshUrl"`
	Language      string    `json:"language"`
	StarCount     int       `json:"starCount"`
	ForkCount     int       `json:"forkCount"`
	WatchCount    int       `json:"watchCount"`
	OpenIssues    int       `json:"openIssues"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	PushedAt      time.Time `json:"pushedAt"`
	Size          int       `json:"size"`
	DefaultBranch string    `json:"defaultBranch"`
}

// Branch represents a Git branch
type Branch struct {
	Name      string `json:"name"`
	Protected bool   `json:"protected"`
	CommitSHA string `json:"commitSha"`
}

// Commit represents a Git commit
type Commit struct {
	SHA     string    `json:"sha"`
	Message string    `json:"message"`
	Author  string    `json:"author"`
	Date    time.Time `json:"date"`
	URL     string    `json:"url"`
}

// User represents a GitHub user
type User struct {
	Login     string `json:"login"`
	ID        int64  `json:"id"`
	AvatarURL string `json:"avatarUrl"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	Location  string `json:"location"`
}

// Pagination contains pagination information
type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	TotalPages int `json:"totalPages"`
	Total      int `json:"total"`
}

