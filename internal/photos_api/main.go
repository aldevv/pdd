package photos_api

import (
	"github.com/gin-gonic/gin"
	"github.com/plant_disease_detection/internal/auth"
	"github.com/plant_disease_detection/internal/db"
	"github.com/plant_disease_detection/internal/handlers"
	"github.com/plant_disease_detection/internal/middleware"
	"github.com/plant_disease_detection/internal/photo_storage"
)

type Server struct {
	router       *gin.Engine
	storage_name string
	address      string
}

type Option func(s *Server)

func NewServer(opts ...Option) *Server {
	s := &Server{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithStorage(storage_name string) Option {
	return func(s *Server) {
		s.storage_name = storage_name
	}
}

func WithAddress(address string) Option {
	return func(s *Server) {
		s.address = address
	}
}

func (s *Server) Serve() {
	db.ConnectDB()
	storage := photo_storage.GetStorage(s.storage_name)

	router := gin.Default()
	router.Use(corsMiddleware())
	router.MaxMultipartMemory = 100 << 20 // 8 MiB

	router.POST("/user", auth.CreateUser)
	router.POST("/login", auth.Login)

	private := router.Group("/api", middleware.Protect)

	private.POST("/upload", storage.SavePhoto)
	private.GET("/photos", handlers.GetPhotos)
	private.GET("/photos/:id", handlers.GetPhoto)
	private.POST("/result", handlers.PostResult)
	private.GET("/result", handlers.GetResult)
	// private.GET("/get_user", auth.GetUser)

	router.Run(s.address)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
