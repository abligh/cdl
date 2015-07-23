package cdl

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type Template map[string]interface{}

type CompiledTemplate struct {
	s map[string]interface{}
}

type Options map[string]interface{}

type Range struct {
	Min int
	Max int
}

type Array struct {
	name string
	r    Range
}

type Requirement struct {
	mandatory bool
	array     bool
	r         Range
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
}

type ValidatorFunc func(obj interface{}) (err *CdlError)

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

func NewError(t int) *CdlError {
	return &CdlError{Type: t}
}

func NewErrorContext(t int, c string) *CdlError {
	return (&CdlError{Type: t}).AddContext(c)
}

func NewErrorContextQuoted(t int, c string) *CdlError {
	return (&CdlError{Type: t}).AddContextQuoted(c)
}

func (e *CdlError) AddContext(c string) *CdlError {
	e.Context = append(e.Context, c)
	return e
}

func (e *CdlError) AddContextQuoted(c string) *CdlError {
	return e.AddContext(fmt.Sprintf("'%s'", c))
}

func (e *CdlError) SetSupplementary(s string) *CdlError {
	e.Supplementary = s
	return e
}

func (r *Range) describeError(value int) string {
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

func (r *Range) contains(value int) bool {
	return (value >= r.Min || r.Min == -1) && (value <= r.Max || r.Max == -1)
}

func makeOptions(optString string) (*Options, *CdlError) {
	opts := make(Options)
	spaceOrBar := func(r rune) bool {
		return unicode.IsSpace(r) || (r == '|')
	}
	for _, o := range strings.FieldsFunc(optString, spaceOrBar) {
		s := regexp.MustCompile("^(\\w+)(.*)$").FindStringSubmatch(o)
		if len(s) < 3 || s[1] == "" {
			return nil, NewErrorContextQuoted(ErrBadOptionValue, o)
		}
		req := Requirement{mandatory: true, array: false, r: Range{-1, -1}}
		if s[2] != "" {
			optslice := regexp.MustCompile("[*+!?]|\\{\\d+,\\d+\\}").FindAllStringSubmatch(s[2], -1)
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
					req.r = Range{1, -1}
					req.array = true
				case c[0] == "*":
					req.array = true
					req.r = Range{0, -1}
				case strings.HasPrefix(c[0], "{"):
					minMax := regexp.MustCompile("^\\{(\\d+),(\\d+)\\}$").FindStringSubmatch(c[0])
					if len(minMax) != 3 {
						return nil, NewErrorContextQuoted(ErrBadRangeOptionModifier, o)
					}
					min, err1 := strconv.Atoi(minMax[1])
					max, err2 := strconv.Atoi(minMax[2])
					if (err1 != nil) || (err2 != nil) || (min > max) {
						return nil, NewErrorContextQuoted(ErrBadRangeOptionModifierValue, o)
					}
					req.array = true
					req.r = Range{min, max}
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
				rng := Range{-1, -1}
				minMax := regexp.MustCompile("^(\\w+)(\\{(\\d+),(\\d+)\\})?$").FindStringSubmatch(arr)
				if len(minMax) != 5 {
					return nil, NewErrorContextQuoted(ErrBadRangeOptionModifier, arr)
				}
				if minMax[3] != "" && minMax[4] != "" {
					min, err1 := strconv.Atoi(minMax[3])
					max, err2 := strconv.Atoi(minMax[4])
					if (err1 != nil) || (err2 != nil) || (min > max) {
						return nil, NewErrorContextQuoted(ErrBadRangeOptionModifierValue, arr)
					}
					rng = Range{min, max}
				}
				ct.s[k] = &Array{name: minMax[1], r: rng}
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
		case *Options:
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

func (ct *CompiledTemplate) validateRange(o interface{}, pos string, r Range) *CdlError {
	slice, ok := o.([]interface{})
	if !ok {
		return NewError(ErrExpectedArray)
	}
	if !r.contains(len(slice)) {
		return NewError(ErrOutOfRange).SetSupplementary(r.describeError(len(slice)))
	}
	for i, v := range slice {
		if err := ct.validateItem(v, pos); err != nil {
			return err.AddContext(fmt.Sprintf("index %d", i))
		}
	}
	return nil
}

func (ct *CompiledTemplate) validateMap(o interface{}, pos string, opts *Options) *CdlError {
	m, ok := o.(map[string]interface{})
	if !ok {
		return NewError(ErrExpectedMap)
	}
	mand := make(map[string]bool)
	for k, v := range *opts {
		switch t := v.(type) {
		case Requirement:
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
			case Requirement:
				if t.array {
					if err := ct.validateRange(v, k, t.r); err != nil {
						return err.AddContextQuoted(k)
					}
				} else {
					if err := ct.validateItem(v, k); err != nil {
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

func (ct *CompiledTemplate) validateItem(o interface{}, pos string) *CdlError {
	if val, ok := ct.s[pos]; !ok {
		return NewError(ErrUnknownKey)
	} else {
		switch t := val.(type) {
		case ValidatorFunc:
			return t(o)
		case *Options:
			return ct.validateMap(o, pos, t)
		case *Array:
			return ct.validateRange(o, t.name, t.r)
		case string:
			ok := false
			switch t {
			case "number":
				switch o.(type) {
				case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
					ok = true
				}
			case "integer":
				switch o.(type) {
				case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
					ok = true
				}
			default:
				otype := fmt.Sprintf("%T", o)
				if otype == t {
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

func (ct *CompiledTemplate) Validate(o interface{}) error {
	if err := ct.validateItem(o, "/"); err != nil {
		return err
	}
	return nil
}
