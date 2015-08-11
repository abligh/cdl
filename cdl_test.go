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
		"/":         "{}apple peach? pear* plum+ raspberry{1,3} strawberry! kiwi{1,4}? guava!{1,2} orange?{2,} mango? blueberry? cherry?",
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
	},
	"integernumberstring": cdl.Template{
		"/": "{}i? n? s? u? w?",
		"n": "number",
		"i": "integer",
		"s": "string",
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
			"w" : 1
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
}

func isOneOrTwo(o interface{}) *cdl.CdlError {
	if v, ok := o.(float64); !ok {
		return cdl.NewError(cdl.ErrBadValue).SetSupplementary("is not a float64")
	} else {
		if v != 1 && v != 2 {
			return cdl.NewError(cdl.ErrBadValue).SetSupplementary("is not 1 or 2")
		}
	}
	return nil
}

func dummy() {
}

func checkCompile(s string, e int) *cdl.CompiledTemplate {
	if t, ok := checkTemplates[s]; !ok {
		log.Fatalf("Cannot find template %s", s)
		return nil
	} else {
		if ct, err := cdl.Compile(t); err != nil {
			if me, ok := err.(*cdl.CdlError); !ok {
				log.Fatalf("Test checkCompile %s Bad error return %T", s, err)
			} else {
				if me.Type != e {
					log.Fatalf("Test checkCompile %s Returned unexpected error - expecting %d got %v", s, e, me)
				}
			}
			return nil
		} else {
			if e != 0 {
				log.Fatalf("Test checkCompile %s was meant to error with %d but didn't", s, e)
			}
			return ct
		}
	}
}

func checkValidate(ct *cdl.CompiledTemplate, s string, e int, c cdl.Configurator) {
	var m interface{}
	if j, ok := checkJsons[s]; !ok {
		log.Fatalf("Cannot find template %s", s)
	} else {
		if err := json.Unmarshal([]byte(j), &m); err != nil {
			log.Fatalf("Test checkJson %s JSON parse error: %v ", s, err)
		}

		if err := ct.Validate(m, c); err != nil {
			if me, ok := err.(*cdl.CdlError); !ok {
				log.Fatalf("Test checkJson %s Bad error return %T", s, err)
			} else {
				if me.Type != e {
					log.Fatalf("Test checkJson %s Returned unexpected error - expecting %d got %v", s, e, me)
				}
			}
		} else {
			if e != 0 {
				log.Fatalf("Test checkJson %s was meant to error with %d but didn't", s, e)
			}
		}
	}
}

func TestCompile(t *testing.T) {
	checkCompile("simple", 0)
	checkCompile("noroot", cdl.ErrMissingRoot)
	checkCompile("badkey", cdl.ErrBadKey)
	checkCompile("array1", 0)
	checkCompile("array2", 0)
	checkCompile("badarray1", cdl.ErrBadRangeOptionModifier)
	checkCompile("badarray2", cdl.ErrBadRangeOptionModifier)
	checkCompile("badarray3", cdl.ErrBadRangeOptionModifier)
	checkCompile("badarray4", cdl.ErrBadRangeOptionModifier)
	checkCompile("badarray5", cdl.ErrBadRangeOptionModifier)
	checkCompile("badvalue", cdl.ErrBadValue)
	checkCompile("validator", 0)
	checkCompile("badvalidator1", cdl.ErrBadValue)
	checkCompile("map", 0)
	checkCompile("badmap1", cdl.ErrBadOptionValue)
	checkCompile("badmap2", cdl.ErrBadOptionModifier)
	checkCompile("badmap3", cdl.ErrBadOptionModifier)
	checkCompile("badmap4", cdl.ErrBadOptionModifier)
	checkCompile("badmap5", cdl.ErrBadOptionModifier)
	checkCompile("badmap6", cdl.ErrBadOptionModifier)
	checkCompile("badmap7", cdl.ErrBadOptionModifier)
	checkCompile("badmap8", cdl.ErrBadRangeOptionModifierValue)
	checkCompile("integernumberstring", 0)
}

func TestValidate(t *testing.T) {
	ct1 := checkCompile("example", 0)

	checkValidate(ct1, "simple1", 0, nil)
	checkValidate(ct1, "simple2", 0, nil)
	checkValidate(ct1, "bad1", cdl.ErrBadType, nil)
	checkValidate(ct1, "bad2", cdl.ErrBadType, nil)
	checkValidate(ct1, "bad3", cdl.ErrBadValue, nil)

	checkValidate(ct1, "mango", 0, nil)
	checkValidate(ct1, "badmango1", cdl.ErrOutOfRange, nil)
	checkValidate(ct1, "badmango2", cdl.ErrOutOfRange, nil)
	checkValidate(ct1, "badmango3", cdl.ErrExpectedMap, nil)
	checkValidate(ct1, "badmango4", cdl.ErrBadKey, nil)

	checkValidate(ct1, "jupiter", 0, nil)
	checkValidate(ct1, "badjupiter1", cdl.ErrExpectedArray, nil)
	checkValidate(ct1, "badjupiter2", cdl.ErrBadKey, nil)
	checkValidate(ct1, "badjupiter3", cdl.ErrExpectedMap, nil)
	checkValidate(ct1, "badjupiter4", cdl.ErrExpectedMap, nil)

	checkValidate(ct1, "blueberry", 0, nil)
	checkValidate(ct1, "badblueberry1", cdl.ErrExpectedMap, nil)
	checkValidate(ct1, "badblueberry2", cdl.ErrExpectedMap, nil)
	checkValidate(ct1, "badblueberry3", cdl.ErrBadKey, nil)
	checkValidate(ct1, "badblueberry4", cdl.ErrMissingMandatory, nil)

	checkValidate(ct1, "cherry", 0, nil)
	checkValidate(ct1, "badcherry1", cdl.ErrBadType, nil)
	checkValidate(ct1, "badcherry2", cdl.ErrBadType, nil)
	checkValidate(ct1, "badcherry3", cdl.ErrBadType, nil)
	ct2 := checkCompile("integernumberstring", 0)

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
	checkValidate(ct2, "integernumberstring", 0, configurator)
	if (n1 != 0.5) || (i1 != 1) || (s1 != "hello") {
		log.Fatalf("Configurator failed: results %d, %f, '%s'", i1, n1, s1)
	}
	checkValidate(ct2, "badintegernumberstring1", cdl.ErrBadType, configurator)
	checkValidate(ct2, "badintegernumberstring2", cdl.ErrBadType, configurator)
	checkValidate(ct2, "badintegernumberstring3", cdl.ErrBadType, configurator)
	checkValidate(ct2, "badintegernumberstring4", cdl.ErrBadType, configurator)
	// tests 5 & 6 will not work as they look at bad values of untyped items for
	// which the configurator is not set up in this test

	var n2 float64
	var i2 int
	var s2 string
	var u2 string
	var w2 float64
	configurator = cdl.Configurator{
		"n": &n2,
		"i": &i2,
		"s": &s2,
		"u": &u2,
		"w": &w2,
	}
	checkValidate(ct2, "integernumberstring", 0, configurator)
	if (n2 != 0.5) || (i2 != 1) || (s2 != "hello") || (u2 != "there") || (w2 != 1) {
		log.Fatalf("Configurator failed: results %d, %f, '%s', '%s', %f", i2, n2, s2, u2, w2)
	}
	checkValidate(ct2, "badintegernumberstring1", cdl.ErrBadType, configurator)
	checkValidate(ct2, "badintegernumberstring2", cdl.ErrBadType, configurator)
	checkValidate(ct2, "badintegernumberstring3", cdl.ErrBadType, configurator)
	checkValidate(ct2, "badintegernumberstring4", cdl.ErrBadType, configurator)
	checkValidate(ct2, "badintegernumberstring5", cdl.ErrBadType, configurator)
	checkValidate(ct2, "badintegernumberstring6", cdl.ErrBadType, configurator)

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
				return cdl.NewError(cdl.ErrBadValue).SetSupplementary("is not a float64")
			} else {
				if v != 1 && v != 2 {
					return cdl.NewError(cdl.ErrBadValue).SetSupplementary("is not 1 or 2")
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
