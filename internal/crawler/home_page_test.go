package crawler

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"testing"

	"github.com/giuliano-macedo/omie-docs/internal/core"
	"github.com/giuliano-macedo/omie-docs/test/mocks"
	"github.com/stretchr/testify/assert"
)

type homePageTestCase struct {
	name string
	url  string
}

func TestGetHomePageCases(t *testing.T) {
	cases := []homePageTestCase{
		{"normal_case", "https://example.com/list-of-apis/"},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			runHomePageTestCase(t, testCase)
		})
	}

}

func TestGetHomeErrorCases(t *testing.T) {
	cases := []homePageTestCase{
		{"no_header_row", ""},
		{"empty_header_row", ""},
		{"header_row_with_multiple_tags", ""},
		{"empty_row", ""},
		{"row_too_many_children", ""},
		{"row_no_href", ""},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			runHomePageTestErrCase(t, testCase)
		})
	}
}

func TestGetHomeFetcherError(t *testing.T) {
	expectedErr := errors.New("fetcher error")
	var fetcherMock = new(mocks.FetcherMock)
	fetcherMock.On("InitialPageUrl").Return("")
	fetcherMock.On("GetInitialPage").Return(nil, expectedErr)

	_, err := GetHomePage(fetcherMock)

	if !assert.Error(t, err) {
		return
	}
	assert.Equal(t, expectedErr, err)
}

func TestGetHomeReadError(t *testing.T) {
	readerErr := errors.New("reader error")
	reader := mocks.NewErrReadCloseMock(readerErr)

	var fetcherMock = new(mocks.FetcherMock)
	fetcherMock.On("InitialPageUrl").Return("")
	fetcherMock.On("GetInitialPage").Return(reader, nil)

	_, err := GetHomePage(fetcherMock)

	if !assert.Error(t, err) {
		return
	}
	assert.ErrorIs(t, err, readerErr)
}

func getHomePageTestSuitdata(caseName string) (docData string, expectedHomePage core.HomePage, err error) {
	docDataBytes, err := os.ReadFile(path.Join("test_data/home_page", caseName, "doc.html"))
	if err != nil {
		return
	}
	docData = string(docDataBytes)

	expectedHomeData, err := os.ReadFile(path.Join("test_data/home_page", caseName, "expected_home.json"))
	if err != nil {
		return
	}

	err = json.Unmarshal(expectedHomeData, &expectedHomePage)

	return
}

func runHomePageTestCase(t *testing.T, testCase homePageTestCase) {
	docData, exptecedHomePage, err := getHomePageTestSuitdata(testCase.name)
	if !assert.NoError(t, err) {
		return
	}

	var fetcherMock = new(mocks.FetcherMock)
	fetcherMock.On("InitialPageUrl").Return(testCase.url)
	fetcherMock.On("GetInitialPage").Return(mocks.NewReadCloseMock(docData), nil)

	home, err := GetHomePage(fetcherMock)

	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, exptecedHomePage, home)
}

func runHomePageTestErrCase(t *testing.T, testCase homePageTestCase) {
	docDataBytes, err := os.ReadFile(path.Join("test_data/home_page/error_cases", testCase.name+".html"))
	if !assert.NoError(t, err) {
		return
	}

	var fetcherMock = new(mocks.FetcherMock)
	fetcherMock.On("InitialPageUrl").Return(testCase.url)
	fetcherMock.On("GetInitialPage").Return(mocks.NewReadCloseMock(string(docDataBytes)), nil)

	_, err = GetHomePage(fetcherMock)

	if !assert.Error(t, err) {
		return
	}
}
