package query

import (
	"reflect"
)

type FieldLinkType uint8

type ModelLink struct {
	FullModelType reflect.Type
	Collection    string
	Resolvers     map[FieldLink]ModelLink
}

const (
	Array FieldLinkType = iota + 1
	Single
	SingleLast
)

type FieldLink struct {
	IdName       string
	LinkIdName   string
	ResolvedName string

	Type FieldLinkType
}
