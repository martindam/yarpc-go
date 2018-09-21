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

package yarpcprotobuf

import (
	"context"

	"github.com/gogo/protobuf/proto"
	yarpc "go.uber.org/yarpc/v2"
)

type unaryHandler struct {
	handle     func(context.Context, proto.Message) (proto.Message, error)
	newRequest func() proto.Message
}

func newUnaryHandler(
	handle func(context.Context, proto.Message) (proto.Message, error),
	newRequest func() proto.Message,
) *unaryHandler {
	return &unaryHandler{handle, newRequest}
}

func (u *unaryHandler) Handle(ctx context.Context, req *yarpc.Request, buf *yarpc.Buffer) (*yarpc.Response, *yarpc.Buffer, error) {
	return nil, nil, nil
}

type streamHandler struct {
	handle func(*ServerStream) error
}

func newStreamHandler(handle func(*ServerStream) error) *streamHandler {
	return &streamHandler{handle}
}

func (s *streamHandler) HandleStream(stream *yarpc.ServerStream) error {
	ctx, call := yarpc.NewInboundCall(stream.Context(), yarpc.DisableResponseHeaders())
	if err := call.ReadFromRequest(stream.Request()); err != nil {
		return err
	}
	protoStream := &ServerStream{
		ctx:    ctx,
		stream: stream,
	}
	return s.handle(protoStream)
}

func getProtoRequest(ctx context.Context, req *yarpc.Request, newRequest func() proto.Message) (context.Context, *yarpc.InboundCall, proto.Message, error) {
	return nil, nil, nil, nil
}

// StreamHandlerParams contains the parameters for creating a new StreamHandler.
type StreamHandlerParams struct {
	Handle func(*ServerStream) error
}

// NewStreamHandler returns a new StreamHandler.
func NewStreamHandler(params StreamHandlerParams) yarpc.StreamHandler {
	return newStreamHandler(params.Handle)
}
