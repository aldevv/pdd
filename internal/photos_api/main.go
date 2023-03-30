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
	router.MaxMultipartMemory = 100 << 20 // 8 MiB

	router.POST("/create_user", auth.CreateUser)
	router.POST("/login", auth.Login)

	private := router.Group("/api", middleware.Protect)

	// TODO: move savephoto to handlers
	private.POST("/upload", storage.SavePhoto)
	private.GET("/photos", handlers.GetPhotos)
	private.GET("/photos/:id", handlers.GetPhoto)
	private.POST("/result/", handlers.PostResult)
	private.GET("/result/", handlers.GetResult)
	// private.GET("/get_user", auth.GetUser)

	router.Run(s.address)
}
