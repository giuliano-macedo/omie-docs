package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/andybalholm/cascadia"
	"github.com/giuliano-macedo/omie-docs/internal/core"
	"github.com/giuliano-macedo/omie-docs/internal/docs_fetcher"
	"golang.org/x/net/html"
)

var (
	pageHeaderSel           = newSelector(".jumbotron div.col-12")
	highlightSel            = newSelector("a.highlight")
	methodSel               = newSelector(".methodItem")
	methodParameterTypeSel  = newSelector(".method-parameter-type>a")
	methodPanelFooterAnchor = newSelector(".panel-footer a")
	preSel                  = newSelector("pre")
	modelSel                = newSelector(".complexTypeItem")
	fieldLengthSel          = newSelector("span.parameter-lenght")
	h3Sel                   = newSelector("h3")
	pSel                    = newSelector("p")
	divPanelSel             = newSelector("div.panel")
	divMethodAreaSel        = newSelector("div.method-example-area")
	tdParameterName         = newSelector("td.parameter-name")
	tdParameterType         = newSelector("td.parameter-type")
	tdParameterDocs         = newSelector("td.parameter-docs")
	fieldFirstSel           = newSelector("div.first")
	fieldMoreSel            = newSelector("div.more>div")
	anchorSel               = newSelector("a")
	tableSel                = newSelector("table")
)

func parsePageHeader(page *core.Page, doc *html.Node) (err error) {
	pageHeader := cascadia.Query(doc, pageHeaderSel)
	if pageHeader == nil {
		return errors.New("couldn't find pageHeader")
	}
	if pageHeader.FirstChild != nil && pageHeader.FirstChild.FirstChild != nil {
		page.Name = pageHeader.FirstChild.FirstChild.Data
	}
	pageLink := cascadia.Query(pageHeader, highlightSel)
	urlString, _ := getAttr(pageLink, "href")
	err = extractPageUrlInfo(urlString, page)

	return
}

func parsePageMethods(doc *html.Node) (methods []core.Method, err error) {
	for _, methodNode := range cascadia.QueryAll(doc, methodSel) {
		var method core.Method

		methodNameNode := cascadia.Query(methodNode, h3Sel)
		methodDescriptionNode := cascadia.Query(methodNode, pSel)
		panelNode := cascadia.Query(methodNode, divPanelSel)
		exampleNode := cascadia.Query(methodNode, divMethodAreaSel)

		if methodNameNode != nil && methodNameNode.FirstChild != nil {
			method.Name = strings.TrimSpace(methodNameNode.FirstChild.Data)
		}
		if methodDescriptionNode != nil && methodDescriptionNode.FirstChild != nil {
			method.Description = strings.TrimSpace(methodDescriptionNode.FirstChild.Data)
			if strings.Contains(method.Description, "DEPRECATED") {
				method.IsDeprecated = true
			}
		}
		if panelNode != nil {
			method.Parameter, err = parsePageMethodParameter(*panelNode)
			if err != nil {
				return methods, err
			}
			method.Return = parsePageMethodReturn(*panelNode)
		}
		if exampleNode != nil {
			method.Example = parsePageMethodExample(*exampleNode)
		}
		methods = append(methods, method)

	}
	return
}

func parsePageMethodExample(exampleNode html.Node) (example map[string]interface{}) {
	pre := cascadia.Query(&exampleNode, preSel)
	text := getNodeText(pre)
	json.Unmarshal([]byte(text), &example)
	return
}

func parsePageMethodReturn(panelNode html.Node) (methodReturn string) {
	footer := cascadia.Query(&panelNode, methodPanelFooterAnchor)
	if footer == nil {
		return
	}
	methodReturn, _ = getAttr(footer, "href")
	return methodReturn[1:]
}

func parsePageMethodParameter(panelNode html.Node) (parameter string, err error) {
	elems := cascadia.QueryAll(&panelNode, methodParameterTypeSel)
	if len(elems) == 0 {
		return
	}
	if len(elems) != 1 {
		return parameter, errors.New(fmt.Sprint("expected to method have only one parameter, got ", len(elems)))
	}
	elem := elems[0]
	parameter, found := getAttr(elem, "href")
	if found {
		parameter = parameter[1:]
	}
	return
}

func parsePageModels(doc *html.Node) (models []core.Model) {
	for _, modelNode := range cascadia.QueryAll(doc, modelSel) {
		anchorNode := cascadia.Query(modelNode, anchorSel)
		fieldsTableNode := cascadia.Query(modelNode, tableSel)

		var model core.Model
		model.Name, _ = getAttr(anchorNode, "name")
		if strings.HasSuffix(model.Name, "Array") {
			continue
		}
		model.Fields = parsePageModelFields(fieldsTableNode)
		models = append(models, model)
	}
	return
}

func parsePageModelFields(fieldsTableNode *html.Node) (fields []core.Field) {
	if fieldsTableNode == nil || fieldsTableNode.FirstChild == nil || fieldsTableNode.FirstChild.FirstChild == nil {
		return
	}
	for row := fieldsTableNode.FirstChild.FirstChild; row != nil; row = row.NextSibling {
		if row.FirstChild == nil {
			// separator
			continue
		}

		nameNode := cascadia.Query(row, tdParameterName)
		typeNode := cascadia.Query(row, tdParameterType)
		docsNode := cascadia.Query(row, tdParameterDocs)

		field := parsePageModelField(nameNode, typeNode, docsNode)
		fields = append(fields, field)
	}
	return
}

func parsePageModelField(nameNode, typeNode, docsNode *html.Node) (field core.Field) {
	if nameNode != nil && nameNode.FirstChild != nil {
		field.Name = nameNode.FirstChild.Data
	}
	field.IsRequired = nodeHasClass(nameNode, "parameter-required")
	if typeNode != nil && typeNode.FirstChild != nil {
		if typeNode.FirstChild.Data == "a" {
			typeName, _ := getAttr(typeNode.FirstChild, "href")
			field.TypeName = typeName[1:]
		} else {
			field.TypeName = strings.TrimSpace(typeNode.FirstChild.Data)
		}
		fieldLengthNode := cascadia.Query(typeNode, fieldLengthSel)
		if fieldLengthNode != nil && fieldLengthNode.FirstChild != nil {
			length, _ := strconv.Atoi(fieldLengthNode.FirstChild.Data)
			field.Length = int64(length)
		}
		field.Type = getFieldType(field.TypeName)
		if field.Type == core.Array {
			field.ElementType = getElementType(field.TypeName)
		}
	}
	if docsNode != nil {
		descriptionNode := cascadia.Query(docsNode, fieldFirstSel)
		if descriptionNode != nil && descriptionNode.FirstChild != nil {
			field.Description = strings.TrimSpace(descriptionNode.FirstChild.Data)
			field.IsDeprecated = strings.Contains(field.Description, "DEPRECATED")
		}

		moreDivNode := cascadia.Query(docsNode, fieldMoreSel)
		field.Tooltip = strings.TrimSpace(getNodeText(moreDivNode))
	}
	return
}

func getElementType(typeName string) string {
	elementType, _ := strings.CutSuffix(typeName, "Array")
	if getFieldType(elementType) != core.Object {
		log.Warning("got an array that element type is not an object", typeName)
	}
	return elementType
}

func extractPageUrlInfo(urlString string, page *core.Page) (err error) {
	page.FullUrl = urlString
	url, err := url.Parse(urlString)
	if err != nil {
		return err
	}
	path := url.Path
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	pathSplit := strings.Split(path, "/")
	if len(pathSplit) < 2 {
		return errors.New(fmt.Sprint("expected path to have at least 2 components, got ", len(pathSplit)))
	}
	page.BasePath = pathSplit[0]
	page.ApiVersion = pathSplit[1]
	page.Endpoint = strings.Join(pathSplit[2:], "/")
	page.BaseUrl = fmt.Sprint(url.Scheme, "://", strings.Join(slices.Concat([]string{url.Host}, pathSplit[:2]), "/"))

	return
}

func getFieldType(typeName string) core.Type {
	switch typeName {
	case "string":
		return core.String
	case "integer":
		return core.Integer
	case "text":
		return core.Text
	case "decimal":
		return core.Decimal
	case "boolean":
		return core.Boolean
	}
	if strings.HasSuffix(typeName, "Array") {
		return core.Array
	}
	return core.Object
}

func GetPage(fetcher docs_fetcher.DocsFetcher, url string) (page core.Page, err error) {
	page.DocsUrl = url
	doc, err := getPage(fetcher, url)
	if err != nil {
		return page, err
	}
	err = parsePageHeader(&page, doc)
	if err != nil {
		return page, err
	}
	page.Methods, err = parsePageMethods(doc)
	if err != nil {
		return page, err
	}
	page.Models = parsePageModels(doc)
	return
}
