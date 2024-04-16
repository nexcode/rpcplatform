/*
 * Copyright 2022 RPCPlatform Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package grpcinject

import (
	"context"
	"strings"
	"time"

	"github.com/nexcode/rpcplatform/errors"
	"github.com/nexcode/rpcplatform/internal/config"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"google.golang.org/grpc"
)

func OpenTelemetry(options interface{}, config config.OpenTelemetryConfig, addr string) error {
	resOptions := []resource.Option{
		resource.WithHost(),
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(semconv.ServiceNameKey.String(config.ServiceName)),
	}

	if addr != "" {
		addr := strings.Split(addr, ":")
		resOptions = append(resOptions,
			resource.WithAttributes(semconv.NetSockHostAddrKey.String(addr[0])),
			resource.WithAttributes(semconv.NetHostPortKey.String(addr[1])),
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := resource.New(ctx, resOptions...)
	cancel()

	if err != nil {
		return err
	}

	tpOptions := []trace.TracerProviderOption{
		trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(config.SampleRate))),
		trace.WithResource(res),
	}

	for _, exporter := range config.Exporters {
		tpOptions = append(tpOptions, trace.WithBatcher(exporter))
	}

	tracerProvider := otelgrpc.WithTracerProvider(
		trace.NewTracerProvider(tpOptions...),
	)

	propagators := otelgrpc.WithPropagators(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)

	switch options := options.(type) {
	case *[]grpc.DialOption:
		*options = append(*options,
			grpc.WithStatsHandler(otelgrpc.NewClientHandler(tracerProvider, propagators)),
		)
	case *[]grpc.ServerOption:
		*options = append(*options,
			grpc.StatsHandler(otelgrpc.NewServerHandler(tracerProvider, propagators)),
		)
	default:
		return errors.ErrGRPCOptionsExpected
	}

	return nil
}
