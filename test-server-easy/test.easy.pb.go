package main

import (
	"github.com/VictoriaMetrics/easyproto"
)

var mp easyproto.MarshalerPool

type TestRequest struct {
	Foo uint64
}

func (r *TestRequest) Marshal(dst []byte) []byte {
	m := mp.Get()
	r.marshal(m.MessageMarshaler())
	dst = m.Marshal(dst)
	mp.Put(m)
	return dst
}

func (r *TestRequest) marshal(mm *easyproto.MessageMarshaler) {
	mm.AppendUint64(1, r.Foo)
}

func (ts *TestRequest) Unmarshal(src []byte) (err error) {
	ts.Foo = 0
	var fc easyproto.FieldContext
	for len(src) > 0 {
		src, _ = fc.NextField(src)
		switch fc.FieldNum {
		case 1:
			ts.Foo, _ = fc.Uint64()
		}
	}
	return nil
}

type TestResponse struct {
	Bar uint64
}

func (r *TestResponse) Marshal(dst []byte) []byte {
	m := mp.Get()
	r.marshal(m.MessageMarshaler())
	dst = m.Marshal(dst)
	mp.Put(m)
	return dst
}

func (r *TestResponse) marshal(mm *easyproto.MessageMarshaler) {
	mm.AppendUint64(1, r.Bar)
}

func (ts *TestResponse) Unmarshal(src []byte) (err error) {
	ts.Bar = 0
	var fc easyproto.FieldContext
	for len(src) > 0 {
		src, _ = fc.NextField(src)
		switch fc.FieldNum {
		case 1:
			ts.Bar, _ = fc.Uint64()
		}
	}
	return nil
}
