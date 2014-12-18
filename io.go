package gomock

import (
	"io"
	"strings"
)

func NewReadCloser(body string) *ReadCloser {
	rc := &ReadCloser{
		Reader: strings.NewReader(body),
	}
	return rc
}

type ReadCloser struct {
	io.Reader
}

func (f ReadCloser) Close() error {
	return nil
}
