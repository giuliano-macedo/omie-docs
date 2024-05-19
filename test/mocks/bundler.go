package mocks

import "github.com/giuliano-macedo/omie-docs/internal/bundler"

type BundlerMock struct {
	CalledBundleArgs bundler.Args
	ReturnError      error
}

func (m *BundlerMock) Bundle(args bundler.Args) error {
	m.CalledBundleArgs = args
	return m.ReturnError
}
