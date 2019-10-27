package main

import (
	"fmt"

	"github.com/mysheep/props"
)

//-----------------------------------------------------------------------------
// Person which implements "object with properties" interface
//-----------------------------------------------------------------------------

type Person struct {
	firstname string
	lastname  string
	age       int
}

func MakeEmtpyPerson() *Person {
	return &Person{}
}

func MakePerson(spec *Person) *Person {
	return &Person{
		firstname: spec.firstname,
		lastname:  spec.lastname,
		age:       spec.age,
	}
}

func (p *Person) Properties() []props.Property {
	props := make([]props.Property, 0)

	props = append(props, MakeNameProp("FirstName", p.firstname))
	props = append(props, MakeNameProp("LastName", p.lastname))
	props = append(props, MakeNumberProp("Age", p.age)) // TODO: Emtpy Value, Default Value

	return props
}

//-----------------------------------------------------------------------------
// Name Property
//-----------------------------------------------------------------------------

type NameProp struct {
	name         string
	value        string
	defaultValue string
	isMandatary  bool
	hasValue     bool
}

func MakeNameProp(name string, value string) props.Property {
	return &NameProp{
		name:        name,
		value:       value,
		isMandatary: false,
		hasValue:    true} // TODO: Emtpy Property, DefaultValue, Mandarty Properties
}

func (n *NameProp) Name() string {
	return n.name
}
func (n *NameProp) Value() interface{} {
	if n.IsMandatary() && !n.HasValue() {
		return n.defaultValue
	}
	return n.value
}

func (n *NameProp) DefaultValue() interface{} {
	return n.defaultValue
}

func (n *NameProp) IsMandatary() bool {
	return n.isMandatary
}

func (n *NameProp) HasValue() bool {
	return n.hasValue
}

func (n *NameProp) String() string {
	return fmt.Sprintf("%15s: [%10s]", n.Name(), n.Value())
}

//-----------------------------------------------------------------------------
// Number Property
//-----------------------------------------------------------------------------

type NumberProp struct {
	name         string
	value        int
	defaultValue string
	isMandatary  bool
	hasValue     bool
}

func MakeNumberProp(name string, value int) props.Property {
	return &NumberProp{
		name:     name,
		value:    value,
		hasValue: true}
}

func (n *NumberProp) Name() string {
	return n.name
}
func (n *NumberProp) Value() interface{} {
	if n.IsMandatary() && !n.HasValue() {
		return n.defaultValue
	}
	return n.value
}

func (n *NumberProp) DefaultValue() interface{} {
	return n.defaultValue
}

func (n *NumberProp) IsMandatary() bool {
	return n.isMandatary
}

func (n *NumberProp) HasValue() bool {
	return n.hasValue
}

func (n *NumberProp) String() string {
	if n.HasValue() {
		return fmt.Sprintf("%15s: [%10d]", n.Name(), n.Value())
	}
	return fmt.Sprintf("%15s: [%10s]", n.Name(), "")

}

//-----------------------------------------------------------------------------
// Main
//-----------------------------------------------------------------------------

func printProps(p props.ObjWithProps) {
	for index, prop := range p.Properties() {
		fmt.Printf("%d %v\n", index, prop.String())
	}
}

// Optional Value
type IntValue struct {
	val      int
	hasValue bool
}

func (i *IntValue) HasValue() bool {
	return i.hasValue
}

func (i *IntValue) Value() int {
	if i.HasValue() == false {
		panic("value is not set")
	}
	return i.val
}

func NewIntVal(val int) IntValue {
	return IntValue{val, true}
}

func EmptyInVal() IntValue {
	return IntValue{hasValue: false}
}

type foo struct {
	name   string
	age    IntValue
	size   IntValue
	weight IntValue
}

func main() {
	p := MakePerson(&Person{firstname: "Noah", lastname: "Wilson", age: 21})
	printProps(p)
	fmt.Println()
	pEmpty := MakeEmtpyPerson()
	printProps(pEmpty)

	f := foo{name: "foo", age: NewIntVal(6), size: EmptyInVal()}
	fmt.Printf("foo %v", f)
}
