package ginx

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/soheilhy/cmux"
	"grpc-template-service/conf"
	"grpc-template-service/core/kernel"
	"grpc-template-service/pkg/colorful"
	"net"
	"net/http"
	"sync"
	"time"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule

	listener net.Listener
	g        *gin.Engine
	httpSrv  *http.Server
	// Gin engine is not stored in the struct since it's typically initialized once and reused.
}

func (m *Mod) Name() string {
	return "ginx"
}

func (m *Mod) Init(hub *kernel.Hub) error {
	// Initialize the Gin router
	m.g = gin.Default()

	// CORS configuration
	corsConf := cors.DefaultConfig()
	corsConf.AllowAllOrigins = true
	corsConf.AllowCredentials = true
	corsConf.AddAllowHeaders("Authorization")
	m.g.Use(cors.New(corsConf))

	// Recovery middleware for panic recovery
	m.g.Use(gin.Recovery())

	// Sentry integration if DSN is provided
	if conf.Get().SentryDsn != "" {
		// Sentry middleware initialization would go here.
		// Note: Gin does not have an official Sentry middleware,
		// so you would need to implement this or find a third-party one.
	}

	// Map the router to the hub or store it as needed
	hub.Map(m.g)

	return nil
}

func (m *Mod) Load(hub *kernel.Hub) error {
	var ginE gin.Engine
	err := hub.Load(&ginE)
	if err != nil {
		return errors.New("can't load gin.Engine from kernel")
	}
	return nil
}

func (m *Mod) Start(hub *kernel.Hub) error {
	var tcpMux cmux.CMux
	err := hub.Load(&tcpMux)
	if err != nil {
		return errors.New("can't load tcpMux from kernel")
	}

	// Prepare the HTTP listener for the Gin router
	httpL := tcpMux.Match(cmux.HTTP1Fast())
	m.listener = httpL

	// Create an HTTP server with the Gin router as its handler
	m.httpSrv = &http.Server{
		Handler: m.g,
	}

	// Start serving requests
	if err := m.httpSrv.Serve(httpL); err != nil && !errors.Is(err, http.ErrServerClosed) {
		hub.Log.Infow("failed to start to listen and serve", "error", err)
		return err
	}
	return nil
}

func (m *Mod) Stop(wg *sync.WaitGroup, ctx context.Context) error {
	defer wg.Done()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Graceful shutdown of the HTTP server
	if err := m.httpSrv.Shutdown(ctx); err != nil {
		fmt.Println(colorful.Yellow("Server forced to shutdown: " + err.Error()))
		return err
	}
	return nil
}
