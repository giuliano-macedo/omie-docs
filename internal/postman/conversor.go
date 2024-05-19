package postman

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/giuliano-macedo/go-postman-collection"
	mod "github.com/giuliano-macedo/omie-docs"
	"github.com/giuliano-macedo/omie-docs/internal/core"
	"github.com/giuliano-macedo/omie-docs/internal/crawler"
)

func mustMarshalJson[T string | []byte](obj interface{}) T {
	data, _ := json.MarshalIndent(obj, "", "    ")
	return T(data)
}

func convertPageToItems(page core.Page) []*postman.Items {
	items := make([]*postman.Items, 0, len(page.Methods))

	exampleModelValueGenerator := core.NewExampleModelValueGenerator([]core.Page{page})
	for _, method := range page.Methods {
		req := convertMethodToRequest(page, method)
		items = append(items, &postman.Items{
			Name:        method.Name,
			Description: &postman.Description{Content: method.Description},
			Request:     req,
			Responses:   convertMethodToResponses(method, exampleModelValueGenerator, req),
		})
	}

	return items
}

func convertMethodToRequest(page core.Page, method core.Method) *postman.Request {
	return &postman.Request{
		URL: &postman.URL{
			Host: []string{"{{baseUrl}}"},
			Path: []string{page.Endpoint, ""},
			Raw:  path.Join("{{baseUrl}}", page.Endpoint) + "/",
		},
		Description: method.Description,
		Method:      "POST",
		Header: []*postman.Header{{
			Key:   "Content-Type",
			Value: "application/json",
		}},
		Body: &postman.Body{
			Mode: "raw",
			Raw: mustMarshalJson[string](
				&map[string]interface{}{
					"app_key":    "{{appKey}}",
					"app_secret": "{{appSecret}}",
					"call":       method.Name,
					"param":      []interface{}{method.Example},
				},
			),
		},
	}
}

func convertMethodToResponses(method core.Method, exampleModelValueGenerator *core.ExampleModelValueGenerator, request *postman.Request) []*postman.Response {
	res := &postman.Response{
		Name:            "Exemplo de resposta",
		Status:          "OK",
		Code:            200,
		OriginalRequest: request,
		Headers: &postman.HeaderList{
			Headers: []*postman.Header{{
				Key:   "Content-Type",
				Value: "application/json",
			}},
		},
	}
	bodyResponse, found := exampleModelValueGenerator.Get(method.Return)
	if found {
		res.Body = mustMarshalJson[string](bodyResponse)
	}
	return []*postman.Response{res}
}

func addEntitiesToGroup(pageByUrl map[string]core.Page, group *postman.Items, entities []core.Entity) {
	for _, entity := range entities {
		entityItemGroup := group.AddItemGroup(entity.Name)

		page := pageByUrl[entity.Url]
		entityItemGroup.Description = &postman.Description{
			Type: "text/markdown",
			Content: strings.Join([]string{
				fmt.Sprint("## ", entity.Name),
				"",
				fmt.Sprintf("[Documentação original](%v)", page.DocsUrl),
				entity.Description,
			}, "\n"),
		}

		entityItemGroup.Items = convertPageToItems(page)
	}
}

func createVariable(name string, value string) *postman.Variable {
	return &postman.Variable{
		Name:  name,
		ID:    name,
		Value: value,
		Type:  "string",
	}
}

func convertToPostman(home core.HomePage, pages []core.Page) postman.Collection {
	var baseUrl string
	if len(pages) != 0 {
		baseUrl = pages[0].BaseUrl
	}

	collection := postman.Collection{
		Info: postman.Info{
			Name: fmt.Sprint("Omie Docs ", mod.Version),
			Description: postman.Description{
				Type:    "text/markdown",
				Content: core.RenderFeatureDescription(nil, core.Postman),
			},
		},
	}
	collection.Info.Schema = "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"

	collection.Variables = []*postman.Variable{
		createVariable("baseUrl", baseUrl),
		createVariable("appKey", ""),
		createVariable("appSecret", ""),
	}

	pageByUrl := crawler.GetPageByUrl(pages)

	for _, feature := range home.Features {
		featureItemGroup := collection.AddItemGroup(feature.Name)
		featureItemGroup.Description = &postman.Description{
			Type:    "text/markdown",
			Content: feature.Description,
		}

		addEntitiesToGroup(pageByUrl, featureItemGroup, feature.MainEntities)

		auxItemGroup := featureItemGroup.AddItemGroup("Cadastros Auxiliares")
		addEntitiesToGroup(pageByUrl, auxItemGroup, feature.AuxiliaryEntites)
	}
	return collection
}
