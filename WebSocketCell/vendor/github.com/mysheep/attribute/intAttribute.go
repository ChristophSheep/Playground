package attribute

import (
	"fmt"
	"math"
)

/*
;; ----------------------------------------------------------------------------
;; Numeric Attribute
;; ----------------------------------------------------------------------------

(deftem (numeric-attr base-attr)  ; like class bar : foo { }

    value           0
    num-type        'int

    min             'empty
    max             'empty
)
*/

type intAttribute struct {
	attribute
	min int // nil
	max int // nil
}

// IntAttribute TODO
type IntAttribute interface {
	Attribute
	Min() int
	Max() int
}

// String - Implement Stringer interface
func (i intAttribute) String() string {
	return fmt.Sprintf("{%v, Min:%v, Max:%v}", i.attribute.String(), i.Min(), i.Max())
}

// CreateIntAttribute TODO
func CreateIntAttribute(name string) IntAttribute {

	attr := createAttribute(name)

	attr.SetValue(0) // TODO: Rethink

	iAttr := intAttribute{
		attribute: attr,
		min:       math.MinInt64,
		max:       math.MaxInt64,
	}

	return iAttr
}

// Min
func (i intAttribute) Min() int {
	return i.min
}

// SetMin
func (i intAttribute) SetMin(min int) {
	i.min = min
}

// Max
func (i intAttribute) Max() int {
	return i.max
}

// SetMax
func (i intAttribute) SetMax(max int) {
	i.max = max
}
