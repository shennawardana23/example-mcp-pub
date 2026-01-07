package model

import "time"

// ServiceRequest represents a request to create/update a service
type ServiceRequest struct {
	Name        string   `json:"name" binding:"required,min=3,max=100"`
	Description string   `json:"description" binding:"required,min=10,max=500"`
	Owner       string   `json:"owner" binding:"required"`
	Type        string   `json:"type" binding:"required,oneof=api library microservice frontend"`
	Status      string   `json:"status" binding:"required,oneof=active deprecated planning"`
	Version     string   `json:"version" binding:"omitempty"`
	Repository  string   `json:"repository" binding:"omitempty,url"`
	DocsURL     string   `json:"docs_url" binding:"omitempty,url"`
	Tags        []string `json:"tags" binding:"omitempty"`
}

// ServiceResponse represents a service response
type ServiceResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Owner       string    `json:"owner"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Version     string    `json:"version,omitempty"`
	Repository  string    `json:"repository,omitempty"`
	DocsURL     string    `json:"docs_url,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ServiceListRequest represents a request to list services
type ServiceListRequest struct {
	Search string `form:"search"`
	Type   string `form:"type"`
	Status string `form:"status"`
	Owner  string `form:"owner"`
	Tag    string `form:"tag"`
	Page   int    `form:"page,default=1"`
	Limit  int    `form:"limit,default=20"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// UserResponse represents a user response
type UserResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// TeamResponse represents a team response
type TeamResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Lead        string    `json:"lead"`
	Members     []string  `json:"members"`
	CreatedAt   time.Time `json:"created_at"`
}

// MetricResponse represents a metric response
type MetricResponse struct {
	ID        int64              `json:"id"`
	ServiceID int64              `json:"service_id"`
	Name      string             `json:"name"`
	Value     float64            `json:"value"`
	Unit      string             `json:"unit"`
	Timestamp time.Time          `json:"timestamp"`
	Labels    map[string]string  `json:"labels,omitempty"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalCount int64       `json:"total_count"`
	TotalPages int         `json:"total_pages"`
}
