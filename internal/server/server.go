package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/Lzrb0x/SmartSchedulingAPI/internal/auth"
	"github.com/Lzrb0x/SmartSchedulingAPI/internal/config"
	"github.com/Lzrb0x/SmartSchedulingAPI/internal/tenant"
)

type Server struct {
	cfg    *config.Config
	db     *sqlx.DB
	engine *gin.Engine
}

func New(cfg *config.Config, db *sqlx.DB) *Server {
	if cfg.Environment == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	s := &Server{
		cfg:    cfg,
		db:     db,
		engine: router,
	}

	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	s.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := s.engine.Group("/api")

	authHandler := auth.NewHandler(s.db, s.cfg.Auth)
	authHandler.RegisterRoutes(api.Group("/auth"))

	tenantGroup := api.Group("")
	tenantGroup.Use(auth.JWTMiddleware(s.cfg.Auth))
	tenantGroup.GET("/tenants/current", tenant.CurrentTenantHandler())
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.HTTP.Host, s.cfg.HTTP.Port)

	srv := &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	log.Printf("server started on %s", addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Println("shutting down server...")
	return srv.Shutdown(ctx)
}
