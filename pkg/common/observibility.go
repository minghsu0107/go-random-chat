package common

import (
	"fmt"
	"net/http"

	"github.com/minghsu0107/go-random-chat/pkg/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	propjaeger "go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

var TracerProvider *tracesdk.TracerProvider

type ObservibilityInjector struct {
	promPort  string
	jaegerUrl string
}

func NewObservibilityInjector(config *config.Config) *ObservibilityInjector {
	return &ObservibilityInjector{
		promPort:  config.Observability.Prometheus.Port,
		jaegerUrl: config.Observability.Tracing.JaegerUrl,
	}
}

func (injector *ObservibilityInjector) Register(service string) error {
	if injector.jaegerUrl != "" {
		err := initTracerProvider(injector.jaegerUrl, service)
		if err != nil {
			return err
		}
		otel.SetTracerProvider(TracerProvider)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propjaeger.Jaeger{}, propagation.Baggage{}))
	}
	if injector.promPort != "" {
		go func() {
			log.Infof("starting prom metrics on PROM_PORT=[%s]", injector.promPort)
			err := http.ListenAndServe(fmt.Sprintf(":%s", injector.promPort), promhttp.Handler())
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	return nil
}

func NewOtelHttpHandler(h http.Handler, operation string) http.Handler {
	httpOptions := []otelhttp.Option{
		otelhttp.WithTracerProvider(otel.GetTracerProvider()),
		otelhttp.WithPropagators(otel.GetTextMapPropagator()),
	}
	return otelhttp.NewHandler(h, operation, httpOptions...)
}

func initTracerProvider(jaegerUrl, service string) error {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerUrl)))
	if err != nil {
		return err
	}
	TracerProvider = tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
		)),
	)
	return nil
}
