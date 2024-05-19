package docs_fetcher

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
)

type CachedHttpClient struct {
	path      string
	client    *http.Client
	cachingFs fs.FS
}

// Create an HttpClient that will cache all OK response bodies to a directory, and read from that directory
// if the file exists
func NewCachedHttpClient(path string, client *http.Client) (httpClient *http.Client, err error) {
	if err = os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, err
	}
	cachedHttpClient := &CachedHttpClient{
		path:      path,
		client:    client,
		cachingFs: os.DirFS(path),
	}
	return &http.Client{
		Transport: cachedHttpClient,
	}, nil
}

func (cachedHttpClient *CachedHttpClient) RoundTrip(req *http.Request) (res *http.Response, err error) {
	url := req.URL.String()
	cacheReader, err := cachedHttpClient.readFromCache(url)
	if err != nil {
		return nil, err
	}
	if cacheReader != nil {
		return &http.Response{
			Body:       cacheReader,
			StatusCode: 200,
			Status:     "OK",
			Request:    req,
		}, nil
	}
	res, err = cachedHttpClient.client.Do(req)
	if err != nil {
		return nil, err
	}
	dataBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	data := string(dataBytes)
	resIsOk := res.StatusCode >= 200 && res.StatusCode <= 299
	if resIsOk {
		if err = cachedHttpClient.writeToCache(url, data); err != nil {
			return nil, err
		}
	}

	res.Body = io.NopCloser(strings.NewReader(data))
	return res, nil
}

func HashUrl(data string) string {
	hashedData := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hashedData[:])
}

func (cachedHttpClient *CachedHttpClient) readFromCache(url string) (io.ReadCloser, error) {
	hashedUrl := HashUrl(url)
	file, err := cachedHttpClient.cachingFs.Open(hashedUrl)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	return file, err
}

func (cachedHttpClient *CachedHttpClient) writeToCache(url string, data string) (err error) {
	hashedUrl := HashUrl(url)

	cacheFile, err := os.Create(path.Join(cachedHttpClient.path, hashedUrl))
	if err != nil {
		return
	}
	defer cacheFile.Close()

	_, err = io.WriteString(cacheFile, data)
	return
}
