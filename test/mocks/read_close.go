package mocks

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ReadCloseMock struct {
	reader      *strings.Reader
	errorOnRead error
	closeCalled bool
}

func (mock *ReadCloseMock) Close() error {
	mock.closeCalled = true
	return nil
}

func (mock *ReadCloseMock) AssertCloseCalled(t *testing.T) bool {
	return assert.Equal(t, true, mock.closeCalled)
}

func (mock *ReadCloseMock) Read(p []byte) (n int, err error) {
	if mock.errorOnRead != nil {
		return 0, mock.errorOnRead
	}
	return mock.reader.Read(p)
}

func NewReadCloseMock(data string) *ReadCloseMock {
	return &ReadCloseMock{reader: strings.NewReader(data)}
}

func NewErrReadCloseMock(errorOnRead error) *ReadCloseMock {
	return &ReadCloseMock{errorOnRead: errorOnRead}
}
