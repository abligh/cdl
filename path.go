package cdl

import (
	"fmt"
	"strings"
)

// type Path is an array of items constituting the path to an item to be checked for configuration
type Path struct {
	items []interface{}
}

func (p *Path) push(o interface{}) Path {
	return Path{items: append(p.items, o)}
}

// func Slice returns a slice of objects representing the path.
//
// The objects may be strings or integers
func (p *Path) Slice() []interface{} {
	return p.items
}

// func StringSlice returns a slice of strings representing a path
func (p *Path) StringSlice() []string {
	ss := make([]string, len(p.items))
	for i, v := range p.items {
		switch s := v.(type) {
		case string:
			ss[i] = s
		case int:
			ss[i] = fmt.Sprintf("%d", s)
		default:
			ss[i] = fmt.Sprintf("%v", s)
		}
	}
	return ss
}

// func String produces a string representation of a path
//
// The path elements are separated by '/'
func (p Path) String() string {
	return "/" + strings.Join(p.StringSlice(), "/")
}
