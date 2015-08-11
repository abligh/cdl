package cdl

import (
	"fmt"
	"sort"
)

// type EnumType represents an enum type within cdl
//
// Each enum type must be initialised exactly once. To initialise use something like
//
//    var myEnumType = cdl.NewEnumType("DEFAULT_VALUE", "ONE_VALUE", "ANOTHER_VALUE")
type EnumType struct {
	toValue  map[string]int
	toString []string
	toText   []string
	items    int
}

// func NewEnumType produces a new EnumType for a given list of enumeration constants
func NewEnumType(values ...string) EnumType {
	var et EnumType = EnumType{items: len(values)}
	et.toValue = make(map[string]int, et.items)
	et.toString = make([]string, et.items)
	et.toText = make([]string, et.items)
	for i, v := range values {
		et.toValue[v] = i
		et.toString[i] = v
	}
	return et
}

// func NewEnumType produces a new EnumType for a given list of enumeration constants
func NewEnumTypeWithText(values map[string]string) EnumType {
	var et EnumType = EnumType{items: len(values)}
	et.toValue = make(map[string]int, et.items)
	et.toString = make([]string, et.items)
	et.toText = make([]string, et.items)
	vals := make([]string, et.items)
	i := 0
	for v, _ := range values {
		vals[i] = v
		i++
	}
	sort.Strings(vals)
	for i, v := range vals {
		et.toValue[v] = i
		et.toString[i] = v
		et.toText[i] = values[v]
		i++
	}
	return et
}

type Enum struct {
	Type  *EnumType
	value int
}

// func String produces the string representation of an Enum
func (e Enum) String() string {
	if e.value >= 0 && e.value < e.Type.items {
		return e.Type.toString[e.value]
	}
	panic(fmt.Sprintf("Bad enum value %d (%d)", e.value, e.Type.items))
}

// func String produces the text representation of an Enum
//
// If no text has been specified, the text is the string representation of the item
func (e Enum) Text() string {
	if e.value >= 0 && e.value < e.Type.items {
		if t := e.Type.toText[e.value]; t != "" {
			return t
		}
		return e.Type.toString[e.value]
	}
	panic(fmt.Sprintf("Bad enum value %d (%d)", e.value, e.Type.items))
}

// func Set sets the value of an Enum to a specific value
//
// returns true if setting the value to v succeeded, else false
func (e *Enum) Set(v string) bool {
	if v, ok := e.Type.toValue[v]; ok {
		e.value = v
		return true
	}
	return false
}

// func Has determines whether an Enum could be set to a value
//
// returns true if the value is valid, else false
func (e *Enum) Has(v string) bool {
	_, ok := e.Type.toValue[v]
	return ok
}

// func Has determines whether an EnumType's instance could be set to a value
//
// returns true if the value is valid, else false
func (et *EnumType) Has(v string) bool {
	_, ok := et.toValue[v]
	return ok
}

// func New creates a new enum value
func (et *EnumType) New(v string) Enum {
	if i, ok := et.toValue[v]; ok {
		return Enum{Type: et, value: i}
	} else {
		panic("Bad enum initialiser " + v)
	}
}
