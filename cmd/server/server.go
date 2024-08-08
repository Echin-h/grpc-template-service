package server

import (
	"fmt"
	"github.com/soheilhy/cmux"
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"
	"grpc-template-service/conf"
	"grpc-template-service/core/kernel"
	"grpc-template-service/core/logx"
	"grpc-template-service/internal/mod/ginx"
	"grpc-template-service/internal/mod/grpcGateway"
	"grpc-template-service/internal/mod/hello"
	"grpc-template-service/pkg/colorful"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var log = logx.NameSpace("cmd.server")

var (
	configYml string
	StartCmd  = &cobra.Command{
		Use:     "server",
		Short:   "Set Application config info",
		Example: "go run main.go server -c ./config.yaml",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info("loading config...")
			conf.LoadConfig(configYml)

			log.Info("init dep...")

			if conf.Get().MODE == "" || conf.Get().MODE == "debug" {
				logx.Init(zapcore.DebugLevel)
			} else {
				logx.Init(zapcore.InfoLevel)
			}
			defer func() {
				if err := recover(); err != nil {
					_ = log.Sync()
				}
			}()
			log.Info("init dep complete")

			log.Info("init kernel...")
			conn, err := net.Listen("tcp", fmt.Sprintf("%s:%s", conf.Get().Host, conf.Get().Port))
			if err != nil {
				log.Fatalw("failed to listen", "error", err)
			}

			tcpMux := cmux.New(conn)
			log.Infow("server start", "port", conf.Get().Port)

			k := kernel.New(
				kernel.Config{
					Listener: conn,
				},
			)
			k.Map(&tcpMux, &conn)
			k.RegMod(
				// TODO: Add your module here
				&grpcGateway.Mod{},
				&ginx.Mod{},
				&hello.Mod{},
			)
			k.Init()

			log.Info("init module...")
			err = k.StartModule()
			log.Info("init module complete")

			log.Info("start server...")
			k.Serve()
			// grpc-gateway start server and the grpc-client is bind with the gateway
			go func() {
				_ = tcpMux.Serve()
			}()

			fmt.Println(colorful.Green("Server run at:"))
			fmt.Printf("-  Local:   http://localhost:%s\n", conf.Get().Port)

			// graceful shutdown
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit
			fmt.Println(colorful.Blue("Shutting down server..."))

			err = k.Stop()
			if err != nil {
				panic(err)
			}
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "", "config file path")
}
