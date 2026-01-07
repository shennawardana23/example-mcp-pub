package entity

import "time"

// Service represents a service in the catalog (database entity)
type Service struct {
	ID          int64
	Name        string
	Description string
	Owner       string
	Type        string
	Status      string
	Version     string
	Repository  string
	DocsURL     string
	Tags        []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CreatedBy   string
	UpdatedBy   string
}

// User represents a user (database entity)
type User struct {
	ID           int64
	Username     string
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	Role         string
	Active       bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Team represents a team (database entity)
type Team struct {
	ID          int64
	Name        string
	Description string
	Lead        string
	Members     []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Metric represents a metric data point (database entity)
type Metric struct {
	ID        int64
	ServiceID int64
	Name      string
	Value     float64
	Unit      string
	Timestamp time.Time
	Labels    map[string]string
}
