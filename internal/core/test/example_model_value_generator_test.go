package core_test

import (
	"testing"

	"github.com/giuliano-macedo/omie-docs/internal/core"
	"github.com/stretchr/testify/require"
)

func TestExampleModelValueGenerator(t *testing.T) {

	testCases := []struct {
		name           string
		page           core.Page
		modelName      string
		expectedOutput map[string]interface{}
	}{
		{
			"primitive types",
			core.Page{
				Models: []core.Model{
					{
						Name: "model1",
						Fields: []core.Field{
							{Name: "field1", Type: core.String},
							{Name: "field2", Type: core.Text},
							{Name: "field3", Type: core.Boolean},
							{Name: "field4", Type: core.Decimal},
							{Name: "field5", Type: core.Integer},
						},
					},
				},
			},
			"model1",
			map[string]interface{}{
				"field1": "string",
				"field2": "text",
				"field3": false,
				"field4": 42.42,
				"field5": 42,
			},
		},
		{
			"objects",
			core.Page{
				Models: []core.Model{
					{
						Name: "model1",
						Fields: []core.Field{
							{Name: "field1", Type: core.String},
							{Name: "field2", Type: core.Object, TypeName: "model2"},
						},
					},
					{
						Name: "model2",
						Fields: []core.Field{
							{Name: "field1", Type: core.String},
							{Name: "field2", Type: core.Text},
						},
					},
				},
			},
			"model1",
			map[string]interface{}{
				"field1": "string",
				"field2": map[string]interface{}{
					"field1": "string",
					"field2": "text",
				},
			},
		},
		{
			"arrays",
			core.Page{
				Models: []core.Model{
					{
						Name: "model1",
						Fields: []core.Field{
							{Name: "field1", Type: core.String},
							{Name: "field2", Type: core.Array, ElementType: "model2"},
						},
					},
					{
						Name: "model2",
						Fields: []core.Field{
							{Name: "field1", Type: core.String},
						},
					},
				},
			},
			"model1",
			map[string]interface{}{
				"field1": "string",
				"field2": []interface{}{
					map[string]interface{}{
						"field1": "string",
					},
					map[string]interface{}{
						"field1": "string",
					},
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			generator := core.NewExampleModelValueGenerator([]core.Page{testCase.page})
			exampleOutput, found := generator.Get(testCase.modelName)
			require.True(t, found)
			require.Equal(t, testCase.expectedOutput, exampleOutput)
		})
	}
}
