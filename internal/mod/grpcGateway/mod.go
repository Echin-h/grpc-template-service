package grpcGateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcAuth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcCtxTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpcOpentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/soheilhy/cmux"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"grpc-template-service/conf"
	"grpc-template-service/core/kernel"
	"grpc-template-service/core/logx"
	"grpc-template-service/internal/mod/grpcGateway/gateway"
	"grpc-template-service/internal/mod/grpcGateway/middleware"
	"grpc-template-service/pkg/colorful"
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

func (m *Mod) Name() string {
	return "grpcGateway"
}

func (m *Mod) PreInit(hub *kernel.Hub) error {
	grpcZap.ReplaceGrpcLoggerV2(logx.NameSpace("grpc").Desugar())
	m.grpc = grpc.NewServer(
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			grpcCtxTags.UnaryServerInterceptor(),
			grpcOpentracing.UnaryServerInterceptor(),
			grpcZap.UnaryServerInterceptor(logx.NameSpace("grpc").Desugar()),
			grpcRecovery.UnaryServerInterceptor(),
			grpcAuth.UnaryServerInterceptor(middleware.AuthInterceptor),
		)), grpc.StatsHandler(otelgrpc.NewServerHandler()))
	reflection.Register(m.grpc)
	hub.Log.Info("init gRPC server success...")
	hub.Map(m.grpc)
	return nil
}

func (m *Mod) PostInit(h *kernel.Hub) error {
	var tcpMux cmux.CMux
	err := h.Load(&tcpMux)
	if err != nil {
		return errors.New("can't load tcpMux from kernel")
	}

	m.grpcL = tcpMux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))

	// 开始初始化grpc-gateway
	opts := []grpc.DialOption{
		//grpc.WithTimeout(10 * time.Second),
		//grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}
	conn, err := grpc.NewClient(fmt.Sprintf("127.0.0.1:%s", conf.Get().Port), opts...)
	if err != nil {
		h.Log.Fatal("gRPC fail to dial: %v", err)
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

	mux := runtime.NewServeMux(
		runtime.WithOutgoingHeaderMatcher(outHeaderFilter),
		runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler),
	)

	var http gin.Engine
	err = h.Load(&http)
	if err != nil {
		return errors.New("can't load jin from kernel")
	}

	http.Any("/v1/*any", func(c *gin.Context) {
		mux.ServeHTTP(c.Writer, c.Request)
	})
	m.gw = &gateway.Gateway{
		Mux:  mux,
		Conn: conn,
	}
	h.Log.Info("init gRPC gateway success...")
	h.Map(m.gw)
	return nil
}

func (m *Mod) Load(h *kernel.Hub) error {
	h.Log.Infow("grpcGateway service Loaded successfully")
	fmt.Println(colorful.Green("grpcGateway service Loaded successfully"))
	return nil
}

func (m *Mod) Start(h *kernel.Hub) error {
	// 初始化grpc
	go func() {
		if err := m.grpc.Serve(m.grpcL); err != nil {
			h.Log.Infow("failed to start to listen and serve", "error", err)
		}
	}()
	fmt.Println(colorful.Green("grpc server run successfully"))
	return nil
}

func (m *Mod) Stop(wg *sync.WaitGroup, _ context.Context) error {
	defer wg.Done()
	m.grpc.GracefulStop()
	return nil
}
