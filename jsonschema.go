// Package jsonschema provides very simple functionality for representing a JSON schema as a
// (nested) struct. This struct can be used with the chat completion "function call" feature.
// For more complicated schemas, it is recommended to use a dedicated JSON schema library
// and/or pass in the schema in []byte format.
package groq

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// DataType is the type of a JSON schema.
// string
type DataType string

const (
	Number  DataType = "number"  // Number is the type of a number.
	String  DataType = "string"  // String is the type of a string.
	Boolean DataType = "boolean" // Boolean is the type of a boolean.
	Integer DataType = "integer" // Integer is the type of an integer.
	Object  DataType = "object"  // Object is the type of an object.
	Array   DataType = "array"   // Array is the type of an array.
	Enum    DataType = "enum"    // Enum is the type of an enum.
	AnyOf   DataType = "anyOf"   // AnyOf is the type of an anyOf.

	Description           = "description"         // Description is the description of the schema.
	DescriptionKey        = "desc"                // DescriptionKey is the description of the schema.
	DescriptionAutoGenFmt = "%s is the %s of %s." // DescriptionAutoGenFmt is the format of the description auto generated.

	Type = "type" // Type is the type of the schema.

	MinItems    = "minItems" // MinItems is the minimum number of items in an array.
	MinItemsKey = "min"      // MinItemsKey is the minimum number of items in an array.

	UniqueItems    = "uniqueItems" // UniqueItems is a boolean value that specifies whether all items in an array must be unique.
	UniqueItemsKey = "unique"      // UniqueItemsKey is a boolean value that specifies whether all items in an array must be unique.

	Required    = "required" // Required is a boolean value that specifies whether a field is required.
	RequiredKey = "req"      // RequiredKey is a boolean value that specifies whether a field is required.
)

var (
	// BaseKeys are the base keys for the base object schema structure definition.
	BaseKeys = []string{
		RequiredKey,
	}
)

// Definition is a struct for describing a JSON Schema.
// It is fairly limited, and you may have better luck using a third-party library.
type Definition struct {
	// Type specifies the data type of the schema.
	Type DataType `json:"type,omitempty"`
	// Description is the description of the schema.
	Description string `json:"description,omitempty"`
	// Enum is used to restrict a value to a fixed set of values. It must be an array with at least
	// one element, where each element is unique. You will probably only use this with strings.
	Enum []string `json:"enum,omitempty"`
	// Properties describes the properties of an object, if the schema type is Object.
	Properties map[string]Definition `json:"properties,omitempty"`
	// Required specifies which properties are required, if the schema type is Object.
	Required []string `json:"required,omitempty"`
	// Items specifies which data type an array contains, if the schema type is Array.
	Items *Definition `json:"items,omitempty"`
	// AdditionalProperties is used to control the handling of properties in an object
	// that are not explicitly defined in the properties section of the schema. example:
	// additionalProperties: true
	// additionalProperties: false
	// additionalProperties: jsonschema.Definition{Type: jsonschema.String}
	AdditionalProperties any `json:"additionalProperties,omitempty"`
}

// MarshalJSON marshals a JSON schema to JSON.
func (d *Definition) MarshalJSON() ([]byte, error) {
	if d.Properties == nil {
		d.Properties = make(map[string]Definition)
	}
	type Alias Definition
	return json.Marshal(struct {
		Alias
	}{
		Alias: (Alias)(*d),
	})
}

// GenerateSchemaForType generates a JSON schema for a given type.
func GenerateSchemaForType(v any) (*Definition, error) {
	return reflectSchema(reflect.TypeOf(v))
}

func reflectSchema(t reflect.Type) (*Definition, error) {
	var d Definition
	switch t.Kind() {
	case reflect.String:
		d.Type = String
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		d.Type = Integer
	case reflect.Float32, reflect.Float64:
		d.Type = Number
	case reflect.Bool:
		d.Type = Boolean
	case reflect.Slice, reflect.Array:
		d.Type = Array
		items, err := reflectSchema(t.Elem())
		if err != nil {
			return nil, err
		}
		d.Items = items
	case reflect.Struct:
		d.Type = Object
		d.AdditionalProperties = false
		object, err := reflectSchemaObject(t)
		if err != nil {
			return nil, err
		}
		d = *object
	case reflect.Ptr:
		definition, err := reflectSchema(t.Elem())
		if err != nil {
			return nil, err
		}
		d = *definition
	case reflect.Invalid,
		reflect.Uintptr,
		reflect.Complex64,
		reflect.Complex128,
		reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.UnsafePointer:
		return nil, fmt.Errorf("unsupported type: %s", t.Kind().String())
	default:
	}
	return &d, nil
}

func reflectSchemaObject(t reflect.Type) (*Definition, error) {
	var d = Definition{
		Type:                 Object,
		AdditionalProperties: false,
	}
	properties := make(map[string]Definition)
	var requiredFields []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		jsonTag := field.Tag.Get("json")
		var required = true
		if jsonTag == "" {
			jsonTag = field.Name
		} else if strings.HasSuffix(jsonTag, ",omitempty") {
			jsonTag = strings.TrimSuffix(jsonTag, ",omitempty")
			required = false
		}

		item, err := reflectSchema(field.Type)
		if err != nil {
			return nil, err
		}
		description := field.Tag.Get("description")
		if description != "" {
			item.Description = description
		}
		properties[jsonTag] = *item

		if s := field.Tag.Get("required"); s != "" {
			required, _ = strconv.ParseBool(s)
		}
		if required {
			requiredFields = append(requiredFields, jsonTag)
		}
	}
	d.Required = requiredFields
	d.Properties = properties
	return &d, nil
}