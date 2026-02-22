package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	"github.com/Lzrb0x/SmartSchedulingAPI/internal/domain"
)

type UserRepository struct {
	db *sqlx.DB
}

type UserWithTenant struct {
	User   domain.User
	Tenant domain.Tenant
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) DB() *sqlx.DB {
	return r.db
}

func (r *UserRepository) CreateTenant(ctx context.Context, exec sqlx.QueryerContext, name string) (*domain.Tenant, error) {
	var tenant domain.Tenant
	query := `INSERT INTO tenants (name) VALUES ($1) RETURNING id, name, status, created_at`
	if err := sqlx.GetContext(ctx, exec, &tenant, query, name); err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, exec sqlx.QueryerContext, user *domain.User) (*domain.User, error) {
	var created domain.User
	query := `
        INSERT INTO users (tenant_id, name, email, password_hash, role, active)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, tenant_id, name, email, password_hash, role, active, created_at
    `
	if err := sqlx.GetContext(ctx, exec, &created, query, user.TenantID, user.Name, user.Email, user.PasswordHash, user.Role, user.Active); err != nil {
		return nil, err
	}
	return &created, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*UserWithTenant, error) {
	var result struct {
		UserID        int64           `db:"user_id"`
		TenantID      int64           `db:"tenant_id"`
		Name          string          `db:"name"`
		Email         string          `db:"email"`
		PasswordHash  string          `db:"password_hash"`
		Role          domain.UserRole `db:"role"`
		Active        bool            `db:"active"`
		CreatedAt     time.Time       `db:"created_at"`
		TenantName    string          `db:"tenant_name"`
		TenantStatus  string          `db:"tenant_status"`
		TenantCreated time.Time       `db:"tenant_created"`
	}

	query := `
        SELECT
            u.id as user_id,
            u.tenant_id,
            u.name,
            u.email,
            u.password_hash,
            u.role,
            u.active,
            u.created_at,
            t.name as tenant_name,
            t.status as tenant_status,
            t.created_at as tenant_created
        FROM users u
        JOIN tenants t ON t.id = u.tenant_id
        WHERE LOWER(u.email) = LOWER($1)
    `

	if err := r.db.GetContext(ctx, &result, query, email); err != nil {
		return nil, err
	}

	return &UserWithTenant{
		User: domain.User{
			ID:           result.UserID,
			TenantID:     result.TenantID,
			Name:         result.Name,
			Email:        result.Email,
			PasswordHash: result.PasswordHash,
			Role:         result.Role,
			Active:       result.Active,
			CreatedAt:    result.CreatedAt,
		},
		Tenant: domain.Tenant{
			ID:        result.TenantID,
			Name:      result.TenantName,
			Status:    result.TenantStatus,
			CreatedAt: result.TenantCreated,
		},
	}, nil
}

func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func IsNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
