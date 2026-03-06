// Package server wires up the Gin HTTP server via Fx.
package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
	"go.uber.org/zap"

	_ "cc.io/arena/docs" // swag-generated docs (import for side-effects)
	"cc.io/arena/internal/middleware"
	log "cc.io/arena/pkg/logging"
)

// Config holds HTTP server configuration.
type Config struct {
	Port int
}

// New constructs and registers a Gin-based HTTP server with the Fx lifecycle.
// The server starts on OnStart and shuts down gracefully on OnStop.
func New(lc fx.Lifecycle, cfg Config, logger *log.Logger) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.AccessLog(logger))

	// Health check – intentionally excluded from access logging.
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			logger.Info("starting HTTP server", zap.String("addr", addr))
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error("HTTP server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("shutting down HTTP server")
			shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			return srv.Shutdown(shutdownCtx)
		},
	})

	return router
}

// Module is an Fx option that provides the Gin router.
func Module(cfg Config) fx.Option {
	return fx.Provide(func(lc fx.Lifecycle, logger *log.Logger) *gin.Engine {
		return New(lc, cfg, logger)
	})
}
