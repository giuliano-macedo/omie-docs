package bruno

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/giuliano-macedo/go-bruno-collection/pkg/bruno"
	mod "github.com/giuliano-macedo/omie-docs"
	"github.com/giuliano-macedo/omie-docs/internal/core"
	"github.com/giuliano-macedo/omie-docs/internal/crawler"
)

const (
	environmentUid = "17b03cda7f1d4b1caf12f"
)

func mustMarshalJson[T string | []byte](obj interface{}) T {
	data, _ := json.MarshalIndent(obj, "", "    ")
	return T(data)
}

func convertPageToItems(page core.Page) []bruno.Item {
	items := make([]bruno.Item, 0, len(page.Methods))

	for _, method := range page.Methods {
		req := convertMethodToRequest(page, method)
		items = append(items, bruno.Item{
			Name:    method.Name,
			Request: req,
			Type:    "http-request",
		})
	}
	return items
}

func convertMethodToRequest(page core.Page, method core.Method) *bruno.Request {
	return &bruno.Request{
		Docs: strings.Join([]string{
			fmt.Sprintf("[Documentação original](%v)", page.DocsUrl),
			"",
			fmt.Sprint("## ", page.EntityName),
			"",
			method.Description,
			"",
		}, "\n"),
		Auth: bruno.Auth{
			Mode: "none",
		},
		URL:    path.Join("{{baseUrl}}", page.Endpoint) + "/",
		Method: "POST",
		Headers: []bruno.Header{{
			Name:    "Content-Type",
			Value:   "application/json",
			Enabled: true,
		}},
		Params: []bruno.Param{},
		Body: bruno.Body{
			Mode: "json",
			Json: mustMarshalJson[string](
				map[string]interface{}{
					"app_key":    "{{appKey}}",
					"app_secret": "{{appSecret}}",
					"call":       method.Name,
					"param":      []interface{}{method.Example},
				},
			),
		},
	}
}

func addEntitiestoItem(pageByUrl map[string]core.Page, rootItem *bruno.Item, entities []core.Entity) {
	rootItem.Items = make([]bruno.Item, len(entities))
	for i, entity := range entities {
		entitiyItem := &rootItem.Items[i]
		entitiyItem.Name = entity.Name
		entitiyItem.Type = "folder"

		page := pageByUrl[entity.Url]
		// Not supported for folders :/
		// entitiyItem.Docs = strings.Join([]string{
		// 	fmt.Sprint("## ", entity.Name),
		// 	"",
		// 	fmt.Sprintf("[Documentação original](%v)", page.DocsUrl),
		// 	entity.Description,
		// }, "\n")

		entitiyItem.Items = convertPageToItems(page)
	}

}

func newEnv(name, value string, secret bool) bruno.EnvironmentVariable {
	return bruno.EnvironmentVariable{
		Name:    name,
		Value:   value,
		Type:    "text",
		Enabled: true,
		Secret:  secret,
	}
}

func ConvertToBruno(home core.HomePage, pages []core.Page) (collection bruno.Collection) {
	var baseUrl string
	if len(pages) != 0 {
		baseUrl = pages[0].BaseUrl
	}

	collection.Name = fmt.Sprint("Omie Docs ", mod.Version)
	collection.Docs = core.RenderFeatureDescription(nil, core.Bruno)
	collection.Version = "1"
	collection.ActiveEnvironmentUid = environmentUid

	collection.Environments = []bruno.Environment{
		{
			Uid:  environmentUid,
			Name: "Default",
			Variables: []bruno.EnvironmentVariable{
				newEnv("baseUrl", baseUrl, false),
				newEnv("appKey", "", true),
				newEnv("appSecret", "", true),
			},
		},
	}

	pageByUrl := crawler.GetPageByUrl(pages)

	collection.Items = make([]bruno.Item, len(home.Features))
	for i, feature := range home.Features {
		rootItem := &collection.Items[i]
		rootItem.Type = "folder"
		// Not supported for folders
		// rootItem.Docs = feature.Description
		rootItem.Name = feature.Name

		addEntitiestoItem(pageByUrl, rootItem, feature.MainEntities)

		auxItems := bruno.Item{Name: "Cadastros Auxiliares", Type: "folder"}
		addEntitiestoItem(pageByUrl, &auxItems, feature.AuxiliaryEntites)
		rootItem.Items = append(rootItem.Items, auxItems)
	}
	return
}
