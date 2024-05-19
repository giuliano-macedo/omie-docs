package docs_fetcher

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

type DocsHttpFetcher struct {
	client     *http.Client
	initialUrl string
}

const defaultInitialUrl = "https://developer.omie.com.br/service-list/"

func NewHttpFetcher(client *http.Client, initialPage string) DocsFetcher {
	if initialPage == "" {
		initialPage = defaultInitialUrl
	}
	return &DocsHttpFetcher{
		client:     client,
		initialUrl: initialPage,
	}
}

func (fetcher *DocsHttpFetcher) GetPage(pageUrl string) (io.ReadCloser, error) {
	resp, err := fetcher.client.Get(pageUrl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprint("Response not OK, ", resp.StatusCode))
	}
	return resp.Body, err
}

func (fetcher *DocsHttpFetcher) InitialPageUrl() string {
	return fetcher.initialUrl
}

func (fetcher *DocsHttpFetcher) GetInitialPage() (io.ReadCloser, error) {
	return fetcher.GetPage(fetcher.initialUrl)
}
