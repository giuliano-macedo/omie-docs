package bundler

import "github.com/giuliano-macedo/omie-docs/internal/core"

type Args struct {
	Pages    []core.Page
	Home     core.HomePage
	FsWriter FSWriter
}

type Bundler interface {
	Bundle(args Args) error
}
