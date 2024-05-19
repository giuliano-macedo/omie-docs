package core

type ExampleModelValueGenerator struct {
	examplesPerModel map[string]interface{}
}

func NewExampleModelValueGenerator(pages []Page) *ExampleModelValueGenerator {
	examplesPerModel := make(map[string]interface{})
	for _, page := range pages {
		generateExamplePerModel(page, &examplesPerModel)
	}
	return &ExampleModelValueGenerator{examplesPerModel}
}

func (generator *ExampleModelValueGenerator) Get(modelName string) (interface{}, bool) {
	res, found := generator.examplesPerModel[modelName]
	return res, found

}

func generateExamplePerModel(page Page, examplesPerModel *map[string]interface{}) {
	modelsByName := make(map[string]Model)
	for _, model := range page.Models {
		modelsByName[model.Name] = model
	}
	for _, model := range page.Models {
		(*examplesPerModel)[model.Name] = resolveModel(modelsByName, model)
	}
}

func resolveModel(modelsByName map[string]Model, model Model) interface{} {
	result := make(map[string]interface{}, len(model.Fields))
	for _, field := range model.Fields {
		var value interface{}
		switch field.Type {
		case String:
			value = "string"
		case Text:
			value = "text"
		case Integer:
			value = 42
		case Decimal:
			value = 42.42
		case Boolean:
			value = false
		case Object:
			subModel, found := modelsByName[field.TypeName]
			if found {
				value = resolveModel(modelsByName, subModel)
			}
		case Array:
			subModel, found := modelsByName[field.ElementType]
			if found {
				valueArray := make([]interface{}, 2)
				for i := range valueArray {
					valueArray[i] = resolveModel(modelsByName, subModel)
				}
				value = valueArray
			}
		}
		if value != nil {
			result[field.Name] = value
		}
	}
	return result
}
