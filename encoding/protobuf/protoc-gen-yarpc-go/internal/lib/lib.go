// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package lib contains the library code for protoc-gen-yarpc-go.
//
// It is split into a separate package so it can be called by the testing package.
package lib

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"go.uber.org/yarpc/internal/protoplugin"
)

const tmpl = `{{$packagePath := .GoPackage.Path}}{{$packageName := .GoPackage.Name}}
// Code generated by protoc-gen-yarpc-go
// source: {{.GetName}}
// DO NOT EDIT!

package {{$packageName}}
{{if .Services}}
import (
	{{range $i := .Imports}}{{if $i.Standard}}{{$i | printf "%s\n"}}{{end}}{{end}}

	{{range $i := .Imports}}{{if not $i.Standard}}{{$i | printf "%s\n"}}{{end}}{{end}}
){{end}}

{{if ne (len .Services) 0}}var _ = ioutil.NopCloser{{end}}

{{range $service := .Services}}
// {{$service.GetName}}YARPCClient is the YARPC client-side interface for the {{$service.GetName}} service.
type {{$service.GetName}}YARPCClient interface {
	{{range $method := unaryMethods $service}}{{$method.GetName}}(context.Context, *{{$method.RequestType.GoType $packagePath}}, ...yarpc.CallOption) (*{{$method.ResponseType.GoType $packagePath}}, error)
	{{end}}{{range $method := onewayMethods $service}}{{$method.GetName}}(context.Context, *{{$method.RequestType.GoType $packagePath}}, ...yarpc.CallOption) (yarpc.Ack, error)
	{{end}}{{range $method := clientStreamingMethods $service}}{{$method.GetName}}(context.Context, ...yarpc.CallOption) ({{$service.GetName}}Service{{$method.GetName}}YARPCClient, error)
	{{end}}{{range $method := serverStreamingMethods $service}}{{$method.GetName}}(context.Context, *{{$method.RequestType.GoType $packagePath}}, ...yarpc.CallOption) ({{$service.GetName}}Service{{$method.GetName}}YARPCClient, error)
	{{end}}{{range $method := clientServerStreamingMethods $service}}{{$method.GetName}}(context.Context, ...yarpc.CallOption) ({{$service.GetName}}Service{{$method.GetName}}YARPCClient, error)
	{{end}}
}

{{range $method := clientStreamingMethods $service}}
// {{$service.GetName}}Service{{$method.GetName}}YARPCClient sends {{$method.RequestType.GoType $packagePath}}s and receives the single {{$method.ResponseType.GoType $packagePath}} when sending is done.
type {{$service.GetName}}Service{{$method.GetName}}YARPCClient interface {
	Context() context.Context
	Send(*{{$method.RequestType.GoType $packagePath}}, ...yarpc.StreamOption) error
	CloseAndRecv(...yarpc.StreamOption) (*{{$method.ResponseType.GoType $packagePath}}, error)
}
{{end}}

{{range $method := serverStreamingMethods $service}}
// {{$service.GetName}}Service{{$method.GetName}}YARPCClient receives {{$method.ResponseType.GoType $packagePath}}s, returning io.EOF when the stream is complete.
type {{$service.GetName}}Service{{$method.GetName}}YARPCClient interface {
	Context() context.Context
	Recv(...yarpc.StreamOption) (*{{$method.ResponseType.GoType $packagePath}}, error)
	CloseSend(...yarpc.StreamOption) error
}
{{end}}

{{range $method := clientServerStreamingMethods $service}}
// {{$service.GetName}}Service{{$method.GetName}}YARPCClient sends {{$method.RequestType.GoType $packagePath}}s and receives {{$method.ResponseType.GoType $packagePath}}s, returning io.EOF when the stream is complete.
type {{$service.GetName}}Service{{$method.GetName}}YARPCClient interface {
	Context() context.Context
	Send(*{{$method.RequestType.GoType $packagePath}}, ...yarpc.StreamOption) error
	Recv(...yarpc.StreamOption) (*{{$method.ResponseType.GoType $packagePath}}, error)
	CloseSend(...yarpc.StreamOption) error
}
{{end}}

// New{{$service.GetName}}YARPCClient builds a new YARPC client for the {{$service.GetName}} service.
func New{{$service.GetName}}YARPCClient(clientConfig transport.ClientConfig, options ...protobuf.ClientOption) {{$service.GetName}}YARPCClient {
	return &_{{$service.GetName}}YARPCCaller{protobuf.NewStreamClient(
		protobuf.ClientParams{
			ServiceName: "{{trimPrefixPeriod $service.FQSN}}",
			ClientConfig: clientConfig,
			Options: options,
		},
	)}
}

// {{$service.GetName}}YARPCServer is the YARPC server-side interface for the {{$service.GetName}} service.
type {{$service.GetName}}YARPCServer interface {
	{{range $method := unaryMethods $service}}{{$method.GetName}}(context.Context, *{{$method.RequestType.GoType $packagePath}}) (*{{$method.ResponseType.GoType $packagePath}}, error)
	{{end}}{{range $method := onewayMethods $service}}{{$method.GetName}}(context.Context, *{{$method.RequestType.GoType $packagePath}}) error
	{{end}}{{range $method := clientStreamingMethods $service}}{{$method.GetName}}({{$service.GetName}}Service{{$method.GetName}}YARPCServer) (*{{$method.ResponseType.GoType $packagePath}}, error)
	{{end}}{{range $method := serverStreamingMethods $service}}{{$method.GetName}}(*{{$method.RequestType.GoType $packagePath}}, {{$service.GetName}}Service{{$method.GetName}}YARPCServer) error
	{{end}}{{range $method := clientServerStreamingMethods $service}}{{$method.GetName}}({{$service.GetName}}Service{{$method.GetName}}YARPCServer) error
	{{end}}
}

{{range $method := clientStreamingMethods $service}}
// {{$service.GetName}}Service{{$method.GetName}}YARPCServer receives {{$method.RequestType.GoType $packagePath}}s.
type {{$service.GetName}}Service{{$method.GetName}}YARPCServer interface {
	Context() context.Context
	Recv(...yarpc.StreamOption) (*{{$method.RequestType.GoType $packagePath}}, error)
}
{{end}}

{{range $method := serverStreamingMethods $service}}
// {{$service.GetName}}Service{{$method.GetName}}YARPCServer sends {{$method.ResponseType.GoType $packagePath}}s.
type {{$service.GetName}}Service{{$method.GetName}}YARPCServer interface {
	Context() context.Context
	Send(*{{$method.ResponseType.GoType $packagePath}}, ...yarpc.StreamOption) error
}
{{end}}

{{range $method := clientServerStreamingMethods $service}}
// {{$service.GetName}}Service{{$method.GetName}}YARPCServer receives {{$method.RequestType.GoType $packagePath}}s and sends {{$method.ResponseType.GoType $packagePath}}.
type {{$service.GetName}}Service{{$method.GetName}}YARPCServer interface {
	Context() context.Context
	Recv(...yarpc.StreamOption) (*{{$method.RequestType.GoType $packagePath}}, error)
	Send(*{{$method.ResponseType.GoType $packagePath}}, ...yarpc.StreamOption) error
}
{{end}}

// Build{{$service.GetName}}YARPCProcedures prepares an implementation of the {{$service.GetName}} service for YARPC registration.
func Build{{$service.GetName}}YARPCProcedures(server {{$service.GetName}}YARPCServer) []transport.Procedure {
	handler := &_{{$service.GetName}}YARPCHandler{server}
	return protobuf.BuildProcedures(
		protobuf.BuildProceduresParams{
			ServiceName: "{{trimPrefixPeriod $service.FQSN}}",
			UnaryHandlerParams: []protobuf.BuildProceduresUnaryHandlerParams{
			{{range $method := unaryMethods $service}}{
					MethodName: "{{$method.GetName}}",
					Handler: protobuf.NewUnaryHandler(
						protobuf.UnaryHandlerParams{
							Handle: handler.{{$method.GetName}},
							NewRequest: new{{$service.GetName}}Service{{$method.GetName}}YARPCRequest,
						},
					),
				},
			{{end}}
			},
			OnewayHandlerParams: []protobuf.BuildProceduresOnewayHandlerParams{
			{{range $method := onewayMethods $service}}{
					MethodName: "{{$method.GetName}}",
					Handler: protobuf.NewOnewayHandler(
						protobuf.OnewayHandlerParams{
							Handle: handler.{{$method.GetName}},
							NewRequest: new{{$service.GetName}}Service{{$method.GetName}}YARPCRequest,
						},
					),
				},
			{{end}}
			},
			StreamHandlerParams: []protobuf.BuildProceduresStreamHandlerParams{
			{{range $method := clientServerStreamingMethods $service}}{
					MethodName: "{{$method.GetName}}",
					Handler: protobuf.NewStreamHandler(
						protobuf.StreamHandlerParams{
							Handle: handler.{{$method.GetName}},
						},
					),
				},
			{{end}}
			{{range $method := serverStreamingMethods $service}}{
					MethodName: "{{$method.GetName}}",
					Handler: protobuf.NewStreamHandler(
						protobuf.StreamHandlerParams{
							Handle: handler.{{$method.GetName}},
						},
					),
				},
			{{end}}
			{{range $method := clientStreamingMethods $service}}{
					MethodName: "{{$method.GetName}}",
					Handler: protobuf.NewStreamHandler(
						protobuf.StreamHandlerParams{
							Handle: handler.{{$method.GetName}},
						},
					),
				},
			{{end}}
			},
		},
	)
}

// Fx{{$service.GetName}}YARPCClientParams defines the input
// for NewFx{{$service.GetName}}YARPCClient. It provides the
// paramaters to get a {{$service.GetName}}YARPCClient in an
// Fx application.
type Fx{{$service.GetName}}YARPCClientParams struct {
	fx.In

	Provider yarpc.ClientConfig
}

// Fx{{$service.GetName}}YARPCClientResult defines the output
// of NewFx{{$service.GetName}}YARPCClient. It provides a
// {{$service.GetName}}YARPCClient to an Fx application.
type Fx{{$service.GetName}}YARPCClientResult struct {
	fx.Out

	Client {{$service.GetName}}YARPCClient

	// We are using an fx.Out struct here instead of just returning a client
	// so that we can add more values or add named versions of the client in
	// the future without breaking any existing code.
}

// NewFx{{$service.GetName}}YARPCClient provides a {{$service.GetName}}YARPCClient
// to an Fx application using the given name for routing.
//
//  fx.Provide(
//    {{$packageName}}.NewFx{{$service.GetName}}YARPCClient("service-name"),
//    ...
//  )
func NewFx{{$service.GetName}}YARPCClient(name string, options ...protobuf.ClientOption) interface{} {
	return func(params Fx{{$service.GetName}}YARPCClientParams) Fx{{$service.GetName}}YARPCClientResult {
		return Fx{{$service.GetName}}YARPCClientResult{
			Client: New{{$service.GetName}}YARPCClient(params.Provider.ClientConfig(name), options...),
		}
	}
}

// Fx{{$service.GetName}}YARPCProceduresParams defines the input
// for NewFx{{$service.GetName}}YARPCProcedures. It provides the
// paramaters to get {{$service.GetName}}YARPCServer procedures in an
// Fx application.
type Fx{{$service.GetName}}YARPCProceduresParams struct {
	fx.In

	Server {{$service.GetName}}YARPCServer
}

// Fx{{$service.GetName}}YARPCProceduresResult defines the output
// of NewFx{{$service.GetName}}YARPCProcedures. It provides
// {{$service.GetName}}YARPCServer procedures to an Fx application.
//
// The procedures are provided to the "yarpcfx" value group.
// Dig 1.2 or newer must be used for this feature to work.
type Fx{{$service.GetName}}YARPCProceduresResult struct {
	fx.Out

	Procedures []transport.Procedure ` + "`" + `group:"yarpcfx"` + "`" + `
	ReflectionMeta reflection.ServerMeta ` + "`" + `group:"yarpcfx"` + "`" + `
}

// NewFx{{$service.GetName}}YARPCProcedures provides {{$service.GetName}}YARPCServer procedures to an Fx application.
// It expects a {{$service.GetName}}YARPCServer to be present in the container.
//
//  fx.Provide(
//    {{$packageName}}.NewFx{{$service.GetName}}YARPCProcedures(),
//    ...
//  )
func NewFx{{$service.GetName}}YARPCProcedures() interface{} {
	return func(params Fx{{$service.GetName}}YARPCProceduresParams) Fx{{$service.GetName}}YARPCProceduresResult {
		return Fx{{$service.GetName}}YARPCProceduresResult{
			Procedures: Build{{$service.GetName}}YARPCProcedures(params.Server),
			ReflectionMeta: reflection.ServerMeta{
				ServiceName: "{{trimPrefixPeriod $service.FQSN}}",
				FileDescriptors: {{ .File.FileDescriptorClosureVarName }},
			},
		}
	}
}

type _{{$service.GetName}}YARPCCaller struct {
	streamClient protobuf.StreamClient
}

{{range $method := unaryMethods $service}}
func (c *_{{$service.GetName}}YARPCCaller) {{$method.GetName}}(ctx context.Context, request *{{$method.RequestType.GoType $packagePath}}, options ...yarpc.CallOption) (*{{$method.ResponseType.GoType $packagePath}}, error) {
	responseMessage, err := c.streamClient.Call(ctx, "{{$method.GetName}}", request, new{{$service.GetName}}Service{{$method.GetName}}YARPCResponse, options...)
	if responseMessage == nil {
		return nil, err
	}
	response, ok := responseMessage.(*{{$method.ResponseType.GoType $packagePath}})
	if !ok {
		return nil, protobuf.CastError(empty{{$service.GetName}}Service{{$method.GetName}}YARPCResponse, responseMessage)
	}
	return response, err
}
{{end}}
{{range $method := onewayMethods $service}}
func (c *_{{$service.GetName}}YARPCCaller) {{$method.GetName}}(ctx context.Context, request *{{$method.RequestType.GoType $packagePath}}, options ...yarpc.CallOption) (yarpc.Ack, error) {
	return c.streamClient.CallOneway(ctx, "{{$method.GetName}}", request, options...)
}
{{end}}
{{range $method := clientStreamingMethods $service}}
func (c *_{{$service.GetName}}YARPCCaller) {{$method.GetName}}(ctx context.Context, options ...yarpc.CallOption) ({{$service.GetName}}Service{{$method.GetName}}YARPCClient, error) {
	stream, err := c.streamClient.CallStream(ctx, "{{$method.GetName}}", options...)
	if err != nil {
		return nil, err
	}
	return &_{{$service.GetName}}Service{{$method.GetName}}YARPCClient{stream: stream}, nil
}
{{end}}
{{range $method := serverStreamingMethods $service}}
func (c *_{{$service.GetName}}YARPCCaller) {{$method.GetName}}(ctx context.Context, request *{{$method.RequestType.GoType $packagePath}}, options ...yarpc.CallOption) ({{$service.GetName}}Service{{$method.GetName}}YARPCClient, error) {
	stream, err := c.streamClient.CallStream(ctx, "{{$method.GetName}}", options...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(request); err != nil {
		return nil, err
	}
	return &_{{$service.GetName}}Service{{$method.GetName}}YARPCClient{stream: stream}, nil
}
{{end}}
{{range $method := clientServerStreamingMethods $service}}
func (c *_{{$service.GetName}}YARPCCaller) {{$method.GetName}}(ctx context.Context, options ...yarpc.CallOption) ({{$service.GetName}}Service{{$method.GetName}}YARPCClient, error) {
	stream, err := c.streamClient.CallStream(ctx, "{{$method.GetName}}", options...)
	if err != nil {
		return nil, err
	}
	return &_{{$service.GetName}}Service{{$method.GetName}}YARPCClient{stream: stream}, nil
}
{{end}}

type _{{$service.GetName}}YARPCHandler struct {
	server {{$service.GetName}}YARPCServer
}

{{range $method := unaryMethods $service}}
func (h *_{{$service.GetName}}YARPCHandler) {{$method.GetName}}(ctx context.Context, requestMessage proto.Message) (proto.Message, error) {
	var request *{{$method.RequestType.GoType $packagePath}}
	var ok bool
	if requestMessage != nil {
		request, ok = requestMessage.(*{{$method.RequestType.GoType $packagePath}})
		if !ok {
			return nil, protobuf.CastError(empty{{$service.GetName}}Service{{$method.GetName}}YARPCRequest, requestMessage)
		}
	}
	response, err := h.server.{{$method.GetName}}(ctx, request)
	if response == nil {
		return nil, err
	}
	return response, err
}
{{end}}
{{range $method := onewayMethods $service}}
func (h *_{{$service.GetName}}YARPCHandler) {{$method.GetName}}(ctx context.Context, requestMessage proto.Message) error {
	var request *{{$method.RequestType.GoType $packagePath}}
	var ok bool
	if requestMessage != nil {
		request, ok = requestMessage.(*{{$method.RequestType.GoType $packagePath}})
		if !ok {
			return protobuf.CastError(empty{{$service.GetName}}Service{{$method.GetName}}YARPCRequest, requestMessage)
		}
	}
	return h.server.{{$method.GetName}}(ctx, request)
}
{{end}}
{{range $method := clientStreamingMethods $service}}
func (h *_{{$service.GetName}}YARPCHandler) {{$method.GetName}}(serverStream *protobuf.ServerStream) error {
	response, err := h.server.{{$method.GetName}}(&_{{$service.GetName}}Service{{$method.GetName}}YARPCServer{serverStream: serverStream})
	if err != nil {
		return err
	}
	return serverStream.Send(response)
}
{{end}}
{{range $method := serverStreamingMethods $service}}
func (h *_{{$service.GetName}}YARPCHandler) {{$method.GetName}}(serverStream *protobuf.ServerStream) error {
	requestMessage, err := serverStream.Receive(new{{$service.GetName}}Service{{$method.GetName}}YARPCRequest)
	if requestMessage == nil {
        return err
    }

	request, ok := requestMessage.(*{{$method.RequestType.GoType $packagePath}})
	if !ok {
		return protobuf.CastError(empty{{$service.GetName}}Service{{$method.GetName}}YARPCRequest, requestMessage)
	}
	return h.server.{{$method.GetName}}(request, &_{{$service.GetName}}Service{{$method.GetName}}YARPCServer{serverStream: serverStream})
}
{{end}}
{{range $method := clientServerStreamingMethods $service}}
func (h *_{{$service.GetName}}YARPCHandler) {{$method.GetName}}(serverStream *protobuf.ServerStream) error {
	return h.server.{{$method.GetName}}(&_{{$service.GetName}}Service{{$method.GetName}}YARPCServer{serverStream: serverStream})
}
{{end}}

{{range $method := clientStreamingMethods $service}}
type _{{$service.GetName}}Service{{$method.GetName}}YARPCClient struct {
	stream *protobuf.ClientStream
}

func (c *_{{$service.GetName}}Service{{$method.GetName}}YARPCClient) Context() context.Context {
	return c.stream.Context()
}

func (c *_{{$service.GetName}}Service{{$method.GetName}}YARPCClient) Send(request *{{$method.RequestType.GoType $packagePath}}, options ...yarpc.StreamOption) error {
	return c.stream.Send(request, options...)
}

func (c *_{{$service.GetName}}Service{{$method.GetName}}YARPCClient) CloseAndRecv(options ...yarpc.StreamOption) (*{{$method.ResponseType.GoType $packagePath}}, error) {
	if err := c.stream.Close(options...); err != nil {
		return nil, err
	}
	responseMessage, err := c.stream.Receive(new{{$service.GetName}}Service{{$method.GetName}}YARPCResponse, options...)
	if responseMessage == nil {
        return nil, err
    }
	response, ok := responseMessage.(*{{$method.ResponseType.GoType $packagePath}})
	if !ok {
		return nil, protobuf.CastError(empty{{$service.GetName}}Service{{$method.GetName}}YARPCResponse, responseMessage)
	}
	return response, err
}
{{end}}

{{range $method := serverStreamingMethods $service}}
type _{{$service.GetName}}Service{{$method.GetName}}YARPCClient struct {
	stream *protobuf.ClientStream
}

func (c *_{{$service.GetName}}Service{{$method.GetName}}YARPCClient) Context() context.Context {
	return c.stream.Context()
}

func (c *_{{$service.GetName}}Service{{$method.GetName}}YARPCClient) Recv(options ...yarpc.StreamOption) (*{{$method.ResponseType.GoType $packagePath}}, error) {
	responseMessage, err := c.stream.Receive(new{{$service.GetName}}Service{{$method.GetName}}YARPCResponse, options...)
	if responseMessage == nil {
        return nil, err
    }
	response, ok := responseMessage.(*{{$method.ResponseType.GoType $packagePath}})
	if !ok {
		return nil, protobuf.CastError(empty{{$service.GetName}}Service{{$method.GetName}}YARPCResponse, responseMessage)
	}
	return response, err
}

func (c *_{{$service.GetName}}Service{{$method.GetName}}YARPCClient) CloseSend(options ...yarpc.StreamOption) error {
	return c.stream.Close(options...)
}
{{end}}

{{range $method := clientServerStreamingMethods $service}}
type _{{$service.GetName}}Service{{$method.GetName}}YARPCClient struct {
	stream *protobuf.ClientStream
}

func (c *_{{$service.GetName}}Service{{$method.GetName}}YARPCClient) Context() context.Context {
	return c.stream.Context()
}

func (c *_{{$service.GetName}}Service{{$method.GetName}}YARPCClient) Send(request *{{$method.RequestType.GoType $packagePath}}, options ...yarpc.StreamOption) error {
	return c.stream.Send(request, options...)
}

func (c *_{{$service.GetName}}Service{{$method.GetName}}YARPCClient) Recv(options ...yarpc.StreamOption) (*{{$method.ResponseType.GoType $packagePath}}, error) {
	responseMessage, err := c.stream.Receive(new{{$service.GetName}}Service{{$method.GetName}}YARPCResponse, options...)
	if responseMessage == nil {
        return nil, err
    }
	response, ok := responseMessage.(*{{$method.ResponseType.GoType $packagePath}})
	if !ok {
		return nil, protobuf.CastError(empty{{$service.GetName}}Service{{$method.GetName}}YARPCResponse, responseMessage)
	}
	return response, err
}

func (c *_{{$service.GetName}}Service{{$method.GetName}}YARPCClient) CloseSend(options ...yarpc.StreamOption) error {
	return c.stream.Close(options...)
}
{{end}}

{{range $method := clientStreamingMethods $service}}
type _{{$service.GetName}}Service{{$method.GetName}}YARPCServer struct {
	serverStream *protobuf.ServerStream
}

func (s *_{{$service.GetName}}Service{{$method.GetName}}YARPCServer) Context() context.Context {
	return s.serverStream.Context()
}

func (s *_{{$service.GetName}}Service{{$method.GetName}}YARPCServer) Recv(options ...yarpc.StreamOption) (*{{$method.RequestType.GoType $packagePath}}, error) {
	requestMessage, err := s.serverStream.Receive(new{{$service.GetName}}Service{{$method.GetName}}YARPCRequest, options...)
	if requestMessage == nil {
        return nil, err
    }
	request, ok := requestMessage.(*{{$method.RequestType.GoType $packagePath}})
	if !ok {
		return nil, protobuf.CastError(empty{{$service.GetName}}Service{{$method.GetName}}YARPCRequest, requestMessage)
	}
	return request, err
}
{{end}}

{{range $method := serverStreamingMethods $service}}
type _{{$service.GetName}}Service{{$method.GetName}}YARPCServer struct {
	serverStream *protobuf.ServerStream
}

func (s *_{{$service.GetName}}Service{{$method.GetName}}YARPCServer) Context() context.Context {
	return s.serverStream.Context()
}

func (s *_{{$service.GetName}}Service{{$method.GetName}}YARPCServer) Send(response *{{$method.ResponseType.GoType $packagePath}}, options ...yarpc.StreamOption) error {
	return s.serverStream.Send(response, options...)
}
{{end}}

{{range $method := clientServerStreamingMethods $service}}
type _{{$service.GetName}}Service{{$method.GetName}}YARPCServer struct {
	serverStream *protobuf.ServerStream
}

func (s *_{{$service.GetName}}Service{{$method.GetName}}YARPCServer) Context() context.Context {
	return s.serverStream.Context()
}

func (s *_{{$service.GetName}}Service{{$method.GetName}}YARPCServer) Recv(options ...yarpc.StreamOption) (*{{$method.RequestType.GoType $packagePath}}, error) {
	requestMessage, err := s.serverStream.Receive(new{{$service.GetName}}Service{{$method.GetName}}YARPCRequest, options...)
	if requestMessage == nil {
        return nil, err
    }
	request, ok := requestMessage.(*{{$method.RequestType.GoType $packagePath}})
	if !ok {
		return nil, protobuf.CastError(empty{{$service.GetName}}Service{{$method.GetName}}YARPCRequest, requestMessage)
	}
	return request, err
}

func (s *_{{$service.GetName}}Service{{$method.GetName}}YARPCServer) Send(response *{{$method.ResponseType.GoType $packagePath}}, options ...yarpc.StreamOption) error {
	return s.serverStream.Send(response, options...)
}
{{end}}

{{range $method := $service.Methods}}
func new{{$service.GetName}}Service{{$method.GetName}}YARPCRequest() proto.Message {
	return &{{$method.RequestType.GoType $packagePath}}{}
}

func new{{$service.GetName}}Service{{$method.GetName}}YARPCResponse() proto.Message {
	return &{{$method.ResponseType.GoType $packagePath}}{}
}
{{end}}
var (
{{range $method := $service.Methods}}
	empty{{$service.GetName}}Service{{$method.GetName}}YARPCRequest = &{{$method.RequestType.GoType $packagePath}}{}
	empty{{$service.GetName}}Service{{$method.GetName}}YARPCResponse = &{{$method.ResponseType.GoType $packagePath}}{}{{end}}
)
{{end}}

var {{ .File.FileDescriptorClosureVarName }} = [][]byte{
	// {{ .Name }}
	{{ encodedFileDescriptor .File }},{{range $dependency := .TransitiveDependencies }}
	// {{ $dependency.Name }}
	{{ encodedFileDescriptor $dependency }},{{end}}
}

{{if .Services}}func init() { {{range $service := .Services}}
	yarpc.RegisterClientBuilder(
		func(clientConfig transport.ClientConfig, structField reflect.StructField) {{$service.GetName}}YARPCClient {
			return New{{$service.GetName}}YARPCClient(clientConfig, protobuf.ClientBuilderOptions(clientConfig, structField)...)
		},
	){{end}}
}{{end}}
`

// Runner is the Runner used for protoc-gen-yarpc-go.
var Runner = protoplugin.NewRunner(
	template.Must(template.New("tmpl").Funcs(
		template.FuncMap{
			"unaryMethods":                 unaryMethods,
			"onewayMethods":                onewayMethods,
			"clientStreamingMethods":       clientStreamingMethods,
			"serverStreamingMethods":       serverStreamingMethods,
			"clientServerStreamingMethods": clientServerStreamingMethods,
			"encodedFileDescriptor":        encodedFileDescriptor,
			"trimPrefixPeriod":             trimPrefixPeriod,
		}).Parse(tmpl)),
	checkTemplateInfo,
	[]string{
		"context",
		"io/ioutil",
		"reflect",
		"github.com/gogo/protobuf/proto",
		"go.uber.org/fx",
		"go.uber.org/yarpc",
		"go.uber.org/yarpc/api/transport",
		"go.uber.org/yarpc/encoding/protobuf",
		"go.uber.org/yarpc/encoding/protobuf/reflection",
	},
	func(file *protoplugin.File) (string, error) {
		name := file.GetName()
		return fmt.Sprintf("%s.pb.yarpc.go", strings.TrimSuffix(name, filepath.Ext(name))), nil
	},
	func(key string, value string) error {
		return nil
	},
)

func checkTemplateInfo(templateInfo *protoplugin.TemplateInfo) error {
	return nil
}

func unaryMethods(service *protoplugin.Service) ([]*protoplugin.Method, error) {
	methods := make([]*protoplugin.Method, 0, len(service.Methods))
	for _, method := range service.Methods {
		if !method.GetClientStreaming() && !method.GetServerStreaming() && method.ResponseType.FQMN() != ".uber.yarpc.Oneway" {
			methods = append(methods, method)
		}
	}
	return methods, nil
}

func onewayMethods(service *protoplugin.Service) ([]*protoplugin.Method, error) {
	methods := make([]*protoplugin.Method, 0, len(service.Methods))
	for _, method := range service.Methods {
		if !method.GetClientStreaming() && !method.GetServerStreaming() && method.ResponseType.FQMN() == ".uber.yarpc.Oneway" {
			methods = append(methods, method)
		}
	}
	return methods, nil
}

func clientStreamingMethods(service *protoplugin.Service) ([]*protoplugin.Method, error) {
	methods := make([]*protoplugin.Method, 0, len(service.Methods))
	for _, method := range service.Methods {
		if method.GetClientStreaming() && !method.GetServerStreaming() {
			methods = append(methods, method)
		}
	}
	return methods, nil
}

func serverStreamingMethods(service *protoplugin.Service) ([]*protoplugin.Method, error) {
	methods := make([]*protoplugin.Method, 0, len(service.Methods))
	for _, method := range service.Methods {
		if !method.GetClientStreaming() && method.GetServerStreaming() {
			methods = append(methods, method)
		}
	}
	return methods, nil
}

func clientServerStreamingMethods(service *protoplugin.Service) ([]*protoplugin.Method, error) {
	methods := make([]*protoplugin.Method, 0, len(service.Methods))
	for _, method := range service.Methods {
		if method.GetClientStreaming() && method.GetServerStreaming() {
			methods = append(methods, method)
		}
	}
	return methods, nil
}

func encodedFileDescriptor(f *protoplugin.File) (string, error) {
	b, err := f.SerializedFileDescriptor()
	if err != nil {
		return "", err
	}

	buf := make([]byte, 0, 0)
	w := bytes.NewBuffer(buf)
	w.WriteString("[]byte{\n")
	for len(b) > 0 {
		n := 16
		if n > len(b) {
			n = len(b)
		}

		s := ""
		for _, c := range b[:n] {
			s += fmt.Sprintf("0x%02x,", c)
		}
		w.WriteString(s)
		w.WriteString("\n")

		b = b[n:]
	}
	w.WriteString("}")
	return w.String(), nil
}

func trimPrefixPeriod(s string) string {
	return strings.TrimPrefix(s, ".")
}
