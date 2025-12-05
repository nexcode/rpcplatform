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
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/nexcode/rpcplatform/internal/attributes"
	etcd "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestNew(t *testing.T) {
	t.Parallel()

	type input struct {
		etcdPrefix string
		options    []PlatformOption
	}

	type expected struct {
		etcdPrefix       *string
		otelServiceName  *string
		otelSampleRate   *float64
		otelExportersLen *int
		clientOptionsLen *int
		serverOptionsLen *int
	}

	tests := []struct {
		name     string
		input    input
		expected expected
	}{
		{
			"Empty etcd prefix",
			input{etcdPrefix: ""},
			expected{etcdPrefix: pointer("")},
		}, {
			"Only slash in etcd prefix",
			input{etcdPrefix: "/"},
			expected{etcdPrefix: pointer("")},
		}, {
			"No slash in etcd prefix",
			input{etcdPrefix: "a"},
			expected{etcdPrefix: pointer("/a")},
		}, {
			"Leading slash in etcd prefix",
			input{etcdPrefix: "/a"},
			expected{etcdPrefix: pointer("/a")},
		}, {
			"Trailing slash in etcd prefix",
			input{etcdPrefix: "a/"},
			expected{etcdPrefix: pointer("/a")},
		}, {
			"Multiple slashes in etcd prefix",
			input{etcdPrefix: "/a/b/"},
			expected{etcdPrefix: pointer("/a/b")},
		}, {
			"Provide additional platform options",
			input{
				options: []PlatformOption{
					PlatformOptions.OpenTelemetry("testName", 0.5, tracetest.NewInMemoryExporter(), tracetest.NewInMemoryExporter()),
					PlatformOptions.ServerOptions(
						ServerOptions.GRPCOptions(grpc.ConnectionTimeout(time.Second)),
						ServerOptions.GRPCOptions(grpc.ConnectionTimeout(time.Second)),
					),
					PlatformOptions.ClientOptions(
						ClientOptions.GRPCOptions(grpc.WithIdleTimeout(time.Second)),
						ClientOptions.GRPCOptions(grpc.WithIdleTimeout(time.Second)),
					),
				},
			},
			expected{
				otelServiceName:  pointer("testName"),
				otelSampleRate:   pointer(0.5),
				otelExportersLen: pointer(2),
				serverOptionsLen: pointer(2),
				clientOptionsLen: pointer(2),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rpcp, err := New(tt.input.etcdPrefix, &etcd.Client{}, tt.input.options...)
			if err != nil {
				t.Fatalf("New() failed: %v", err)
			}

			if rpcp.etcdClient == nil {
				t.Error("etcdClient is nil")
			}

			if tt.expected.etcdPrefix != nil {
				if rpcp.etcdPrefix != *tt.expected.etcdPrefix {
					t.Errorf("etcdPrefix = %v, want: %v", rpcp.etcdPrefix, *tt.expected.etcdPrefix)
				}
			}

			if tt.expected.otelServiceName != nil {
				if rpcp.config.OpenTelemetry.ServiceName != *tt.expected.otelServiceName {
					t.Errorf("OpenTelemetry.ServiceName = %v, want: %v", rpcp.config.OpenTelemetry.ServiceName, *tt.expected.otelServiceName)
				}
			}

			if tt.expected.otelSampleRate != nil {
				if rpcp.config.OpenTelemetry.SampleRate != *tt.expected.otelSampleRate {
					t.Errorf("OpenTelemetry.SampleRate = %v, want: %v", rpcp.config.OpenTelemetry.SampleRate, *tt.expected.otelSampleRate)
				}
			}

			if tt.expected.otelExportersLen != nil {
				if len(rpcp.config.OpenTelemetry.Exporters) != *tt.expected.otelExportersLen {
					t.Errorf("OpenTelemetry.Exporters length = %v, want: %v", len(rpcp.config.OpenTelemetry.Exporters), *tt.expected.otelExportersLen)
				}
			}

			if tt.expected.serverOptionsLen != nil {
				if len(rpcp.config.ServerOptions) != *tt.expected.serverOptionsLen {
					t.Errorf("ServerOptions length = %v, want: %v", len(rpcp.config.ServerOptions), *tt.expected.serverOptionsLen)
				}
			}

			if tt.expected.clientOptionsLen != nil {
				if len(rpcp.config.ClientOptions) != *tt.expected.clientOptionsLen {
					t.Errorf("ClientOptions length = %v, want: %v", len(rpcp.config.ClientOptions), *tt.expected.clientOptionsLen)
				}
			}
		})
	}
}

func TestLookup(t *testing.T) {
	t.Parallel()

	etcdClient := getEtcdClient(t)
	defer etcdClient.Close()

	rpcp, err := New("rpcplatform", etcdClient)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	attrs := NewAttributes()
	attrs.BalancerWeight = 10
	attrs.BalancerPriority = 20

	serverName := "testLookup"
	publicAddr := "1.2.3.4:56789"

	server, err := rpcp.NewServer(serverName, "localhost:",
		ServerOptions.Attributes(attrs), ServerOptions.PublicAddr(publicAddr),
	)

	if err != nil {
		t.Fatalf("NewServer() failed: %v", err)
	}

	defer server.Server().Stop()

	go func() {
		if err := server.Serve(); err != nil {
			t.Errorf("Serve() failed: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	lookupChan, err := rpcp.Lookup(ctx, serverName, true)
	if err != nil {
		t.Fatalf("Lookup() failed: %v", err)
	}

	for infoMap := range lookupChan {
		info, ok := infoMap[server.ID()]
		if !ok {
			continue
		}

		if info.Address != publicAddr {
			t.Errorf("Address = %v, want: %v", info.Address, publicAddr)
		}

		if !slices.Equal(attributes.Values(info.Attributes), attributes.Values(attrs)) {
			t.Errorf("Attributes = %+v, want: %+v", info.Attributes, attrs)
		}

		return
	}

	t.Errorf("channel closed by timeout or unexpectedly")
}

func TestNewClient(t *testing.T) {
	t.Parallel()

	insecureTransport := PlatformOptions.ClientOptions(
		ClientOptions.GRPCOptions(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)

	type input struct {
		etcdPrefix      string
		target          string
		platformOptions []PlatformOption
		clientOptions   []ClientOption
	}

	type expected struct {
		target           *string
		maxActiveServers *int
		grpcOptionsLen   *int
	}

	tests := []struct {
		name     string
		input    input
		expected expected
	}{
		{
			"Target name without etcd prefix",
			input{
				target:          "testNewServer",
				platformOptions: []PlatformOption{insecureTransport},
			},
			expected{
				target: pointer("/testNewServer/"),
			},
		}, {
			"Target name with etcd prefix",
			input{
				etcdPrefix:      "rpcplatform",
				target:          "testNewServer",
				platformOptions: []PlatformOption{insecureTransport},
			},
			expected{
				target: pointer("/rpcplatform/testNewServer/"),
			},
		}, {
			"Provide MaxActiveServers option",
			input{
				target:          "testNewServer",
				platformOptions: []PlatformOption{insecureTransport},
				clientOptions: []ClientOption{
					ClientOptions.MaxActiveServers(10),
				},
			},
			expected{
				maxActiveServers: pointer(10),
			},
		}, {
			"Provide options that gRPC relies on",
			input{
				target: "testNewServer",
				platformOptions: []PlatformOption{
					insecureTransport,
					PlatformOptions.OpenTelemetry("testName", 0.5, tracetest.NewInMemoryExporter()),
					PlatformOptions.ClientOptions(
						ClientOptions.GRPCOptions(grpc.WithIdleTimeout(time.Second), grpc.WithIdleTimeout(time.Second)),
					),
				},
				clientOptions: []ClientOption{
					ClientOptions.GRPCOptions(grpc.WithIdleTimeout(time.Second), grpc.WithIdleTimeout(time.Second)),
				},
			},
			expected{
				grpcOptionsLen: pointer(6 + 2), // NewClient adds 2 additional options
			},
		},
	}

	etcdClient := getEtcdClient(t)
	t.Cleanup(func() { etcdClient.Close() })

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rpcp, err := New(tt.input.etcdPrefix, etcdClient, tt.input.platformOptions...)
			if err != nil {
				t.Fatalf("New() failed: %v", err)
			}

			client, err := rpcp.NewClient(tt.input.target, tt.input.clientOptions...)
			if err != nil {
				t.Fatalf("NewClient() failed: %v", err)
			}

			if client.ID() == "" {
				t.Error("client ID is empty")
			}

			if client.Client() == nil {
				t.Error("grpc.Client is nil")
			}

			if tt.expected.target != nil {
				if client.target != *tt.expected.target {
					t.Errorf("target = %v, want: %v", client.target, *tt.expected.target)
				}
			}

			if tt.expected.maxActiveServers != nil {
				if client.config.MaxActiveServers != *tt.expected.maxActiveServers {
					t.Errorf("MaxActiveServers = %v, want: %v", client.config.MaxActiveServers, *tt.expected.maxActiveServers)
				}
			}

			if tt.expected.grpcOptionsLen != nil {
				if len(client.config.GRPCOptions) != *tt.expected.grpcOptionsLen {
					t.Errorf("GRPCOptions length = %v, want: %v", len(client.config.GRPCOptions), *tt.expected.grpcOptionsLen)
				}
			}
		})
	}
}

func TestNewServer(t *testing.T) {
	t.Parallel()

	type input struct {
		etcdPrefix      string
		name            string
		addr            string
		platformOptions []PlatformOption
		serverOptions   []ServerOption
	}

	type expected struct {
		name           *string
		addr           bool
		publicAddr     *string
		attributes     *Attributes
		grpcOptionsLen *int
	}

	tests := []struct {
		name     string
		input    input
		expected expected
	}{
		{
			"Server name without etcd prefix",
			input{
				name: "testNewServer",
			},
			expected{
				name: pointer("/testNewServer"),
			},
		}, {
			"Server name with etcd prefix",
			input{
				etcdPrefix: "rpcplatform",
				name:       "testNewServer",
			},
			expected{
				name: pointer("/rpcplatform/testNewServer"),
			},
		}, {
			"Listen on 127.0.0.1 with port 0",
			input{
				name: "testNewServer",
				addr: "127.0.0.1:0",
			},
			expected{
				addr: true,
			},
		}, {
			"Provide PublicAddr option",
			input{
				name: "testNewServer",
				serverOptions: []ServerOption{
					ServerOptions.PublicAddr("1.2.3.4:56789"),
				},
			},
			expected{
				publicAddr: pointer("1.2.3.4:56789"),
			},
		}, {
			"Provide Attributes option",
			input{
				name: "testNewServer",
				serverOptions: []ServerOption{
					ServerOptions.Attributes(&Attributes{
						BalancerWeight:   10,
						BalancerPriority: 20,
					}),
				},
			},
			expected{
				attributes: &Attributes{
					BalancerWeight:   10,
					BalancerPriority: 20,
				},
			},
		}, {
			"Provide options that gRPC relies on",
			input{
				name: "testNewServer",
				platformOptions: []PlatformOption{
					PlatformOptions.OpenTelemetry("testName", 0.5, tracetest.NewInMemoryExporter()),
					PlatformOptions.ServerOptions(
						ServerOptions.GRPCOptions(grpc.ConnectionTimeout(time.Second), grpc.ConnectionTimeout(time.Second)),
					),
				},
				serverOptions: []ServerOption{
					ServerOptions.GRPCOptions(grpc.ConnectionTimeout(time.Second), grpc.ConnectionTimeout(time.Second)),
				},
			},
			expected{
				grpcOptionsLen: pointer(5),
			},
		},
	}

	etcdClient := getEtcdClient(t)
	t.Cleanup(func() { etcdClient.Close() })

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rpcp, err := New(tt.input.etcdPrefix, etcdClient, tt.input.platformOptions...)
			if err != nil {
				t.Fatalf("New() failed: %v", err)
			}

			server, err := rpcp.NewServer(tt.input.name, tt.input.addr, tt.input.serverOptions...)
			if err != nil {
				t.Fatalf("NewServer() failed: %v", err)
			}

			if server.listener == nil {
				t.Error("listener is nil")
			}

			if server.ID() == "" {
				t.Error("server ID is empty")
			}

			if server.Server() == nil {
				t.Error("grpc.Server is nil")
			}

			if tt.expected.name != nil {
				if server.name != *tt.expected.name {
					t.Errorf("name = %v, want: %v", server.name, *tt.expected.name)
				}
			}

			if tt.expected.publicAddr != nil {
				if server.config.PublicAddr != *tt.expected.publicAddr {
					t.Errorf("PublicAddr = %v, want: %v", server.config.PublicAddr, *tt.expected.publicAddr)
				}
			}

			if tt.expected.attributes != nil {
				if !slices.Equal(attributes.Values(server.config.Attributes), attributes.Values(tt.expected.attributes)) {
					t.Errorf("Attributes = %+v, want: %+v", server.config.Attributes, tt.expected.attributes)
				}
			}

			if tt.expected.grpcOptionsLen != nil {
				if len(server.config.GRPCOptions) != *tt.expected.grpcOptionsLen {
					t.Errorf("GRPCOptions length = %v, want: %v", len(server.config.GRPCOptions), *tt.expected.grpcOptionsLen)
				}
			}

			if tt.expected.addr {
				addr := strings.TrimSuffix(tt.input.addr, ":0") + ":"

				if !strings.HasPrefix(server.listener.Addr().String(), addr) {
					t.Errorf("listen addr = %v, want: %v", server.listener.Addr().String(), addr)
				}
			}

			go func() {
				time.Sleep(time.Second)
				server.Server().Stop()
			}()

			if err := server.Serve(); err != nil {
				t.Errorf("Serve() failed: %v", err)
			}
		})
	}
}

func getEtcdClient(t *testing.T) *etcd.Client {
	etcdAddr := os.Getenv("ETCD_ADDR")
	if etcdAddr == "" {
		t.Skipf("environment variable ETCD_ADDR is empty")
	}

	etcdClient, err := etcd.New(etcd.Config{
		Endpoints:   []string{etcdAddr},
		DialTimeout: 10 * time.Second,
	})

	if err != nil {
		t.Fatalf("etcd client creation failed: %v", err)
	}

	return etcdClient
}

func pointer[T any](v T) *T {
	return &v
}
