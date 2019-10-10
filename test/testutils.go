package test

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

type NoOpWriteCloser struct {
	io.Writer
}

func (no *NoOpWriteCloser) Close() error {
	return nil
}

type NoOpReadCloser struct {
	io.Reader
}

func (no *NoOpReadCloser) Close() error {
	return nil
}

type ErroringReader struct {
	Error error
}

func (reader *ErroringReader) Read(p []byte) (n int, err error) {
	return 0, reader.Error
}

func ExpectError(t *testing.T, expect error, actual error) {
	if actual != nil && expect == nil {
		t.Error(fmt.Sprintf("expecting 'nil' error but got '%s' error", actual))
		t.Fail()
	} else if actual == nil && expect != nil {
		t.Error(fmt.Sprintf("expecting '%s' error but got 'nil' error", expect))
		t.Fail()
	} else if actual != nil && expect != nil && actual.Error() != expect.Error() {
		t.Error(fmt.Sprintf("expecting '%s' error but got '%s' error", expect, actual))
		t.Fail()
	}
}

func ExpectErrorLike(t *testing.T, expect error, actual error) {
	if actual != nil && expect == nil {
		t.Error(fmt.Sprintf("expecting 'nil' error but got '%s' error", actual))
		t.Fail()
	} else if actual == nil && expect != nil {
		t.Error(fmt.Sprintf("expecting '%s' error but got 'nil' error", expect))
		t.Fail()
	} else if actual != nil && expect != nil && !strings.Contains(actual.Error(), expect.Error()) {
		t.Error(fmt.Sprintf("expecting '%s' error to contain '%s' error but it does not.", actual, expect))
		t.Fail()
	}
}

func ExpectString(t *testing.T, expect string, actual string) {
	if expect != actual {
		t.Error(fmt.Sprintf("expecting string '%s' but got '%s'", expect, actual))
		t.Fail()
	}
}
