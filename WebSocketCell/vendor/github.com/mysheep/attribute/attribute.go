package attribute

import "fmt"

/*
;; ----------------------------------------------------------------------------
;; Base Attribute
;; ----------------------------------------------------------------------------

(deftem base-attr

    name            "unknown"
    value           'empty

    is-readable     't
    is-writeable    't
    is-storeable    't
)
*/

type any interface{}

/*
Attribute - Base attribute
*/
type attribute struct {
	name        string
	value       any
	isReadable  bool
	isWriteable bool
	isStoreable bool
}

/*
Attribute - Interface of attribute
*/
type Attribute interface {
	Name() string
	Value() any
	IsReadable() bool
	IsWriteable() bool
	IsStoreable() bool
}

func createAttribute(name string) attribute {
	attr := attribute{
		name:        name,
		isReadable:  true,
		isWriteable: true,
		isStoreable: true,
	}
	return attr
}

// Empty attribute
var Empty = attribute{
	name:        "Empty",
	isReadable:  false,
	isWriteable: false,
	isStoreable: false,
}

// String
func (a attribute) String() string {
	return fmt.Sprintf("Name:\"%v\", Value:%v, Readable:%v, Writeable:%v, Storeable:%v", a.Name(), a.Value(), a.IsReadable(), a.IsWriteable(), a.IsStoreable())
}

// Name
func (a attribute) Name() string {
	return a.name
}

// SetName
func (a attribute) SetName(name string) {
	a.name = name
}

// Value
func (a attribute) Value() any {
	return a.value
}

// GetValue
func (a attribute) SetValue(value any) {
	a.value = value
}

// IsReadable
func (a attribute) IsReadable() bool {
	return a.isReadable
}

// SetReadable
func (a attribute) SetReadable(isReadable bool) {
	a.isReadable = isReadable
}

// IsWriteable
func (a attribute) IsWriteable() bool {
	return a.isWriteable
}

// SetWriteable
func (a attribute) SetWriteable(isWriteable bool) {
	a.isWriteable = isWriteable
}

// IsStoreable
func (a attribute) IsStoreable() bool {
	return a.isStoreable
}

// SetStoreable
func (a attribute) SetStoreable(isStoreable bool) {
	a.isStoreable = isStoreable
}
