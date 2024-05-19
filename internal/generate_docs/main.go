package generate_docs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	// "github.com/giuliano-macedo/omie-docs/internal/bruno"
	"github.com/giuliano-macedo/omie-docs/internal/bruno"
	"github.com/giuliano-macedo/omie-docs/internal/bundler"
	"github.com/giuliano-macedo/omie-docs/internal/crawler"
	"github.com/giuliano-macedo/omie-docs/internal/docs_fetcher"
	openapi_bundler "github.com/giuliano-macedo/omie-docs/internal/openapi/bundler"
	"github.com/giuliano-macedo/omie-docs/internal/postman"
)

func handleError(e error) {
	if e != nil {
		panic(e)
	}
}

func saveStructToFile(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "    ")

	handleError(err)

	typeName := fmt.Sprintf("%T", data)

	fileName := "types_dump/" + typeName + ".json"
	err = os.WriteFile(fileName, jsonData, 0644)
	handleError(err)
}

type Args struct {
	OutputDir                 string
	NumberOfWorkers           int
	InitialPage               string
	DumpIntermediaryDataTypes bool
	CachingDirectory          string
	HttpClient                *http.Client
	PrettifyJson              bool
}

func Run(args Args) {
	var (
		httpClient *http.Client
		err        error
	)
	switch {
	case args.HttpClient != nil:
		httpClient = args.HttpClient
	case args.CachingDirectory != "":
		httpClient, err = docs_fetcher.NewCachedHttpClient(args.CachingDirectory, http.DefaultClient)
	default:
		httpClient = http.DefaultClient
	}
	handleError(err)

	fetcher := docs_fetcher.NewHttpFetcher(httpClient, args.InitialPage)
	home, err := crawler.GetHomePage(fetcher)
	handleError(err)
	if args.DumpIntermediaryDataTypes {
		saveStructToFile(home)
	}

	pages, err := crawler.GetPages(fetcher, home, args.NumberOfWorkers)
	handleError(err)

	if args.DumpIntermediaryDataTypes {
		saveStructToFile(pages)
	}
	fsWriter, err := bundler.NewOsFileWriter(args.OutputDir, args.PrettifyJson)
	handleError(err)

	bundlerArgs := bundler.Args{
		Pages:    pages,
		Home:     home,
		FsWriter: fsWriter,
	}

	runner := bundler.Runner{Args: bundlerArgs}

	err = runner.Run(
		openapi_bundler.NewOpenApiBundler(),
		postman.NewPostmanBundler(),
		bruno.NewBrunoBundler(),
	)

	handleError(err)
}
