package props

import "fmt"

// see packages "flag"

// Value is the interface to the dynamic value stored in a flag.
// (The default value is represented as a string.)
//
// If a Value has an IsBoolFlag() bool method returning true,
// the command-line parser makes -name equivalent to -name=true
// rather than using the next command-line argument.
//
// Set is called once, in command line order, for each flag present.
// The flag package may call the String method with a zero-valued receiver,
// such as a nil pointer.

/*
type Value interface {
	String() string
	Set(string) error
}
*/

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
