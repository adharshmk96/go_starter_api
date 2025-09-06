package infra

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func SetupOtelSDK(ctx context.Context) (func(context.Context) error, error) {
	var shutdownFuncs []func(context.Context) error
	var err error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}
	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// propagators are used to propagate the trace context and baggage across the different services.
	initPropagators()

	// tracer provider is used to create and manage the tracers.
	tp, err := newTracerProvider(ctx)
	if err != nil {
		handleErr(err)
	}
	otel.SetTracerProvider(tp)
	shutdownFuncs = append(shutdownFuncs, tp.Shutdown)

	// metric provider is used to create and manage the metrics.
	mp, err := newMetricProvider(ctx)
	if err != nil {
		handleErr(err)
	}
	otel.SetMeterProvider(mp)
	shutdownFuncs = append(shutdownFuncs, mp.Shutdown)

	// logger provider is used to create and manage the loggers.
	lp, err := newLoggerProvider(ctx)
	if err != nil {
		handleErr(err)
	}
	global.SetLoggerProvider(lp)
	shutdownFuncs = append(shutdownFuncs, lp.Shutdown)

	return shutdown, nil
}

func initPropagators() {
	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)
}

func newTracerProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	traceExporter, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithInsecure(), // Use HTTP instead of HTTPS
	)
	if err != nil {
		return nil, err
	}

	// Create resource with service name and version
	res, err := newResource()
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)

	return tp, nil
}

func newMetricProvider(ctx context.Context) (*sdkmetric.MeterProvider, error) {
	metricExporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithInsecure(), // Use HTTP instead of HTTPS
	)
	if err != nil {
		return nil, err
	}

	// Create resource with service name and version
	res, err := newResource()
	if err != nil {
		return nil, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(res),
	)

	return mp, nil
}

func newLoggerProvider(ctx context.Context) (*sdklog.LoggerProvider, error) {
	logExporter, err := otlploghttp.New(ctx, otlploghttp.WithInsecure())
	if err != nil {
		return nil, err
	}

	// Create resource with service name and version
	res, err := newResource()
	if err != nil {
		return nil, err
	}

	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(
			sdklog.NewBatchProcessor(logExporter),
		),
		sdklog.WithResource(res),
	)

	return loggerProvider, nil
}

// newResource creates a new OpenTelemetry resource with service name and version
func newResource() (*resource.Resource, error) {
	return resource.New(
		context.Background(),
		resource.WithAttributes(
			// Service name - this will replace "unknown_service:main"
			semconv.ServiceNameKey.String("servicehub-api"),
			// Service version
			semconv.ServiceVersionKey.String("1.0.0"),
			// Service namespace (optional)
			semconv.ServiceNamespaceKey.String("knullsoft"),
			// Additional attributes
			semconv.DeploymentEnvironmentKey.String("development"),
		),
	)
}
