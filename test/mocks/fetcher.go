package mocks

import (
	"io"

	"github.com/stretchr/testify/mock"
)

type FetcherMock struct {
	mock.Mock
}

func argsGetReaderCloser(args mock.Arguments, index int) io.ReadCloser {
	res := args.Get(index)
	if res == nil {
		return nil
	}
	return res.(io.ReadCloser)
}

func (m *FetcherMock) GetInitialPage() (io.ReadCloser, error) {
	args := m.Called()
	return argsGetReaderCloser(args, 0), args.Error(1)
}

func (m *FetcherMock) GetPage(url string) (io.ReadCloser, error) {
	args := m.Called(url)
	return argsGetReaderCloser(args, 0), args.Error(1)
}

func (m *FetcherMock) InitialPageUrl() string {
	args := m.Called()
	return args.String(0)
}
