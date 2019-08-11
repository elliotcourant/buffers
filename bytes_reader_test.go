package buffers

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBytesReader_NextBytes(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var nilBytes []byte = nil
		buf := NewBytesBuffer()
		buf.Append(nilBytes...)
		result := buf.Bytes()
		fmt.Println(hex.Dump(result))
		rdr := NewBytesReader(result)
		resultNilBytes := rdr.NextBytes()
		assert.Nil(t, resultNilBytes)
	})
}
