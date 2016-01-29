package authn

import (
	"bytes"
	"io"
)

// Rewinder is a wrapper for an io.ReadCloser() that adds the ability to rewind
// the incoming data stream.  It does this by caching everything that's read
// from it, and starting over when the Rewind() method is called, answering
// read requests from its cache until the cache is exhausted, and then pulling
// "live" data from the Reader.
type Rewinder interface {
	Rewind()
	io.ReadCloser
}

type rewinder struct {
	buffer bytes.Buffer
	reader io.Reader
	read   int
}

// Rewind rewinds the input stream in the object, so that the next Read attempt
// will return data starting at the first byte that was ever read.
func (r *rewinder) Rewind() {
	r.read = 0
}

// Close fakes closing the reader.
func (r *rewinder) Close() error {
	r.Rewind()
	return nil
}

func (r *rewinder) Read(p []byte) (n int, err error) {
	// If we have enough data to satisfy the read attempt from the buffer,
	// just return the data and increment our read offset.
	if r.read+len(p) < r.buffer.Len() {
		n, err = bytes.NewReader(r.buffer.Bytes()[r.read:]).Read(p)
		if n > 0 {
			r.read += n
		}
		return n, err
	}
	// We don't have enough data.  Read what we've already buffered first.
	n2, err2 := bytes.NewReader(r.buffer.Bytes()[r.read:]).Read(p)
	if n2 > 0 {
		r.read += n2
	}
	if err2 != nil && err2 != io.EOF {
		return n2, err2
	}
	// Try to read some more.
	n3, err3 := r.reader.Read(p[n2:])
	if n3 > 0 {
		r.buffer.Write(p[n2 : n2+n3])
		r.read += n3
	}
	return n2 + n3, err3
}

func makeRewinder(reader io.Reader) Rewinder {
	return &rewinder{reader: reader}
}
