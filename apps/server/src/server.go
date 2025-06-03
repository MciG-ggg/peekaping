package main

import (
	"net/http"
	"peekaping/src/config"
	"peekaping/src/modules/auth"
	"peekaping/src/modules/monitor"
	"peekaping/src/modules/notification"
	"peekaping/src/modules/proxy"
	"peekaping/src/modules/setting"
	"peekaping/src/modules/websocket"

	_ "peekaping/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// @Summary      Get server version
// @Description  Returns the current server version
// @Tags         System
// @Produce      json
// @Success      200  {object}  map[string]string  "{"version": "1.2.3"}"
// @Router       /version [get]
func versionHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"version": Version})
}

// @Summary      Get server health
// @Description  Returns the current server health
// @Tags         System
// @Produce      json
// @Success      200  {object}  map[string]string  "{"status": "success"}"
// @Router       /health [get]
func healthHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

type Server struct {
	router *gin.Engine
	cfg    *config.Config
}

func ProvideServer(
	logger *zap.SugaredLogger,
	cfg *config.Config,
	monitorRoute *monitor.MonitorRoute,
	monitorController *monitor.MonitorController,
	authRoute *auth.Route,
	authController *auth.Controller,
	wsServer *websocket.Server,
	notificationRoute *notification.Route,
	notificationController *notification.Controller,
	proxyRoute *proxy.Route,
	proxyController *proxy.Controller,
	settingRoute *setting.Route,
	settingController *setting.Controller,
) *Server {
	server := gin.Default()
	// server := gin.New()

	server.RedirectTrailingSlash = false

	// CORS configuration
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: true,
	}))

	// server.Use(LogMiddleware(logger))

	router := server.Group("/api/v1")
	router.GET("/health", healthHandler)
	router.GET("/version", versionHandler)

	// Connect routes
	monitorRoute.ConnectRoute(router, monitorController)
	authRoute.ConnectRoute(router, authController)
	notificationRoute.ConnectRoute(router, notificationController)
	proxyRoute.ConnectRoute(router, proxyController)
	settingRoute.ConnectRoute(router, settingController)

	// Swagger routes
	url := ginSwagger.URL("/swagger/doc.json")
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// WebSocket route
	server.GET("/socket.io/*f", func(c *gin.Context) {
		wsServer.ServeHTTP(c.Writer, c.Request)
	})
	server.POST("/socket.io/*f", func(c *gin.Context) {
		wsServer.ServeHTTP(c.Writer, c.Request)
	})

	return &Server{router: server, cfg: cfg}
}
