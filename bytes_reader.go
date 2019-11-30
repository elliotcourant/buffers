package buffers

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

type BytesReader interface {
	NextByte() byte
	NextBytes() []byte
	NextString() string
	NextShortString() string
	NextError() error
	NextInt8() int8
	NextInt16() int16
	NextInt32() int32
	NextInt64() int64
	NextUint8() uint8
	NextUint16() uint16
	NextUint32() uint32
	NextUint64() uint64
	NextBool() bool
	NextReflection(kind reflect.Kind) interface{}
}

func NewBytesReader(src []byte) BytesReader {
	return &bytesReader{
		data:   src,
		offset: 0,
	}
}

type bytesReader struct {
	data   []byte
	offset uint32
}

func (b *bytesReader) NextByte() byte {
	i := b.data[b.offset]
	b.offset++
	return i
}

func (b *bytesReader) NextBytes() []byte {
	length := b.NextInt32()
	if length == -1 {
		return nil
	}
	i := b.data[b.offset : b.offset+uint32(length)]
	b.offset += uint32(length)
	return i
}

func (b *bytesReader) NextShortBytes() []byte {
	length := b.NextUint8()
	if length == 0 {
		return nil
	}
	i := b.data[b.offset : b.offset+uint32(length)]
	b.offset += uint32(length)
	return i
}

func (b *bytesReader) NextError() error {
	eBytes := b.NextBytes()
	if eBytes == nil {
		return nil
	} else {
		return errors.New(string(eBytes))
	}
}

func (b *bytesReader) NextString() string {
	return string(b.NextBytes())
}

func (b *bytesReader) NextShortString() string {
	return string(b.NextShortBytes())
}

func (b *bytesReader) NextBool() bool {
	t := b.NextByte()
	return t == 1
}

func (b *bytesReader) NextInt8() int8 {
	return int8(b.NextUint8())
}

func (b *bytesReader) NextInt16() int16 {
	return int16(b.NextUint16())
}

func (b *bytesReader) NextInt32() int32 {
	return int32(b.NextUint32())
}

func (b *bytesReader) NextInt64() int64 {
	return int64(b.NextUint64())
}

func (b *bytesReader) NextUint8() uint8 {
	return b.NextByte()
}

func (b *bytesReader) NextUint16() uint16 {
	i := b.data[b.offset : b.offset+Uint16Size]
	b.offset += Uint16Size
	return binary.BigEndian.Uint16(i)
}

func (b *bytesReader) NextUint32() uint32 {
	i := b.data[b.offset : b.offset+Uint32Size]
	b.offset += Uint32Size
	return binary.BigEndian.Uint32(i)
}

func (b *bytesReader) NextUint64() uint64 {
	i := b.data[b.offset : b.offset+Uint64Size]
	b.offset += Uint64Size
	return binary.BigEndian.Uint64(i)
}

func (b *bytesReader) NextReflection(kind reflect.Kind) interface{} {
	switch kind {

	case reflect.String:
		return b.NextString()

	case reflect.Bool:
		return b.NextBool()

	case reflect.Uint8:
		return b.NextUint8()
	case reflect.Uint16:
		return b.NextUint16()
	case reflect.Uint:
		return uint(b.NextUint32())
	case reflect.Uint32:
		return b.NextUint32()
	case reflect.Uint64:
		return b.NextUint64()

	case reflect.Int8:
		return b.NextInt8()
	case reflect.Int16:
		return b.NextInt16()
	case reflect.Int:
		return int(b.NextInt32())
	case reflect.Int32:
		return b.NextInt32()
	case reflect.Int64:
		return b.NextInt64()

	default:
		panic(fmt.Sprintf("cannot read %s", kind))
	}
}
