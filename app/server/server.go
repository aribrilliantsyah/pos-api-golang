package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"pos-api/app/middlewares"
	"pos-api/app/routes"
	dbCon "pos-api/db/sqlc"
	"pos-api/util/config"
	"pos-api/util/swagger"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Server struct {
	engine *gin.Engine
	db     *dbCon.Queries
	ctx    context.Context
}

func NewServer(config config.Config) *Server {
	ctx := context.TODO()

	conn, err := sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	db := dbCon.New(conn)
	fmt.Println("PostgreSql connected successfully...")

	server := &Server{
		engine: gin.Default(),
		db:     db,
		ctx:    ctx,
	}

	// Initialize Swagger
	swagger.Initialize(server.engine)

	// Set up routes
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	prefix := "/api/v1/"
	s.engine.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "welcome to api, docs: /docs/index.html",
		})
	})

	// Public routes
	group := s.engine.Group(prefix)
	routes.SetupAuthRoutes(s.db, s.ctx, group)

	// Protected routes with AuthMiddleware
	protected := s.engine.Group(prefix)
	protected.Use(middlewares.AuthMiddleware(s.db, s.ctx))
	routes.SetupUserRoutes(s.db, s.ctx, protected)
	routes.SetupCategoryRoutes(s.db, s.ctx, protected)
	routes.SetupCustomerRoutes(s.db, s.ctx, protected)
	routes.SetupProductRoutes(s.db, s.ctx, protected)

	// Handle 404
	s.engine.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": fmt.Sprintf("The specified route %s not found", ctx.Request.URL),
		})
	})
}

func (s *Server) Run() error {
	config, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}
	return s.engine.Run(":" + config.ServerAddress)
}
