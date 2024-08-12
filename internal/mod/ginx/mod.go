package ginx

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/soheilhy/cmux"
	"grpc-template-service/conf"
	"grpc-template-service/core/kernel"
	"grpc-template-service/internal/mod/ginx/middleware"
	"grpc-template-service/pkg/colorful"
	"net"
	"net/http"
	"sync"
	"time"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule // 请为所有Module引入UnimplementedModule

	listener net.Listener
	g        *gin.Engine
	httpSrv  *http.Server
}

func (m *Mod) Name() string {
	return "ginx"
}

func (m *Mod) Init(hub *kernel.Hub) error {
	m.g = gin.New()
	// TODO: cors

	if conf.Get().MODE == "" {
		gin.SetMode("debug")
	} else {
		gin.SetMode(conf.Get().MODE)
	}

	m.g.Use(
		gin.Recovery(),
		gin.Logger(),
		//otelgin.Middleware("grpc-template-service"),
		middleware.Trace(),
	)

	//if conf.Get().SentryDsn != "" {
	//	m.j.Use(sentryjin.New(sentryjin.Options{Repanic: true}))
	//}

	hub.Map(m.g)
	return nil
}

func (m *Mod) Load(hub *kernel.Hub) error {
	var ginE gin.Engine
	err := hub.Load(&ginE)
	if err != nil {
		return errors.New("can't load jin.Engine from kernel")
	}
	fmt.Println(colorful.Green("jin.Engine Loaded successfully"))
	return nil
}

func (m *Mod) Start(hub *kernel.Hub) error {
	var tcpMux cmux.CMux
	err := hub.Load(&tcpMux)
	if err != nil {
		return errors.New("can't load tcpMux from kernel")
	}

	httpL := tcpMux.Match(cmux.HTTP1Fast())
	m.listener = httpL
	m.httpSrv = &http.Server{
		Handler: m.g,
	}

	if err := m.httpSrv.Serve(httpL); err != nil && !errors.Is(err, http.ErrServerClosed) {
		hub.Log.Infow("failed to start to listen and serve", "error", err)
	}
	return nil
}

func (m *Mod) Stop(wg *sync.WaitGroup, ctx context.Context) error {
	defer wg.Done()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := m.httpSrv.Shutdown(ctx); err != nil {
		fmt.Println(colorful.Yellow("Server forced to shutdown: " + err.Error()))
		return err
	}
	return nil
}
