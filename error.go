package cdl

import (
	"fmt"
	"strings"
)

type CdlError struct {
	Type          Enum
	Supplementary string
	Context       []string
}

// var ErrorEnum is the Enum containing cdl errors.
var (
	ErrorEnum = NewEnumTypeWithText(map[string]string{
		"ErrInternal":                    "Internal error",
		"ErrMissingRoot":                 "No root key in template",
		"ErrBadOptionValue":              "Bad option value",
		"ErrBadRangeOptionModifier":      "Bad range option modifer",
		"ErrBadRangeOptionModifierValue": "Bad range option modifier value",
		"ErrBadOptionModifier":           "Bad option modifier",
		"ErrBadKey":                      "Bad key",
		"ErrBadValue":                    "Bad value",
		"ErrUnknownKey":                  "Unknown key",
		"ErrExpectedMap":                 "Expected map",
		"ErrExpectedArray":               "Expected array",
		"ErrOutOfRange":                  "Number of array items outside permissible range",
		"ErrBadType":                     "Bad type",
		"ErrMissingMandatory":            "Missing mandatory key",
		"ErrBadConfigurator":             "Bad configurator",
		"ErrBadEnumValue":                "Bad option",
	})
)

// func Error implements the Error() function of the error interface.
//
// An error string is returned in context.
func (e CdlError) Error() string {
	main := e.Type.Text()
	if e.Supplementary != "" {
		main = fmt.Sprintf("%s; %s", main, e.Supplementary)
	}
	if len(e.Context) == 0 {
		return fmt.Sprintf("%s (code %s)", main, e.Type.String())
	} else {
		return fmt.Sprintf("%s (code %s) near %s", main, e.Type.String(), strings.Join(e.Context, " at "))
	}
}

// func NewError returns a new CdlError of a given type.
//
// The type should be a type starting with `Err` in the constants section.
func NewError(t string) *CdlError {
	return &CdlError{Type: ErrorEnum.New(t)}
}

// func NewErrorContext creates a new CdlError with the specified context string.
//
// The type should be a type starting with `Err` in the constants section.
func NewErrorContext(t string, c string) *CdlError {
	return (&CdlError{Type: ErrorEnum.New(t)}).AddContext(c)
}

// func NewErrorContext creates a new CdlError with the specified context string.
//
// The type should be a type starting with `Err` in the constants section.
// The context string will be quoted.
func NewErrorContextQuoted(t string, c string) *CdlError {
	return (&CdlError{Type: ErrorEnum.New(t)}).AddContextQuoted(c)
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
