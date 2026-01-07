package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/shennawardana23/example-mcp-pub/internal/app/entity"
)

// PostgresServiceRepository implements ServiceRepository for PostgreSQL
type PostgresServiceRepository struct {
	db *sql.DB
}

// NewPostgresServiceRepository creates a new PostgreSQL service repository
func NewPostgresServiceRepository(db *sql.DB) *PostgresServiceRepository {
	return &PostgresServiceRepository{db: db}
}

// Create creates a new service
func (r *PostgresServiceRepository) Create(ctx context.Context, service *entity.Service) error {
	query := `
		INSERT INTO services (name, description, owner, type, status, version, repository, docs_url, tags, created_by, updated_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id`

	now := time.Now()
	service.CreatedAt = now
	service.UpdatedAt = now

	err := r.db.QueryRowContext(ctx, query,
		service.Name,
		service.Description,
		service.Owner,
		service.Type,
		service.Status,
		service.Version,
		service.Repository,
		service.DocsURL,
		pq.Array(service.Tags),
		service.CreatedBy,
		service.UpdatedBy,
		service.CreatedAt,
		service.UpdatedAt,
	).Scan(&service.ID)

	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	return nil
}

// GetByID retrieves a service by ID
func (r *PostgresServiceRepository) GetByID(ctx context.Context, id int64) (*entity.Service, error) {
	query := `
		SELECT id, name, description, owner, type, status, version, repository, docs_url, tags, created_by, updated_by, created_at, updated_at
		FROM services
		WHERE id = $1`

	service := &entity.Service{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&service.ID,
		&service.Name,
		&service.Description,
		&service.Owner,
		&service.Type,
		&service.Status,
		&service.Version,
		&service.Repository,
		&service.DocsURL,
		pq.Array(&service.Tags),
		&service.CreatedBy,
		&service.UpdatedBy,
		&service.CreatedAt,
		&service.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("service not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	return service, nil
}

// List retrieves services with filtering
func (r *PostgresServiceRepository) List(ctx context.Context, filters ServiceFilters) ([]entity.Service, int64, error) {
	// Build query with filters
	whereConditions := []string{}
	args := []interface{}{}
	argPos := 1

	if filters.Search != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d)", argPos, argPos))
		args = append(args, "%"+filters.Search+"%")
		argPos++
	}

	if filters.Type != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("type = $%d", argPos))
		args = append(args, filters.Type)
		argPos++
	}

	if filters.Status != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("status = $%d", argPos))
		args = append(args, filters.Status)
		argPos++
	}

	if filters.Owner != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("owner = $%d", argPos))
		args = append(args, filters.Owner)
		argPos++
	}

	if filters.Tag != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("$%d = ANY(tags)", argPos))
		args = append(args, filters.Tag)
		argPos++
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM services %s", whereClause)
	var totalCount int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count services: %w", err)
	}

	// Get services
	query := fmt.Sprintf(`
		SELECT id, name, description, owner, type, status, version, repository, docs_url, tags, created_by, updated_by, created_at, updated_at
		FROM services
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argPos, argPos+1)

	args = append(args, filters.Limit, filters.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list services: %w", err)
	}
	defer rows.Close()

	services := []entity.Service{}
	for rows.Next() {
		var service entity.Service
		err := rows.Scan(
			&service.ID,
			&service.Name,
			&service.Description,
			&service.Owner,
			&service.Type,
			&service.Status,
			&service.Version,
			&service.Repository,
			&service.DocsURL,
			pq.Array(&service.Tags),
			&service.CreatedBy,
			&service.UpdatedBy,
			&service.CreatedAt,
			&service.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan service: %w", err)
		}
		services = append(services, service)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating services: %w", err)
	}

	return services, totalCount, nil
}

// Update updates an existing service
func (r *PostgresServiceRepository) Update(ctx context.Context, service *entity.Service) error {
	query := `
		UPDATE services
		SET name = $1, description = $2, owner = $3, type = $4, status = $5, version = $6,
		    repository = $7, docs_url = $8, tags = $9, updated_by = $10, updated_at = $11
		WHERE id = $12`

	service.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		service.Name,
		service.Description,
		service.Owner,
		service.Type,
		service.Status,
		service.Version,
		service.Repository,
		service.DocsURL,
		pq.Array(service.Tags),
		service.UpdatedBy,
		service.UpdatedAt,
		service.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("service not found")
	}

	return nil
}

// Delete deletes a service
func (r *PostgresServiceRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM services WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("service not found")
	}

	return nil
}
