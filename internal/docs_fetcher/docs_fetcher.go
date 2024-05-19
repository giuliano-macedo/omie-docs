package docs_fetcher

import "io"

type DocsFetcher interface {
	GetInitialPage() (io.ReadCloser, error)
	GetPage(url string) (io.ReadCloser, error)
	InitialPageUrl() string
}
