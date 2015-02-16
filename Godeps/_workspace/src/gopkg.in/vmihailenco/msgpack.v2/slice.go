package msgpack

import (
	"fmt"
	"io"
	"reflect"
)

func (e *Encoder) encodeBytesLen(l int) error {
	switch {
	case l < 256:
		if err := e.write1(bin8Code, uint64(l)); err != nil {
			return err
		}
	case l < 65536:
		if err := e.write2(bin16Code, uint64(l)); err != nil {
			return err
		}
	default:
		if err := e.write4(bin32Code, uint64(l)); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) encodeStrLen(l int) error {
	switch {
	case l < 32:
		if err := e.w.WriteByte(fixStrLowCode | uint8(l)); err != nil {
			return err
		}
	case l < 256:
		if err := e.write1(str8Code, uint64(l)); err != nil {
			return err
		}
	case l < 65536:
		if err := e.write2(str16Code, uint64(l)); err != nil {
			return err
		}
	default:
		if err := e.write4(str32Code, uint64(l)); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) EncodeString(v string) error {
	if err := e.encodeStrLen(len(v)); err != nil {
		return err
	}
	return e.writeString(v)
}

func (e *Encoder) EncodeBytes(v []byte) error {
	if v == nil {
		return e.EncodeNil()
	}
	if err := e.encodeBytesLen(len(v)); err != nil {
		return err
	}
	return e.write(v)
}

func (e *Encoder) EncodeSliceLen(l int) error {
	switch {
	case l < 16:
		if err := e.w.WriteByte(fixArrayLowCode | byte(l)); err != nil {
			return err
		}
	case l < 65536:
		if err := e.write2(array16Code, uint64(l)); err != nil {
			return err
		}
	default:
		if err := e.write4(array32Code, uint64(l)); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) encodeStringSlice(s []string) error {
	if s == nil {
		return e.EncodeNil()
	}
	if err := e.EncodeSliceLen(len(s)); err != nil {
		return err
	}
	for _, v := range s {
		if err := e.EncodeString(v); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) encodeSlice(v reflect.Value) error {
	if v.IsNil() {
		return e.EncodeNil()
	}
	return e.encodeArray(v)
}

func (e *Encoder) encodeArray(v reflect.Value) error {
	l := v.Len()
	if err := e.EncodeSliceLen(l); err != nil {
		return err
	}
	for i := 0; i < l; i++ {
		if err := e.EncodeValue(v.Index(i)); err != nil {
			return err
		}
	}
	return nil
}

func (d *Decoder) DecodeBytesLen() (int, error) {
	c, err := d.r.ReadByte()
	if err != nil {
		return 0, err
	}
	if c == nilCode {
		return -1, nil
	} else if c >= fixStrLowCode && c <= fixStrHighCode {
		return int(c & fixStrMask), nil
	}
	switch c {
	case str8Code, bin8Code:
		n, err := d.uint8()
		return int(n), err
	case str16Code, bin16Code:
		n, err := d.uint16()
		return int(n), err
	case str32Code, bin32Code:
		n, err := d.uint32()
		return int(n), err
	}
	return 0, fmt.Errorf("msgpack: invalid code %x decoding bytes length", c)
}

func (d *Decoder) DecodeBytes() ([]byte, error) {
	n, err := d.DecodeBytesLen()
	if err != nil {
		return nil, err
	}
	if n == -1 {
		return nil, nil
	}
	b := make([]byte, n)
	_, err = io.ReadFull(d.r, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (d *Decoder) bytesValue(value reflect.Value) error {
	v, err := d.DecodeBytes()
	if err != nil {
		return err
	}
	value.SetBytes(v)
	return nil
}

func (d *Decoder) DecodeString() (string, error) {
	n, err := d.DecodeBytesLen()
	if err != nil {
		return "", err
	}
	if n == -1 {
		return "", nil
	}
	b, err := d.readN(n)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (d *Decoder) stringValue(value reflect.Value) error {
	v, err := d.DecodeString()
	if err != nil {
		return err
	}
	value.SetString(v)
	return nil
}

func (d *Decoder) DecodeSliceLen() (int, error) {
	c, err := d.r.ReadByte()
	if err != nil {
		return 0, err
	}
	if c == nilCode {
		return -1, nil
	} else if c >= fixArrayLowCode && c <= fixArrayHighCode {
		return int(c & fixArrayMask), nil
	}
	switch c {
	case array16Code:
		n, err := d.uint16()
		return int(n), err
	case array32Code:
		n, err := d.uint32()
		return int(n), err
	}
	return 0, fmt.Errorf("msgpack: invalid code %x decoding array length", c)
}

func (d *Decoder) decodeIntoStrings(sp *[]string) error {
	n, err := d.DecodeSliceLen()
	if err != nil {
		return err
	}
	if n == -1 {
		return nil
	}
	s := *sp
	if s == nil || len(s) < n {
		s = make([]string, n)
	}
	for i := 0; i < n; i++ {
		v, err := d.DecodeString()
		if err != nil {
			return err
		}
		s[i] = v
	}
	*sp = s
	return nil
}

func (d *Decoder) DecodeSlice() ([]interface{}, error) {
	n, err := d.DecodeSliceLen()
	if err != nil {
		return nil, err
	}

	if n == -1 {
		return nil, nil
	}

	s := make([]interface{}, n)
	for i := 0; i < n; i++ {
		v, err := d.DecodeInterface()
		if err != nil {
			return nil, err
		}
		s[i] = v
	}

	return s, nil
}

func (d *Decoder) sliceValue(v reflect.Value) error {
	n, err := d.DecodeSliceLen()
	if err != nil {
		return err
	}

	if n == -1 {
		v.Set(reflect.Zero(v.Type()))
		return nil
	}

	if v.Len() < n || (v.Kind() == reflect.Slice && v.IsNil()) {
		v.Set(reflect.MakeSlice(v.Type(), n, n))
	}

	for i := 0; i < n; i++ {
		sv := v.Index(i)
		if err := d.DecodeValue(sv); err != nil {
			return err
		}
	}

	return nil
}

func (d *Decoder) strings() ([]string, error) {
	n, err := d.DecodeSliceLen()
	if err != nil {
		return nil, err
	}

	if n == -1 {
		return nil, nil
	}

	ss := make([]string, n)
	for i := 0; i < n; i++ {
		s, err := d.DecodeString()
		if err != nil {
			return nil, err
		}
		ss[i] = s
	}

	return ss, nil
}

func (d *Decoder) stringsValue(value reflect.Value) error {
	ss, err := d.strings()
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(ss))
	return nil
}
