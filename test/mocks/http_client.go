package mocks

import (
	"errors"
	"io"
	"net/http"
	"strings"
)

type HttpClientMock struct {
	responses    map[string]*http.Response
	requestCount map[string]int
}

func (httpClientMock *HttpClientMock) RoundTrip(req *http.Request) (*http.Response, error) {
	url := req.URL.String()
	response, ok := httpClientMock.responses[url]
	if !ok {
		return nil, errors.New("Unexpected URL: " + url)
	}
	httpClientMock.requestCount[url]++
	return response, nil

}

func NewHttpClientMock() *HttpClientMock {
	return &HttpClientMock{
		responses:    map[string]*http.Response{},
		requestCount: map[string]int{},
	}
}

func (httpClientMock *HttpClientMock) ToHttpClient() *http.Client {
	return &http.Client{
		Transport: httpClientMock,
	}
}

func (httpClientMock *HttpClientMock) OnRequest(url, responseData string, statusCode int) {
	httpClientMock.responses[url] = &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(responseData)),
	}
}

func (httpClientMock *HttpClientMock) RequestCount(url string) int {
	return httpClientMock.requestCount[url]
}
