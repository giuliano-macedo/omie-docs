package bruno

import (
	"github.com/giuliano-macedo/go-bruno-collection/pkg/bruno"
	"github.com/giuliano-macedo/omie-docs/internal/bundler"
	"github.com/giuliano-macedo/omie-docs/internal/core"
)

type BrunonBundler struct{}

func (b *BrunonBundler) Bundle(args bundler.Args) error {
	collection := ConvertToBruno(args.Home, args.Pages)
	return writeBrucollection(args.FsWriter, collection)
}

func writeBrucollection(fSWriter bundler.FSWriter, collection bruno.Collection) error {
	file, err := fSWriter.Create(core.BrunoCollectionName)
	if err != nil {
		return err
	}
	defer file.Close()
	return bruno.CreateBruTar(collection, file)
}

func NewBrunoBundler() bundler.Bundler {
	return &BrunonBundler{}
}
