package main

import (
	"flag"
	"runtime"

	"github.com/giuliano-macedo/omie-docs/internal/generate_docs"
)

func main() {
	outputDir := flag.String("output", "bundle", "Output dir to save openapi definitions and swagger ui static files")
	cachingDirectory := flag.String("caching_dir", "", "Directory to use as http caching (empty to disable caching)")
	numberOfWorkers := flag.Int("workers", runtime.NumCPU(), "Number of workers to fetch and process data")
	initialPage := flag.String("initial_page", "", "Initial URL to crawl (empty to use default)")
	dumpIntermediaryDataTypes := flag.Bool("intermediary", false, "Enables dumping intermediary crawling results as json in directory ./cache")
	prettifyJson := flag.Bool("prettify", false, "Prettify openapi schmema json files")
	flag.Parse()

	generate_docs.Run(generate_docs.Args{
		OutputDir:                 *outputDir,
		InitialPage:               *initialPage,
		DumpIntermediaryDataTypes: *dumpIntermediaryDataTypes,
		CachingDirectory:          *cachingDirectory,
		NumberOfWorkers:           *numberOfWorkers,
		PrettifyJson:              *prettifyJson,
	})
}
