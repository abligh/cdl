package cdl_test

import (
	"encoding/json"
	"fmt"
	"github.com/abligh/cdl"
	"log"
	"testing"
)

type checkTemplate map[string]cdl.Template
type checkJson map[string]string

var fruitPart = cdl.NewEnumType("flesh", "pips", "rind")

var checkTemplates checkTemplate = checkTemplate{
	"simple": cdl.Template{
		"/":   "foo",
		"bar": "int",
	},
	"noroot": cdl.Template{
		"x": "foo",
	},
	"badkey": cdl.Template{
		"/": "foo",
		"+": "foo",
	},
	"array1": cdl.Template{
		"/": "[]foo",
	},
	"array2": cdl.Template{
		"/": "[]foo{1,3}",
	},
	"badarray1": cdl.Template{
		"/": "[]",
	},
	"badarray2": cdl.Template{
		"/": "[]!",
	},
	"badarray3": cdl.Template{
		"/": "[]foo{3}",
	},
	"badarray4": cdl.Template{
		"/": "[]foo{3,a}",
	},
	"badarray5": cdl.Template{
		"/": "[]foo {1,3}",
	},
	"badvalue": cdl.Template{
		"/": 1,
	},
	"validator": cdl.Template{
		"/": isOneOrTwo,
	},
	"badvalidator1": cdl.Template{
		"/": dummy,
	},
	"map": cdl.Template{
		"/":     "{}apple peach? pear* plum+ raspberry{1,3} strawberry! kiwi{1,4}? guava!{1,2} orange?{2,}",
		"apple": "int",
		"peach": isOneOrTwo,
	},
	"badmap1": cdl.Template{
		"/": "{}/a",
	},
	"badmap2": cdl.Template{
		"/": "{}a/a",
	},
	"badmap3": cdl.Template{
		"/": "{}apple/",
	},
	"badmap4": cdl.Template{
		"/": "{}apple{1",
	},
	"badmap5": cdl.Template{
		"/": "{}apple{-1,-1}",
	},
	"badmap6": cdl.Template{
		"/": "{}apple{a,1}",
	},
	"badmap7": cdl.Template{
		"/": "{}apple{1,a}",
	},
	"badmap8": cdl.Template{
		"/": "{}apple{3,1}",
	},
	"example": cdl.Template{
		"/":         "{}apple peach? pear* plum+ raspberry{1,3} strawberry! kiwi{1,4}? guava!{1,2} orange?{2,} mango? blueberry? cherry? tangerine?",
		"apple":     "float64",
		"peach":     "number",
		"pear":      "string",
		"plum":      isOneOrTwo,
		"mango":     "[]planet{2,4}",
		"planet":    "{}earth venus? jupiter?",
		"jupiter":   "[]gods",
		"gods":      "{}thor? odin?",
		"blueberry": "{}red yellow?",
		"cherry":    "ipport",
		"tangerine": fruitPart,
	},
	"integernumberstring": cdl.Template{
		"/": "{}i? n? s? u? w? e? f?",
		"n": "number",
		"i": "integer",
		"s": "string",
		"e": fruitPart,
	},
}

var checkJsons checkJson = checkJson{
	"simple1": `
		{
			"apple" : 3,
			"pear" : [],
			"plum" : [ 1 ],
			"raspberry" : [ "a", "b" ],
			"strawberry" : "here",
			"guava": [ "c", "d" ]
		}
	`,
	"simple2": `
		{
			"apple" : 3,
			"peach" : 4.2,
			"pear" : [ "astring" ],
			"plum" : [ 1, 2 ],
			"raspberry" : [ "a", "b", "c" ],
			"strawberry" : "here",
			"kiwi" : [ 1, 2, 3, 4 ],
			"orange" : [ 1, 2, 3, 4, 5 ],
			"guava": [ "d" ]
		}
	`,
	"bad1": `
		{
			"apple" : "notmeanttobeastring",
			"pear" : [],
			"plum" : [ 1 ],
			"raspberry" : [ "a", "b" ],
			"strawberry" : "here",
			"guava": [ "c", "d" ]
		}
	`,
	"bad2": `
		{
			"apple" : 3,
			"pear" : [ 1 ],
			"plum" : [ 1 ],
			"raspberry" : [ "a", "b" ],
			"strawberry" : "here",
			"guava": [ "c", "d" ]
		}
	`,
	"bad3": `
		{
			"apple" : 3,
			"pear" : [],
			"plum" : [ 4 ],
			"raspberry" : [ "a", "b" ],
			"strawberry" : "here",
			"guava": [ "c", "d" ]
		}
	`,
	"mango": `
		{
			"apple" : 3,
			"pear" : [],
			"plum" : [ 1 ],
			"raspberry" : [ "a", "b" ],
			"strawberry" : "here",
			"guava": [ "c", "d" ],
			"mango": [ {"earth" : 1}, {"earth" : 1, "venus" : 1} ]
		}
	`,
	"badmango1": `
		{
			"apple" : 3,
			"pear" : [],
			"plum" : [ 1 ],
			"raspberry" : [ "a", "b" ],
			"strawberry" : "here",
			"guava": [ "c", "d" ],
			"mango": [ {"earth" : 1} ]
		}
	`,
	"badmango2": `
		{
			"apple" : 3,
			"pear" : [],
			"plum" : [ 1 ],
			"raspberry" : [ "a", "b" ],
			"strawberry" : "here",
			"guava": [ "c", "d" ],
			"mango": [ {"earth" : 1}, {"earth" : 1}, {"earth" : 1}, {"earth" : 1}, {"earth" : 1} ]
		}
	`,
	"badmango3": `
		{
			"apple" : 3,
			"pear" : [],
			"plum" : [ 1 ],
			"raspberry" : [ "a", "b" ],
			"strawberry" : "here",
			"guava": [ "c", "d" ],
			"mango": [ 1, 2 ]
		}
	`,
	"badmango4": `
		{
			"apple" : 3,
			"pear" : [],
			"plum" : [ 1 ],
			"raspberry" : [ "a", "b" ],
			"strawberry" : "here",
			"guava": [ "c", "d" ],
			"mango": [ { "foo" : "bar"}, { "foo" : "bar"} ]
		}
	`,
	"jupiter": `
		{
			"apple" : 3,
			"pear" : [],
			"plum" : [ 1 ],
			"raspberry" : [ "a", "b" ],
			"strawberry" : "here",
			"guava": [ "c", "d" ],
			"mango": [
				{
					"earth" : 1
				},
				{
					"earth" : 1,
					"venus" : 1,
					"jupiter": [
						{"thor" : 1},
						{"odin" : 1}
					]
				}
			]
		}
	`,
	"badjupiter1": `
		{
			"apple" : 3,
			"pear" : [],
			"plum" : [ 1 ],
			"raspberry" : [ "a", "b" ],
			"strawberry" : "here",
			"guava": [ "c", "d" ],
			"mango": [
				{
					"earth" : 1
				},
				{
					"earth" : 1,
					"venus" : 1,
					"jupiter": 1
				}
			]
		}
	`,
	"badjupiter2": `
		{
			"apple" : 3,
			"pear" : [],
			"plum" : [ 1 ],
			"raspberry" : [ "a", "b" ],
			"strawberry" : "here",
			"guava": [ "c", "d" ],
			"mango": [
				{
					"earth" : 1
				},
				{
					"earth" : 1,
					"venus" : 1,
					"jupiter": [
						{"wotan" : 1},
						{"odin" : 1}
					]
				}
			]
		}
	`,
	"badjupiter3": `
		{
			"apple" : 3,
			"pear" : [],
			"plum" : [ 1 ],
			"raspberry" : [ "a", "b" ],
			"strawberry" : "here",
			"guava": [ "c", "d" ],
			"mango": [
				{
					"earth" : 1
				},
				{
					"earth" : 1,
					"venus" : 1,
					"jupiter": [
						1
					]
				}
			]
		}
	`,
	"badjupiter4": `
		{
			"apple" : 3,
			"pear" : [],
			"plum" : [ 1 ],
			"raspberry" : [ "a", "b" ],
			"strawberry" : "here",
			"guava": [ "c", "d" ],
			"mango": [
				{
					"earth" : 1
				},
				{
					"earth" : 1,
					"venus" : 1,
					"jupiter": [
						[ 1 ]
					]
				}
			]
		}
	`,
	"blueberry": `
		{
			"apple" : 3,
			"pear" : [],
			"plum" : [ 1 ],
			"raspberry" : [ "a", "b" ],
			"strawberry" : "here",
			"guava": [ "c", "d" ],
			"blueberry": { "red" : 1 }
		}
	`,
	"badblueberry1": `
		{
			"apple" : 3,
			"pear" : [],
			"plum" : [ 1 ],
			"raspberry" : [ "a", "b" ],
			"strawberry" : "here",
			"guava": [ "c", "d" ],
			"blueberry": 1
		}
	`,
	"badblueberry2": `
		{
			"apple" : 3,
			"pear" : [],
			"plum" : [ 1 ],
			"raspberry" : [ "a", "b" ],
			"strawberry" : "here",
			"guava": [ "c", "d" ],
			"blueberry": [ 1 ]
		}
	`,
	"badblueberry3": `
		{
			"apple" : 3,
			"pear" : [],
			"plum" : [ 1 ],
			"raspberry" : [ "a", "b" ],
			"strawberry" : "here",
			"guava": [ "c", "d" ],
			"blueberry": { "green" : 1 }
		}
	`,
	"badblueberry4": `
		{
			"apple" : 3,
			"pear" : [],
			"plum" : [ 1 ],
			"raspberry" : [ "a", "b" ],
			"strawberry" : "here",
			"guava": [ "c", "d" ],
			"blueberry": { "yellow" : 1 }
		}
	`,
	"integernumberstring": `
		{
			"i" : 1,
			"n" : 0.5,
			"s" : "hello",
			"u" : "there",
			"w" : 1,
			"e" : "rind",
			"f" : "rind"
		}
	`,
	"badintegernumberstring1": `
		{
			"i" : 1.1,
			"n" : 0.5,
			"s" : "hello",
			"u" : "there",
			"w" : 1
		}
	`,
	"badintegernumberstring2": `
		{
			"i" : "a string",
			"n" : 0.5,
			"s" : "hello",
			"u" : "there",
			"w" : 1
		}
	`,
	"badintegernumberstring3": `
		{
			"i" : 1,
			"n" : "a string",
			"s" : "hello",
			"u" : "there",
			"w" : 1
		}
	`,
	"badintegernumberstring4": `
		{
			"i" : 1,
			"n" : 0.5,
			"s" : 37,
			"u" : "there",
			"w" : 1
		}
	`,
	"badintegernumberstring5": `
		{
			"i" : 1,
			"n" : 0.5,
			"s" : "hello",
			"u" : 2,
			"w" : 1
		}
	`,
	"badintegernumberstring6": `
		{
			"i" : 1,
			"n" : 0.5,
			"s" : "hello",
			"u" : "there",
			"w" : "notanint"
		}
	`,
	"badintegernumberstring7": `
		{
			"i" : 1,
			"n" : 0.5,
			"s" : "hello",
			"u" : "there",
			"w" : 1,
			"e" : "cerebralcortex"
		}
	`,
	"badintegernumberstring8": `
		{
			"i" : 1,
			"n" : 0.5,
			"s" : "hello",
			"u" : "there",
			"w" : 1,
			"f" : "cerebralcortex"
		}
	`,
	"badintegernumberstring9": `
		{
			"i" : 1,
			"n" : 0.5,
			"s" : "hello",
			"u" : "there",
			"w" : 1,
			"e" : 1
		}
	`,
	"badintegernumberstring10": `
		{
			"i" : 1,
			"n" : 0.5,
			"s" : "hello",
			"u" : "there",
			"w" : 1,
			"f" : 1
		}
	`,
	"cherry": `
	{
		"apple" : 3,
		"pear" : [],
		"plum" : [ 1 ],
		"raspberry" : [ "a", "b" ],
		"strawberry" : "here",
		"guava": [ "c", "d" ],
		"cherry": "127.0.0.1:1234"
	}
	`,
	"badcherry1": `
	{
		"apple" : 3,
		"pear" : [],
		"plum" : [ 1 ],
		"raspberry" : [ "a", "b" ],
		"strawberry" : "here",
		"guava": [ "c", "d" ],
		"cherry": 1234
	}
	`,
	"badcherry2": `
	{
		"apple" : 3,
		"pear" : [],
		"plum" : [ 1 ],
		"raspberry" : [ "a", "b" ],
		"strawberry" : "here",
		"guava": [ "c", "d" ],
		"cherry": "thisisnotahostportpair"
	}
	`,
	"badcherry3": `
	{
		"apple" : 3,
		"pear" : [],
		"plum" : [ 1 ],
		"raspberry" : [ "a", "b" ],
		"strawberry" : "here",
		"guava": [ "c", "d" ],
		"cherry": "127.0.0.1"
	}
	`,
	"tangerine": `
	{
		"apple" : 3,
		"pear" : [],
		"plum" : [ 1 ],
		"raspberry" : [ "a", "b" ],
		"strawberry" : "here",
		"guava": [ "c", "d" ],
		"tangerine": "pips"
	}
	`,
	"badtangerine1": `
	{
		"apple" : 3,
		"pear" : [],
		"plum" : [ 1 ],
		"raspberry" : [ "a", "b" ],
		"strawberry" : "here",
		"guava": [ "c", "d" ],
		"tangerine": "cerebralcortex"
	}
	`,
	"badtangerine2": `
	{
		"apple" : 3,
		"pear" : [],
		"plum" : [ 1 ],
		"raspberry" : [ "a", "b" ],
		"strawberry" : "here",
		"guava": [ "c", "d" ],
		"tangerine": 7
	}
	`,
}

func isOneOrTwo(o interface{}) *cdl.CdlError {
	if v, ok := o.(float64); !ok {
		return cdl.NewError("ErrBadValue").SetSupplementary("is not a float64")
	} else {
		if v != 1 && v != 2 {
			return cdl.NewError("ErrBadValue").SetSupplementary("is not 1 or 2")
		}
	}
	return nil
}

func dummy() {
}

func checkCompile(s string, e string) *cdl.CompiledTemplate {
	if t, ok := checkTemplates[s]; !ok {
		log.Fatalf("Cannot find template %s", s)
		return nil
	} else {
		if ct, err := cdl.Compile(t); err != nil {
			if me, ok := err.(*cdl.CdlError); !ok {
				log.Fatalf("Test checkCompile %s Bad error return %T", s, err)
			} else {
				if me.Type.String() != e {
					log.Fatalf("Test checkCompile %s Returned unexpected error - expecting '%s' got %v; %s", s, e, me.Type.String(), me.Error())
				}
			}
			return nil
		} else {
			if e != "" {
				log.Fatalf("Test checkCompile %s was meant to error with '%s' but didn't", s, e)
			}
			return ct
		}
	}
}

func checkValidate(ct *cdl.CompiledTemplate, s string, e string, c cdl.Configurator) {
	var m interface{}
	if j, ok := checkJsons[s]; !ok {
		log.Fatalf("Test checkValidate Cannot find template %s", s)
	} else {
		if err := json.Unmarshal([]byte(j), &m); err != nil {
			log.Fatalf("Test checkValidate %s JSON parse error: %v ", s, err)
		}

		if err := ct.Validate(m, c); err != nil {
			if me, ok := err.(*cdl.CdlError); !ok {
				log.Fatalf("Test checkValidate %s Bad error return %T", s, err)
			} else {
				if me.Type.String() != e {
					log.Fatalf("Test checkValidate %s Returned unexpected error - expecting '%s' got %v; %s", s, e, me.Type.String(), me.Error())
				}
			}
		} else {
			if e != "" {
				log.Fatalf("Test checkValidate %s was meant to error with '%s' but didn't", s, e)
			}
		}
	}
}

func TestCompile(t *testing.T) {
	checkCompile("simple", "")
	checkCompile("noroot", "ErrMissingRoot")
	checkCompile("badkey", "ErrBadKey")
	checkCompile("array1", "")
	checkCompile("array2", "")
	checkCompile("badarray1", "ErrBadRangeOptionModifier")
	checkCompile("badarray2", "ErrBadRangeOptionModifier")
	checkCompile("badarray3", "ErrBadRangeOptionModifier")
	checkCompile("badarray4", "ErrBadRangeOptionModifier")
	checkCompile("badarray5", "ErrBadRangeOptionModifier")
	checkCompile("badvalue", "ErrBadValue")
	checkCompile("validator", "")
	checkCompile("badvalidator1", "ErrBadValue")
	checkCompile("map", "")
	checkCompile("badmap1", "ErrBadOptionValue")
	checkCompile("badmap2", "ErrBadOptionModifier")
	checkCompile("badmap3", "ErrBadOptionModifier")
	checkCompile("badmap4", "ErrBadOptionModifier")
	checkCompile("badmap5", "ErrBadOptionModifier")
	checkCompile("badmap6", "ErrBadOptionModifier")
	checkCompile("badmap7", "ErrBadOptionModifier")
	checkCompile("badmap8", "ErrBadRangeOptionModifierValue")
	checkCompile("integernumberstring", "")
}

func TestValidate(t *testing.T) {
	ct1 := checkCompile("example", "")

	checkValidate(ct1, "simple1", "", nil)
	checkValidate(ct1, "simple2", "", nil)
	checkValidate(ct1, "bad1", "ErrBadType", nil)
	checkValidate(ct1, "bad2", "ErrBadType", nil)
	checkValidate(ct1, "bad3", "ErrBadValue", nil)

	checkValidate(ct1, "mango", "", nil)
	checkValidate(ct1, "badmango1", "ErrOutOfRange", nil)
	checkValidate(ct1, "badmango2", "ErrOutOfRange", nil)
	checkValidate(ct1, "badmango3", "ErrExpectedMap", nil)
	checkValidate(ct1, "badmango4", "ErrBadKey", nil)

	checkValidate(ct1, "jupiter", "", nil)
	checkValidate(ct1, "badjupiter1", "ErrExpectedArray", nil)
	checkValidate(ct1, "badjupiter2", "ErrBadKey", nil)
	checkValidate(ct1, "badjupiter3", "ErrExpectedMap", nil)
	checkValidate(ct1, "badjupiter4", "ErrExpectedMap", nil)

	checkValidate(ct1, "blueberry", "", nil)
	checkValidate(ct1, "badblueberry1", "ErrExpectedMap", nil)
	checkValidate(ct1, "badblueberry2", "ErrExpectedMap", nil)
	checkValidate(ct1, "badblueberry3", "ErrBadKey", nil)
	checkValidate(ct1, "badblueberry4", "ErrMissingMandatory", nil)

	checkValidate(ct1, "cherry", "", nil)
	checkValidate(ct1, "badcherry1", "ErrBadType", nil)
	checkValidate(ct1, "badcherry2", "ErrBadType", nil)
	checkValidate(ct1, "badcherry3", "ErrBadType", nil)

	checkValidate(ct1, "tangerine", "", nil)
	checkValidate(ct1, "badtangerine1", "ErrBadEnumValue", nil)
	checkValidate(ct1, "badtangerine2", "ErrBadType", nil)

	ct2 := checkCompile("integernumberstring", "")

	var n1 float64
	var i1 int
	var s1 string
	configurator := cdl.Configurator{
		"n": func(o interface{}, p cdl.Path) *cdl.CdlError {
			n1 = o.(float64)
			return nil
		},
		"i": func(o interface{}, p cdl.Path) *cdl.CdlError {
			i1 = o.(int)
			return nil
		},
		"s": func(o interface{}, p cdl.Path) *cdl.CdlError {
			s1 = o.(string)
			return nil
		},
	}
	checkValidate(ct2, "integernumberstring", "", configurator)
	if (n1 != 0.5) || (i1 != 1) || (s1 != "hello") {
		log.Fatalf("Configurator failed: results %d, %f, '%s'", i1, n1, s1)
	}
	checkValidate(ct2, "badintegernumberstring1", "ErrBadType", configurator)
	checkValidate(ct2, "badintegernumberstring2", "ErrBadType", configurator)
	checkValidate(ct2, "badintegernumberstring3", "ErrBadType", configurator)
	checkValidate(ct2, "badintegernumberstring4", "ErrBadType", configurator)
	// tests 5 onwards will not work as they look at bad values of untyped items for
	// which the configurator is not set up in this test

	var n2 float64
	var i2 int
	var s2 string
	var u2 string
	var w2 float64
	var e2 = fruitPart.New("flesh")
	var f2 = fruitPart.New("flesh")
	configurator = cdl.Configurator{
		"n": &n2,
		"i": &i2,
		"s": &s2,
		"u": &u2,
		"w": &w2,
		"e": &e2,
		"f": &f2,
	}
	checkValidate(ct2, "integernumberstring", "", configurator)
	if (n2 != 0.5) || (i2 != 1) || (s2 != "hello") || (u2 != "there") || (w2 != 1) || (e2.String() != "rind") || (f2.String() != "rind") {
		log.Fatalf("Configurator failed: results %d, %f, '%s', '%s', %f, '%s', '%s'", i2, n2, s2, u2, w2, e2, f2)
	}
	checkValidate(ct2, "badintegernumberstring1", "ErrBadType", configurator)
	checkValidate(ct2, "badintegernumberstring2", "ErrBadType", configurator)
	checkValidate(ct2, "badintegernumberstring3", "ErrBadType", configurator)
	checkValidate(ct2, "badintegernumberstring4", "ErrBadType", configurator)
	checkValidate(ct2, "badintegernumberstring5", "ErrBadType", configurator)
	checkValidate(ct2, "badintegernumberstring6", "ErrBadType", configurator)
	checkValidate(ct2, "badintegernumberstring7", "ErrBadEnumValue", configurator)
	checkValidate(ct2, "badintegernumberstring8", "ErrBadEnumValue", configurator)
	checkValidate(ct2, "badintegernumberstring9", "ErrBadType", configurator)
	checkValidate(ct2, "badintegernumberstring10", "ErrBadType", configurator)
}

func Example_cdlCompile() {

	// here's our template
	template := cdl.Template{
		"/":     "{}apple peach? pear* plum+ raspberry{1,3} strawberry! kiwi{1,4}? guava!{1,2} orange?{2,31}",
		"apple": "float64",
	}

	if ct, err := cdl.Compile(template); err != nil {
		log.Fatalf("Error on compile: %v", err)
	} else {

		// use ct here
		_ = ct
	}

	fmt.Println("Success!")
	// Output: Success!
}

func Example_cdlValidate() {

	// here's our template
	template := cdl.Template{
		"/":     "{}apple peach? pear* plum+ raspberry{1,3} strawberry! kiwi{1,4}? guava!{1,2} orange?{2,31}",
		"apple": "float64",
		"peach": func(o interface{}) *cdl.CdlError {
			if v, ok := o.(float64); !ok {
				return cdl.NewError("ErrBadValue").SetSupplementary("is not a float64")
			} else {
				if v != 1 && v != 2 {
					return cdl.NewError("ErrBadValue").SetSupplementary("is not 1 or 2")
				}
			}
			return nil
		},
	}

	if ct, err := cdl.Compile(template); err != nil {
		log.Fatalf("Error on compile: %v", err)
	} else {

		var strawberry string

		// here's our configurator
		configurator := cdl.Configurator{

			// First an easy example using a pointer
			"strawberry": &strawberry,

			// Now a more complex example using a string
			"apple": func(o interface{}, p cdl.Path) *cdl.CdlError {
				fmt.Printf("Apple is %1.0f - ", o.(float64))
				return nil
			},
		}

		// Unmarshal some JSON
		var m interface{}

		j := `
		     {
				"apple" : 3,
				"pear" : [],
				"peach" : 2,
				"plum" : [ 1 ],
				"raspberry" : [ "a", "b" ],
				"strawberry" : "here",
				"guava": [ "c", "d" ]
		     }`

		if err := json.Unmarshal([]byte(j), &m); err != nil {
			log.Fatalf("Cannot unmarshal JSON: %v", err)
		}

		// Validate it
		if err := ct.Validate(m, configurator); err != nil {
			log.Fatalf("Validation error: %v", err)
		}

		if strawberry != "here" {
			log.Fatal("Strawberry variable not set correctly")
		}

		fmt.Println("Success!")
		// Output: Apple is 3 - Success!
	}
}
