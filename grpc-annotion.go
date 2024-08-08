package main

//
//import (
//"context"
//"errors"
//"fmt"
//"github.com/gin-gonic/gin"
//grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
//grpcAuth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
//grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
//grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
//grpcCtxTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
//grpcOpentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
//"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
//"github.com/soheilhy/cmux"
//"google.golang.org/grpc"
//"google.golang.org/grpc/reflection"
//"grpc-template-service/conf"
//"grpc-template-service/core/kernel"
//"grpc-template-service/core/logx"
//"grpc-template-service/internal/mod/grpcGateway/gateway"
//"grpc-template-service/internal/mod/grpcGateway/middleware"
//"net"
//"sync"
//)
//
//var _ kernel.Module = (*Mod)(nil)
//
//type Mod struct {
//	kernel.UnimplementedModule
//
//	grpcL net.Listener
//	grpc  *grpc.Server
//	gw    *gateway.Gateway
//}
//
//func (m *Mod) Name() string { return "grpcGateway" }
//
//// PreInit the grpc-gateway module with a grpc server into injector
//func (m *Mod) PreInit(hub *kernel.Hub) error {
//	grpcZap.ReplaceGrpcLoggerV2(logx.NameSpace("grpc").Desugar())
//	m.grpc = grpc.NewServer(grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
//		grpcCtxTags.UnaryServerInterceptor(),
//		grpcOpentracing.UnaryServerInterceptor(),
//		grpcZap.UnaryServerInterceptor(logx.NameSpace("grpc").Desugar()),
//		grpcRecovery.UnaryServerInterceptor(),
//		grpcAuth.UnaryServerInterceptor(middleware.AuthInterceptor),
//	)))
//	reflection.Register(m.grpc)
//	hub.Log.Info("init gRPC server success...")
//	hub.Map(m.grpc)
//	return nil
//}
//
//// PostInit the grpc-gateway module with a grpc client, it is important to note that the gateway is created here
//func (m *Mod) PostInit(hub *kernel.Hub) error {
//	var tcpMux cmux.CMux
//	err := hub.Load(&tcpMux)
//	if err != nil {
//		return errors.New("failed to load tcpMux")
//	}
//
//	// create a net.listener for grpc server
//	m.grpcL = tcpMux.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
//
//	// 开始初始化 grpc client
//	//opts := []grpc.DialOption{
//	//	grpc.WithTimeout(10 * time.Second),
//	//	//grpc.WithBlock(),
//	//	grpc.WithTransportCredentials(insecure.NewCredentials()),
//	//}
//
//	conn, err := grpc.NewClient(fmt.Sprintf("0.0.0.0:%s", conf.Get().Port), grpc.WithInsecure())
//	if err != nil {
//		hub.Log.Fatalw("grpc fail to Dail : %v", err)
//	}
//
//	// 开始初始化 grpc-gateway
//	//var allowedHeaders = map[string]struct{}{
//	//	"x-request-id": {}, // 还没用到 后续做追踪
//	//}
//	//outHeaderFilter := func(s string) (string, bool) {
//	//	if _, isAllowed := allowedHeaders[s]; isAllowed {
//	//		return strings.ToUpper(s), true
//	//	}
//	//	return s, false
//	//}
//	// serverMux is a grpc-gateway multiplexer that serves the provided gRPC server.
//	mux := runtime.NewServeMux()
//	//runtime.WithOutgoingHeaderMatcher(outHeaderFilter),
//	//runtime.WithHealthzEndpoint(""),  健康检查
//	//runtime.WithOutgoingHeaderMatcher() 这个方法可以用来过滤header
//	//runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true}),为了实现将响应封装在固定格式json.data中
//	// .......
//
//	// 将grpc gateway 嵌入到 gin 中
//	var http gin.Engine
//	err = hub.Load(&http)
//	if err != nil {
//		return errors.New("failed to load gin from kernel")
//	}
//
//	http.Any("/v1/*any", func(c *gin.Context) { mux.ServeHTTP(c.Writer, c.Request) })
//
//	m.gw = &gateway.Gateway{
//		Mux:  mux,
//		Conn: conn,
//	}
//
//	hub.Log.Info("init gRPC gateway success...")
//	hub.Map(m.gw)
//	return nil
//}
//
//func (m *Mod) Start(hub *kernel.Hub) error {
//	go func() {
//		if err := m.grpc.Serve(m.grpcL); err != nil {
//			hub.Log.Infow("grpc server failed to serve", "error", err)
//		}
//	}()
//	return nil
//}
//
//func (m *Mod) Stop(wg *sync.WaitGroup, _ context.Context) error {
//	defer wg.Done()
//	m.grpc.GracefulStop()
//	return nil
//}
