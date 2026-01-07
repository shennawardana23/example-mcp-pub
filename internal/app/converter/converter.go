package converter

import (
	"github.com/shennawardana23/example-mcp-pub/internal/app/entity"
	"github.com/shennawardana23/example-mcp-pub/internal/app/model"
)

// ServiceConverter handles conversions between Service entity and models
type ServiceConverter struct{}

// NewServiceConverter creates a new service converter
func NewServiceConverter() *ServiceConverter {
	return &ServiceConverter{}
}

// ToResponse converts a Service entity to ServiceResponse model
func (c *ServiceConverter) ToResponse(service *entity.Service) *model.ServiceResponse {
	if service == nil {
		return nil
	}

	return &model.ServiceResponse{
		ID:          service.ID,
		Name:        service.Name,
		Description: service.Description,
		Owner:       service.Owner,
		Type:        service.Type,
		Status:      service.Status,
		Version:     service.Version,
		Repository:  service.Repository,
		DocsURL:     service.DocsURL,
		Tags:        service.Tags,
		CreatedAt:   service.CreatedAt,
		UpdatedAt:   service.UpdatedAt,
	}
}

// ToResponseList converts a slice of Service entities to ServiceResponse models
func (c *ServiceConverter) ToResponseList(services []entity.Service) []model.ServiceResponse {
	responses := make([]model.ServiceResponse, 0, len(services))
	for _, service := range services {
		if resp := c.ToResponse(&service); resp != nil {
			responses = append(responses, *resp)
		}
	}
	return responses
}

// ToEntity converts a ServiceRequest model to Service entity
func (c *ServiceConverter) ToEntity(req *model.ServiceRequest) *entity.Service {
	if req == nil {
		return nil
	}

	return &entity.Service{
		Name:        req.Name,
		Description: req.Description,
		Owner:       req.Owner,
		Type:        req.Type,
		Status:      req.Status,
		Version:     req.Version,
		Repository:  req.Repository,
		DocsURL:     req.DocsURL,
		Tags:        req.Tags,
	}
}

// UserConverter handles conversions between User entity and models
type UserConverter struct{}

// NewUserConverter creates a new user converter
func NewUserConverter() *UserConverter {
	return &UserConverter{}
}

// ToResponse converts a User entity to UserResponse model
func (c *UserConverter) ToResponse(user *entity.User) *model.UserResponse {
	if user == nil {
		return nil
	}

	return &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}

// TeamConverter handles conversions between Team entity and models
type TeamConverter struct{}

// NewTeamConverter creates a new team converter
func NewTeamConverter() *TeamConverter {
	return &TeamConverter{}
}

// ToResponse converts a Team entity to TeamResponse model
func (c *TeamConverter) ToResponse(team *entity.Team) *model.TeamResponse {
	if team == nil {
		return nil
	}

	return &model.TeamResponse{
		ID:          team.ID,
		Name:        team.Name,
		Description: team.Description,
		Lead:        team.Lead,
		Members:     team.Members,
		CreatedAt:   team.CreatedAt,
	}
}

// MetricConverter handles conversions between Metric entity and models
type MetricConverter struct{}

// NewMetricConverter creates a new metric converter
func NewMetricConverter() *MetricConverter {
	return &MetricConverter{}
}

// ToResponse converts a Metric entity to MetricResponse model
func (c *MetricConverter) ToResponse(metric *entity.Metric) *model.MetricResponse {
	if metric == nil {
		return nil
	}

	return &model.MetricResponse{
		ID:        metric.ID,
		ServiceID: metric.ServiceID,
		Name:      metric.Name,
		Value:     metric.Value,
		Unit:      metric.Unit,
		Timestamp: metric.Timestamp,
		Labels:    metric.Labels,
	}
}
