package openapi_bundler

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/giuliano-macedo/omie-docs/internal/bundler"
)

const (
	title            = "Omie OpenAPI"
	description      = "Documentação OpenAPI não oficial da Omie"
	swaggerUiVersion = "5.17.2"
)

func swaggerUiDistFile(path string) string {
	return fmt.Sprint("https://unpkg.com/swagger-ui-dist@", swaggerUiVersion, "/", path)
}

//go:embed templates/index.html.go.templ
var indexTemplate string

//go:embed static/main.js
var mainJs string

func generateIndex(urlConfigs []UrlConfig, fsWriter bundler.FSWriter) error {
	tmpl, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		return err
	}

	SwaggerUrls, err := json.Marshal(urlConfigs)
	if err != nil {
		return err
	}

	ctx := IndexContext{
		Title:       title,
		Description: description,
		Css:         swaggerUiDistFile("swagger-ui.css"),
		Icon16:      swaggerUiDistFile("favicon-16x16.png"),
		Icon32:      swaggerUiDistFile("favicon-32x32.png"),
		JsBundle:    swaggerUiDistFile("swagger-ui-bundle.js"),
		JsPreset:    swaggerUiDistFile("swagger-ui-standalone-preset.js"),
		SwaggerUrls: string(SwaggerUrls),
		MainJs:      removeJsComments(mainJs),
	}

	file, err := fsWriter.Create("index.html")
	if err != nil {
		return err
	}
	defer file.Close()

	if err := tmpl.Execute(file, ctx); err != nil {
		return err
	}

	return nil
}
