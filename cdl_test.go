package cdl_test

import (
	"encoding/json"
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
		"/":         "{}apple peach? pear* plum+ raspberry{1,3} strawberry! kiwi{1,4}? guava!{1,2} orange?{2,} mango? blueberry?",
		"apple":     "float64",
		"peach":     "number",
		"pear":      "string",
		"plum":      isOneOrTwo,
		"mango":     "[]planet{2,4}",
		"planet":    "{}earth venus? jupiter?",
		"jupiter":   "[]gods",
		"gods":      "{}thor? odin?",
		"blueberry": "{}red yellow?",
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
	if ct, err := cdl.Compile(checkTemplates[s]); err != nil {
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

func checkValidate(ct *cdl.CompiledTemplate, s string, e int) {
	var m interface{}

	if err := json.Unmarshal([]byte(checkJsons[s]), &m); err != nil {
		log.Fatalf("Test checkJson %s JSON parse error: %v ", s, err)
	}

	if err := ct.Validate(m); err != nil {
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
}

func TestValidate(t *testing.T) {
	ct1 := checkCompile("example", 0)

	checkValidate(ct1, "simple1", 0)
	checkValidate(ct1, "simple2", 0)
	checkValidate(ct1, "bad1", cdl.ErrBadType)
	checkValidate(ct1, "bad2", cdl.ErrBadType)
	checkValidate(ct1, "bad3", cdl.ErrBadValue)

	checkValidate(ct1, "mango", 0)
	checkValidate(ct1, "badmango1", cdl.ErrOutOfRange)
	checkValidate(ct1, "badmango2", cdl.ErrOutOfRange)
	checkValidate(ct1, "badmango3", cdl.ErrExpectedMap)
	checkValidate(ct1, "badmango4", cdl.ErrBadKey)

	checkValidate(ct1, "jupiter", 0)
	checkValidate(ct1, "badjupiter1", cdl.ErrExpectedArray)
	checkValidate(ct1, "badjupiter2", cdl.ErrBadKey)
	checkValidate(ct1, "badjupiter3", cdl.ErrExpectedMap)
	checkValidate(ct1, "badjupiter4", cdl.ErrExpectedMap)

	checkValidate(ct1, "blueberry", 0)
	checkValidate(ct1, "badblueberry1", cdl.ErrExpectedMap)
	checkValidate(ct1, "badblueberry2", cdl.ErrExpectedMap)
	checkValidate(ct1, "badblueberry3", cdl.ErrBadKey)
	checkValidate(ct1, "badblueberry4", cdl.ErrMissingMandatory)
}
