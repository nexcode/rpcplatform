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

package options

import (
	"github.com/nexcode/rpcplatform/internal/config"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Platform struct{}

// ClientOptions sets default options for all new clients.
// Local client options can override these global settings.
func (Platform) ClientOptions(options ...func(*config.Client)) func(*config.Platform) {
	return func(c *config.Platform) {
		c.ClientOptions = append(c.ClientOptions, options...)
	}
}

// ServerOptions sets default options for all new servers.
// Local server options can override these global settings.
func (Platform) ServerOptions(options ...func(*config.Server)) func(*config.Platform) {
	return func(c *config.Platform) {
		c.ServerOptions = append(c.ServerOptions, options...)
	}
}

// OpenTelemetry configures OpenTelemetry tracing for clients and servers.
func (Platform) OpenTelemetry(serviceName string, sampleRate float64, exporters ...trace.SpanExporter) func(*config.Platform) {
	return func(c *config.Platform) {
		c.OpenTelemetry = &config.OpenTelemetry{
			ServiceName: serviceName,
			SampleRate:  sampleRate,
			Exporters:   exporters,
		}
	}
}
