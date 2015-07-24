// Package cdl provides a configuration definition language for Go
// 
// There are several ways to import configuration files into Go, from
// unmarshalling your own JSON using `encoding/json` to projects such as
// Viper (http://spf13.com/project/viper) which will read the configuration
// in JSON, YAML, TOML or from etcd, the command line, or possibly other
// sources. All of these have in common that in the end they produce
// something that looks in Go like
// 
//     map[string]interface{}
// 
// However, these have the issue that they don't validate the configuration.
// If the supplied configuration omits mandatory keys, or puts in extra
// ones, or does both by misspelling one key, or puts the right bit of
// configuration at the wrong level, or any other error is made that
// doesn't actually prevent the JSON, YAML etc. parsing, then it's left
// for you to detect that manually in your program.
// 
// cdl makes this all much easier. Simply supply a cdl template, compile
// it (once) using
//     ct := cdl.Compile(...)
// then validate using
//     err := ct.Validate(object)
// 
// If the validation fails, you will get an `error` return with a context
// that will allow a user to discover the error in his file.
// 
// Templates
// 
// cdl templates are themselves a
//     map[string]interface{}
// 
// but are flat (i.e. only one level deep). The key represents a point in
// in the hierarchy to parse, and the value specifies what may appear at
// that point. The value is normally a `string`, but may be a pointer
// to a validation function you supply.
// 
// For example:
//     template := cdl.Template{
// 		"/":     "{}apple peach? lemon",
// 		"apple": "float64",
// 		"peach": isOneOrTwo,
// 	}
// 
// Here:
//  * The root level is specified to be a map ('`{}`'), which may consist of the elements
//    `apple`, `peach` and `lemon`.
//  
//  * There must be an `apple` element and `lemon` element, but the `peach` element is optional
//  
//  * The `apple` element must be a `float64`
//  
//  * In order to validate the `peach` element, your own validator function (`isOneOrTwo`)
//     is called. If this returns a `cdl.CdlError`, that error will be passed to the user (as
//     an `error`). If it returns `nil`, then validation will continue.
//   
//  * There is no validation at all on `peach`
// 
// Let's take a more complicated example:
// 	template := cdl.Template{
// 		"/":          "{}apple peach? pear* plum+ raspberry{1,3} strawberry! kiwi{1,4}? guava!{1,2} orange?{2,31}",
// 		"apple":      "float64",
// 		"peach":      isOneOrTwo,
// 		"strawberry": "[nectarine]{1,3}",
// 		"nectarine":  "string",
// 		"raspberry":  "string"
// 	}
// 
// Here we have allowed in the root level:
//  * `strawberry`: The `!` indicates it is mandatory; this is the default, so the `!` is unnecessary.
//     Each `strawberry` must be an array of `nectarine` with between 1 and 3 components, and each
//     `nectarine` must be a `string`.
//  
//  * `rasbperry`: This is a shorthand for writing the same thing as above, i.e. an array of between
//     1 and 3 `raspberry`, each of which must be a string.
//  
//  * `pear`: An array of zero or more items. Note the empty array must be there (if the array itself is
//     optional, write `pear?*`).
//  
//  * `plum`: An array of one or more entries.
// 
//  * `kiwi`: An optionally present array of between 1 and 4 entries
//  
//  * `guava`: A mandatory array of between 1 and 2 entries
// 
// Template syntax in detail
// 
// 1. Each key must either be `/` (for the root key) or consist of word characters (i.e. matching `\w+` in regexp terms)
// 
// 2. Each key must have a value, which may be either a validator function, or
// a validation instruction in the form of a `string`
// 
// 3. A validator function is a function with the signature
//    func(obj interface{}) (err *CdlError)`
// 
// 4. Each validation instruction may be either
//   * The Go name of a type (not a slice), e.g. `bool`, `string` etc. (in quotes as it's a `string`)
//   * An array specifier, having a form beginning `[]`
//   * A map specifier, having a form beginning `{}`
// 
// 5. An array specifier has the form `[]key` optionally followed by a range specifier
//   * The key (`key` above) consists of word characters.
//   * The key need not be specified within the template (if it isn't, no validation will be done on it).
// 
// 6. A range specifier takes the form
//   * `{n,m}` (meaning between `n` and `m`) or
//   * `{n,}` (meaning at least `n`).
// 
// 7. A map specifier has the form `{}` followed by zero or more space-separated map elements
// 
// 8. A map element consists of a key (`key`) followed by zero or more modifiers
//   * The key consists of word characters.
//   * The key need not be specified within the template (if it isn't, no validation will be done on it).
// 
// 9. Permitted modifiers are:
//   * `?` means the key is optional
//   * `!` means the key is mandatory (the default)
//   * `*` means the key is an array of 0 or more elements
//   * `+` means the key is an array of 1 or more elements
//   * A *range specifier* (see above), i.e.
//     * `{n,m}` (meaning between `n` and `m`) or
//     * `{n,}` (meaning at least `n`)
// 
// Validator Functions
// 
// Where the validator is passed, it is a function with signature:
//     func (o interface{}) *cdl.CdlError
// 
// Here's an example showing how it can return an error and send supplementary data back to the user.
// Note that cdl itself will add the appropriate context.
// 
//     func isOneOrTwo(o interface{}) *cdl.CdlError {
//     	if v, ok := o.(float64); !ok {
//     		return cdl.NewError(cdl.ErrBadValue).SetSupplementary("is not a float64")
//     	} else {
//     		if v != 1 && v != 2 {
//     			return cdl.NewError(cdl.ErrBadValue).SetSupplementary("is not 1 or 2")
//     		}
//     	}
//     	return nil
//     }
//
package cdl
