package postman

import (
	"github.com/giuliano-macedo/omie-docs/internal/bundler"
	"github.com/giuliano-macedo/omie-docs/internal/core"
)

type PostmanBundler struct{}

func (p *PostmanBundler) Bundle(args bundler.Args) error {
	collection := convertToPostman(args.Home, args.Pages)
	return args.FsWriter.SaveJsonFile(core.PostmanCollectionName, collection)
}

func NewPostmanBundler() bundler.Bundler {
	return &PostmanBundler{}
}
