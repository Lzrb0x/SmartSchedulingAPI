package auth

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"github.com/Lzrb0x/SmartSchedulingAPI/internal/config"
	"github.com/Lzrb0x/SmartSchedulingAPI/internal/database"
	"github.com/Lzrb0x/SmartSchedulingAPI/internal/domain"
)

var (
	ErrInvalidCredentials = errors.New("invalid_credentials")
	ErrEmailInUse         = errors.New("email_in_use")
)

type Service struct {
	repo *database.UserRepository
	cfg  config.AuthConfig
}

func NewService(db *sqlx.DB, cfg config.AuthConfig) *Service {
	return &Service{
		repo: database.NewUserRepository(db),
		cfg:  cfg,
	}
}

type RegisterInput struct {
	TenantName string
	Name       string
	Email      string
	Password   string
	Role       domain.UserRole
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (*domain.User, *domain.Tenant, error) {
	if input.Role != domain.RoleOwner {
		return nil, nil, errors.New("only owner registration supported")
	}

	tx, err := s.repo.DB().BeginTxx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	tenant, err := s.repo.CreateTenant(ctx, tx, input.TenantName)
	if err != nil {
		return nil, nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	user := &domain.User{
		TenantID:     tenant.ID,
		Name:         input.Name,
		Email:        strings.ToLower(input.Email),
		PasswordHash: string(hash),
		Role:         input.Role,
		Active:       true,
	}

	created, err := s.repo.CreateUser(ctx, tx, user)
	if err != nil {
		if database.IsUniqueViolation(err) {
			return nil, nil, ErrEmailInUse
		}
		return nil, nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}

	return created, tenant, nil
}

type LoginOutput struct {
	Token  string
	User   domain.User
	Tenant domain.Tenant
}

func (s *Service) Login(ctx context.Context, email, password string) (*LoginOutput, error) {
	userWithTenant, err := s.repo.GetUserByEmail(ctx, strings.ToLower(email))
	if err != nil {
		if database.IsNoRows(err) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userWithTenant.User.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	claims := &Claims{
		TenantID: userWithTenant.Tenant.ID,
		UserID:   userWithTenant.User.ID,
		Role:     userWithTenant.User.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    s.cfg.JWTIssuer,
			Subject:   strconv.FormatInt(userWithTenant.User.ID, 10),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		Token:  signed,
		User:   userWithTenant.User,
		Tenant: userWithTenant.Tenant,
	}, nil
}
