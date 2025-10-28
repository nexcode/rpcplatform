/*
 * Copyright 2025 RPCPlatform Authors
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

package rpcplatform

import (
	"context"
	"errors"
	"log"
	"net"
	"strconv"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"google.golang.org/grpc/stats"
)

func (p *RPCPlatform) openTelemetry(instanceID string, localAddr net.Addr, publicAddr string) (stats.Handler, error) {
	resOptions := []resource.Option{
		resource.WithHost(),
		resource.WithOS(),
		resource.WithContainer(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(semconv.ServiceName(p.config.OpenTelemetry.ServiceName)),
		resource.WithAttributes(semconv.ServiceInstanceID(instanceID)),
	}

	if localAddr != nil {
		host, port, err := net.SplitHostPort(localAddr.String())
		if err != nil {
			return nil, err
		}

		portInt, err := strconv.Atoi(port)
		if err != nil {
			return nil, err
		}

		resOptions = append(resOptions,
			resource.WithAttributes(semconv.NetworkTransportKey.String(localAddr.Network())),
			resource.WithAttributes(semconv.NetworkLocalAddress(host)),
			resource.WithAttributes(semconv.NetworkLocalPort(portInt)),
		)

		if publicAddr != localAddr.String() {
			host, port, err = net.SplitHostPort(publicAddr)
			if err != nil {
				return nil, err
			}

			portInt, err = strconv.Atoi(port)
			if err != nil {
				return nil, err
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

	if errors.Is(err, resource.ErrPartialResource) || errors.Is(err, resource.ErrSchemaURLConflict) {
		log.Println(err)
	} else if err != nil {
		return nil, err
	}

	tpOptions := []trace.TracerProviderOption{
		trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(p.config.OpenTelemetry.SampleRate))),
		trace.WithResource(res),
	}

	for _, exporter := range p.config.OpenTelemetry.Exporters {
		tpOptions = append(tpOptions, trace.WithBatcher(exporter))
	}

	tracerProvider := otelgrpc.WithTracerProvider(
		trace.NewTracerProvider(tpOptions...),
	)

	propagators := otelgrpc.WithPropagators(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)

	if localAddr != nil {
		return otelgrpc.NewServerHandler(tracerProvider, propagators), nil
	}

	return otelgrpc.NewClientHandler(tracerProvider, propagators), nil
}
