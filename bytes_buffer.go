package buffers

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
)

const (
	Uint8Size  = 1
	Uint16Size = 2
	Uint32Size = 4
	Uint64Size = 8
)

type BytesBuffer interface {
	Append(bytes ...byte)
	AppendRaw(bytes []byte)
	AppendByte(byte byte)
	AppendString(str string)
	AppendShortString(str string)
	AppendError(err error)
	AppendUint8(item uint8)
	AppendUint16(item uint16)
	AppendUint32(item uint32)
	AppendUint64(item uint64)
	AppendInt16(item int16)
	AppendInt32(item int32)
	AppendInt64(item int64)
	AppendBool(item bool)
	AppendNil32()
	AppendBuffer(buffer BytesBuffer)
	AppendReflection(value reflect.Value)

	Bytes() []byte
}

func NewBytesBuffer() BytesBuffer {
	return &bytesBuffer{
		buf: make([]byte, 0),
	}
}

type bytesBuffer struct {
	buf []byte
}

func (b *bytesBuffer) AppendBuffer(buffer BytesBuffer) {
	b.Append(buffer.Bytes()...)
}

func (b *bytesBuffer) AppendByte(byte byte) {
	b.AppendUint8(byte)
}

func (b *bytesBuffer) Append(bytes ...byte) {
	if bytes == nil {
		b.AppendInt32(-1)
		return
	}
	b.AppendInt32(int32(len(bytes)))
	b.buf = append(b.buf, bytes...)
}

func (b *bytesBuffer) AppendRaw(bytes []byte) {
	b.buf = append(b.buf, bytes...)
}

func (b *bytesBuffer) AppendShort(bytes ...byte) {
	if bytes == nil {
		b.AppendUint8(0)
		return
	}
	if len(bytes) > math.MaxUint8 {
		panic("too many bytes for short array")
	}
	b.AppendUint8(uint8(len(bytes)))
	b.buf = append(b.buf, bytes...)
}

func (b *bytesBuffer) AppendError(err error) {
	if err == nil {
		b.Append(nil...)
	} else {
		b.AppendString(err.Error())
	}
}

func (b *bytesBuffer) AppendString(str string) {
	b.Append([]byte(str)...)
}

func (b *bytesBuffer) AppendShortString(str string) {
	b.AppendShort([]byte(str)...)
}

func (b *bytesBuffer) AppendUint8(item uint8) {
	b.buf = append(b.buf, item)
}

func (b *bytesBuffer) AppendUint16(item uint16) {
	wp := len(b.buf)
	b.buf = append(b.buf, 0, 0)
	binary.BigEndian.PutUint16(b.buf[wp:], item)
}

func (b *bytesBuffer) AppendUint32(item uint32) {
	wp := len(b.buf)
	b.buf = append(b.buf, 0, 0, 0, 0)
	binary.BigEndian.PutUint32(b.buf[wp:], item)
}

func (b *bytesBuffer) AppendUint64(item uint64) {
	wp := len(b.buf)
	b.buf = append(b.buf, 0, 0, 0, 0, 0, 0, 0, 0)
	binary.BigEndian.PutUint64(b.buf[wp:], item)
}

func (b *bytesBuffer) AppendInt8(item int8) {
	b.AppendUint8(uint8(item))
}

func (b *bytesBuffer) AppendInt16(item int16) {
	b.AppendUint16(uint16(item))
}

func (b *bytesBuffer) AppendInt32(item int32) {
	b.AppendUint32(uint32(item))
}

func (b *bytesBuffer) AppendInt64(item int64) {
	b.AppendUint64(uint64(item))
}

func (b *bytesBuffer) AppendBool(item bool) {
	switch item {
	case true:
		b.AppendUint8(1)
	default:
		b.AppendUint8(0)
	}
}

func (b *bytesBuffer) AppendNil32() {
	b.AppendInt32(-1)
}

func (b *bytesBuffer) Bytes() []byte {
	return b.buf
}

func (b *bytesBuffer) AppendReflection(value reflect.Value) {
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	switch value.Kind() {

	case reflect.String:
		b.AppendString(value.String())

	case reflect.Bool:
		b.AppendBool(value.Bool())

	case reflect.Uint8:
		b.AppendUint8(uint8(value.Uint()))
	case reflect.Uint16:
		b.AppendUint16(uint16(value.Uint()))
	case reflect.Uint:
		b.AppendUint32(uint32(value.Uint()))
	case reflect.Uint32:
		b.AppendUint32(uint32(value.Uint()))
	case reflect.Uint64:
		b.AppendUint64(value.Uint())

	case reflect.Int8:
		b.AppendInt8(int8(value.Int()))
	case reflect.Int16:
		b.AppendInt16(int16(value.Int()))
	case reflect.Int:
		b.AppendInt32(int32(value.Int()))
	case reflect.Int32:
		b.AppendInt32(int32(value.Int()))
	case reflect.Int64:
		b.AppendInt64(value.Int())

	default:
		panic(fmt.Sprintf("cannot append %s", value.Kind()))
	}
}
