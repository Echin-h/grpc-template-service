package otel

import (
	"context"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"grpc-template-service/conf"
	"grpc-template-service/core/kernel"
	"grpc-template-service/pkg/colorful"
	"sync"
	"time"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule
	shutdown []func(ctx context.Context) error
	tracer   trace.Tracer
}

func (m *Mod) Name() string { return "otel" }

func (m *Mod) Init(hub *kernel.Hub) error {
	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(conf.Get().Otel.ServiceName),
		),
	)
	if err != nil {
		hub.Log.Errorw("failed to merge resource", "error", err)
		return err
	}

	provider := sdktrace.NewTracerProvider(sdktrace.WithResource(r))
	otel.SetTracerProvider(provider)
	exp, err := newOTLPExporter(context.Background())
	if err != nil {
		hub.Log.Errorw("failed to create exporter", "error", err)
		return err
	}

	bsp := sdktrace.NewBatchSpanProcessor(exp)
	provider.RegisterSpanProcessor(bsp)

	m.shutdown = append(m.shutdown, provider.Shutdown)

	m.tracer = otel.GetTracerProvider().Tracer("grpc-template-service")
	if m.tracer == nil {
		hub.Log.Errorw("failed to get tracer", "error", err)
		return err
	}
	hub.Map(&m.tracer)

	return nil
}

func (m *Mod) Load(h *kernel.Hub) error {
	if m.shutdown == nil {
		return errors.New("otel shutdown not initialized")
	}
	var tc trace.Tracer
	if err := h.Load(&tc); err != nil {
		return errors.New("can't load tracer from kernel")
	}
	fmt.Println(colorful.Green("otel loaded successfully"))
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Printf("-  Local:   http://localhost:%v\n", 3000)
	}()

	return nil
}

func (m *Mod) Stop(wg *sync.WaitGroup, ctx context.Context) error {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if m.shutdown != nil {
		for _, fn := range m.shutdown {
			fn(ctx)
		}
	}

	return nil
}
