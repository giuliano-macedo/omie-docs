package bundler

import "io"

type FSWriter interface {
	Create(name string) (io.WriteCloser, error)
	MkdirAll(name string) error
	SaveJsonFile(name string, value interface{}) error
}
