package repository

import (
	"context"

	"github.com/shennawardana23/example-mcp-pub/internal/app/entity"
)

// ServiceRepository defines the interface for service data access
type ServiceRepository interface {
	Create(ctx context.Context, service *entity.Service) error
	GetByID(ctx context.Context, id int64) (*entity.Service, error)
	List(ctx context.Context, filters ServiceFilters) ([]entity.Service, int64, error)
	Update(ctx context.Context, service *entity.Service) error
	Delete(ctx context.Context, id int64) error
}

// ServiceFilters represents service filtering options
type ServiceFilters struct {
	Search string
	Type   string
	Status string
	Owner  string
	Tag    string
	Offset int
	Limit  int
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int64) error
}

// TeamRepository defines the interface for team data access
type TeamRepository interface {
	Create(ctx context.Context, team *entity.Team) error
	GetByID(ctx context.Context, id int64) (*entity.Team, error)
	List(ctx context.Context, offset, limit int) ([]entity.Team, int64, error)
	Update(ctx context.Context, team *entity.Team) error
	Delete(ctx context.Context, id int64) error
}

// MetricRepository defines the interface for metric data access
type MetricRepository interface {
	Create(ctx context.Context, metric *entity.Metric) error
	GetByServiceID(ctx context.Context, serviceID int64, limit int) ([]entity.Metric, error)
	GetLatest(ctx context.Context, serviceID int64, name string) (*entity.Metric, error)
}
