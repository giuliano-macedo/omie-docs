package crawler

import (
	"errors"
	"strings"

	"github.com/andybalholm/cascadia"
	"github.com/giuliano-macedo/omie-docs/internal/core"
	"github.com/giuliano-macedo/omie-docs/internal/docs_fetcher"
	"golang.org/x/net/html"
)

var (
	tableRowSel = newSelector(".table tr")
)

func parseHomePageHeader(columns []html.Node) (feature *core.Feature, err error) {
	if columns == nil {
		return nil, errors.New("empty node")
	}
	if len(columns) != 1 {
		return nil, errors.New("len columns not equal to 1")
	}
	feature = new(core.Feature)

	column := columns[0]
	children := getChildrenNodes(&column)
	strong := children[0]
	paragraph := children[1]
	if strong.FirstChild != nil {
		feature.Name = strong.FirstChild.Data
	}
	if paragraph.FirstChild != nil {
		feature.Description = strings.TrimSpace(paragraph.FirstChild.Data)
	}
	return
}

func parseHomePageRow(columns []html.Node) (entity core.Entity, err error) {
	if columns == nil {
		return entity, errors.New("empty node")
	}
	if len(columns) != 3 {
		return entity, errors.New("len columns not equal to 3")
	}
	column1 := columns[0]
	column2 := columns[1]
	column3 := columns[2]

	entityLink := column1.FirstChild
	entityUrl, found := getAttr(entityLink, "href")
	if !found {
		return entity, errors.New("couldn't find href in entityUrl")
	}

	entity.Url = entityUrl
	entity.Name = entityLink.FirstChild.Data
	if column2.FirstChild != nil {
		entity.Description = strings.TrimSpace(column2.FirstChild.Data)
	}
	if column3.FirstChild != nil && column3.FirstChild.FirstChild != nil {
		entity.Version = column3.FirstChild.FirstChild.Data
	}
	return
}
func GetHomePage(fetcher docs_fetcher.DocsFetcher) (home core.HomePage, err error) {
	pageReader, err := fetcher.GetInitialPage()
	if err != nil {
		return home, err
	}
	home.DocsUrl = fetcher.InitialPageUrl()
	doc, err := parsePage(pageReader)
	if err != nil {
		return home, err
	}
	var feature *core.Feature
	var entities *[]core.Entity

	for _, row := range cascadia.QueryAll(doc, tableRowSel) {
		is_header := nodeHasClass(row, "active")
		columns := getChildrenNodes(row)

		if is_header {
			if feature != nil {
				home.Features = append(home.Features, *feature)
			}
			feature, err = parseHomePageHeader(columns)
			if err != nil {
				return home, err
			}
			entities = &feature.MainEntities
			continue
		}
		if feature == nil {
			return home, errors.New("unexpected error, rows sequence did not start with header")
		}

		if len(columns) == 1 {
			entities = &feature.AuxiliaryEntites
			continue
		}

		entity, err := parseHomePageRow(columns)
		if err != nil {
			return home, err
		}
		*entities = append(*entities, entity)
	}
	if feature != nil {
		home.Features = append(home.Features, *feature)
	}
	return home, err
}
