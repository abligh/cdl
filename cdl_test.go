package cdl

import (
	"encoding/json"
	"log"
	"testing"
)

type checkTemplate map[string]Template
type checkJson map[string]string

var checkTemplates checkTemplate = checkTemplate{
	"simple": Template{
		"/":   "foo",
		"bar": "int",
	},
	"noroot": Template{
		"x": "foo",
	},
	"badkey": Template{
		"/": "foo",
		"+": "foo",
	},
	"array1": Template{
		"/": "[]foo",
	},
	"array2": Template{
		"/": "[]foo{1,3}",
	},
	"badarray1": Template{
		"/": "[]",
	},
	"badarray2": Template{
		"/": "[]!",
	},
	"badarray3": Template{
		"/": "[]foo{3}",
	},
	"badarray4": Template{
		"/": "[]foo{3,a}",
	},
	"badarray5": Template{
		"/": "[]foo {1,3}",
	},
	"badvalue": Template{
		"/": 1,
	},
	"validator": Template{
		"/": isOneOrTwo,
	},
	"badvalidator1": Template{
		"/": dummy,
	},
	"map": Template{
		"/":     "{}apple peach? pear* plum+ raspberry{1,3} strawberry! kiwi{1,4}? guava!{1,2} orange?{2,31}",
		"apple": "int",
		"peach": isOneOrTwo,
	},
	"badmap1": Template{
		"/": "{}/a",
	},
	"badmap2": Template{
		"/": "{}a/a",
	},
	"badmap3": Template{
		"/": "{}apple/",
	},
	"badmap4": Template{
		"/": "{}apple{1",
	},
	"badmap5": Template{
		"/": "{}apple{-1,-1}",
	},
	"badmap6": Template{
		"/": "{}apple{a,1}",
	},
	"badmap7": Template{
		"/": "{}apple{1,a}",
	},
	"badmap8": Template{
		"/": "{}apple{3,1}",
	},
	"example": Template{
		"/":         "{}apple peach? pear* plum+ raspberry{1,3} strawberry! kiwi{1,4}? guava!{1,2} orange?{2,31} mango? blueberry?",
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

func isOneOrTwo(o interface{}) *CdlError {
	if v, ok := o.(float64); !ok {
		return NewError(ErrBadValue).SetSupplementary("is not a float64")
	} else {
		if v != 1 && v != 2 {
			return NewError(ErrBadValue).SetSupplementary("is not 1 or 2")
		}
	}
	return nil
}

func dummy() {
}

func checkCompile(s string, e int) *CompiledTemplate {
	if ct, err := Compile(checkTemplates[s]); err != nil {
		if me, ok := err.(*CdlError); !ok {
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

func checkValidate(ct *CompiledTemplate, s string, e int) {
	var m interface{}

	if err := json.Unmarshal([]byte(checkJsons[s]), &m); err != nil {
		log.Fatalf("Test checkJson %s JSON parse error: %v ", s, err)
	}

	if err := ct.Validate(m); err != nil {
		if me, ok := err.(*CdlError); !ok {
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
	checkCompile("noroot", ErrMissingRoot)
	checkCompile("badkey", ErrBadKey)
	checkCompile("array1", 0)
	checkCompile("array2", 0)
	checkCompile("badarray1", ErrBadRangeOptionModifier)
	checkCompile("badarray2", ErrBadRangeOptionModifier)
	checkCompile("badarray3", ErrBadRangeOptionModifier)
	checkCompile("badarray4", ErrBadRangeOptionModifier)
	checkCompile("badarray5", ErrBadRangeOptionModifier)
	checkCompile("badvalue", ErrBadValue)
	checkCompile("validator", 0)
	checkCompile("badvalidator1", ErrBadValue)
	checkCompile("map", 0)
	checkCompile("badmap1", ErrBadOptionValue)
	checkCompile("badmap2", ErrBadOptionModifier)
	checkCompile("badmap3", ErrBadOptionModifier)
	checkCompile("badmap4", ErrBadOptionModifier)
	checkCompile("badmap5", ErrBadOptionModifier)
	checkCompile("badmap6", ErrBadOptionModifier)
	checkCompile("badmap7", ErrBadOptionModifier)
	checkCompile("badmap8", ErrBadRangeOptionModifierValue)
}

func TestValidate(t *testing.T) {
	ct1 := checkCompile("example", 0)

	checkValidate(ct1, "simple1", 0)
	checkValidate(ct1, "simple2", 0)
	checkValidate(ct1, "bad1", ErrBadType)
	checkValidate(ct1, "bad2", ErrBadType)
	checkValidate(ct1, "bad3", ErrBadValue)

	checkValidate(ct1, "mango", 0)
	checkValidate(ct1, "badmango1", ErrOutOfRange)
	checkValidate(ct1, "badmango2", ErrOutOfRange)
	checkValidate(ct1, "badmango3", ErrExpectedMap)
	checkValidate(ct1, "badmango4", ErrBadKey)

	checkValidate(ct1, "jupiter", 0)
	checkValidate(ct1, "badjupiter1", ErrExpectedArray)
	checkValidate(ct1, "badjupiter2", ErrBadKey)
	checkValidate(ct1, "badjupiter3", ErrExpectedMap)
	checkValidate(ct1, "badjupiter4", ErrExpectedMap)

	checkValidate(ct1, "blueberry", 0)
	checkValidate(ct1, "badblueberry1", ErrExpectedMap)
	checkValidate(ct1, "badblueberry2", ErrExpectedMap)
	checkValidate(ct1, "badblueberry3", ErrBadKey)
	checkValidate(ct1, "badblueberry4", ErrMissingMandatory)
}
