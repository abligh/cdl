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
//     err := ct.Validate(object, nil)
//
// If the validation fails, you will get an `error` return with a context
// that will allow a user to discover the error in his file.
//
// So what was that `nil` parameter to `cdt.Validate` about? cdl also
// permits you to pass a configurator in, so that you can store the values
// retrieved in appropriate places.
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
//  * There must be an `apple` element and `lemon` element, but the `peach` element
//    is optional.
//
//  * The `apple` element must be a `float64`
//
//  * In order to validate the `peach` element, your own validator function (`isOneOrTwo`)
//     is called. If this returns a `cdl.CdlError`, that error will be passed to the user
//     (as an `error`). If it returns `nil`, then validation will continue.
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
//  * `strawberry`: The `!` indicates it is mandatory; this is the default, so the `!`
//     is unnecessary. Each `strawberry` must be an array of `nectarine` with between
//      1 and 3 components, and each `nectarine` must be a `string`.
//
//  * `rasbperry`: This is a shorthand for writing the same thing as above, i.e. an
//     array of between 1 and 3 `raspberry`, each of which must be a string.
//
//  * `pear`: An array of zero or more items. Note the empty array must be there (if
//     the array itself is optional, write `pear?*`).
//
//  * `plum`: An array of one or more entries.
//
//  * `kiwi`: An optionally present array of between 1 and 4 entries
//
//  * `guava`: A mandatory array of between 1 and 2 entries
//
// Template syntax in detail
//
// 1. Each key must either be `/` (for the root key) or consist of word characters
// (i.e. matching `\w+` in regexp terms)
//
// 2. Each key must have a value, which may be either a validator function, or
// a validation instruction in the form of a `string`
//
// 3. A validator function is a function with the signature
//    func(obj interface{}) (err *CdlError)`
//
// 4. Each validation instruction may be either
//   * The Go name of a type (not a slice), e.g. `bool`, `string` etc. (in quotes as
//     it's a `string`)
//   * A pseudotype (e.g. `number`, `integer`) - see below
//   * An array specifier, having a form beginning `[]`
//   * A map specifier, having a form beginning `{}`
//
// 5. Each pseudotype may be either
//   * The word `number` which indicates any numerical type (not `bool`)
//   * The word `integer` which indicates any numerical type where the value is an
//     integer (useful for parsing JSON with `json/encoding` which presents these as
//     `float64`)
//   * The word `ipport` for an IP port pair which is successfully decoded by
//     `net.SplitHostPort`
//
// 6. An array specifier has the form `[]key` optionally followed by a range specifier
//   * The key (`key` above) consists of word characters.
//   * The key need not be specified within the template (if it isn't, no validation
//     will be done on it).
//
// 7. A range specifier takes the form
//   * `{n,m}` (meaning between `n` and `m`) or
//   * `{n,}` (meaning at least `n`).
//
// 8. A map specifier has the form `{}` followed by zero or more space-separated
//    map elements
//
// 9. A map element consists of a key (`key`) followed by zero or more modifiers
//   * The key consists of word characters.
//   * The key need not be specified within the template (if it isn't, no validation
//     will be done on it).
//
// 10. Permitted modifiers are:
//   * `?` means the key is optional
//   * `!` means the key is mandatory (the default)
//   * `*` means the key is an array of 0 or more elements
//   * `+` means the key is an array of 1 or more elements
//   * A range specifier (see above), i.e.
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
// Configurators
//
// A cdl configurator may optionally be passed to the `Validate` function. The
// configurator allows you to consume the configuration in your program now
// you know that is validated.
//
// The configurator consists of a map of keys to items. Each item should
// be either
//   * a pointer to the variable to be set; or
//   * a pointer to a configuration function.
//
// If a pointer to a variable is used, the variable must be of the same
// type as the item in the configuration, or an error will be issued;
// therefore as a type check is performed here, it is unnecessary in this
// case to require a specific type in the template. If a specific
// type is require, a type check is done twice. Certain pseudo-types
// being required will cause a type conversion:
//
// 1. If you required the pseudo-type `number`, you will be always be given a `float64`
//
// 2. If you required the pseudo-type `integer`, you will always be given an `int`
//
// If a pointer configuration function is used, it has a `ConfiguratorFunc` type
// (or a function with a similar signature), which looks like this:
//
//     type ConfiguratorFunc func(obj interface{}, path Path) (err *CdlError)
//
// This function is guaranteed to be called for each item in the tree
// (if it's key is present in the configurator) after it and all of
// its children have been validated. It may return an error (just like
// a validator function).
//
// The object passed will be the validated object from the configuration
// tree. It is guaranteed to be of the correct type, which means the type
// you asked for save for the following exceptions:
//
// 1. If you asked for the pseudo-type `number`, you will always be given a `float64`.
//
// 2. If you asked for the pseudo-type `integer`, you will always be given an `int`.
//
// As a trivial example:
//
//     var i int
//     err := ct.Validate(object, cdl.Configurator{
//         "i": func(o interface{}, p cdl.Path) *cdl.CdlError {
// 		        i = o.(int)
// 		        return nil
// 	        },
//     })
//
// Here the parameter named `"i"` in the template will be stored in
// variable `i`.
package cdl
