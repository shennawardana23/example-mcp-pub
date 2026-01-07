package service

import (
	"context"
	"fmt"
	"time"

	"github.com/shennawardana23/example-mcp-pub/internal/app/converter"
	"github.com/shennawardana23/example-mcp-pub/internal/app/entity"
	"github.com/shennawardana23/example-mcp-pub/internal/app/model"
	"github.com/shennawardana23/example-mcp-pub/internal/app/repository"
	"github.com/redis/go-redis/v9"
)

// ServiceService defines the interface for service business logic
type ServiceService interface {
	Create(ctx context.Context, req *model.ServiceRequest, username string) (*model.ServiceResponse, error)
	GetByID(ctx context.Context, id int64) (*model.ServiceResponse, error)
	List(ctx context.Context, req *model.ServiceListRequest) (*model.PaginatedResponse, error)
	Update(ctx context.Context, id int64, req *model.ServiceRequest, username string) (*model.ServiceResponse, error)
	Delete(ctx context.Context, id int64) error
}

// serviceService implements ServiceService
type serviceService struct {
	repo      repository.ServiceRepository
	converter *converter.ServiceConverter
	redis     *redis.Client
}

// NewServiceService creates a new service service
func NewServiceService(repo repository.ServiceRepository, redis *redis.Client) ServiceService {
	return &serviceService{
		repo:      repo,
		converter: converter.NewServiceConverter(),
		redis:     redis,
	}
}

// Create creates a new service
func (s *serviceService) Create(ctx context.Context, req *model.ServiceRequest, username string) (*model.ServiceResponse, error) {
	service := s.converter.ToEntity(req)
	service.CreatedBy = username
	service.UpdatedBy = username

	if err := s.repo.Create(ctx, service); err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}

	// Invalidate list cache
	s.invalidateListCache(ctx)

	return s.converter.ToResponse(service), nil
}

// GetByID retrieves a service by ID with caching
func (s *serviceService) GetByID(ctx context.Context, id int64) (*model.ServiceResponse, error) {
	cacheKey := fmt.Sprintf("service:%d", id)

	// Try cache first
	if s.redis != nil {
		var cached model.ServiceResponse
		err := s.redis.Get(ctx, cacheKey).Scan(&cached)
		if err == nil {
			return &cached, nil
		}
	}

	// Get from database
	service, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	response := s.converter.ToResponse(service)

	// Cache for 1 hour
	if s.redis != nil {
		s.redis.Set(ctx, cacheKey, response, time.Hour)
	}

	return response, nil
}

// List retrieves services with filtering and pagination
func (s *serviceService) List(ctx context.Context, req *model.ServiceListRequest) (*model.PaginatedResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}

	offset := (req.Page - 1) * req.Limit

	filters := repository.ServiceFilters{
		Search: req.Search,
		Type:   req.Type,
		Status: req.Status,
		Owner:  req.Owner,
		Tag:    req.Tag,
		Offset: offset,
		Limit:  req.Limit,
	}

	services, totalCount, err := s.repo.List(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	responses := s.converter.ToResponseList(services)
	totalPages := int((totalCount + int64(req.Limit) - 1) / int64(req.Limit))

	return &model.PaginatedResponse{
		Data:       responses,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

// Update updates an existing service
func (s *serviceService) Update(ctx context.Context, id int64, req *model.ServiceRequest, username string) (*model.ServiceResponse, error) {
	// Get existing service
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	// Update fields
	existing.Name = req.Name
	existing.Description = req.Description
	existing.Owner = req.Owner
	existing.Type = req.Type
	existing.Status = req.Status
	existing.Version = req.Version
	existing.Repository = req.Repository
	existing.DocsURL = req.DocsURL
	existing.Tags = req.Tags
	existing.UpdatedBy = username

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("failed to update service: %w", err)
	}

	// Invalidate caches
	s.invalidateCache(ctx, id)
	s.invalidateListCache(ctx)

	return s.converter.ToResponse(existing), nil
}

// Delete deletes a service
func (s *serviceService) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}

	// Invalidate caches
	s.invalidateCache(ctx, id)
	s.invalidateListCache(ctx)

	return nil
}

// invalidateCache invalidates service cache
func (s *serviceService) invalidateCache(ctx context.Context, id int64) {
	if s.redis != nil {
		cacheKey := fmt.Sprintf("service:%d", id)
		s.redis.Del(ctx, cacheKey)
	}
}

// invalidateListCache invalidates list cache
func (s *serviceService) invalidateListCache(ctx context.Context) {
	if s.redis != nil {
		// Could implement a more sophisticated cache invalidation strategy
		keys, _ := s.redis.Keys(ctx, "services:list:*").Result()
		if len(keys) > 0 {
			s.redis.Del(ctx, keys...)
		}
	}
}
