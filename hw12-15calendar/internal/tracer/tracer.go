package tracer

import (
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type Conf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type InfoService struct {
	Name    string
	Version string
}

func TraceInit(conf *Conf, service *InfoService) error {
	jaegerEndpoint := fmt.Sprintf("http://%s:%d/api/traces", conf.Host, conf.Port)

	// Create and configure the Jaeger exporter
	jaegerExporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(jaegerEndpoint),
		),
	)
	if err != nil {
		return err
	}

	// Create and configure the TracerProvider exporter using the
	// newly-created exporters.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(jaegerExporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(service.Name),
			semconv.ServiceVersion(service.Version),
		)),
	)

	// Now we can register tp as the otel trace provider.
	otel.SetTracerProvider(tp)

	return nil
}
