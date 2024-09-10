package groq

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// ToSnakeCase converts the provided string into snake case using dashes.
// This is useful for Schema IDs and definitions to be coherent with
// common JSON Schema examples.
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}-${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}-${2}")
	return strings.ToLower(snake)
}

// NewProperties is a helper method to instantiate a new properties ordered
// map.
func NewProperties() *orderedmap.OrderedMap[string, *Schema] {
	return orderedmap.New[string, *Schema]()
}

// Version is the JSON Schema version.
var Version = "https://json-schema.org/draft/2020-12/schema"

// Schema represents a JSON Schema object type.
// RFC draft-bhutton-json-schema-00 section 4.3
type Schema struct {
	// RFC draft-bhutton-json-schema-00
	Version     string      `json:"$schema,omitempty"`     // section 8.1.1
	ID          ID          `json:"$id,omitempty"`         // section 8.2.1
	Anchor      string      `json:"$anchor,omitempty"`     // section 8.2.2
	Ref         string      `json:"$ref,omitempty"`        // section 8.2.3.1
	DynamicRef  string      `json:"$dynamicRef,omitempty"` // section 8.2.3.2
	Definitions Definitions `json:"$defs,omitempty"`       // section 8.2.4
	// Comments specifies a comment for the schema as
	// specified RFC draft-bhutton-json-schema-00 section 8.3
	//
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-8.3
	//
	// The value of this field MUST be a string.  Implementations MUST NOT
	// present this string to end users.  Tools for editing schemas SHOULD
	// support displaying and editing this field.  The value of this
	// field MAY be used in debug or error output which is intended for
	// developers making use of schemas.
	//
	// Schema vocabularies SHOULD allow "$comment" within any object
	// containing vocabulary fields.  Implementations MAY assume
	// "$comment" is allowed unless the vocabulary specifically forbids it.
	// Vocabularies MUST NOT specify any effect of "$comment" beyond what is
	// described in this specification.
	//
	// Tools that translate other media types or programming languages to
	// and from application/schema+json MAY choose to convert that media
	// type or programming language's native comments to or from "$comment"
	// values.  The behavior of such translation when both native comments
	// and "$comment" properties are present is implementation-dependent.
	//
	// Implementations MAY strip "$comment" values at any point during
	// processing.  In particular, this allows for shortening schemas when
	// the size of deployed schemas is a concern.
	//
	// Implementations MUST NOT take any other action based on the presence,
	// absence, or contents of "$comment" properties.  In particular, the
	// value of "$comment" MUST NOT be collected as an annotation result.
	Comments string `json:"$comment,omitempty"` // section 8.3
	// AllOf specifies that the schema is an all of of the schema as
	// specifified RFC draft-bhutton-json-schema-00 section 10.2.1
	//
	// section 10.2.1.1
	//
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-10.2.1
	//
	// The value of this field MUST be an array.  Elements in the array
	// MUST be objects.  Each object MUST be a valid JSON Schema.
	//
	// An instance validates successfully against this field if it
	// validates successfully against all schemas defined by "allOf".
	//
	// Omitting this field has the same behavior as an empty array.
	AllOf []*Schema `json:"allOf,omitempty"` // section 10.2.1.1
	AnyOf []*Schema `json:"anyOf,omitempty"` // section 10.2.1.2
	OneOf []*Schema `json:"oneOf,omitempty"` // section 10.2.1.3
	Not   *Schema   `json:"not,omitempty"`   // section 10.2.1.4
	// RFC draft-bhutton-json-schema-00 section 10.2.2 (Apply sub-schemas conditionally)
	If               *Schema            `json:"if,omitempty"`               // section 10.2.2.1
	Then             *Schema            `json:"then,omitempty"`             // section 10.2.2.2
	Else             *Schema            `json:"else,omitempty"`             // section 10.2.2.3
	DependentSchemas map[string]*Schema `json:"dependentSchemas,omitempty"` // section 10.2.2.4
	// RFC draft-bhutton-json-schema-00 section 10.3.1 (arrays)
	PrefixItems []*Schema `json:"prefixItems,omitempty"` // section 10.3.1.1
	Items       *Schema   `json:"items,omitempty"`       // section 10.3.1.2  (replaces additionalItems)
	Contains    *Schema   `json:"contains,omitempty"`    // section 10.3.1.3
	// RFC draft-bhutton-json-schema-00 section 10.3.2 (sub-schemas)
	Properties           *orderedmap.OrderedMap[string, *Schema] `json:"properties,omitempty"`           // section 10.3.2.1
	PatternProperties    map[string]*Schema                      `json:"patternProperties,omitempty"`    // section 10.3.2.2
	AdditionalProperties *Schema                                 `json:"additionalProperties,omitempty"` // section 10.3.2.3
	PropertyNames        *Schema                                 `json:"propertyNames,omitempty"`        // section 10.3.2.4
	// RFC draft-bhutton-json-schema-validation-00, section 6
	Type             string      `json:"type,omitempty"`             // section 6.1.1
	Enum             []any       `json:"enum,omitempty"`             // section 6.1.2
	Const            any         `json:"const,omitempty"`            // section 6.1.3
	MultipleOf       json.Number `json:"multipleOf,omitempty"`       // section 6.2.1
	Maximum          json.Number `json:"maximum,omitempty"`          // section 6.2.2
	ExclusiveMaximum json.Number `json:"exclusiveMaximum,omitempty"` // section 6.2.3
	Minimum          json.Number `json:"minimum,omitempty"`          // section 6.2.4
	ExclusiveMinimum json.Number `json:"exclusiveMinimum,omitempty"` // section 6.2.5
	MaxLength        *uint64     `json:"maxLength,omitempty"`        // section 6.3.1
	// MinLength specifies the minimum length of the string as specified in
	// section 6.3.2 of RFC
	// draft-bhutton-json-schema-validation-00.
	//
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.3.2
	//
	// The value of this field MUST be a non-negative integer.
	//
	// A string instance is valid against "minLength" if its length is
	// greater than, or equal to, the value of this field.
	//
	// Omitting this field has the same behavior as a value of 0.
	MinLength *uint64 `json:"minLength,omitempty"` // section 6.3.2
	// Pattern specifies the regular expression pattern of the schema as
	// specified in section 6.3.3 of RFC
	// draft-bhutton-json-schema-validation-00.
	//
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.3.3
	//
	// The value of this field MUST be a string.  This string SHOULD be a
	// valid regular expression, according to the ECMA-262 regular
	// expression dialect.
	//
	// A string instance is considered valid if the regular expression
	// matches the instance successfully.  Recall: regular expressions are
	// not implicitly anchored.
	Pattern string `json:"pattern,omitempty"` // section 6.3.3
	// MaxItems specifies the maximum number of items in the array as
	// specified in section 6.4.1 of RFC
	// draft-bhutton-json-schema-validation-00.
	//
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.4.1
	//
	// The value of this field MUST be a non-negative integer.
	//
	// An array instance is valid against "maxItems" if its size is less
	// than, or equal to, the value of this field.
	//
	// Omitting this field has the same behavior as a value of an
	// implementation-defined number.
	MaxItems *uint64 `json:"maxItems,omitempty"` // section 6.4.1
	// MinItems specifies the minimum number of items in the array as
	// specified in section 6.4.2 of RFC
	// draft-bhutton-json-schema-validation-00.
	//
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.4.2
	//
	// The value of this field MUST be a non-negative integer.
	//
	// An array instance is valid against "minItems" if its size is greater
	// than, or equal to, the value of this field.
	//
	// Omitting this field has the same behavior as a value of 0.
	MinItems *uint64 `json:"minItems,omitempty"` // section 6.4.2
	// UniqueItems specifies that the instance array is unique as specified
	// in section 6.4.3 of RFC
	// draft-bhutton-json-schema-validation-00.
	//
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.4.3
	//
	// The value of this field MUST be a boolean.
	//
	// If this field has boolean value false, the instance validates
	// successfully.  If it has boolean value true, the instance validates
	// successfully if all of its elements are unique.
	UniqueItems bool `json:"uniqueItems,omitempty"` // section 6.4.3
	// MaxContains specifies the maximum number of items in the array as
	// specified in section 6.4.4 of RFC
	// draft-bhutton-json-schema-validation-00.
	//
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.4.4
	//
	// The value of this field MUST be a non-negative integer.
	//
	// If "contains" is not present within the same schema object, then this
	// field has no effect.
	//
	// An instance array is valid against "maxContains" in two ways,
	// depending on the form of the annotation result of an adjacent
	// "contains" [json-schema] field.  The first way is if the annotation
	// result is an array and the length of that array is less than or equal
	// to the "maxContains" value.  The second way is if the annotation
	// result is a boolean "true" and the instance array length is less than
	// or equal to the "maxContains" value.
	MaxContains *uint64 `json:"maxContains,omitempty"`
	// MinContains specifies the minimum number of items in the array as
	// specified in section 6.4.5 of RFC
	// draft-bhutton-json-schema-validation-00.
	//
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.4.5
	//
	// The value of this field MUST be a non-negative integer.
	//
	// If "contains" is not present within the same schema object, then this
	// field has no effect.
	//
	// An instance array is valid against "minContains" in two ways,
	// depending on the form of the annotation result of an adjacent
	// "contains" [json-schema] field.  The first way is if the annotation
	// result is an array and the length of that array is greater than or
	// equal to the "minContains" value.  The second way is if the
	// annotation result is a boolean "true" and the instance array length
	// is greater than or equal to the "minContains" value.
	//
	// A value of 0 is allowed, but is only useful for setting a range of
	// occurrences from 0 to the value of "maxContains".  A value of 0 with
	// no "maxContains" causes "contains" to always pass validation.
	MinContains *uint64 `json:"minContains,omitempty"`
	// MaxProperties specifies the maximum number of properties of the
	// schema as specified in section 6.5.1 of RFC
	// draft-bhutton-json-schema-validation-00.
	//
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.5.1
	//
	// The value of this field MUST be a non-negative integer.
	//
	// An object instance is valid against "maxProperties" if its number of
	// properties is less than, or equal to, the value of this field.
	//
	// Omitting this field has the same behavior as a value of an
	// implementation-defined number.
	MaxProperties *uint64 `json:"maxProperties,omitempty"`
	// MinProperties specifies the minimum number of properties of the
	// schema as specifiied in section 6.5.2 of RFC
	// draft-bhutton-json-schema-validation-00.
	//
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.5.2
	//
	// The value of this field MUST be a non-negative integer.
	//
	// An object instance is valid against "minProperties" if its number of
	// properties is greater than, or equal to, the value of this field.
	//
	// Omitting this field has the same behavior as a value of 0.
	MinProperties *uint64 `json:"minProperties,omitempty"`
	// Required specifies the required properties of the schema as
	// specified in section 6.5.3 of RFC
	// draft-bhutton-json-schema-validation-00.
	//
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.5.4
	//
	// The value of this field MUST be an object.  Properties in this
	// object, if any, MUST be arrays.  Elements in each array, if any, MUST
	// be strings, and MUST be unique.
	//
	// This field specifies properties that are required if a specific
	// other property is present.  Their requirement is dependent on the
	// presence of the other property.
	//
	// Validation succeeds if, for each name that appears in both the
	// instance and as a name within this field's value, every item in the
	// corresponding array is also the name of a property in the instance.
	//
	// Omitting this field has the same behavior as an empty object.
	Required []string `json:"required,omitempty"`
	// DependentRequired is the dependent required of the schema.
	//
	// section 6.5.4
	//
	// url: https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.5.4
	//
	// The value of this field MUST be an object.  Properties in this
	// object, if any, MUST be arrays.  Elements in each array, if any, MUST
	// be strings, and MUST be unique.
	//
	// This field specifies properties that are required if a specific
	// other property is present.  Their requirement is dependent on the
	// presence of the other property.
	//
	// Validation succeeds if, for each name that appears in both the
	// instance and as a name within this field's value, every item in the
	// corresponding array is also the name of a property in the instance.
	//
	// Omitting this field has the same behavior as an empty object.
	DependentRequired map[string][]string `json:"dependentRequired,omitempty"` // section 6.5.4
	// Format specifies the format of the schema as specified in section
	// 7.3 of RFC draft-bhutton-json-schema-validation-00.
	//
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-7.3
	//
	// TODO: add type for format and all the possible formats
	//
	// The value of this field MUST be a string.  Implementations that
	// use a subset of JSON as their input format, such as JSON Hyper-Schema
	// or JSON Schema Hyper-Schema, MAY implement validation against
	// meta-schemas that define format-specific fields that describe
	// additional constraints beyond those specified herein.
	Format string `json:"format,omitempty"` // RFC draft-bhutton-json-schema-validation-00, section 7
	// RFC draft-bhutton-json-schema-validation-00, section 8
	// ContentEncoding specifies the content encoding of the schema as
	// specified in section 8.3 of RFC
	// draft-bhutton-json-schema-validation-00.
	//
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-8.3
	//
	// If the instance value is a string, this property defines that the
	// string SHOULD be interpreted as binary data and decoded using the
	// encoding named by this property.
	//
	// Possible values indicating base 16, 32, and 64 encodings with several
	// variations are listed in RFC 4648 [RFC4648].  Additionally, sections
	// 6.7 and 6.8 of RFC 2045 [RFC2045] provide encodings used in MIME.  As
	// "base64" is defined in both RFCs, the definition from RFC 4648 SHOULD
	// be assumed unless the string is specifically intended for use in a
	// MIME context.  Note that all of these encodings result in strings
	// consisting only of 7-bit ASCII characters.  Therefore, this field
	// has no meaning for strings containing characters outside of that
	// range.
	//
	// If this field is absent, but "contentMediaType" is present, this
	// indicates that the encoding is the identity encoding, meaning that no
	// transformation was needed in order to represent the content in a
	// UTF-8 string.
	ContentEncoding string `json:"contentEncoding,omitempty"` // section 8.3
	// ContentMediaType specifies the content media type of the schema as
	// specified in section 8.4 of RFC
	// draft-bhutton-json-schema-validation-00.
	//
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-8.4
	//
	// If the instance is a string, this property indicates the media type
	// of the contents of the string.  If "contentEncoding" is present, this
	// property describes the decoded string.
	//
	// The value of this property MUST be a string, which MUST be a media
	// type, as defined by RFC 2046 [RFC2046].
	ContentMediaType string `json:"contentMediaType,omitempty"` // section 8.4
	// ContentSchema specifies the content schema of the schema as
	// specified in section 8.5 of RFC
	// draft-bhutton-json-schema-validation-00.
	//
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-8.5
	//
	//
	// If the instance is a string, and if "contentMediaType" is present,
	// this property contains a schema which describes the structure of the
	// string.
	//
	// This field MAY be used with any media type that can be mapped into
	// JSON Schema's data model.
	//
	// The value of this property MUST be a valid JSON schema.  It SHOULD be
	// ignored if "contentMediaType" is not present.
	ContentSchema *Schema `json:"contentSchema,omitempty"` // section 8.5
	// RFC draft-bhutton-json-schema-validation-00, section 9
	Title       string `json:"title,omitempty"`       // section 9.1
	Description string `json:"description,omitempty"` // section 9.1
	Default     any    `json:"default,omitempty"`     // section 9.2
	Deprecated  bool   `json:"deprecated,omitempty"`  // section 9.3
	ReadOnly    bool   `json:"readOnly,omitempty"`    // section 9.4
	WriteOnly   bool   `json:"writeOnly,omitempty"`   // section 9.4
	Examples    []any  `json:"examples,omitempty"`    // section 9.5

	Extras map[string]any `json:"-"`

	// Special boolean representation of the Schema - section 4.3.2
	boolean *bool
}

var (
	// TrueSchema defines a schema with a true value
	TrueSchema = &Schema{boolean: &[]bool{true}[0]}
	// FalseSchema defines a schema with a false value
	FalseSchema = &Schema{boolean: &[]bool{false}[0]}
)

// Definitions hold schema definitions.
//
// http://json-schema.org/latest/json-schema-validation.html#rfc.section.5.26
//
// RFC draft-wright-json-schema-validation-00, section 5.26
type Definitions map[string]*Schema

// ID represents a Schema ID type which should always be a URI.
// See draft-bhutton-json-schema-00 section 8.2.1
type ID string

// EmptyID is used to explicitly define an ID with no value.
const EmptyID ID = ""

// Validate is used to check if the ID looks like a proper schema.
// This is done by parsing the ID as a URL and checking it has all the
// relevant parts.
func (id ID) Validate() error {
	u, err := url.Parse(id.String())
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if u.Hostname() == "" {
		return errors.New("missing hostname")
	}
	if !strings.Contains(u.Hostname(), ".") {
		return errors.New("hostname does not look valid")
	}
	if u.Path == "" {
		return errors.New("path is expected")
	}
	if u.Scheme != "https" && u.Scheme != "http" {
		return errors.New("unexpected schema")
	}
	return nil
}

// Anchor sets the anchor part of the schema URI.
func (id ID) Anchor(name string) ID {
	b := id.Base()
	return ID(b.String() + "#" + name)
}

// Def adds or replaces a definition identifier.
func (id ID) Def(name string) ID {
	b := id.Base()
	return ID(b.String() + "#/$defs/" + name)
}

// Add appends the provided path to the id, and removes any
// anchor data that might be there.
func (id ID) Add(path string) ID {
	b := id.Base()
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return ID(b.String() + path)
}

// Base removes any anchor information from the schema
func (id ID) Base() ID {
	s := id.String()
	i := strings.LastIndex(s, "#")
	if i != -1 {
		s = s[0:i]
	}
	s = strings.TrimRight(s, "/")
	return ID(s)
}

// String provides string version of ID
func (id ID) String() string {
	return string(id)
}
