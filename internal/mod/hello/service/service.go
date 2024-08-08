package service

import (
	"context"
	helloV1 "github.com/Echin-h/grpc-template-proto/gen/proto/hello/v1"
	"go.uber.org/zap"
)

type S struct {
	helloV1.UnimplementedGreeterServiceServer
	Log *zap.SugaredLogger
}

var _ helloV1.GreeterServiceServer = (*S)(nil)

func (s *S) SayHello(_ context.Context, in *helloV1.SayHelloRequest) (*helloV1.SayHelloResponse, error) {
	return &helloV1.SayHelloResponse{Message: in.Name + " world"}, nil
}

func (s *S) SayHelloAgain(_ context.Context, in *helloV1.SayHelloAgainRequest) (*helloV1.SayHelloAgainResponse, error) {
	return &helloV1.SayHelloAgainResponse{Message: in.Name + " world again"}, nil
}

func (s *S) SayHelloThirdly(_ context.Context, in *helloV1.SayHelloThirdlyRequest) (*helloV1.SayHelloThirdlyResponse, error) {
	return &helloV1.SayHelloThirdlyResponse{Message: in.Name + " world thirdly"}, nil
}

//
