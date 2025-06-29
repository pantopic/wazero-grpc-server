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
