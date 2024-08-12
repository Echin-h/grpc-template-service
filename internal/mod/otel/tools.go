package otel

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	oteltrace "go.opentelemetry.io/otel/sdk/trace"
	"grpc-template-service/conf"
)

func newOTLPExporter(ctx context.Context) (oteltrace.SpanExporter, error) {
	// Change default HTTPS -> HTTP
	insecureOpt := otlptracehttp.WithInsecure()

	// Update default OTLP reciver endpoint
	endPoint := fmt.Sprintf("%s:%s", conf.Get().Otel.AgentHost, conf.Get().Otel.AgentPort)
	endpointOpt := otlptracehttp.WithEndpoint(endPoint)
	return otlptracehttp.New(ctx, insecureOpt, endpointOpt)
}
