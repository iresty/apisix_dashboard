package server2

import (
	"fmt"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/apache/apisix-dashboard/api/internal/config"
	"github.com/apache/apisix-dashboard/api/internal/filter"
)

type Server struct {
	initialized bool
	config      config.Config
	gin         *gin.Engine
}

func NewServer(cfg config.Config) *Server {
	r := gin.New()

	s := &Server{
		gin:    r,
		config: cfg,
	}

	return s
}

func (s *Server) setupMiddlewares() {
	r := s.gin

	staticPath := "./html/"
	r.Use(
		filter.IPFilter(s.config.Security),
		filter.InvalidRequest(),
		filter.Authentication(s.config.Authentication),
		gzip.Gzip(gzip.DefaultCompression),
		filter.CORS(s.config.Security),
		filter.RequestId(),
		filter.SchemaCheck(),
		filter.RecoverHandler(),
		//static.Serve("/", static.ServeFileSystem()), //TODO
	)
	//r.StaticFS("/", http.FS)

	r.NoRoute(func(c *gin.Context) {
		c.File(fmt.Sprintf("%s/index.html", staticPath))
	})

}

func (s *Server) setupHandler() {

}

func (s *Server) Listen() {
	if !s.initialized {
		s.setupMiddlewares()
		s.setupHandler()
	}
}
