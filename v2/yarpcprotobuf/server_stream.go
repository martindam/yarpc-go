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

// ServerStream is a protobuf-specific server stream.
type ServerStream struct {
	ctx    context.Context
	stream *yarpc.ServerStream
}

// Context returns the context of the stream.
func (s *ServerStream) Context() context.Context {
	return s.ctx
}

// Receive will receive a protobuf message from the server stream.
func (s *ServerStream) Receive(newMessage func() proto.Message, options ...yarpc.StreamOption) (proto.Message, error) {
	return nil, nil
}

// Send will send a protobuf message to the server stream.
func (s *ServerStream) Send(message proto.Message, options ...yarpc.StreamOption) error {
	return nil
}
