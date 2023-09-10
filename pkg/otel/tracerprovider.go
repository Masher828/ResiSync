package otel

import (
	pkg_constants "ResiSync/pkg/constants"
	"ResiSync/pkg/logger"
	"context"
	"os"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

var tracingLevel map[string]int = map[string]int{
	pkg_constants.TracingLevelInfo:     1,
	pkg_constants.TracingLevelDebug:    0,
	pkg_constants.TracingLevelCritical: 2}

func InitTracer(appName string) (*trace.TracerProvider, error) {

	if !viper.GetBool(pkg_constants.ConfigTracingEnabled) {
		return nil, nil
	}
	tracerProvider, err := GetTracerProvider(appName)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tracerProvider)
	return tracerProvider, nil
}

func GetJaegarClient() (trace.SpanExporter, error) {

	tracer, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(viper.GetString(pkg_constants.ConfigJaegarUrlCollectorKey))))

	return tracer, err
}

func GetTracerProvider(appName string) (*trace.TracerProvider, error) {

	log := logger.GetAppStartupLogger()

	exporter, err := GetJaegarClient()
	if err != nil {
		log.Error("Error while getting tracer exporter", zap.Error(err))
		return nil, err
	}

	serviceName := appName + "__" + os.Getenv(pkg_constants.EnvEnvironment)

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Error("Error while creating resource", zap.Error(err))
		return nil, err
	}

	return trace.NewTracerProvider(trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(resources)), nil
}

func IsTracingEnabled() bool {
	return viper.GetBool(pkg_constants.ConfigTracingEnabled)
}

func OtelTracingLevel() string {
	return viper.GetString(pkg_constants.ConfigTracingLevel)
}

func TracingLevel(level string) int {
	return tracingLevel[level]
}
