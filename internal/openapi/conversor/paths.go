package openapi_conversor

import (
	"fmt"
	"slices"

	"github.com/giuliano-macedo/omie-docs/internal/core"
	"github.com/swaggest/openapi-go/openapi3"
)

func convertAndAddMethodsToPaths(page core.Page, methods []core.Method, mapOfPathItemValues *map[string]openapi3.PathItem) {
	for _, method := range methods {
		pathName := fmt.Sprint("/", page.Endpoint, "/", "#", method.Name)
		(*mapOfPathItemValues)[pathName] = convertMethodToPath(page, method)
	}
}

func convertMethodToPath(page core.Page, method core.Method) (pathItem openapi3.PathItem) {
	pathItem.MapOfOperationValues = make(map[string]openapi3.Operation, 1)
	op := openapi3.Operation{
		Tags:        []string{page.EntityName},
		Description: &method.Description,
		Deprecated:  &method.IsDeprecated,
		RequestBody: &openapi3.RequestBodyOrRef{
			RequestBody: &openapi3.RequestBody{
				Content: map[string]openapi3.MediaType{
					"application/json": {
						Schema: &openapi3.SchemaOrRef{
							Schema: convertMethodToSchema(page, method),
						},
					},
				},
			},
		},
	}
	if method.Return != "" {
		op.Responses = openapi3.Responses{
			MapOfResponseOrRefValues: map[string]openapi3.ResponseOrRef{
				"200": {
					Response: &openapi3.Response{
						Content: map[string]openapi3.MediaType{
							"application/json": {
								Schema: &openapi3.SchemaOrRef{
									SchemaReference: &openapi3.SchemaReference{
										Ref: getFieldRef(page.Endpoint, method.Return, true),
									},
								},
							},
						},
					},
				},
			},
		}
	} else {
		op.Responses = openapi3.Responses{
			MapOfResponseOrRefValues: map[string]openapi3.ResponseOrRef{
				"200": {
					Response: &openapi3.Response{
						Content: map[string]openapi3.MediaType{
							"application/json": {
								Schema: &openapi3.SchemaOrRef{
									Schema: &openapi3.Schema{
										Type:  newPtr(openapi3.SchemaTypeString),
										Title: newPtr("Não documentado"),
									},
								},
							},
						},
					},
				},
			},
		}
	}
	pathItem.MapOfOperationValues["post"] = op
	return
}

func convertMethodToSchema(page core.Page, method core.Method) (schema *openapi3.Schema) {
	schema = new(openapi3.Schema)
	schema.Type = newPtr(openapi3.SchemaTypeObject)
	schema.Properties = map[string]openapi3.SchemaOrRef{
		"call": {
			Schema: &openapi3.Schema{
				Title: newPtr("call"),
				Enum:  []interface{}{method.Name},
				Type:  newPtr(openapi3.SchemaTypeString),
			},
		},
		"app_key": {
			Schema: &openapi3.Schema{
				Title:   newPtr("app_key"),
				Example: newInterface(newPtr("{{appKey}}")),
				Type:    newPtr(openapi3.SchemaTypeString),
			},
		},
		"app_secret": {
			Schema: &openapi3.Schema{
				Title:   newPtr("app_secret"),
				Example: newInterface(newPtr("{{appSecret}}")),
				Type:    newPtr(openapi3.SchemaTypeString),
			},
		},
		"param": {
			Schema: &openapi3.Schema{
				Title:   newPtr("param"),
				Example: newInterface([]map[string]interface{}{method.Example}),
				Type:    newPtr(openapi3.SchemaTypeArray),
				Items: &openapi3.SchemaOrRef{
					SchemaReference: &openapi3.SchemaReference{
						Ref: getFieldRef(page.Endpoint, method.Parameter, true),
					},
				},
			},
		},
	}
	schema.Required = make([]string, 0, len(schema.Properties))
	for field := range schema.Properties {
		schema.Required = append(schema.Required, field)
	}
	slices.Sort(schema.Required)
	return
}
