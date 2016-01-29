package authn

import (
	"bytes"
	"io"
	"testing"
)

func TestRewinder(t *testing.T) {
	s := "abcdefghijklmnopqrstuvwxyz"
	rewinder := makeRewinder(bytes.NewBufferString(s))
	for i := 1; i < len(s); i++ {
		b := bytes.NewBuffer([]byte{})
		j := 0
		for j < len(s) {
			p := make([]byte, i)
			n, err := rewinder.Read(p)
			b.Write(p[:n])
			j += n
			if n < len(p) {
				if err == io.EOF {
					break
				}
				t.Fatalf("Unexpected read error at %d (%d): %v", i, j, err)
				return
			}
		}
		if j != b.Len() || j == 0 {
			t.Fatalf("Unexpected length %d (expected %d)", j, b.Len())
			return
		}
		if !bytes.Equal(b.Bytes()[:j], bytes.NewBufferString(s).Bytes()[:j]) {
			t.Fatalf("Mismatch on data pass %d/%d)", i, j)
			return
		}
		rewinder.Rewind()
	}
}
