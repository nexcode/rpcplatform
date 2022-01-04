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

package options

import (
	"github.com/nexcode/rpcplatform/internal/config"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

func ClientOptions(options ...grpc.DialOption) Option {
	return newFuncOption(func(c *config.Config) {
		c.GRPCOptions.Client = options
	})
}

func ServerOptions(options ...grpc.ServerOption) Option {
	return newFuncOption(func(c *config.Config) {
		c.GRPCOptions.Server = options
	})
}

func OpenTelemetry(serviceName string, sampleRate float64, exporters ...trace.SpanExporter) Option {
	return newFuncOption(func(c *config.Config) {
		c.OpenTelemetry = config.OpenTelemetryConfig{
			ServiceName: serviceName,
			SampleRate:  sampleRate,
			Exporters:   exporters,
		}
	})
}
