package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	v1 "github.com/LeonardoForconesi/go-template/internal/adapter/http/v1"
	"github.com/LeonardoForconesi/go-template/internal/server/middleware"
)

func NewRouter(log *zap.Logger, h *v1.UserHandlers) *gin.Engine {
	r := gin.New()
	// Middlewares globales
	r.Use(middleware.RequestID())
	r.Use(middleware.Recovery(log))
	r.Use(middleware.Logger(log))
	r.Use(middleware.CORS())
	r.Use(middleware.Timeout(5 * time.Second))

	api := r.Group("/api")
	v := api.Group("/v1")
	h.Register(v)

	// Health
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	// si el dia de mañana quisiera otras versiones
	//v2 := api.Group("/v2")
	//otherHandlers.Register(v2)

	return r
}
