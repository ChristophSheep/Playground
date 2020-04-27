package attribute

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

/*
NumAttribute - Numeric Attribute
*/
type numAttribute struct {
	Min int // nil
	Max int // nil
	Attribute
}

/*
NumAttribute - Interface of attribute
see https://golang.org/pkg/math/big/
*/
type NumAttribute interface {
	Attribute
	NumType() string
	Min() interface{}
	Max() interface{}
}
