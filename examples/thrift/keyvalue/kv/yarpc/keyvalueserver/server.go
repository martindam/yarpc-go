// Code generated by thriftrw-plugin-yarpc
// @generated

// Copyright (c) 2016 Uber Technologies, Inc.
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

package keyvalueserver

import (
	"context"
	"go.uber.org/thriftrw/wire"
	"go.uber.org/yarpc/api/transport"
	"go.uber.org/yarpc/examples/thrift/keyvalue/kv"
	"go.uber.org/yarpc/encoding/thrift"
	"go.uber.org/yarpc"
)

// Interface is the server-side interface for the KeyValue service.
type Interface interface {
	GetValue(
		ctx context.Context,
		reqMeta yarpc.ReqMeta,
		Key *string,
	) (string, yarpc.ResMeta, error)

	SetValue(
		ctx context.Context,
		reqMeta yarpc.ReqMeta,
		Key *string,
		Value *string,
	) (yarpc.ResMeta, error)
}

// New prepares an implementation of the KeyValue service for
// registration.
//
// 	handler := KeyValueHandler{}
// 	dispatcher.Register(keyvalueserver.New(handler))
func New(impl Interface, opts ...thrift.RegisterOption) []transport.Registrant {
	h := handler{impl}
	service := thrift.Service{
		Name: "KeyValue",
		Methods: map[string]thrift.UnaryHandler{
			"getValue": thrift.UnaryHandlerFunc(h.GetValue),
			"setValue": thrift.UnaryHandlerFunc(h.SetValue),
		},
		OnewayMethods: map[string]thrift.OnewayHandler{},
	}
	return thrift.BuildRegistrants(service, opts...)
}

type handler struct{ impl Interface }

func (h handler) GetValue(
	ctx context.Context,
	reqMeta yarpc.ReqMeta,
	body wire.Value,
) (thrift.Response, error) {
	var args kv.KeyValue_GetValue_Args
	if err := args.FromWire(body); err != nil {
		return thrift.Response{}, err
	}

	success, resMeta, err := h.impl.GetValue(ctx, reqMeta, args.Key)

	hadError := err != nil
	result, err := kv.KeyValue_GetValue_Helper.WrapResponse(success, err)

	var response thrift.Response
	if err == nil {
		response.IsApplicationError = hadError
		response.Meta = resMeta
		response.Body = result
	}
	return response, err
}

func (h handler) SetValue(
	ctx context.Context,
	reqMeta yarpc.ReqMeta,
	body wire.Value,
) (thrift.Response, error) {
	var args kv.KeyValue_SetValue_Args
	if err := args.FromWire(body); err != nil {
		return thrift.Response{}, err
	}

	resMeta, err := h.impl.SetValue(ctx, reqMeta, args.Key, args.Value)

	hadError := err != nil
	result, err := kv.KeyValue_SetValue_Helper.WrapResponse(err)

	var response thrift.Response
	if err == nil {
		response.IsApplicationError = hadError
		response.Meta = resMeta
		response.Body = result
	}
	return response, err
}