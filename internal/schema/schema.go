package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/conneroisu/groq-go/internal/omap"
)

const (
	// version is the JSON Schema version.
	version = "https://json-schema.org/draft/2020-12/schema"
	// EmptyID is used to explicitly define an ID with no value.
	EmptyID schemaID = ""
)

// ReflectSchema returns a schema from a value.
func ReflectSchema(a any) (*Schema, error) {
	r := &reflector{}
	schema := r.ReflectFromType(reflect.TypeOf(a))
	return schema, nil
}

// Available Go defined types for JSON Schema Validation.
//
// https://datatracker.ietf.org/doc/html/draft-wright-json-schema-validation-00#section-7.3
//
// RFC draft-wright-json-schema-validation-00, section 7.3
var (
	trueSchema                    = &Schema{boolean: &[]bool{true}[0]}
	falseSchema                   = &Schema{boolean: &[]bool{false}[0]}
	timeType                      = reflect.TypeOf(time.Time{}) // date-time RFC section 7.3.1
	ipType                        = reflect.TypeOf(net.IP{})    // ipv4 and ipv6 RFC section 7.3.4, 7.3.5
	uriType                       = reflect.TypeOf(url.URL{})   // uri RFC section 7.3.6
	byteSliceType                 = reflect.TypeOf([]byte(nil))
	rawMessageType                = reflect.TypeOf(json.RawMessage{})
	customType                    = reflect.TypeOf((*customSchemaImpl)(nil)).Elem()
	extendType                    = reflect.TypeOf((*extendSchemaImpl)(nil)).Elem()
	customStructGetFieldDocString = reflect.TypeOf((*customSchemaGetFieldDocString)(nil)).Elem()
	protoEnumType                 = reflect.TypeOf((*protoEnum)(nil)).Elem()
	matchFirstCap                 = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap                   = regexp.MustCompile("([a-z0-9])([A-Z])")
	customAliasSchema             = reflect.TypeOf((*aliasSchemaImpl)(nil)).Elem()
	customPropertyAliasSchema     = reflect.TypeOf((*propertyAliasSchemaImpl)(nil)).
					Elem()
)

type (
	// Schema represents a JSON Schema object type.
	// RFC draft-bhutton-json-Schema-00 section 4.3
	Schema struct {
		// RFC draft-bhutton-json-schema-00
		// Version is the version of the schema as specified in section 8.1.1 of RFC
		// draft-bhutton-json-schema-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-8.1.1
		//
		// The value of this field MUST be a string.  This string SHOULD be a
		// version number of the schema.
		Version string `json:"$schema,omitempty"`
		// ID is the ID of the schema as specified in section 8.2.1 of RFC
		// draft-bhutton-json-schema-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-8.2.1
		//
		// The value of this field MUST be a string.  This string SHOULD be a
		// URI.
		ID schemaID `json:"$id,omitempty"`
		// Anchor is the anchor of the schema as specified in section 8.2.2 of RFC
		// draft-bhutton-json-schema-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-8.2.2
		//
		// The value of this field MUST be a string.  This string SHOULD be a
		// valid URI fragment.
		Anchor string `json:"$anchor,omitempty"`
		// Ref is the ref of the schema as specified in section 8.2.3.1 of RFC
		// draft-bhutton-json-schema-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-8.2.3.1
		//
		// The value of this field MUST be a string.  This string SHOULD be a
		// valid URI.
		Ref string `json:"$ref,omitempty"`
		// DynamicRef is the dynamic ref of the schema as specified in section 8.2.3.2 of RFC
		// draft-bhutton-json-schema-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-8.2.3.2
		//
		// The value of this field MUST be a string.  This string SHOULD be a
		// valid URI.
		DynamicRef string `json:"$dynamicRef,omitempty"`
		// Definitions is the definitions of the schema as specified in section 8.2.4 of RFC
		// draft-bhutton-json-schema-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-8.2.4
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
		// instance and as a name within this field's value, the child instance
		// for that name successfully validates against the corresponding schema.
		//
		// The annotation result of this field is the set of instance property
		// names matched by this field.
		//
		// Omitting this field has the same assertion behavior as an empty
		// object.
		Definitions schemaDefinitions `json:"$defs,omitempty"`
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
		Comments string `json:"$comment,omitempty"`
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
		AllOf []*Schema `json:"allOf,omitempty"`
		// AnyOf is the any of of the schema as specified in section 10.2.1.2
		// of RFC draft-bhutton-json-schema-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-10.2.1.2
		//
		// The value of this field MUST be an array.  Elements in the array
		// MUST be objects.  Each object MUST be a valid JSON Schema.
		//
		// An instance validates successfully against this field if it
		// validates successfully against at least one schema defined by
		// "anyOf".
		//
		// Omitting this field has the same behavior as an empty array.
		AnyOf []*Schema `json:"anyOf,omitempty"`
		// OneOf is the one of of the schema as specified in section 10.2.1.3
		// of RFC draft-bhutton-json-schema-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-10.2.1.3
		//
		// The value of this field MUST be an array.  Elements in the array
		// MUST be objects.  Each object MUST be a valid JSON Schema.
		//
		// An instance validates successfully against this field if it
		// validates successfully against exactly one schema defined by
		// "oneOf".
		//
		// Omitting this field has the same behavior as an empty array.
		OneOf []*Schema `json:"oneOf,omitempty"`
		// Not is the not of the schema as specified in section 10.2.1.4 of
		// RFC draft-bhutton-json-schema-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-10.2.1.4
		//
		// The value of this field MUST be an object.  This object MUST be a
		// valid JSON Schema.
		//
		// An instance validates successfully against this field if it
		// validates successfully against the schema defined by "not".
		//
		// Omitting this field has the same behavior as an empty object.
		Not *Schema `json:"not,omitempty"`
		// RFC draft-bhutton-json-schema-00 section 10.2.2 (Apply sub-schemas conditionally)
		If *Schema `json:"if,omitempty"` // section 10.2.2.1
		// Then is the then of the schema as specified in section 10.2.2.2 of RFC
		// draft-bhutton-json-schema-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-10.2.2.2
		//
		// The value of this field MUST be an object.  This object MUST be a
		// valid JSON Schema.
		//
		// An instance validates successfully against this field if it
		// validates successfully against the schema defined by "then".
		//
		// Omitting this field has the same behavior as an empty object.
		Then *Schema `json:"then,omitempty"` // section 10.2.2.2
		// Else is the else of the schema as specified in section 10.2.2.3 of RFC
		// draft-bhutton-json-schema-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-10.2.2.3
		//
		// The value of this field MUST be an object.  This object MUST be a
		// valid JSON Schema.
		//
		// An instance validates successfully against this field if it
		// validates successfully against the schema defined by "else".
		//
		// Omitting this field has the same behavior as an empty object.
		Else *Schema `json:"else,omitempty"` // section 10.2.2.3
		// DependentSchemas is the dependent schemas of the schema as specified in section 10.2.2.4 of RFC
		// draft-bhutton-json-schema-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-10.2.2.4
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
		DependentSchemas map[string]*Schema `json:"dependentSchemas,omitempty"` // section 10.2.2.4
		// PrefixItems is the prefix items of the schema as specified in section 10.3.1.1 of RFC
		// draft-bhutton-json-schema-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-10.3.1.1
		//
		// The value of this field MUST be an array.  Elements in the array,
		// if any, MUST be valid JSON Schemas.
		//
		// An array instance is valid against "prefixItems" if its length is
		// greater than or equal to the value of "minItems" and if each item
		// in the instance array is valid against the schema defined by the
		// corresponding item in "prefixItems".
		//
		// Omitting this field has the same behavior as an empty array.
		PrefixItems []*Schema `json:"prefixItems,omitempty"` // section 10.3.1.1
		// Items is the items of the schema as specified in section 10.3.1.2 of RFC
		// draft-bhutton-json-schema-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-10.3.1.2
		//
		// The value of "items" MUST be a valid JSON Schema.
		//
		// This keyword applies its subschema to all instance elements at
		// indexes greater than the length of the "prefixItems" array in the
		// same schema object, as reported by the annotation result of that
		// "prefixItems" keyword.  If no such annotation result exists, "items"
		// applies its subschema to all instance array elements.  [[CREF11: Note
		// that the behavior of "items" without "prefixItems" is identical to
		// that of the schema form of "items" in prior drafts.  When
		// "prefixItems" is present, the behavior of "items" is identical to the
		// former "additionalItems" keyword.  ]]
		//
		// If the "items" subschema is applied to any positions within the
		// instance array, it produces an annotation result of boolean true,
		// indicating that all remaining array elements have been evaluated
		// against this keyword's subschema.
		//
		// Omitting this keyword has the same assertion behavior as an empty
		// schema.
		Items *Schema `json:"items,omitempty"` // section 10.3.1.2  (replaces additionalItems)
		// Contains is the contains of the schema as specified in section 10.3.1.3 of RFC
		// draft-bhutton-json-schema-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-10.3.1.3
		//
		//
		// The value of this field MUST be a valid JSON Schema.
		//
		// An array instance is valid against "contains" if at least one of its
		// elements is valid against the given schema.  The subschema MUST be
		// applied to every array element even after the first match has been
		// found, in order to collect annotations for use by other fields.
		// This is to ensure that all possible annotations are collected.
		Contains *Schema `json:"contains,omitempty"` // section 10.3.1.3
		// RFC draft-bhutton-json-schema-00 section 10.3.2 (sub-schemas)
		// Properties are the properties of the schema as specified in section 10.3.2.1 of RFC
		// draft-bhutton-json-schema-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-10.3.2.1
		//
		// The value of "properties" MUST be an object.  Each value of this
		// object MUST be a valid JSON Schema.
		//
		// Validation succeeds if, for each name that appears in both the
		// instance and as a name within this field's value, the child
		// instance for that name successfully validates against the
		// corresponding schema.
		//
		// The annotation result of this field is the set of instance property
		// names matched by this field.
		//
		// Omitting this field has the same assertion behavior as an empty
		// object.
		Properties *omap.OrderedMap[string, *Schema] `json:"properties,omitempty"`
		// PatternProperties are the pattern properties of the schema as specified in section 10.3.2.2 of RFC
		// draft-bhutton-json-schema-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-10.3.2.2
		//
		// The value of "patternProperties" MUST be an object.  Each property
		// name of this object SHOULD be a valid regular expression, according
		// to the ECMA-262 regular expression dialect.  Each property value of
		// this object MUST be a valid JSON Schema.
		//
		// Validation succeeds if, for each instance name that matches any
		// regular expressions that appear as a property name in this field's
		// value, the child instance for that name successfully validates
		// against each schema that corresponds to a matching regular
		// expression.
		//
		// The annotation result of this field is the set of instance property
		// names matched by this field.
		//
		// Omitting this field has the same assertion behavior as an empty
		// object.
		PatternProperties map[string]*Schema `json:"patternProperties,omitempty"` // section 10.3.2.2
		// AdditionalProperties is the additional properties of the schema as
		// specified in section 10.3.2.3 of RFC
		// draft-bhutton-json-schema-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-10.3.2.3
		//
		//
		// The value of "additionalProperties" MUST be a valid JSON Schema.
		//
		// The behavior of this field depends on the presence and annotation
		// results of "properties" and "patternProperties" within the same
		// schema object.  Validation with "additionalProperties" applies only
		// to the child values of instance names that do not appear in the
		// annotation results of either "properties" or "patternProperties".
		//
		// For all such properties, validation succeeds if the child instance
		// validates against the "additionalProperties" schema.
		//
		// The annotation result of this field is the set of instance property
		// names validated by this field's subschema.
		//
		// Omitting this field has the same assertion behavior as an empty
		// schema.
		//
		// Implementations MAY choose to implement or optimize this field in
		// another way that produces the same effect, such as by directly
		// checking the names in "properties" and the patterns in
		// "patternProperties" against the instance property set.
		// Implementations that do not support annotation collection MUST do so.
		AdditionalProperties *Schema `json:"additionalProperties,omitempty"` // section 10.3.2.3
		// PropertyNames is the property names of the schema as specified in
		// section 10.3.2.4 of RFC
		// draft-bhutton-json-schema-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-10.3.2.4
		//
		// The value of this field MUST be an object.  Properties in this
		// object, if any, MUST be arrays.  Elements in each array, if any,
		// MUST be strings, and MUST be unique.
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
		PropertyNames *Schema `json:"propertyNames,omitempty"` // section 10.3.2.4
		// Type is the type of the schema as specified in section 6.1.1 of
		// RFC draft-bhutton-json-schema-validation-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.1.1
		//
		// The value of this field MUST be a string.  This string SHOULD be a
		// valid JSON Schema type.
		//
		// Omitting this field has the same behavior as an empty string.
		Type string `json:"type,omitempty"` // section 6.1.1
		// Enum is the enum of the schema as specified in section 6.1.2 of
		// RFC draft-bhutton-json-schema-validation-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.1.2
		//
		// The value of this field MUST be an array.  Elements in the array,
		// if any, MUST be unique.
		//
		// A numeric instance is valid against "enum" if its value is equal
		// to one of the values in the array.
		//
		// Omitting this field has the same behavior as an empty array.
		Enum []any `json:"enum,omitempty"`
		// Const is the const of the schema as specified in section 6.1.3 of
		// RFC draft-bhutton-json-schema-validation-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.1.3
		//
		// The value of this field MUST be an instance of the data type
		// defined by the "type" field.
		//
		// A numeric instance is valid against "const" if its value is equal
		// to the value of this field.
		//
		// Omitting this field has the same behavior as an empty value.
		Const any `json:"const,omitempty"`
		// MultipleOf specifies the multiple of the schema as specified in
		// section 6.2.1 of RFC draft-bhutton-json-schema-validation-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.2.1
		//
		// The value of this field MUST be a JSON number, representing an
		// instance of the data type defined by the "type" field.
		//
		// A numeric instance is valid against "multipleOf" if the result of
		// the division of the instance by this field's value leaves no
		// remainder.
		//
		// Omitting this field has the same behavior as an empty value.
		MultipleOf json.Number `json:"multipleOf,omitempty"`
		// Maximum is the maximum of the schema as specified in section 6.2.2
		// of RFC draft-bhutton-json-schema-validation-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.2.2
		//
		// The value of this field MUST be a JSON number, representing an
		// instance of the data type defined by the "type" field.
		//
		// A numeric instance is valid against "maximum" if it has a value
		// less than the value of "exclusiveMaximum" and it has a value
		// greater than the value of "minimum".
		//
		// Omitting this field has the same behavior as an empty value.
		Maximum json.Number `json:"maximum,omitempty"`
		// ExclusiveMaximum is the exclusive maximum of the schema as specified
		// in section 6.2.3 of RFC
		// draft-bhutton-json-schema-validation-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.2.3
		//
		// The value of this field MUST be a JSON number, representing an
		// instance of the data type defined by the "type" field.
		//
		// A numeric instance is valid against "exclusiveMaximum" if it has a
		// value less than the value of "minimum" and it has a value greater
		// than the value of "exclusiveMinimum".
		//
		// Omitting this field has the same behavior as an empty value.
		ExclusiveMaximum json.Number `json:"exclusiveMaximum,omitempty"`
		// Minimum is the minimum of the schema as specified in section 6.2.4
		// of RFC draft-bhutton-json-schema-validation-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.2.4
		//
		// The value of this field MUST be a JSON number, representing an
		// instance of the data type defined by the "type" field.
		//
		// A numeric instance is valid against "minimum" if it has a value
		// greater than the value of "exclusiveMinimum" and it has a value
		// less than the value of "maximum".
		//
		// Omitting this field has the same behavior as an empty value.
		Minimum json.Number `json:"minimum,omitempty"` // section 6.2.4
		// ExclusiveMinimum is the exclusive minimum of the schema as specified
		// in section 6.2.5 of RFC
		// draft-bhutton-json-schema-validation-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.2.5
		//
		// The value of this field MUST be a JSON number, representing an
		// instance of the data type defined by the "type" field.
		//
		// A numeric instance is valid against "exclusiveMinimum" if it has a
		// value less than the value of "minimum" and it has a value greater
		// than the value of "exclusiveMaximum".
		//
		// Omitting this field has the same behavior as an empty value.
		ExclusiveMinimum json.Number `json:"exclusiveMinimum,omitempty"` // section 6.2.5
		// MaxLength specifies the maximum length of the string as specified in
		// section 6.3.1 of RFC
		// draft-bhutton-json-schema-validation-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.3.1
		//
		// The value of this field MUST be a non-negative integer.
		//
		// A string instance is valid against "maxLength" if its length is
		// less than, or equal to, the value of this field.
		//
		// Omitting this field has the same behavior as a value of an
		// implementation-defined number.
		MaxLength *uint64 `json:"maxLength,omitempty"` // section 6.3.1
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
		DependentRequired map[string][]string `json:"dependentRequired,omitempty"`
		// Format specifies the format of the schema as specified in section
		// 7.3 of RFC draft-bhutton-json-schema-validation-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-7.3
		//
		// The value of this field MUST be a string.  Implementations that
		// use a subset of JSON as their input format, such as JSON Hyper-Schema
		// or JSON Schema Hyper-Schema, MAY implement validation against
		// meta-schemas that define format-specific fields that describe
		// additional constraints beyond those specified herein.
		Format string `json:"format,omitempty"`
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
		ContentEncoding string `json:"contentEncoding,omitempty"`
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
		ContentMediaType string `json:"contentMediaType,omitempty"`
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
		ContentSchema *Schema `json:"contentSchema,omitempty"`
		// Title is the title of the schema as specified in section 9.1 of
		// RFC draft-bhutton-json-schema-validation-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-9.1
		//
		// The value of this field MUST be a string.  This string SHOULD be a
		// short description of the schema.  The value of this field SHOULD be
		// true when the instance described by this schema is a boolean.
		//
		// Omitting this field has the same behavior as an empty string.
		Title string `json:"title,omitempty"`
		// Description is the description of the schema as specified in section 9.1
		// of RFC draft-bhutton-json-schema-validation-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-9.1
		//
		// The value of this field MUST be a string.  This string SHOULD be a
		// description of the schema.  The value of this field SHOULD be
		// true when the instance described by this schema is a boolean.
		//
		// Omitting this field has the same behavior as an empty string.
		Description string `json:"description,omitempty"`
		// Default is the default of the schema as specified in section 9.2 of
		// RFC draft-bhutton-json-schema-validation-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-9.2
		//
		// The value of this field MUST be an instance of the data type defined
		// by the "type" field.  This instance SHOULD be used as the default
		// value of the instance if the instance is undefined or its value is
		// equal to null.
		//
		// Omitting this field has the same behavior as an empty value.
		Default any `json:"default,omitempty"`
		// Deprecated is the deprecated of the schema as specified in section 9.3
		// of RFC draft-bhutton-json-schema-validation-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-9.3
		//
		// The value of this field MUST be a boolean.  This boolean SHOULD be
		// true when the instance described by this schema is deprecated.
		//
		// Omitting this field has the same behavior as false.
		Deprecated bool `json:"deprecated,omitempty"`
		// ReadOnly is the read only of the schema as specified in section 9.4
		// of RFC draft-bhutton-json-schema-validation-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-9.4
		//
		// The value of this field MUST be a boolean.  This boolean SHOULD be
		// true when the instance described by this schema is read only.
		//
		// Omitting this field has the same behavior as false.
		ReadOnly bool `json:"readOnly,omitempty"` // section 9.4
		// WriteOnly is the write only of the schema as specified in section 9.4
		// of RFC draft-bhutton-json-schema-validation-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-9.4
		//
		// The value of this field MUST be a boolean.  This boolean SHOULD be
		// true when the instance described by this schema is write only.
		//
		// Omitting this field has the same behavior as false.
		WriteOnly bool `json:"writeOnly,omitempty"`
		// Examples is the examples of the schema as specified in section 9.5
		// of RFC draft-bhutton-json-schema-validation-00.
		//
		// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-9.5
		//
		// The value of this field MUST be an array.  Elements in this array,
		// if any, MUST be valid against the "items" schema that describes
		// the type of the array.
		//
		// Omitting this field has the same behavior as an empty array.
		Examples []any `json:"examples,omitempty"`
		// Extras holds additional information about the schema.
		//
		// This field is not part of the official JSON Schema specification.
		Extras map[string]any `json:"-"`
		// Special boolean representation of the Schema - section 4.3.2
		boolean *bool
	}
	// schemaDefinitions hold schema schemaDefinitions.
	//
	// http://json-schema.org/latest/json-schema-validation.html#rfc.section.5.26
	//
	// RFC draft-wright-json-schema-validation-00, section 5.26
	schemaDefinitions map[string]*Schema
	// schemaID represents a Schema schemaID type which should always be a URI.
	// See draft-bhutton-json-schema-00 section 8.2.1
	schemaID string
	// customSchemaImpl is used to detect if the type provides it's own
	// custom Schema Type definition to use instead. Very useful for situations
	// where there are custom JSON Marshal and Unmarshal methods.
	customSchemaImpl interface {
		JSONSchema() *Schema
	}
	// Function to be run after the schema has been generated.
	// this will let you modify a schema afterwards
	extendSchemaImpl interface {
		JSONSchemaExtend(*Schema)
	}
	// If the object to be reflected defines a `JSONSchemaAlias` method, its type will
	// be used instead of the original type.
	aliasSchemaImpl interface {
		JSONSchemaAlias() any
	}
	// If an object to be reflected defines a `JSONSchemaPropertyAlias` method,
	// it will be called for each property to determine if another object
	// should be used for the contents.
	propertyAliasSchemaImpl interface {
		JSONSchemaProperty(prop string) any
	}
	// customSchemaGetFieldDocString
	customSchemaGetFieldDocString interface {
		GetFieldDocString(fieldName string) string
	}
	customGetFieldDocString func(fieldName string) string
)

// Go code generated from protobuf enum types should fulfil this interface.
type protoEnum interface {
	EnumDescriptor() ([]byte, []int)
}

func appendUniqueString(base []string, value string) []string {
	for _, v := range base {
		if v == value {
			return base
		}
	}
	return append(base, value)
}
func (t *Schema) fieldsFromTags(
	f reflect.StructField,
	parent *Schema,
	propertyName string,
) {
	t.Description = f.Tag.Get("jsonschema_description")
	tags := splitOnUnescapedCommas(f.Tag.Get("jsonschema"))
	tags = t.genericfields(tags, parent, propertyName)
	switch t.Type {
	case "string":
		t.stringfields(tags)
	case "number":
		t.numericalfields(tags)
	case "integer":
		t.numericalfields(tags)
	case "array":
		t.arrayfields(tags)
	case "boolean":
		t.booleanfields(tags)
	}
	extras := strings.Split(f.Tag.Get("jsonschema_extras"), ",")
	t.extrafields(extras)
}

// genericfields reads struct tags for generic keywords
func (t *Schema) genericfields(
	tags []string,
	parent *Schema,
	propertyName string,
) []string {
	unprocessed := make([]string, 0, len(tags))
	for _, tag := range tags {
		nameValue := strings.SplitN(tag, "=", 2)
		if len(nameValue) == 2 {
			name, val := nameValue[0], nameValue[1]
			switch name {
			case "title":
				t.Title = val
			case "description":
				t.Description = val
			case "type":
				t.Type = val
			case "anchor":
				t.Anchor = val
			case "oneof_required":
				var typeFound *Schema
				for i := range parent.OneOf {
					if parent.OneOf[i].Title == nameValue[1] {
						typeFound = parent.OneOf[i]
					}
				}
				if typeFound == nil {
					typeFound = &Schema{
						Title:    nameValue[1],
						Required: []string{},
					}
					parent.OneOf = append(parent.OneOf, typeFound)
				}
				typeFound.Required = append(typeFound.Required, propertyName)
			case "anyof_required":
				var typeFound *Schema
				for i := range parent.AnyOf {
					if parent.AnyOf[i].Title == nameValue[1] {
						typeFound = parent.AnyOf[i]
					}
				}
				if typeFound == nil {
					typeFound = &Schema{
						Title:    nameValue[1],
						Required: []string{},
					}
					parent.AnyOf = append(parent.AnyOf, typeFound)
				}
				typeFound.Required = append(typeFound.Required, propertyName)
			case "oneof_ref":
				subSchema := t
				if t.Items != nil {
					subSchema = t.Items
				}
				if subSchema.OneOf == nil {
					subSchema.OneOf = make([]*Schema, 0, 1)
				}
				subSchema.Ref = ""
				refs := strings.Split(nameValue[1], ";")
				for _, r := range refs {
					subSchema.OneOf = append(subSchema.OneOf, &Schema{
						Ref: r,
					})
				}
			case "oneof_type":
				if t.OneOf == nil {
					t.OneOf = make([]*Schema, 0, 1)
				}
				t.Type = ""
				types := strings.Split(nameValue[1], ";")
				for _, ty := range types {
					t.OneOf = append(t.OneOf, &Schema{
						Type: ty,
					})
				}
			case "anyof_ref":
				subSchema := t
				if t.Items != nil {
					subSchema = t.Items
				}
				if subSchema.AnyOf == nil {
					subSchema.AnyOf = make([]*Schema, 0, 1)
				}
				subSchema.Ref = ""
				refs := strings.Split(nameValue[1], ";")
				for _, r := range refs {
					subSchema.AnyOf = append(subSchema.AnyOf, &Schema{
						Ref: r,
					})
				}
			case "anyof_type":
				if t.AnyOf == nil {
					t.AnyOf = make([]*Schema, 0, 1)
				}
				t.Type = ""
				types := strings.Split(nameValue[1], ";")
				for _, ty := range types {
					t.AnyOf = append(t.AnyOf, &Schema{
						Type: ty,
					})
				}
			default:
				unprocessed = append(unprocessed, tag)
			}
		}
	}
	return unprocessed
}

// read struct tags for boolean type fields
func (t *Schema) booleanfields(tags []string) {
	for _, tag := range tags {
		nameValue := strings.Split(tag, "=")
		if len(nameValue) != 2 {
			continue
		}
		name, val := nameValue[0], nameValue[1]
		if name == "default" {
			if val == "true" {
				t.Default = true
				continue
			}
			if val == "false" {
				t.Default = false
				continue
			}
		}
	}
}

// read struct tags for string type fields
func (t *Schema) stringfields(tags []string) {
	for _, tag := range tags {
		nameValue := strings.SplitN(tag, "=", 2)
		if len(nameValue) == 2 {
			name, val := nameValue[0], nameValue[1]
			switch name {
			case "minLength":
				t.MinLength = parseUint(val)
			case "maxLength":
				t.MaxLength = parseUint(val)
			case "pattern":
				t.Pattern = val
			case "format":
				t.Format = val
			case "readOnly":
				i, _ := strconv.ParseBool(val)
				t.ReadOnly = i
			case "writeOnly":
				i, _ := strconv.ParseBool(val)
				t.WriteOnly = i
			case "default":
				t.Default = val
			case "example":
				t.Examples = append(t.Examples, val)
			case "enum":
				t.Enum = append(t.Enum, val)
			}
		}
	}
}

// read struct tags for numerical type fields
func (t *Schema) numericalfields(tags []string) {
	for _, tag := range tags {
		nameValue := strings.Split(tag, "=")
		if len(nameValue) == 2 {
			name, val := nameValue[0], nameValue[1]
			switch name {
			case "multipleOf":
				t.MultipleOf, _ = toJSONNumber(val)
			case "minimum":
				t.Minimum, _ = toJSONNumber(val)
			case "maximum":
				t.Maximum, _ = toJSONNumber(val)
			case "exclusiveMaximum":
				t.ExclusiveMaximum, _ = toJSONNumber(val)
			case "exclusiveMinimum":
				t.ExclusiveMinimum, _ = toJSONNumber(val)
			case "default":
				if num, ok := toJSONNumber(val); ok {
					t.Default = num
				}
			case "example":
				if num, ok := toJSONNumber(val); ok {
					t.Examples = append(t.Examples, num)
				}
			case "enum":
				if num, ok := toJSONNumber(val); ok {
					t.Enum = append(t.Enum, num)
				}
			}
		}
	}
}

// read struct tags for array type fields
func (t *Schema) arrayfields(tags []string) {
	var defaultValues []any
	unprocessed := make([]string, 0, len(tags))
	for _, tag := range tags {
		nameValue := strings.Split(tag, "=")
		if len(nameValue) == 2 {
			name, val := nameValue[0], nameValue[1]
			switch name {
			case "minItems":
				t.MinItems = parseUint(val)
			case "maxItems":
				t.MaxItems = parseUint(val)
			case "uniqueItems":
				t.UniqueItems = true
			case "default":
				defaultValues = append(defaultValues, val)
			case "format":
				t.Items.Format = val
			case "pattern":
				t.Items.Pattern = val
			default:
				unprocessed = append(
					unprocessed,
					tag,
				) // left for further processing by underlying type
			}
		}
	}
	if len(defaultValues) > 0 {
		t.Default = defaultValues
	}
	if len(unprocessed) == 0 {
		return
	}
	switch t.Items.Type {
	case "string":
		t.Items.stringfields(unprocessed)
	case "number":
		t.Items.numericalfields(unprocessed)
	case "integer":
		t.Items.numericalfields(unprocessed)
	case "array":
		// explicitly don't support traversal for the [][]..., as it's unclear where the array tags belong
	case "boolean":
		t.Items.booleanfields(unprocessed)
	}
}
func (t *Schema) extrafields(tags []string) {
	for _, tag := range tags {
		nameValue := strings.SplitN(tag, "=", 2)
		if len(nameValue) == 2 {
			t.setExtra(nameValue[0], nameValue[1])
		}
	}
}
func (t *Schema) setExtra(key, val string) {
	if t.Extras == nil {
		t.Extras = map[string]any{}
	}
	if existingVal, ok := t.Extras[key]; ok {
		switch existingVal := existingVal.(type) {
		case string:
			t.Extras[key] = []string{existingVal, val}
		case []string:
			t.Extras[key] = append(existingVal, val)
		case int:
			t.Extras[key], _ = strconv.Atoi(val)
		case bool:
			t.Extras[key] = (val == "true" || val == "t")
		}
		return
	}
	switch key {
	case "minimum":
		t.Extras[key], _ = strconv.Atoi(val)
	default:
		var x any
		if val == "true" {
			x = true
		} else if val == "false" {
			x = false
		} else {
			x = val
		}
		t.Extras[key] = x
	}
}
func requiredFromJSONTags(tags []string, val *bool) {
	if ignoredByJSONTags(tags) {
		return
	}
	for _, tag := range tags[1:] {
		if tag == "omitempty" {
			*val = false
			return
		}
	}
	*val = true
}
func requiredFromJSONSchemaTags(tags []string, val *bool) {
	if ignoredByJSONSchemaTags(tags) {
		return
	}
	for _, tag := range tags {
		if tag == "required" {
			*val = true
		}
	}
}
func nullableFromJSONSchemaTags(tags []string) bool {
	if ignoredByJSONSchemaTags(tags) {
		return false
	}
	for _, tag := range tags {
		if tag == "nullable" {
			return true
		}
	}
	return false
}
func ignoredByJSONTags(tags []string) bool {
	return tags[0] == "-"
}
func ignoredByJSONSchemaTags(tags []string) bool {
	return tags[0] == "-"
}
func inlinedByJSONTags(tags []string) bool {
	for _, tag := range tags[1:] {
		if tag == "inline" {
			return true
		}
	}
	return false
}

// toJSONNumber converts string to *json.Number.
// It'll aso return whether the number is valid.
func toJSONNumber(s string) (json.Number, bool) {
	num := json.Number(s)
	if _, err := num.Int64(); err == nil {
		return num, true
	}
	if _, err := num.Float64(); err == nil {
		return num, true
	}
	return json.Number(""), false
}
func parseUint(num string) *uint64 {
	val, err := strconv.ParseUint(num, 10, 64)
	if err != nil {
		return nil
	}
	return &val
}
func (r *reflector) fieldNameTag() string {
	if r.FieldNameTag != "" {
		return r.FieldNameTag
	}
	return "json"
}
func (r *reflector) reflectFieldName(
	f reflect.StructField,
) (string, bool, bool, bool) {
	jsonTagString := f.Tag.Get(r.fieldNameTag())
	jsonTags := strings.Split(jsonTagString, ",")
	if ignoredByJSONTags(jsonTags) {
		return "", false, false, false
	}
	schemaTags := strings.Split(f.Tag.Get("jsonschema"), ",")
	if ignoredByJSONSchemaTags(schemaTags) {
		return "", false, false, false
	}
	var required bool
	if !r.RequiredFromJSONSchemaTags {
		requiredFromJSONTags(jsonTags, &required)
	}
	requiredFromJSONSchemaTags(schemaTags, &required)
	nullable := nullableFromJSONSchemaTags(schemaTags)
	if f.Anonymous && jsonTags[0] == "" {
		// As per JSON Marshal rules, anonymous structs are inherited
		if f.Type.Kind() == reflect.Struct {
			return "", true, false, false
		}
		// As per JSON Marshal rules, anonymous pointer to structs are inherited
		if f.Type.Kind() == reflect.Ptr &&
			f.Type.Elem().Kind() == reflect.Struct {
			return "", true, false, false
		}
	}
	// As per JSON Marshal rules, inline nested structs that have `inline` tag.
	if inlinedByJSONTags(jsonTags) {
		return "", true, false, false
	}
	// Try to determine the name from the different combos
	name := f.Name
	if jsonTags[0] != "" {
		name = jsonTags[0]
	}
	if !f.Anonymous && f.PkgPath != "" {
		// field not anonymous and not export has no export name
		name = ""
		return name, false, required, nullable
	}
	if r.KeyNamer != nil {
		name = r.KeyNamer(name)
		return name, false, required, nullable
	}
	return name, true, required, nullable
}

// UnmarshalJSON is used to parse a schema object or boolean.
func (t *Schema) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("true")) {
		*t = *trueSchema
		return nil
	} else if bytes.Equal(data, []byte("false")) {
		*t = *falseSchema
		return nil
	}
	type SchemaAlt Schema
	aux := &struct {
		*SchemaAlt
	}{
		SchemaAlt: (*SchemaAlt)(t),
	}
	return json.Unmarshal(data, aux)
}

// MarshalJSON is used to serialize a schema object or boolean.
func (t *Schema) MarshalJSON() ([]byte, error) {
	if t.boolean != nil {
		if *t.boolean {
			return []byte("true"), nil
		}
		return []byte("false"), nil
	}
	if reflect.DeepEqual(&Schema{}, t) {
		// Don't bother returning empty schemas
		return []byte("true"), nil
	}
	type SchemaAlt Schema
	b, err := json.Marshal((*SchemaAlt)(t))
	if err != nil {
		return nil, err
	}
	if len(t.Extras) == 0 {
		return b, nil
	}
	m, err := json.Marshal(t.Extras)
	if err != nil {
		return nil, err
	}
	if len(b) == 2 {
		return m, nil
	}
	b[len(b)-1] = ','
	return append(b, m[1:]...), nil
}
func (r *reflector) typeName(t reflect.Type) string {
	if r.Namer != nil {
		if name := r.Namer(t); name != "" {
			return name
		}
	}
	return t.Name()
}

// Split on commas that are not preceded by `\`.
//
// This way, we prevent splitting regexes.
func splitOnUnescapedCommas(tagString string) []string {
	ret := make([]string, 0)
	separated := strings.Split(tagString, ",")
	ret = append(ret, separated[0])
	i := 0
	for _, nextTag := range separated[1:] {
		if len(ret[i]) == 0 {
			ret = append(ret, nextTag)
			i++
			continue
		}
		if ret[i][len(ret[i])-1] == '\\' {
			ret[i] = ret[i][:len(ret[i])-1] + "," + nextTag
			continue
		}
		ret = append(ret, nextTag)
		i++
	}
	return ret
}
func fullyQualifiedTypeName(t reflect.Type) string {
	return t.PkgPath() + "." + t.Name()
}

// ToSnakeCase converts the provided string into snake case using dashes.
// This is useful for Schema IDs and definitions to be coherent with
// common JSON Schema examples.
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}-${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}-${2}")
	return strings.ToLower(snake)
}

// newProperties is a helper method to instantiate a new properties ordered
// map.
func newProperties() *omap.OrderedMap[string, *Schema] {
	return omap.New[string, *Schema]()
}

// Validate is used to check if the ID looks like a proper schema.
// This is done by parsing the ID as a URL and checking it has all the
// relevant parts.
func (i schemaID) Validate() error {
	u, err := url.Parse(string(i))
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if u.Hostname() == "" {
		return fmt.Errorf("missing hostname: %s", u.Hostname())
	}
	if !strings.Contains(u.Hostname(), ".") {
		return fmt.Errorf("hostname does not look valid: %s", u.Hostname())
	}
	if u.Path == "" {
		return fmt.Errorf("path is expected: %s", u.Path)
	}
	if u.Scheme != "https" && u.Scheme != "http" {
		return fmt.Errorf("unexpected schema: %s", u.Scheme)
	}
	return nil
}

// Anchor sets the anchor part of the schema URI.
func (i schemaID) Anchor(name string) schemaID {
	b := i.Base()
	return schemaID(string(b) + "#" + name)
}

// Def adds or replaces a definition identifier.
func (i schemaID) Def(name string) schemaID {
	b := i.Base()
	return schemaID(string(b) + "#/$defs/" + name)
}

// Add appends the provided path to the id, and removes any
// anchor data that might be there.
func (i schemaID) Add(path string) schemaID {
	b := i.Base()
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return schemaID(string(b) + path)
}

// Base removes any anchor information from the schema
func (i schemaID) Base() schemaID {
	s := string(i)
	li := strings.LastIndex(s, "#")
	if li != -1 {
		s = s[0:li]
	}
	s = strings.TrimRight(s, "/")
	return schemaID(s)
}
