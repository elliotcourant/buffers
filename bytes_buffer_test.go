package buffers

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestBytesBuffer_Append(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		emptyBytes := make([]byte, 0)
		buf := NewBytesBuffer()
		buf.Append(emptyBytes...)
		result := buf.Bytes()
		fmt.Println(hex.Dump(result))
	})

	t.Run("nil", func(t *testing.T) {
		var nilBytes []byte = nil
		buf := NewBytesBuffer()
		buf.Append(nilBytes...)
		result := buf.Bytes()
		fmt.Println(hex.Dump(result))
	})
}
