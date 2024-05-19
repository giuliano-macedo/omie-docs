package core

import "slices"

type HomePage struct {
	DocsUrl  string
	Features []Feature
}

type Feature struct {
	Name             string
	Description      string
	MainEntities     []Entity
	AuxiliaryEntites []Entity
}

func (feature *Feature) AllEntities() []Entity {
	return slices.Concat(feature.MainEntities, feature.AuxiliaryEntites)
}

type Entity struct {
	Name        string
	Description string
	Version     string
	Url         string
}

type Page struct {
	EntityName string
	Name       string
	DocsUrl    string
	Endpoint   string
	BaseUrl    string
	FullUrl    string
	BasePath   string
	ApiVersion string
	Methods    []Method
	Models     []Model
}
type Method struct {
	Name        string
	Description string
	// Model name
	Parameter string
	// model name
	Return       string
	Example      map[string]interface{}
	IsDeprecated bool
}
type Model struct {
	Name   string
	Fields []Field
}

type Type int

const (
	String Type = iota
	Integer
	Text
	Decimal
	Boolean
	Object
	Array
)

type Field struct {
	Name       string
	Tooltip    string
	Type       Type
	IsRequired bool
	TypeName   string
	// only for array
	ElementType  string
	Length       int64
	Description  string
	IsDeprecated bool
}
