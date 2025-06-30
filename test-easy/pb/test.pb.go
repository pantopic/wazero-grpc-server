// Code generated manually. Feel free to edit.
// source: test.proto

package pb

import (
	"github.com/VictoriaMetrics/easyproto"
)

type Message interface {
	Unmarshal([]byte) error
	Marshal([]byte) []byte
}

var mp easyproto.MarshalerPool

type TestRequest struct {
	Foo uint64
}

func (tr *TestRequest) Marshal(dst []byte) []byte {
	m := mp.Get()
	tr.marshal(m.MessageMarshaler())
	dst = m.Marshal(dst)
	mp.Put(m)
	return dst
}

func (tr *TestRequest) marshal(mm *easyproto.MessageMarshaler) {
	mm.AppendUint64(1, tr.Foo)
}

func (tr *TestRequest) Unmarshal(src []byte) (err error) {
	tr.Foo = 0
	var fc easyproto.FieldContext
	for len(src) > 0 {
		src, _ = fc.NextField(src)
		switch fc.FieldNum {
		case 1:
			tr.Foo, _ = fc.Uint64()
		}
	}
	return nil
}

type TestResponse struct {
	Bar uint64
}

func (tr *TestResponse) Marshal(dst []byte) []byte {
	m := mp.Get()
	tr.marshal(m.MessageMarshaler())
	dst = m.Marshal(dst)
	mp.Put(m)
	return dst
}

func (tr *TestResponse) marshal(mm *easyproto.MessageMarshaler) {
	mm.AppendUint64(1, tr.Bar)
}

func (tr *TestResponse) Unmarshal(src []byte) (err error) {
	tr.Bar = 0
	var fc easyproto.FieldContext
	for len(src) > 0 {
		src, _ = fc.NextField(src)
		switch fc.FieldNum {
		case 1:
			tr.Bar, _ = fc.Uint64()
		}
	}
	return nil
}

type RetestRequest struct {
	Bar uint64
}

func (tr *RetestRequest) Marshal(dst []byte) []byte {
	m := mp.Get()
	tr.marshal(m.MessageMarshaler())
	dst = m.Marshal(dst)
	mp.Put(m)
	return dst
}

func (tr *RetestRequest) marshal(mm *easyproto.MessageMarshaler) {
	mm.AppendUint64(1, tr.Bar)
}

func (tr *RetestRequest) Unmarshal(src []byte) (err error) {
	tr.Bar = 0
	var fc easyproto.FieldContext
	for len(src) > 0 {
		src, _ = fc.NextField(src)
		switch fc.FieldNum {
		case 1:
			tr.Bar, _ = fc.Uint64()
		}
	}
	return nil
}

type RetestResponse struct {
	Foo uint64
}

func (tr *RetestResponse) Marshal(dst []byte) []byte {
	m := mp.Get()
	tr.marshal(m.MessageMarshaler())
	dst = m.Marshal(dst)
	mp.Put(m)
	return dst
}

func (tr *RetestResponse) marshal(mm *easyproto.MessageMarshaler) {
	mm.AppendUint64(1, tr.Foo)
}

func (tr *RetestResponse) Unmarshal(src []byte) (err error) {
	tr.Foo = 0
	var fc easyproto.FieldContext
	for len(src) > 0 {
		src, _ = fc.NextField(src)
		switch fc.FieldNum {
		case 1:
			tr.Foo, _ = fc.Uint64()
		}
	}
	return nil
}

type ClientStreamRequest struct {
	Foo2 uint64
}

func (tr *ClientStreamRequest) Marshal(dst []byte) []byte {
	m := mp.Get()
	tr.marshal(m.MessageMarshaler())
	dst = m.Marshal(dst)
	mp.Put(m)
	return dst
}

func (tr *ClientStreamRequest) marshal(mm *easyproto.MessageMarshaler) {
	mm.AppendUint64(2, tr.Foo2)
}

func (tr *ClientStreamRequest) Unmarshal(src []byte) (err error) {
	tr.Foo2 = 0
	var fc easyproto.FieldContext
	for len(src) > 0 {
		src, _ = fc.NextField(src)
		switch fc.FieldNum {
		case 2:
			tr.Foo2, _ = fc.Uint64()
		}
	}
	return nil
}

type ClientStreamResponse struct {
	Bar2 uint64
}

func (tr *ClientStreamResponse) Marshal(dst []byte) []byte {
	m := mp.Get()
	tr.marshal(m.MessageMarshaler())
	dst = m.Marshal(dst)
	mp.Put(m)
	return dst
}

func (tr *ClientStreamResponse) marshal(mm *easyproto.MessageMarshaler) {
	mm.AppendUint64(3, tr.Bar2)
}

func (tr *ClientStreamResponse) Unmarshal(src []byte) (err error) {
	tr.Bar2 = 0
	var fc easyproto.FieldContext
	for len(src) > 0 {
		src, _ = fc.NextField(src)
		switch fc.FieldNum {
		case 3:
			tr.Bar2, _ = fc.Uint64()
		}
	}
	return nil
}

type ServerStreamRequest struct {
	Foo3 uint64
}

func (tr *ServerStreamRequest) Marshal(dst []byte) []byte {
	m := mp.Get()
	tr.marshal(m.MessageMarshaler())
	dst = m.Marshal(dst)
	mp.Put(m)
	return dst
}

func (tr *ServerStreamRequest) marshal(mm *easyproto.MessageMarshaler) {
	mm.AppendUint64(4, tr.Foo3)
}

func (tr *ServerStreamRequest) Unmarshal(src []byte) (err error) {
	tr.Foo3 = 0
	var fc easyproto.FieldContext
	for len(src) > 0 {
		src, _ = fc.NextField(src)
		switch fc.FieldNum {
		case 4:
			tr.Foo3, _ = fc.Uint64()
		}
	}
	return nil
}

type ServerStreamResponse struct {
	Bar3 uint64
}

func (tr *ServerStreamResponse) Marshal(dst []byte) []byte {
	m := mp.Get()
	tr.marshal(m.MessageMarshaler())
	dst = m.Marshal(dst)
	mp.Put(m)
	return dst
}

func (tr *ServerStreamResponse) marshal(mm *easyproto.MessageMarshaler) {
	mm.AppendUint64(5, tr.Bar3)
}

func (tr *ServerStreamResponse) Unmarshal(src []byte) (err error) {
	tr.Bar3 = 0
	var fc easyproto.FieldContext
	for len(src) > 0 {
		src, _ = fc.NextField(src)
		switch fc.FieldNum {
		case 5:
			tr.Bar3, _ = fc.Uint64()
		}
	}
	return nil
}
