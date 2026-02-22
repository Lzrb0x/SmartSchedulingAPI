package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"

	"github.com/Lzrb0x/SmartSchedulingAPI/internal/config"
	"github.com/Lzrb0x/SmartSchedulingAPI/internal/domain"
	"github.com/Lzrb0x/SmartSchedulingAPI/internal/tenant"
)

type Handler struct {
	db  *sqlx.DB
	cfg config.AuthConfig
}

func NewHandler(db *sqlx.DB, cfg config.AuthConfig) *Handler {
	return &Handler{db: db, cfg: cfg}
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
	Token string `json:"token"`
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

	// TODO: replace with real user lookup
	claims := &Claims{
		TenantID: 1,
		UserID:   1,
		Role:     domain.RoleOwner,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    h.cfg.JWTIssuer,
			Subject:   "1",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, loginResponse{Token: signed})
}

func (h *Handler) register(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"status": "registered"})
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
