package openapi_conversor

import (
	"fmt"

	"github.com/giuliano-macedo/omie-docs/internal/core"
	"github.com/swaggest/openapi-go/openapi3"
)

func convertAndAddModelsToSchemas(page core.Page, models []core.Model, mapOfSchemaOrRefValues *map[string]openapi3.SchemaOrRef) {
	for _, model := range models {
		(*mapOfSchemaOrRefValues)[getFieldRef(page.Endpoint, model.Name, false)] = openapi3.SchemaOrRef{
			Schema: convertModelToSchema(page, model),
		}
	}
}

func convertModelToSchema(page core.Page, model core.Model) (schema *openapi3.Schema) {
	schema = new(openapi3.Schema)
	schema.Type = newPtr(openapi3.SchemaTypeObject)
	properties := make(map[string]openapi3.SchemaOrRef, len(model.Fields))
	for _, field := range model.Fields {
		properties[field.Name] = convertFieldToSchema(field, page.Endpoint)
	}
	schema.Properties = properties
	schema.Required = getFieldRequiredProperties(model.Fields)
	return
}

func getFieldRequiredProperties(fields []core.Field) []string {
	requiredFields := make([]string, 0)
	for _, field := range fields {
		if field.IsRequired {
			requiredFields = append(requiredFields, field.Name)
		}
	}
	if len(requiredFields) == 0 {
		return nil
	}
	return requiredFields
}

func convertFieldToSchema(field core.Field, endpoint string) openapi3.SchemaOrRef {
	schemaType := newPtr(convertFieldTypeToSchemaType(field.Type))
	if *schemaType == openapi3.SchemaTypeObject {
		return openapi3.SchemaOrRef{
			SchemaReference: &openapi3.SchemaReference{
				Ref: getFieldRef(endpoint, field.TypeName, true),
			},
		}
	}
	schema := new(openapi3.Schema)
	schema.Type = schemaType
	schema.Title = &field.Name

	if field.Tooltip != "" {
		schema.Description = newPtr(fmt.Sprint(field.Description, "<br />", field.Tooltip))
	} else {
		schema.Description = &field.Description
	}
	if field.Length != 0 {
		schema.MaxLength = &field.Length
	}
	if field.IsDeprecated {
		schema.Deprecated = &field.IsDeprecated
	}

	if *schemaType == openapi3.SchemaTypeArray {
		schema.Items = &openapi3.SchemaOrRef{
			SchemaReference: &openapi3.SchemaReference{
				Ref: getFieldRef(endpoint, field.ElementType, true),
			},
		}
	}

	return openapi3.SchemaOrRef{Schema: schema}
}

func convertFieldTypeToSchemaType(fieldType core.Type) openapi3.SchemaType {
	switch fieldType {
	case core.String, core.Text:
		return openapi3.SchemaTypeString
	case core.Integer, core.Decimal:
		return openapi3.SchemaTypeNumber
	case core.Boolean:
		return openapi3.SchemaTypeBoolean
	case core.Array:
		return openapi3.SchemaTypeArray
	default: // case crawler.Object or anything else
		return openapi3.SchemaTypeObject
	}
}
