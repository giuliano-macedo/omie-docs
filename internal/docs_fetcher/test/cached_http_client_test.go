package docs_fetcher_test

import (
	"os"
	"path"
	"testing"

	"github.com/giuliano-macedo/omie-docs/internal/docs_fetcher"
	"github.com/giuliano-macedo/omie-docs/test/mocks"
	"github.com/stretchr/testify/require"
)

func TestCachedHttpClient(t *testing.T) {
	dirTemp, err := os.MkdirTemp("", "cached_http_client")
	require.NoError(t, err)

	url, content := "https://example.com", "example page"
	httpClient := mocks.NewHttpClientMock()
	httpClient.OnRequest(url, content, 200)

	cachedHttpClient, err := docs_fetcher.NewCachedHttpClient(dirTemp, httpClient.ToHttpClient())
	require.NoError(t, err)

	t.Run("Uncached request", func(t *testing.T) {
		res, err := cachedHttpClient.Get(url)
		require.NoError(t, err)

		require.Equal(t, readRes(t, res.Body), content)
		require.Equal(t, httpClient.RequestCount(url), 1)
		require.Equal(t, readFileInTemp(t, dirTemp, docs_fetcher.HashUrl(url)), content)
		expectNumberOfFilesInFolder(t, dirTemp, 1)
	})
	t.Run("cached request", func(t *testing.T) {
		res, err := cachedHttpClient.Get(url)
		require.NoError(t, err)

		require.Equal(t, readRes(t, res.Body), content)
		require.Equal(t, httpClient.RequestCount(url), 1)
	})
}

func readFileInTemp(t *testing.T, tmpDir, fname string) string {
	data, err := os.ReadFile(path.Join(tmpDir, fname))
	require.NoError(t, err)
	return string(data)
}

func expectNumberOfFilesInFolder(t *testing.T, dir string, expectedCount int) {
	count := 0
	files, err := os.ReadDir(dir)
	require.NoError(t, err)
	for _, file := range files {
		if !file.IsDir() {
			count++
		}
	}
	require.Equal(t, expectedCount, count)
}
