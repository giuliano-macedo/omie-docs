package openapi_conversor

import (
	"fmt"

	mod "github.com/giuliano-macedo/omie-docs"
	"github.com/giuliano-macedo/omie-docs/internal/core"
	"github.com/swaggest/openapi-go/openapi3"
)

const (
	Title          = "Omie OpenApi"
	OpenapiVersion = "3.0.3"
)

func ConvertToOpenApi(homepage core.HomePage, pages []core.Page, feature *core.Feature) (spec openapi3.Spec) {
	spec.Openapi = OpenapiVersion
	spec.Info.Version = mod.Version
	spec.ExternalDocs = &openapi3.ExternalDocumentation{
		Description: newPtr("Documentação original das APIs"),
		URL:         homepage.DocsUrl,
	}
	if feature != nil {
		spec.Info.Title = fmt.Sprint(Title, " - ", feature.Name)
	} else {
		spec.Info.Title = Title
	}
	if len(pages) != 0 {
		spec.Servers = []openapi3.Server{
			{URL: pages[0].BaseUrl, Description: newPtr("Servidor de produção")},
		}
	}
	spec.Info.Description = newPtr(core.RenderFeatureDescription(feature, core.OpenAPI))

	mapOfPathItemValues, mapOfSchemaOrRefValues := newMapOfPathItemsAndSchmemas(pages)

	for _, page := range pages {
		convertAndAddMethodsToPaths(page, page.Methods, &mapOfPathItemValues)
		convertAndAddModelsToSchemas(page, page.Models, &mapOfSchemaOrRefValues)
	}
	spec.Tags = convertToTags(pages, homepage)
	spec.Paths.MapOfPathItemValues = mapOfPathItemValues
	spec.Components = &openapi3.Components{
		Schemas: &openapi3.ComponentsSchemas{
			MapOfSchemaOrRefValues: mapOfSchemaOrRefValues,
		},
	}
	return
}

func convertToTags(pages []core.Page, homepage core.HomePage) []openapi3.Tag {
	entitiesByName := make(map[string]core.Entity)
	for _, feature := range homepage.Features {
		for _, entity := range feature.AllEntities() {
			entitiesByName[entity.Name] = entity
		}
	}

	tags := make([]openapi3.Tag, 0, len(pages))
	for _, page := range pages {
		entity := entitiesByName[page.EntityName]
		tags = append(tags, openapi3.Tag{
			Name:        entity.Name,
			Description: &entity.Description,
			ExternalDocs: &openapi3.ExternalDocumentation{
				Description: newPtr("Documentação original"),
				URL:         page.DocsUrl,
			},
		})
	}
	return tags
}

func newMapOfPathItemsAndSchmemas(pages []core.Page) (map[string]openapi3.PathItem, map[string]openapi3.SchemaOrRef) {
	estimatedSizeOfMethods := len(pages)
	estimatedSizeOfComponents := len(pages)
	if len(pages) != 0 {
		estimatedSizeOfMethods *= len(pages[0].Methods)
		estimatedSizeOfComponents *= len(pages[0].Models)
	}
	mapOfPathItemValues := make(map[string]openapi3.PathItem, estimatedSizeOfMethods)
	mapOfSchemaOrRefValues := make(map[string]openapi3.SchemaOrRef, estimatedSizeOfComponents)
	return mapOfPathItemValues, mapOfSchemaOrRefValues
}
