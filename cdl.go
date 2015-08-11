package cdl

import (
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// type Template is a user-provided uncompiled template.
//
// See the overview for how these work.
type Template map[string]interface{}

// type Configurator is a map of Configurator functions
type Configurator map[string]interface{}

// type CompiledTemplate is a compiled template.
//
// It is opaque to the user in operations.
type CompiledTemplate struct {
	s map[string]interface{}
}

type options map[string]interface{}

type optrange struct {
	Min int
	Max int
}

type array struct {
	name string
	r    optrange
}

type requirement struct {
	mandatory bool
	array     bool
	r         optrange
}

type CdlError struct {
	Type          int
	Supplementary string
	Context       []string
}

const (
	_                              = iota
	ErrInternal                    = iota
	ErrMissingRoot                 = iota
	ErrBadOptionValue              = iota
	ErrBadRangeOptionModifier      = iota
	ErrBadRangeOptionModifierValue = iota
	ErrBadOptionModifier           = iota
	ErrBadKey                      = iota
	ErrBadValue                    = iota
	ErrUnknownKey                  = iota
	ErrExpectedMap                 = iota
	ErrExpectedArray               = iota
	ErrOutOfRange                  = iota
	ErrBadType                     = iota
	ErrMissingMandatory            = iota
	ErrBadConfigurator             = iota
)

var ErrorMap map[int]string = map[int]string{
	ErrInternal:                    "Internal error",
	ErrMissingRoot:                 "No root key in template",
	ErrBadOptionValue:              "Bad option value",
	ErrBadRangeOptionModifier:      "Bad range option modifer",
	ErrBadRangeOptionModifierValue: "Bad range option modifier value",
	ErrBadOptionModifier:           "Bad option modifier",
	ErrBadKey:                      "Bad key",
	ErrBadValue:                    "Bad value",
	ErrUnknownKey:                  "Unknown key",
	ErrExpectedMap:                 "Expected map",
	ErrExpectedArray:               "Expected array",
	ErrOutOfRange:                  "Number of array items outside permissible range",
	ErrBadType:                     "Bad type",
	ErrMissingMandatory:            "Missing mandatory key",
	ErrBadConfigurator:             "Bad configurator",
}

// type ValidatorFunc allows user specified validation functions to be passed to cdl.
type ValidatorFunc func(obj interface{}) (err *CdlError)

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

// type ConfiguratorFunc allows user specified configurator functions to be passed to cdl.
type ConfiguratorFunc func(obj interface{}, path Path) (err *CdlError)

// func Error implements the Error() function of the error interface.
//
// An error string is returned in context.
func (e CdlError) Error() string {
	main := ErrorMap[e.Type]
	if e.Supplementary != "" {
		main = fmt.Sprintf("%s; %s", main, e.Supplementary)
	}
	if len(e.Context) == 0 {
		return fmt.Sprintf("%s (code %d)", main, e.Type)
	} else {
		return fmt.Sprintf("%s (code %d) near %s", main, e.Type, strings.Join(e.Context, " at "))
	}
}

// func NewError returns a new CdlError of a given type.
//
// The type should be a type starting with `Err` in the constants section.
func NewError(t int) *CdlError {
	return &CdlError{Type: t}
}

// func NewErrorContext creates a new CdlError with the specified context string.
//
// The type should be a type starting with `Err` in the constants section.
func NewErrorContext(t int, c string) *CdlError {
	return (&CdlError{Type: t}).AddContext(c)
}

// func NewErrorContext creates a new CdlError with the specified context string.
//
// The type should be a type starting with `Err` in the constants section.
// The context string will be quoted.
func NewErrorContextQuoted(t int, c string) *CdlError {
	return (&CdlError{Type: t}).AddContextQuoted(c)
}

// func AddContext adds the specified context to an existing cdl error.
func (e *CdlError) AddContext(c string) *CdlError {
	e.Context = append(e.Context, c)
	return e
}

// func AddContextQuoted adds the specified context to an existing cdl error.
//
// The context will be quoted.
func (e *CdlError) AddContextQuoted(c string) *CdlError {
	return e.AddContext(fmt.Sprintf("'%s'", c))
}

// func SetSupplementary adds the specified supplementary data to an existing cdl error.
func (e *CdlError) SetSupplementary(s string) *CdlError {
	e.Supplementary = s
	return e
}

func (r *optrange) describeError(value int) string {
	min := r.Min
	if min < 0 {
		min = 0
	}
	if r.Max < 0 {
		return fmt.Sprintf("got %d, expecting at least %d", value, min)
	} else {
		return fmt.Sprintf("got %d, expecting between %d and %d", value, min, r.Max)
	}
}

func (r *optrange) contains(value int) bool {
	return (value >= r.Min || r.Min == -1) && (value <= r.Max || r.Max == -1)
}

func makeOptions(optString string) (*options, *CdlError) {
	opts := make(options)
	spaceOrBar := func(r rune) bool {
		return unicode.IsSpace(r) || (r == '|')
	}
	for _, o := range strings.FieldsFunc(optString, spaceOrBar) {
		s := regexp.MustCompile("^(\\w+)(.*)$").FindStringSubmatch(o)
		if len(s) < 3 || s[1] == "" {
			return nil, NewErrorContextQuoted(ErrBadOptionValue, o)
		}
		req := requirement{mandatory: true, array: false, r: optrange{-1, -1}}
		if s[2] != "" {
			optslice := regexp.MustCompile("[*+!?]|\\{\\d+,\\d*\\}").FindAllStringSubmatch(s[2], -1)
			if len(optslice) == 0 {
				return nil, NewErrorContextQuoted(ErrBadOptionModifier, o)
			}
			for _, c := range optslice {
				if len(c) != 1 {
					return nil, NewErrorContextQuoted(ErrBadOptionModifier, o)
				}
				switch {
				case c[0] == "?":
					req.mandatory = false
				case c[0] == "!":
					req.mandatory = true
				case c[0] == "+":
					req.r = optrange{1, -1}
					req.array = true
				case c[0] == "*":
					req.array = true
					req.r = optrange{0, -1}
				case strings.HasPrefix(c[0], "{"):
					minMax := regexp.MustCompile("^\\{(\\d+),(\\d*)\\}$").FindStringSubmatch(c[0])
					if len(minMax) != 3 {
						return nil, NewErrorContextQuoted(ErrBadRangeOptionModifier, o)
					}
					min, err1 := strconv.Atoi(minMax[1])
					if err1 != nil {
						return nil, NewErrorContextQuoted(ErrBadRangeOptionModifierValue, o)
					}
					max := -1
					if minMax[2] != "" {
						max, err2 := strconv.Atoi(minMax[2])
						if (err2 != nil) || (min > max) {
							return nil, NewErrorContextQuoted(ErrBadRangeOptionModifierValue, o)
						}
					}
					req.array = true
					req.r = optrange{min, max}
				default:
					return nil, NewErrorContextQuoted(ErrBadOptionModifier, o)
				}
			}
		}
		opts[s[1]] = req
	}

	return &opts, nil
}

func newCompiledTemplate() *CompiledTemplate {
	return &CompiledTemplate{s: make(map[string]interface{})}
}

// func Compile compiles a specified cdl template.
func Compile(t Template) (*CompiledTemplate, error) {
	ct := newCompiledTemplate()
	for k, v := range t {
		if match, err := regexp.MatchString("^(/|(\\w+))?$", k); !match || err != nil {
			return nil, NewErrorContextQuoted(ErrBadKey, k)
		}
		switch t := v.(type) {
		case string:
			if t == "" {
				t = "/"
			}
			switch {
			case strings.HasPrefix(t, "{}"):
				if o, err := makeOptions(strings.TrimPrefix(t, "{}")); err != nil {
					return nil, err.AddContextQuoted(k)
				} else {
					ct.s[k] = o
				}
			case strings.HasPrefix(t, "[]"):
				arr := strings.TrimPrefix(t, "[]")
				rng := optrange{-1, -1}
				minMax := regexp.MustCompile("^(\\w+)(\\{(\\d+),(\\d*)\\})?$").FindStringSubmatch(arr)
				if len(minMax) != 5 {
					return nil, NewErrorContextQuoted(ErrBadRangeOptionModifier, arr)
				}
				if minMax[3] != "" {
					min, err1 := strconv.Atoi(minMax[3])
					if err1 != nil {
						return nil, NewErrorContextQuoted(ErrBadRangeOptionModifierValue, arr)
					}
					max := -1
					if minMax[4] != "" {
						var err2 error
						max, err2 = strconv.Atoi(minMax[4])
						if (err2 != nil) || (min > max) {
							return nil, NewErrorContextQuoted(ErrBadRangeOptionModifierValue, arr)
						}
					}
					rng = optrange{min, max}
				}
				ct.s[k] = &array{name: minMax[1], r: rng}
			default:
				ct.s[k] = t
			}
		case ValidatorFunc:
			ct.s[k] = t
		case func(interface{}) *CdlError: // in case they didn't cast it
			ct.s[k] = ValidatorFunc(t)
		default:
			return nil, NewErrorContextQuoted(ErrBadValue, fmt.Sprintf("%T", t)).AddContextQuoted(k)
		}
	}
	for _, v := range ct.s {
		switch t := v.(type) {
		case *options:
			for optk, _ := range *t {
				if _, ok := ct.s[optk]; !ok {
					ct.s[optk] = 0 // autodiscovered
				}
			}
		}
	}
	if _, ok := ct.s["/"]; !ok {
		return nil, NewError(ErrMissingRoot)
	}
	return ct, nil
}

// MustCompile is like Compile but panics if the expression cannot be parsed.
// It simplifies safe initialization of global variables holding compiled templates
func MustCompile(t Template) *CompiledTemplate {
	ct, error := Compile(t)
	if error != nil {
		panic(`cdl: Compile failed: ` + error.Error())
	}
	return ct
}

func (ct *CompiledTemplate) validateRange(o interface{}, pos string, r optrange, configurator Configurator, path Path) *CdlError {
	slice, ok := o.([]interface{})
	if !ok {
		return NewError(ErrExpectedArray)
	}
	if !r.contains(len(slice)) {
		return NewError(ErrOutOfRange).SetSupplementary(r.describeError(len(slice)))
	}
	for i, v := range slice {
		if err := ct.validateAndConfigureItem(v, pos, configurator, path.push(i)); err != nil {
			return err.AddContext(fmt.Sprintf("index %d", i))
		}
	}
	return nil
}

func (ct *CompiledTemplate) validateMap(o interface{}, pos string, opts *options, configurator Configurator, path Path) *CdlError {
	m, ok := o.(map[string]interface{})
	if !ok {
		return NewError(ErrExpectedMap)
	}
	mand := make(map[string]bool)
	for k, v := range *opts {
		switch t := v.(type) {
		case requirement:
			if t.mandatory {
				mand[k] = true
			}
		}
	}
	for k, v := range m {
		if o, ok := (*opts)[k]; !ok {
			return NewErrorContextQuoted(ErrBadKey, k)
		} else {
			switch t := o.(type) {
			case requirement:
				if t.array {
					if err := ct.validateRange(v, k, t.r, configurator, path.push(k)); err != nil {
						return err.AddContextQuoted(k)
					}
				} else {
					if err := ct.validateAndConfigureItem(v, k, configurator, path.push(k)); err != nil {
						return err.AddContextQuoted(k)
					}
				}
				if t.mandatory {
					delete(mand, k)
				}
			}
		}
	}
	if len(mand) != 0 {
		missing := make([]string, len(mand))
		i := 0
		for k := range mand {
			missing[i] = fmt.Sprintf("'%s'", k)
			i++
		}
		return NewError(ErrMissingMandatory).SetSupplementary(fmt.Sprintf("missing %s", strings.Join(missing, ", ")))
	}
	return nil
}

func (ct *CompiledTemplate) validateItem(o interface{}, pos string, configurator Configurator, path Path) *CdlError {
	if val, ok := ct.s[pos]; !ok {
		return NewError(ErrUnknownKey)
	} else {
		switch t := val.(type) {
		case ValidatorFunc:
			return t(o)
		case *options:
			return ct.validateMap(o, pos, t, configurator, path)
		case *array:
			return ct.validateRange(o, t.name, t.r, configurator, path)
		case string:
			ok := false
			switch t {
			case "number":
				switch o.(type) {
				case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
					ok = true
				}
			case "integer":
				switch n := o.(type) {
				case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
					ok = true
				case float64:
					if n == float64(int(n)) {
						ok = true
					}
				case float32:
					if n == float32(int(n)) {
						ok = true
					}
				}
			case "ipport":
				switch n := o.(type) {
				case string:
					if _, _, err := net.SplitHostPort(n); err == nil {
						ok = true
					}
				}
			default:
				if reflect.TypeOf(o).String() == t {
					ok = true
				}
			}
			if !ok {
				return NewError(ErrBadType).SetSupplementary(fmt.Sprintf("got %T expected %s", o, t))
			}
		case int:
			// autodiscovered
		default:
			return NewError(ErrInternal).SetSupplementary(fmt.Sprintf("type is neither validator func nor options: %T", val))
		}
	}
	return nil
}

func assign(ptr interface{}, obj interface{}) *CdlError {
	p := reflect.ValueOf(ptr)

	switch p.Kind() {
	case reflect.Ptr:
		v := p.Elem()
		if v.Type() != reflect.TypeOf(obj) {
			return NewError(ErrBadType).SetSupplementary(fmt.Sprintf("at configuration got %s expected %s",
				v.Type().String(),
				reflect.TypeOf(obj).String()))
		}
		v.Set(reflect.ValueOf(obj))
		return nil
	default:
		return NewError(ErrBadConfigurator).SetSupplementary("got object that is not a pointer")
	}
}

func (ct *CompiledTemplate) validateAndConfigureItem(o interface{}, pos string, configurator Configurator, path Path) *CdlError {
	if err := ct.validateItem(o, pos, configurator, path); err != nil {
		return err
	}
	if configurator != nil {
		if cnf, ok := configurator[pos]; ok && (cnf != nil) {
			if val, ok := ct.s[pos]; !ok {
				return NewError(ErrUnknownKey)
			} else {
				v := o
				switch t := val.(type) {
				case string:
					switch t {
					case "number":
						switch n := o.(type) {
						// Go unhelpfully does not allow casting with a multiple case type assertion
						case int:
							v = float64(n)
						case int8:
							v = float64(n)
						case int16:
							v = float64(n)
						case int32:
							v = float64(n)
						case int64:
							v = float64(n)
						case uint:
							v = float64(n)
						case uint8:
							v = float64(n)
						case uint16:
							v = float64(n)
						case uint32:
							v = float64(n)
						case uint64:
							v = float64(n)
						case float32:
							v = float64(n)
						case float64:
							v = float64(n)
						}
					case "integer":
						switch n := o.(type) {
						// Go unhelpfully does not allow casting with a multiple case type assertion
						case int:
							v = int(n)
						case int8:
							v = int(n)
						case int16:
							v = int(n)
						case int32:
							v = int(n)
						case int64:
							v = int(n)
						case uint:
							v = int(n)
						case uint8:
							v = int(n)
						case uint16:
							v = int(n)
						case uint32:
							v = int(n)
						case uint64:
							v = int(n)
						case float32:
							v = int(n)
						case float64:
							v = int(n)
						}
					}
				}
				switch t := cnf.(type) {
				case ConfiguratorFunc:
					return t(v, path)
				case func(interface{}, Path) *CdlError: // in case they didn't cast it
					return t(v, path)
				default:
					if reflect.ValueOf(cnf).Kind() == reflect.Ptr {
						if err := assign(cnf, v); err != nil {
							return err
						}
					} else {
						return NewError(ErrBadConfigurator).SetSupplementary("got unknown type")
					}
				}
			}
		}
	}
	return nil
}

// func Validate validates an object against a cdl template.
//
// Optionally a configurator may be passed. This can be nil if you do not need configurator functions calling
func (ct *CompiledTemplate) Validate(o interface{}, configurator Configurator) error {
	path := Path{}
	if err := ct.validateAndConfigureItem(o, "/", configurator, path); err != nil {
		return err
	}
	return nil
}
