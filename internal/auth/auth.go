package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"

	"github.com/Lzrb0x/SmartSchedulingAPI/internal/config"
	"github.com/Lzrb0x/SmartSchedulingAPI/internal/domain"
	"github.com/Lzrb0x/SmartSchedulingAPI/internal/tenant"
)

type Handler struct {
	service *Service
}

func NewHandler(db *sqlx.DB, cfg config.AuthConfig) *Handler {
	return &Handler{service: NewService(db, cfg)}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/login", h.login)
	r.POST("/register", h.register)
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Token  string        `json:"token"`
	User   domain.User   `json:"user"`
	Tenant domain.Tenant `json:"tenant"`
}

type Claims struct {
	TenantID int64           `json:"tenant_id"`
	UserID   int64           `json:"user_id"`
	Role     domain.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type UserContext struct {
	ID       int64
	TenantID int64
	Role     domain.UserRole
}

const userContextKey = "current_user"

func setCurrentUser(c *gin.Context, user UserContext) {
	c.Set(userContextKey, user)
}

func CurrentUser(c *gin.Context) (UserContext, bool) {
	val, ok := c.Get(userContextKey)
	if !ok {
		return UserContext{}, false
	}
	user, ok := val.(UserContext)
	return user, ok
}

func (h *Handler) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	resp, err := h.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

	c.JSON(http.StatusOK, loginResponse{Token: resp.Token, User: resp.User, Tenant: resp.Tenant})
}

type registerRequest struct {
	TenantName string          `json:"tenant_name" binding:"required"`
	Name       string          `json:"name" binding:"required"`
	Email      string          `json:"email" binding:"required,email"`
	Password   string          `json:"password" binding:"required,min=6"`
	Role       domain.UserRole `json:"role" binding:"required"`
}

func (h *Handler) register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	user, tenant, err := h.service.Register(ctx, RegisterInput{
		TenantName: req.TenantName,
		Name:       req.Name,
		Email:      req.Email,
		Password:   req.Password,
		Role:       req.Role,
	})
	if err != nil {
		status := http.StatusBadRequest
		message := err.Error()
		if errors.Is(err, ErrEmailInUse) {
			message = "email already in use"
		} else if !errors.Is(err, ErrEmailInUse) && !errors.Is(err, ErrInvalidCredentials) {
			status = http.StatusInternalServerError
			message = "registration failed"
		}
		c.JSON(status, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"tenant": tenant,
		"user":   user,
	})
}

func JWTMiddleware(cfg config.AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		if tokenString == header {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		tenant.SetTenant(c, claims.TenantID)
		setCurrentUser(c, UserContext{ID: claims.UserID, TenantID: claims.TenantID, Role: claims.Role})
		c.Next()
	}
}
