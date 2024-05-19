package docs_fetcher_test

import (
	"io"
	"testing"

	"github.com/giuliano-macedo/omie-docs/internal/docs_fetcher"
	"github.com/giuliano-macedo/omie-docs/test/mocks"
	"github.com/stretchr/testify/require"
)

func TestHttpFetcher(t *testing.T) {
	httpClient := mocks.NewHttpClientMock()
	httpClient.OnRequest("https://developer.omie.com.br/service-list/", "initial page", 200)
	httpClient.OnRequest("https://example.com", "example page", 200)

	httpFetcher := docs_fetcher.NewHttpFetcher(httpClient.ToHttpClient(), "")

	t.Run("GetInitialPage", func(t *testing.T) {
		res, err := httpFetcher.GetInitialPage()
		require.NoError(t, err)
		require.Equal(t, readRes(t, res), "initial page")
	})

	t.Run("GetPage", func(t *testing.T) {
		res, err := httpFetcher.GetPage("https://example.com")
		require.NoError(t, err)
		require.Equal(t, readRes(t, res), "example page")
	})

	t.Run("InitialPageUrl", func(t *testing.T) {
		require.Equal(t, httpFetcher.InitialPageUrl(), "https://developer.omie.com.br/service-list/")
	})
}

func TestHttpFetcherErrors(t *testing.T) {
	t.Run("non 2xx errors", func(t *testing.T) {
		httpClient := mocks.NewHttpClientMock()
		httpClient.OnRequest("https://developer.omie.com.br/service-list/", "initial page", 400)

		httpFetcher := docs_fetcher.NewHttpFetcher(httpClient.ToHttpClient(), "")
		res, err := httpFetcher.GetInitialPage()
		require.Nil(t, res)
		require.Error(t, err)
	})
	t.Run("request error", func(t *testing.T) {
		httpClient := mocks.NewHttpClientMock()
		httpFetcher := docs_fetcher.NewHttpFetcher(httpClient.ToHttpClient(), "")
		res, err := httpFetcher.GetInitialPage()
		require.Nil(t, res)
		require.Error(t, err)
	})
}

func readRes(t *testing.T, res io.ReadCloser) string {
	defer res.Close()
	bodyData, err := io.ReadAll(res)
	require.NoError(t, err)
	return string(bodyData)
}
