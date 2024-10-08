package service

import (
	"context"
	helloV1 "github.com/Echin-h/grpc-template-proto/gen/proto/hello/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"grpc-template-service/internal/mod/hello/dao"
	model2 "grpc-template-service/internal/mod/hello/model"
)

type S struct {
	helloV1.UnimplementedGreeterServiceServer
	Log *zap.SugaredLogger
}

var _ helloV1.GreeterServiceServer = (*S)(nil)

func (s *S) SayHello(ctx context.Context, in *helloV1.SayHelloRequest) (*helloV1.SayHelloResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	_, span := otel.Tracer("grpc-template-service").
		Start(ctx, "SayHello",
			trace.WithAttributes(
				attribute.String("name", in.GetName()),
				attribute.StringSlice("client-id", md.Get("client-id")),
				attribute.StringSlice("user-id", md.Get("user-id")),
			),
		)
	defer span.End()
	s.Log.Info("SayHello==========================================================")
	var model = model2.Hello{
		ID:        9,
		Name:      in.Name,
		Age:       77,
		Email:     "xx5xx",
		Telephone: "1299999",
	}
	tx := dao.Get().WithContext(ctx).Create(&model)
	if tx.Error != nil {
		return nil, status.Error(500, tx.Error.Error())
	}
	return &helloV1.SayHelloResponse{Message: in.Name + " world"}, nil
}

func (s *S) SayHelloAgain(_ context.Context, in *helloV1.SayHelloAgainRequest) (*helloV1.SayHelloAgainResponse, error) {
	return &helloV1.SayHelloAgainResponse{Message: in.Name + " world again"}, nil
}

func (s *S) SayHelloThirdly(_ context.Context, in *helloV1.SayHelloThirdlyRequest) (*helloV1.SayHelloThirdlyResponse, error) {
	return &helloV1.SayHelloThirdlyResponse{Message: in.Name + " world thirdly"}, nil
}
