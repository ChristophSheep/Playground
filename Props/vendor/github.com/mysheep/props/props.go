package props

import "fmt"

//-----------------------------------------------------------------------------
// Interfaces
//-----------------------------------------------------------------------------

type Property interface {
	Name() string

	Value() interface{}
	DefaultValue() interface{}

	IsMandatary() bool
	HasValue() bool

	fmt.Stringer // ToString in other languages
}

type ObjWithProps interface {
	Properties() []Property
}
