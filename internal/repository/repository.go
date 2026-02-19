package repository

import (
	"context"
	"database/sql"

	"cadastral-service/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CreateQuery is create a new request
func (r *Repository) CreateQuery(ctx context.Context, query *models.Query) error {
	queryStr := `
		INSERT INTO queries (id, cadastral_number, latitude, longitude, status, user_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, queryStr,
		query.ID,
		query.CadastralNumber,
		query.Latitude,
		query.Longitude,
		query.Status,
		query.UserID,
		query.CreatedAt,
	)

	return err
}

// UpdateQuery is update a status of request
func (r *Repository) UpdateQuery(ctx context.Context, id string, status string, result *bool) error {
	queryStr := `
		UPDATE queries 
		SET status = $1, result = $2, completed_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, queryStr, status, result, id)
	return err
}

// GetQueries is return list of request
func (r *Repository) GetQueries(ctx context.Context, userID, page, limit string) ([]models.Query, error) {
	var queries []models.Query
	var queryStr string
	var args []interface{}

	if userID != "" {
		queryStr = `
			SELECT id, cadastral_number, latitude, longitude, status, result, user_id, created_at, completed_at
			FROM queries
			WHERE user_id = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{userID, limit, (page - 1) * limit}
	} else {
		queryStr = `
			SELECT id, cadastral_number, latitude, longitude, status, result, user_id, created_at, completed_at
			FROM queries
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`
		args = []interface{}{limit, (page - 1) * limit}
	}

	rows, err := r.db.QueryContext(ctx, queryStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var q models.Query
		err := rows.Scan(
			&q.ID,
			&q.CadastralNumber,
			&q.Latitude,
			&q.Longitude,
			&q.Status,
			&q.Result,
			&q.UserID,
			&q.CreatedAt,
			&q.CompletedAt,
		)
		if err != nil {
			return nil, err
		}
		queries = append(queries, q)
	}

	return queries, nil
}

// GetQueriesByCadastral is return request by cadastral number
func (r *Repository) GetQueriesByCadastral(ctx context.Context, cadastralNumber, userID string) ([]models.Query, error) {
	var queries []models.Query
	var queryStr string
	var args []interface{}

	if userID != "" {
		queryStr = `
			SELECT id, cadastral_number, latitude, longitude, status, result, user_id, created_at, completed_at
			 FROM queries
			WHERE cadastral_number = $1 AND user_id = $2
			ORDER BY created_at DESC
		`
		args = []interface{}{cadastralNumber, userID}
	} else {
		queryStr = `
			SELECT id, cadastral_number, latitude, longitude, status, result, user_id, created_at, completed_at
			FROM queries
			WHERE cadastral_number = $1
			ORDER BY created_at DESC
		`
		args = []interface{}{cadastralNumber}
	}

	rows, err := r.db.QueryContext(ctx, queryStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var q models.Query
		err := rows.Scan(
			&q.ID,
			&q.CadastralNumber,
			&q.Latitude,
			&q.Longitude,
			&q.Status,
			&q.Result,
			&q.UserID,
			&q.CreatedAt,
			&q.CompletedAt,
		)
		if err != nil {
			return nil, err
		}
		queries = append(queries, q)
	}

	return queries, nil
}

// CreateUser is create a new user
func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	queryStr := `
		INSERT INTO users (id, username, password_hash, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, queryStr,
		user.ID,
		user.Username,
		user.PasswordHash,
		user.CreatedAt,
	)

	return err
}

// GetUserByUsername is return user by name
func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	queryStr := `
		SELECT id, username, password_hash, created_at
		FROM users
		WHERE username = $1
	`

	var user models.User
	err := r.db.QueryRowContext(ctx, queryStr, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
