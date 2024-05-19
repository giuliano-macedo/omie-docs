package core

import (
	"bytes"
	_ "embed"
	"text/template"

	mod "github.com/giuliano-macedo/omie-docs"
)

//go:embed description.md.go.templ
var descriptionTemplate string

func RenderFeatureDescription(feature *Feature, collectionType CollectionType) string {
	tmpl, _ := template.New("FeatureDescription").Parse(descriptionTemplate)

	var buff bytes.Buffer
	var ctx struct {
		Description                   string
		IsOpenAPI, IsPostman, IsBruno bool
		CollectionName                string
		ProjectUrl                    string
		ExternalCollectionNames       []string
	}
	ctx.ProjectUrl = mod.Url
	ctx.ExternalCollectionNames = []string{
		PostmanCollectionName,
		BrunoCollectionName,
	}
	if feature != nil {
		ctx.Description = feature.Description
	}
	switch collectionType {
	case OpenAPI:
		ctx.IsOpenAPI = true
	case Postman:
		ctx.IsPostman = true
	case Bruno:
		ctx.IsBruno = true
	}
	ctx.CollectionName = collectionType.String()

	_ = tmpl.Execute(&buff, ctx)
	return buff.String()
}
