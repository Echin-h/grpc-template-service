package ginx

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/soheilhy/cmux"
	"grpc-template-service/core/kernel"
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
	//corsConf := cors.DefaultConfig()
	//corsConf.AllowAllOrigins = true
	//corsConf.AllowCredentials = true
	//corsConf.AddAllowHeaders("Authorization")
	//m.g.Use(cors.New(corsConf))

	hub.Map(m.g)
	return nil
}

func (m *Mod) Load(hub *kernel.Hub) error { return nil }

func (m *Mod) Start(hub *kernel.Hub) error {
	var tcpMux cmux.CMux
	err := hub.Load(&tcpMux)
	if err != nil {
		return errors.New("can't load tcpMux from kernel")
	}

	httpL := tcpMux.Match(cmux.HTTP1Fast())
	m.listener = httpL

	m.httpSrv = &http.Server{Handler: m.g}

	if err = m.httpSrv.Serve(httpL); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
		return err
	}
	return nil
}
