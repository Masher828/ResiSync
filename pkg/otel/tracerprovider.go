package otel

import (
	"ResiSync/pkg/constants"
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

func InitTracer(appName string) (*trace.TracerProvider, error) {

	if !viper.GetBool(constants.ConfigTracingEnabled) {
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

	tracer, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(viper.GetString(constants.ConfigJaegarUrlCollectorKey))))

	return tracer, err
}

func GetTracerProvider(appName string) (*trace.TracerProvider, error) {

	log := logger.GetAppStartupLogger()

	exporter, err := GetJaegarClient()
	if err != nil {
		log.Error("Error while getting tracer exporter", zap.Error(err))
		return nil, err
	}

	serviceName := appName + "__" + os.Getenv(constants.EnvEnvironment)

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.Name", serviceName),
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
