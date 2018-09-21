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
	yarpc "go.uber.org/yarpc/v2"
	"go.uber.org/yarpc/v2/yarpcprocedure"
)

// BuildProceduresParams contains the parameters for BuildProcedures.
type BuildProceduresParams struct {
	ServiceName         string
	UnaryHandlerParams  []BuildProceduresUnaryHandlerParams
	StreamHandlerParams []BuildProceduresStreamHandlerParams
}

// BuildProceduresUnaryHandlerParams contains the parameters for a UnaryHandler for BuildProcedures.
type BuildProceduresUnaryHandlerParams struct {
	MethodName string
	Handler    yarpc.UnaryHandler
}

// BuildProceduresStreamHandlerParams contains the parameters for a StreamHandler for BuildProcedures.
type BuildProceduresStreamHandlerParams struct {
	MethodName string
	Handler    yarpc.StreamHandler
}

// BuildProcedures builds the transport.Procedures.
func BuildProcedures(params BuildProceduresParams) []yarpc.Procedure {
	procedures := make([]yarpc.Procedure, 0, 2*(len(params.UnaryHandlerParams)+len(params.StreamHandlerParams)))
	for _, unaryHandlerParams := range params.UnaryHandlerParams {
		procedures = append(
			procedures,
			yarpc.Procedure{
				Name:        yarpcprocedure.ToName(params.ServiceName, unaryHandlerParams.MethodName),
				HandlerSpec: yarpc.NewUnaryHandlerSpec(unaryHandlerParams.Handler),
				Encoding:    protoEncoding,
			},
			yarpc.Procedure{
				Name:        yarpcprocedure.ToName(params.ServiceName, unaryHandlerParams.MethodName),
				HandlerSpec: yarpc.NewUnaryHandlerSpec(unaryHandlerParams.Handler),
				Encoding:    jsonEncoding,
			},
		)
	}
	for _, streamHandlerParams := range params.StreamHandlerParams {
		procedures = append(
			procedures,
			yarpc.Procedure{
				Name:        yarpcprocedure.ToName(params.ServiceName, streamHandlerParams.MethodName),
				HandlerSpec: yarpc.NewStreamHandlerSpec(streamHandlerParams.Handler),
				Encoding:    protoEncoding,
			},
			yarpc.Procedure{
				Name:        yarpcprocedure.ToName(params.ServiceName, streamHandlerParams.MethodName),
				HandlerSpec: yarpc.NewStreamHandlerSpec(streamHandlerParams.Handler),
				Encoding:    jsonEncoding,
			},
		)
	}
	return procedures
}
