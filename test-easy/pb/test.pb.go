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

type TestBytesRequest struct {
	Key []byte
	Val []byte
}

func (tr *TestBytesRequest) Marshal(dst []byte) []byte {
	m := mp.Get()
	tr.marshal(m.MessageMarshaler())
	dst = m.Marshal(dst)
	mp.Put(m)
	return dst
}

func (tr *TestBytesRequest) marshal(mm *easyproto.MessageMarshaler) {
	mm.AppendBytes(1, tr.Key)
	mm.AppendBytes(2, tr.Val)
}

func (tr *TestBytesRequest) Unmarshal(src []byte) (err error) {
	tr.Key = tr.Key[:0]
	tr.Val = tr.Val[:0]
	var fc easyproto.FieldContext
	for len(src) > 0 {
		src, _ = fc.NextField(src)
		switch fc.FieldNum {
		case 1:
			tr.Key, _ = fc.Bytes()
		case 2:
			tr.Val, _ = fc.Bytes()
		}
	}
	return nil
}

type TestBytesResponse struct {
	Code uint64
	Data []byte
}

func (tr *TestBytesResponse) Marshal(dst []byte) []byte {
	m := mp.Get()
	tr.marshal(m.MessageMarshaler())
	dst = m.Marshal(dst)
	mp.Put(m)
	return dst
}

func (tr *TestBytesResponse) marshal(mm *easyproto.MessageMarshaler) {
	mm.AppendUint64(1, tr.Code)
	mm.AppendBytes(2, tr.Data)
}

func (tr *TestBytesResponse) Unmarshal(src []byte) (err error) {
	tr.Code = 0
	tr.Data = tr.Data[:0]
	var fc easyproto.FieldContext
	for len(src) > 0 {
		src, _ = fc.NextField(src)
		switch fc.FieldNum {
		case 1:
			tr.Code, _ = fc.Uint64()
		case 2:
			tr.Data, _ = fc.Bytes()
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

type BidirectionalStreamRequest struct {
	Foo4 uint64
}

func (tr *BidirectionalStreamRequest) Marshal(dst []byte) []byte {
	m := mp.Get()
	tr.marshal(m.MessageMarshaler())
	dst = m.Marshal(dst)
	mp.Put(m)
	return dst
}

func (tr *BidirectionalStreamRequest) marshal(mm *easyproto.MessageMarshaler) {
	mm.AppendUint64(6, tr.Foo4)
}

func (tr *BidirectionalStreamRequest) Unmarshal(src []byte) (err error) {
	tr.Foo4 = 0
	var fc easyproto.FieldContext
	for len(src) > 0 {
		src, _ = fc.NextField(src)
		switch fc.FieldNum {
		case 6:
			tr.Foo4, _ = fc.Uint64()
		}
	}
	return nil
}

type BidirectionalStreamResponse struct {
	Bar4 uint64
}

func (tr *BidirectionalStreamResponse) Marshal(dst []byte) []byte {
	m := mp.Get()
	tr.marshal(m.MessageMarshaler())
	dst = m.Marshal(dst)
	mp.Put(m)
	return dst
}

func (tr *BidirectionalStreamResponse) marshal(mm *easyproto.MessageMarshaler) {
	mm.AppendUint64(7, tr.Bar4)
}

func (tr *BidirectionalStreamResponse) Unmarshal(src []byte) (err error) {
	tr.Bar4 = 0
	var fc easyproto.FieldContext
	for len(src) > 0 {
		src, _ = fc.NextField(src)
		switch fc.FieldNum {
		case 7:
			tr.Bar4, _ = fc.Uint64()
		}
	}
	return nil
}
