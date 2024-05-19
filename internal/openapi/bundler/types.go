package openapi_bundler

type OpenApiBundler struct{}

type UrlConfig struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type IndexContext struct {
	Title       string
	Description string
	Css         string
	Icon32      string
	Icon16      string
	JsBundle    string
	JsPreset    string
	SwaggerUrls string
	MainJs      string
}
