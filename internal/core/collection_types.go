package core

type CollectionType uint

const (
	OpenAPI CollectionType = iota
	Postman
	Bruno
)

var collectionTypeNames = []string{OpenAPI: "OpenAPI", Postman: "Postman", Bruno: "Bruno"}

func (collectionType CollectionType) String() string {
	return collectionTypeNames[collectionType]
}
