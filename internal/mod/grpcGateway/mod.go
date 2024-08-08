package grpcGateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc-template-service/conf"
	"grpc-template-service/core/kernel"
	"grpc-template-service/internal/mod/grpcGateway/gateway"
	"net"
	"strings"
	"sync"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule

	grpcL net.Listener
	grpc  *grpc.Server
	gw    *gateway.Gateway
}

func (m *Mod) Name() string { return "grpcGateway" }

// PreInit the grpc-gateway module with a grpc server into injector
func (m *Mod) PreInit(hub *kernel.Hub) error {
	m.grpc = grpc.NewServer()
	reflection.Register(m.grpc)
	hub.Map(m.grpc)
	return nil
}

// PostInit the grpc-gateway module with a grpc client, it is important to note that the gateway is created here
func (m *Mod) PostInit(hub *kernel.Hub) error {
	var tcpMux cmux.CMux
	err := hub.Load(&tcpMux)
	if err != nil {
		return errors.New("failed to load tcpMux")
	}
	m.grpcL = tcpMux.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	//------------------------------------------------------------------------------------------------------------client
	conn, err := grpc.NewClient(fmt.Sprintf("0.0.0.0:%s", conf.Get().Port), grpc.WithInsecure())
	if err != nil {
		hub.Log.Fatalw("grpc fail to Dail : %v", err)
	}

	var allowedHeaders = map[string]struct{}{
		"x-request-id": {}, // 还没用到 后续做追踪
	}

	outHeaderFilter := func(s string) (string, bool) {
		if _, isAllowed := allowedHeaders[s]; isAllowed {
			return strings.ToUpper(s), true
		}
		return s, false
	}

	// serverMux is a grpc-gateway multiplexer that serves the provided gRPC server.
	mux := runtime.NewServeMux(
		runtime.WithOutgoingHeaderMatcher(outHeaderFilter),
		runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler),
	)

	var http gin.Engine
	err = hub.Load(&http)
	if err != nil {
		return errors.New("failed to load gin from kernel")
	}

	http.Any("/v1/*any", func(c *gin.Context) { mux.ServeHTTP(c.Writer, c.Request) })

	m.gw = &gateway.Gateway{Mux: mux, Conn: conn}
	hub.Map(m.gw)
	return nil
}

func (m *Mod) Start(hub *kernel.Hub) error {
	go func() {
		if err := m.grpc.Serve(m.grpcL); err != nil {
			hub.Log.Infow("grpc server failed to serve", "error", err)
		}
	}()
	return nil
}

func (m *Mod) Stop(wg *sync.WaitGroup, _ context.Context) error {
	defer wg.Done()
	m.grpc.GracefulStop()
	return nil
}
