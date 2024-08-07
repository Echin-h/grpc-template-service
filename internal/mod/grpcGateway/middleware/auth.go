package middleware

import (
	"context"
	grpcAuth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"grpc-template-service/core/logx"
	"grpc-template-service/pkg/auth"
	"grpc-template-service/pkg/ctxKey"
)

// grpc-go Middlewares: interceptors, helpers and utilities.
// It is a perfect way to implement common patterns:
// auth, logging, tracing, metrics, validation, retries, rate limiting and more
// https://github.com/grpc-ecosystem/go-grpc-middleware
// and this file is aimed to add costumed auth

func AuthInterceptor(ctx context.Context) (context.Context, error) {
	token, err := grpcAuth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return ctx, nil
	}

	jwtClaims, err := auth.ParseToken(token)
	if err != nil {
		logx.NameSpace("grpc.middleware.auth").Infof("Error parsing token: %v", err)
		return ctx, nil
	}

	newCtx := context.WithValue(ctx, ctxKey.UID, jwtClaims.Info.UID)
	return newCtx, nil

}
