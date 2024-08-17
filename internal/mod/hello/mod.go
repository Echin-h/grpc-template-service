package hello

import (
	"errors"
	"fmt"
	helloV1 "github.com/Echin-h/grpc-template-proto/gen/proto/hello/v1"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"grpc-template-service/core/kernel"
	"grpc-template-service/internal/mod/grpcGateway/gateway"
	"grpc-template-service/internal/mod/hello/dao"
	"grpc-template-service/internal/mod/hello/service"
	"grpc-template-service/pkg/colorful"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule
}

func (m *Mod) Name() string { return "hello" }

func (m *Mod) Load(hub *kernel.Hub) error {
	var gw gateway.Gateway
	if hub.Load(&gw) != nil {
		return errors.New("can't load gateway from kernel")
	}
	var GRPC grpc.Server
	if hub.Load(&GRPC) != nil {
		return errors.New("can't load gRPC server from kernel")
	}
	var db *gorm.DB
	if hub.Load(&db) != nil {
		return errors.New("can't load gorm from kernel")
	}
	if err := dao.Init(db); err != nil {
		return err
	}

	helloV1.RegisterGreeterServiceServer(&GRPC, &service.S{
		Log: hub.Log.Named("hello.service"),
	})
	err := gw.Register(helloV1.RegisterGreeterServiceHandler)
	fmt.Println(colorful.Green("hello service Loaded successfully"))
	if err != nil {
		hub.Log.Fatalw("failed to register", "error", err)
	}

	return nil
}
