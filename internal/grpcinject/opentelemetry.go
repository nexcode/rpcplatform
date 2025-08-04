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
	"net"
	"strconv"
	"time"

	"github.com/nexcode/rpcplatform/internal/config"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"google.golang.org/grpc"
)

func OpenTelemetry(config *config.Config, localAddr net.Addr, publicAddr string) error {
	resOptions := []resource.Option{
		resource.WithHost(),
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(semconv.ServiceName(config.OpenTelemetry.ServiceName)),
	}

	if localAddr != nil {
		host, port, err := net.SplitHostPort(localAddr.String())
		if err != nil {
			return err
		}

		portInt, err := strconv.Atoi(port)
		if err != nil {
			return err
		}

		resOptions = append(resOptions,
			resource.WithAttributes(semconv.NetworkTransportKey.String(localAddr.Network())),
			resource.WithAttributes(semconv.NetworkLocalAddress(host)),
			resource.WithAttributes(semconv.NetworkLocalPort(portInt)),
		)

		if publicAddr != localAddr.String() {
			host, port, err = net.SplitHostPort(publicAddr)
			if err != nil {
				return err
			}

			portInt, err = strconv.Atoi(port)
			if err != nil {
				return err
			}
		}

		resOptions = append(resOptions,
			resource.WithAttributes(semconv.ServerAddress(host)),
			resource.WithAttributes(semconv.ServerPort(portInt)),
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := resource.New(ctx, resOptions...)
	cancel()

	if err != nil {
		return err
	}

	tpOptions := []trace.TracerProviderOption{
		trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(config.OpenTelemetry.SampleRate))),
		trace.WithResource(res),
	}

	for _, exporter := range config.OpenTelemetry.Exporters {
		tpOptions = append(tpOptions, trace.WithBatcher(exporter))
	}

	tracerProvider := otelgrpc.WithTracerProvider(
		trace.NewTracerProvider(tpOptions...),
	)

	propagators := otelgrpc.WithPropagators(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)

	if localAddr != nil {
		config.GRPCOptions.Server = append(config.GRPCOptions.Server,
			grpc.StatsHandler(otelgrpc.NewServerHandler(tracerProvider, propagators)),
		)
	} else {
		config.GRPCOptions.Client = append(config.GRPCOptions.Client,
			grpc.WithStatsHandler(otelgrpc.NewClientHandler(tracerProvider, propagators)),
		)
	}

	return nil
}
