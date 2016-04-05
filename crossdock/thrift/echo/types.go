// Code generated by thriftrw

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

package echo

import (
	"fmt"
	"github.com/thriftrw/thriftrw-go/wire"
	"strings"
)

type Ping struct{ Beep string }

func (v *Ping) ToWire() wire.Value {
	var fields [1]wire.Field
	i := 0
	fields[i] = wire.Field{ID: 1, Value: wire.NewValueString(v.Beep)}
	i++
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]})
}

func (v *Ping) FromWire(w wire.Value) error {
	var err error
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TBinary {
				v.Beep, err = field.Value.GetString(), error(nil)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (v *Ping) String() string {
	var fields [1]string
	i := 0
	fields[i] = fmt.Sprintf("Beep: %v", v.Beep)
	i++
	return fmt.Sprintf("Ping{%v}", strings.Join(fields[:i], ", "))
}

type Pong struct{ Boop string }

func (v *Pong) ToWire() wire.Value {
	var fields [1]wire.Field
	i := 0
	fields[i] = wire.Field{ID: 1, Value: wire.NewValueString(v.Boop)}
	i++
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]})
}

func (v *Pong) FromWire(w wire.Value) error {
	var err error
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TBinary {
				v.Boop, err = field.Value.GetString(), error(nil)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (v *Pong) String() string {
	var fields [1]string
	i := 0
	fields[i] = fmt.Sprintf("Boop: %v", v.Boop)
	i++
	return fmt.Sprintf("Pong{%v}", strings.Join(fields[:i], ", "))
}