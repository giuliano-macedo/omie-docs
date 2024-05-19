package crawler

import (
	"io"
	"strings"

	"github.com/andybalholm/cascadia"
	"github.com/giuliano-macedo/omie-docs/internal/core"
	"github.com/giuliano-macedo/omie-docs/internal/docs_fetcher"
	"golang.org/x/net/html"
)

func GetPageByUrl(pages []core.Page) map[string]core.Page {
	pageByUrl := make(map[string]core.Page, len(pages))
	for _, page := range pages {
		pageByUrl[page.FullUrl] = page
	}
	return pageByUrl
}

func newSelector(query string) cascadia.Sel {
	sel, err := cascadia.Parse(query)
	if err != nil {
		panic(err)
	}
	return sel
}

func normalizeHtml(doc *html.Node) {
	var nodesToremove []*html.Node
	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		for child := node.LastChild; child != nil; child = child.PrevSibling {
			traverse(child)
		}

		isEmptyNode := node.Type == html.TextNode && strings.TrimSpace(node.Data) == ""
		if isEmptyNode && node.Parent != nil {
			nodesToremove = append(nodesToremove, node)
		}
	}

	traverse(doc)
	for _, node := range nodesToremove {
		node.Parent.RemoveChild(node)
	}
}

func parsePage(pageReader io.ReadCloser) (*html.Node, error) {
	defer pageReader.Close()
	doc, err := html.Parse(pageReader)
	if err != nil {
		return nil, err
	}
	normalizeHtml(doc)
	return doc, err
}

func getPage(fetcher docs_fetcher.DocsFetcher, url string) (*html.Node, error) {
	pageReader, err := fetcher.GetPage(url)
	if err != nil {
		return nil, err
	}
	return parsePage(pageReader)
}

func getChildrenNodes(node *html.Node) (nodes []html.Node) {
	if node == nil {
		return
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		nodes = append(nodes, *child)
	}
	return
}

func getAttr(node *html.Node, attrName string) (string, bool) {
	if node == nil {
		return "", false
	}
	for _, attr := range node.Attr {
		if attr.Key == attrName {
			return attr.Val, true
		}
	}
	return "", false
}

func nodeHasClass(node *html.Node, class string) bool {
	classStr, found := getAttr(node, "class")
	if !found {
		return false
	}
	return strings.Contains(classStr, class)
}

func getNodeText(node *html.Node) string {
	var stringBuilder strings.Builder

	var traverse func(node *html.Node)
	traverse = func(node *html.Node) {
		if node == nil {
			return
		}
		if node.Type == html.TextNode {
			stringBuilder.WriteString(node.Data + "\n")
		}
		for children := node.FirstChild; children != nil; children = children.NextSibling {
			traverse(children)
		}
	}
	traverse(node)
	return stringBuilder.String()
}
