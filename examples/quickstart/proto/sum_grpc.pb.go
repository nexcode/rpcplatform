//
// Copyright 2022 RPCPlatform Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.30.2
// source: sum.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Sum_Sum_FullMethodName = "/Sum/Sum"
)

// SumClient is the client API for Sum service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SumClient interface {
	Sum(ctx context.Context, in *SumRequest, opts ...grpc.CallOption) (*SumResponse, error)
}

type sumClient struct {
	cc grpc.ClientConnInterface
}

func NewSumClient(cc grpc.ClientConnInterface) SumClient {
	return &sumClient{cc}
}

func (c *sumClient) Sum(ctx context.Context, in *SumRequest, opts ...grpc.CallOption) (*SumResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SumResponse)
	err := c.cc.Invoke(ctx, Sum_Sum_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SumServer is the server API for Sum service.
// All implementations must embed UnimplementedSumServer
// for forward compatibility.
type SumServer interface {
	Sum(context.Context, *SumRequest) (*SumResponse, error)
	mustEmbedUnimplementedSumServer()
}

// UnimplementedSumServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSumServer struct{}

func (UnimplementedSumServer) Sum(context.Context, *SumRequest) (*SumResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Sum not implemented")
}
func (UnimplementedSumServer) mustEmbedUnimplementedSumServer() {}
func (UnimplementedSumServer) testEmbeddedByValue()             {}

// UnsafeSumServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SumServer will
// result in compilation errors.
type UnsafeSumServer interface {
	mustEmbedUnimplementedSumServer()
}

func RegisterSumServer(s grpc.ServiceRegistrar, srv SumServer) {
	// If the following call pancis, it indicates UnimplementedSumServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Sum_ServiceDesc, srv)
}

func _Sum_Sum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SumRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SumServer).Sum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Sum_Sum_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SumServer).Sum(ctx, req.(*SumRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Sum_ServiceDesc is the grpc.ServiceDesc for Sum service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Sum_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Sum",
	HandlerType: (*SumServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Sum",
			Handler:    _Sum_Sum_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sum.proto",
}
