package openapi_bundler

import (
	"fmt"
	"path"

	"github.com/giuliano-macedo/omie-docs/internal/bundler"
	"github.com/giuliano-macedo/omie-docs/internal/core"
	"github.com/giuliano-macedo/omie-docs/internal/crawler"
	conversor "github.com/giuliano-macedo/omie-docs/internal/openapi/conversor"
)

func GetSpecFilePath(feature *core.Feature) string {
	if feature != nil {
		return path.Join("openapi", fmt.Sprint(formatFeatureName(feature.Name), ".openapi.json"))
	}
	return path.Join("openapi", "all.openapi.json")
}

func (bundler *OpenApiBundler) Bundle(args bundler.Args) error {
	pageByUrl := crawler.GetPageByUrl(args.Pages)

	urlConfigs := make([]UrlConfig, 0, len(args.Home.Features)+1)
	addUrlConfig := func(name, url string) {
		urlConfigs = append(urlConfigs, UrlConfig{
			Url:  url,
			Name: name,
		})
	}
	err := args.FsWriter.MkdirAll("openapi")
	if err != nil {
		return err
	}
	for _, feature := range args.Home.Features {
		allEntities := feature.AllEntities()
		pagesByFeature := make([]core.Page, 0, len(allEntities))
		for _, entity := range allEntities {
			pagesByFeature = append(pagesByFeature, pageByUrl[entity.Url])
		}

		openApiPath := GetSpecFilePath(&feature)
		featureSpec := conversor.ConvertToOpenApi(args.Home, pagesByFeature, &feature)
		if err = args.FsWriter.SaveJsonFile(openApiPath, featureSpec); err != nil {
			return err
		}
		addUrlConfig(feature.Name, openApiPath)
	}
	openApiPath := GetSpecFilePath(nil)
	allSpec := conversor.ConvertToOpenApi(args.Home, args.Pages, nil)
	if err = args.FsWriter.SaveJsonFile(openApiPath, allSpec); err != nil {
		return err
	}
	addUrlConfig("Todas APIs", openApiPath)
	if err = generateIndex(urlConfigs, args.FsWriter); err != nil {
		return err
	}
	return nil
}

func NewOpenApiBundler() bundler.Bundler {
	return &OpenApiBundler{}
}
