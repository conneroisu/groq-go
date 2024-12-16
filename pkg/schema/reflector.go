package schema

import "reflect"

type (

	// A reflector reflects values into a Schema.
	reflector struct {
		// BaseSchemaID defines the URI that will be used as a base to determine
		// Schema IDs for models. For example, a base Schema ID of `
		// https://conneroh.com/schemas` when defined with a struct called
		// `User{}`, will result in a schema with an ID set to
		// `https://conneroh.com/schemas/user`.
		//
		// If no `BaseSchemaID` is provided, we'll take the type's complete
		// package path and use that as a base instead. Set `Anonymous` to try
		// if you do not want to include a schema ID.
		BaseSchemaID schemaID
		// Anonymous when true will hide the auto-generated Schema ID and
		// provide what is known as an "anonymous schema". As a rule, this is
		// not recommended.
		Anonymous bool
		// AssignAnchor when true will use the original struct's name as an
		// anchor inside every definition, including the root schema. These can
		// be useful for having a reference to the original struct's name in
		// CamelCase instead of the snake-case used
		// by default for URI compatibility.
		//
		// Anchors do not appear to be widely used out in the wild, so at this
		// time the anchors themselves will not be used inside generated schema.
		AssignAnchor bool
		// AllowAdditionalProperties will cause the Reflector to generate a
		// schema without additionalProperties set to 'false' for all struct
		// types. This means the presence of additional keys in JSON objects
		// will not cause validation to fail. Note said additional keys will
		// simply be dropped when the validated JSON is unmarshaled.
		AllowAdditionalProperties bool
		// RequiredFromJSONSchemaTags will cause the Reflector to generate a
		// schema that requires any key tagged with `jsonschema:required`,
		// overriding the default of requiring any key *not* tagged with
		// `json:,omitempty`.
		RequiredFromJSONSchemaTags bool
		// Do not reference definitions. This will remove the top-level $defs
		// map and instead cause the entire structure of types to be output in
		// one tree. The list of type definitions (`$defs`) will not be
		// included.
		DoNotReference bool
		// ExpandedStruct when true will include the reflected type's definition
		// in the root as opposed to a definition with a reference.
		ExpandedStruct bool
		// FieldNameTag will change the tag used to get field names. json tags
		// are used by default.
		FieldNameTag string
		// IgnoredTypes defines a slice of types that should be ignored in the
		// schema, switching to just allowing additional properties instead.
		IgnoredTypes []any
		// Lookup allows a function to be defined that will provide a custom
		// mapping of types to Schema IDs. This allows existing schema documents
		// to be referenced by their ID instead of being embedded into the
		// current schema definitions. Reflected types will never be pointers,
		// only underlying elements.
		Lookup func(reflect.Type) schemaID
		// Mapper is a function that can be used to map custom Go types to
		// jsonschema schemas.
		Mapper func(reflect.Type) *Schema
		// Namer allows customizing of type names. The default is to use the
		// type's name provided by the reflect package.
		Namer func(reflect.Type) string
		// KeyNamer allows customizing of key names.
		// The default is to use the key's name as is, or the json tag if
		// present.
		//
		// If a json tag is present, KeyNamer will receive the tag's name as an
		// argument, not the original key name.
		KeyNamer func(string) string
		// AdditionalFields allows adding structfields for a given type
		AdditionalFields func(reflect.Type) []reflect.StructField
		// CommentMap is a dictionary of fully qualified go types and fields to
		// comment strings that will be used if a description has not already
		// been provided in the tags. Types and fields are added to the package
		// path using "." as a separator.
		//
		// Type descriptions should be defined like:
		//
		//   map[string]string{"github.com/conneroisu/groq.Reflector": "A Reflector reflects values into a Schema."}
		//
		// And Fields defined as:
		//
		//   map[string]string{"github.com/conneroisu/groq.Reflector.DoNotReference": "Do not reference definitions."}
		//
		// See also: AddGoComments
		CommentMap map[string]string
	}
)

// Reflect reflects to Schema from a value.
func (r *reflector) Reflect(v any) *Schema {
	return r.ReflectFromType(reflect.TypeOf(v))
}

// ReflectFromType generates root schema
func (r *reflector) ReflectFromType(t reflect.Type) *Schema {
	if t.Kind() == reflect.Ptr {
		t = t.Elem() // re-assign from pointer
	}
	name := r.typeName(t)
	s := new(Schema)
	definitions := schemaDefinitions{}
	s.Definitions = definitions
	bs := r.reflectTypeToSchemaWithID(definitions, t)
	if r.ExpandedStruct {
		*s = *definitions[name]
		delete(definitions, name)
	} else {
		*s = *bs
	}
	// Attempt to set the schema ID
	if !r.Anonymous && s.ID == EmptyID {
		baseSchemaID := r.BaseSchemaID
		if baseSchemaID == EmptyID {
			i := schemaID("https://" + t.PkgPath())
			if err := i.Validate(); err == nil {
				// it's okay to silently ignore URL errors
				baseSchemaID = i
			}
		}
		if baseSchemaID != EmptyID {
			s.ID = baseSchemaID.Add(ToSnakeCase(name))
		}
	}
	s.Version = version
	if !r.DoNotReference {
		s.Definitions = definitions
	}
	return s
}

// SetBaseSchemaID is a helper use to be able to set the reflectors base
// schema ID from a string as opposed to then ID instance.
func (r *reflector) SetBaseSchemaID(identifier string) {
	r.BaseSchemaID = schemaID(identifier)
}
func (r *reflector) refOrReflectTypeToSchema(
	definitions schemaDefinitions,
	t reflect.Type,
) *Schema {
	id := r.lookupID(t)
	if id != EmptyID {
		return &Schema{
			Ref: string(id),
		}
	}
	// Already added to definitions?
	if def := r.refDefinition(definitions, t); def != nil {
		return def
	}
	return r.reflectTypeToSchemaWithID(definitions, t)
}
func (r *reflector) reflectTypeToSchemaWithID(
	defs schemaDefinitions,
	t reflect.Type,
) *Schema {
	s := r.reflectTypeToSchema(defs, t)
	if s != nil {
		if r.Lookup != nil {
			identifier := r.Lookup(t)
			if identifier != EmptyID {
				s.ID = identifier
			}
		}
	}
	return s
}
func (r *reflector) reflectTypeToSchema(
	definitions schemaDefinitions,
	t reflect.Type,
) *Schema {
	// only try to reflect non-pointers
	if t.Kind() == reflect.Ptr {
		return r.refOrReflectTypeToSchema(definitions, t.Elem())
	}
	// Check if the there is an alias method that provides an object
	// that we should use instead of this one.
	if t.Implements(customAliasSchema) {
		v := reflect.New(t)
		o := v.Interface().(aliasSchemaImpl)
		t = reflect.TypeOf(o.JSONSchemaAlias())
		return r.refOrReflectTypeToSchema(definitions, t)
	}
	// Do any pre-definitions exist?
	if r.Mapper != nil {
		if t := r.Mapper(t); t != nil {
			return t
		}
	}
	if rt := r.reflectCustomSchema(definitions, t); rt != nil {
		return rt
	}
	// Prepare a base to which details can be added
	st := new(Schema)
	// jsonpb will marshal protobuf enum options as either strings or integers.
	// It will unmarshal either.
	if t.Implements(protoEnumType) {
		st.OneOf = []*Schema{
			{Type: "string"},
			{Type: "integer"},
		}
		return st
	}
	// Defined format types for JSON Schema Validation
	// RFC draft-wright-json-schema-validation-00, section 7.3
	// TODO email RFC section 7.3.2, hostname RFC section 7.3.3, uriref RFC section 7.3.7
	if t == ipType {
		// TODO differentiate ipv4 and ipv6 RFC section 7.3.4, 7.3.5
		st.Type = "string"
		st.Format = "ipv4"
		return st
	}
	switch t.Kind() {
	case reflect.Struct:
		r.reflectStruct(definitions, t, st)
	case reflect.Slice, reflect.Array:
		r.reflectSliceOrArray(definitions, t, st)
	case reflect.Map:
		r.reflectMap(definitions, t, st)
	case reflect.Interface:
		// empty
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
		st.Type = "integer"
	case reflect.Float32, reflect.Float64:
		st.Type = "number"
	case reflect.Bool:
		st.Type = "boolean"
	case reflect.String:
		st.Type = "string"
	default:
		panic("unsupported type " + t.String())
	}
	r.reflectSchemaExtend(definitions, t, st)
	// Always try to reference the definition which may have just been
	// created
	if def := r.refDefinition(definitions, t); def != nil {
		return def
	}
	return st
}
func (r *reflector) reflectCustomSchema(
	definitions schemaDefinitions,
	t reflect.Type,
) *Schema {
	if t.Kind() == reflect.Ptr {
		return r.reflectCustomSchema(definitions, t.Elem())
	}
	if t.Implements(customType) {
		o := reflect.New(t).Interface().(customSchemaImpl)
		st := o.JSONSchema()
		r.addDefinition(definitions, t, st)
		if ref := r.refDefinition(definitions, t); ref != nil {
			return ref
		}
		return st
	}
	return nil
}
func (r *reflector) reflectSchemaExtend(
	definitions schemaDefinitions,
	t reflect.Type,
	s *Schema,
) *Schema {
	if t.Implements(extendType) {
		o := reflect.New(t).Interface().(extendSchemaImpl)
		o.JSONSchemaExtend(s)
		if ref := r.refDefinition(definitions, t); ref != nil {
			return ref
		}
	}
	return s
}
func (r *reflector) reflectSliceOrArray(
	definitions schemaDefinitions,
	t reflect.Type,
	st *Schema,
) {
	if t == rawMessageType {
		return
	}
	r.addDefinition(definitions, t, st)
	if st.Description == "" {
		st.Description = r.lookupComment(t, "")
	}
	if t.Kind() == reflect.Array {
		l := uint64(t.Len())
		st.MinItems = &l
		st.MaxItems = &l
	}
	if t.Kind() == reflect.Slice && t.Elem() == byteSliceType.Elem() {
		st.Type = "string"
		st.ContentEncoding = "base64"
		return
	}
	st.Type = "array"
	st.Items = r.refOrReflectTypeToSchema(definitions, t.Elem())
}
func (r *reflector) reflectMap(
	definitions schemaDefinitions,
	t reflect.Type,
	st *Schema,
) {
	r.addDefinition(definitions, t, st)
	st.Type = "object"
	if st.Description == "" {
		st.Description = r.lookupComment(t, "")
	}
	switch t.Key().Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		st.PatternProperties = map[string]*Schema{
			"^[0-9]+$": r.refOrReflectTypeToSchema(definitions, t.Elem()),
		}
		st.AdditionalProperties = falseSchema
		return
	}
	if t.Elem().Kind() != reflect.Interface {
		st.AdditionalProperties = r.refOrReflectTypeToSchema(
			definitions,
			t.Elem(),
		)
	}
}

// Reflects a struct to a JSON Schema type.
func (r *reflector) reflectStruct(
	definitions schemaDefinitions,
	t reflect.Type,
	s *Schema,
) {
	// Handle special types
	switch t {
	case timeType: // date-time RFC section 7.3.1
		s.Type = "string"
		s.Format = "date-time"
		return
	case uriType: // uri RFC section 7.3.6
		s.Type = "string"
		s.Format = "uri"
		return
	}
	r.addDefinition(definitions, t, s)
	s.Type = "object"
	s.Properties = newProperties()
	s.Description = r.lookupComment(t, "")
	if r.AssignAnchor {
		s.Anchor = t.Name()
	}
	if !r.AllowAdditionalProperties && s.AdditionalProperties == nil {
		s.AdditionalProperties = falseSchema
	}
	ignored := false
	for _, it := range r.IgnoredTypes {
		if reflect.TypeOf(it) == t {
			ignored = true
			break
		}
	}
	if !ignored {
		r.reflectStructFields(s, definitions, t)
	}
}

func (r *reflector) reflectStructFields(
	st *Schema,
	definitions schemaDefinitions,
	t reflect.Type,
) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}
	var getFieldDocString customGetFieldDocString
	if t.Implements(customStructGetFieldDocString) {
		v := reflect.New(t)
		o := v.Interface().(customSchemaGetFieldDocString)
		getFieldDocString = o.GetFieldDocString
	}
	customPropertyMethod := func(string) any {
		return nil
	}
	if t.Implements(customPropertyAliasSchema) {
		v := reflect.New(t)
		o := v.Interface().(propertyAliasSchemaImpl)
		customPropertyMethod = o.JSONSchemaProperty
	}
	handleField := func(f reflect.StructField) {
		name, shouldEmbed, required, nullable := r.reflectFieldName(f)
		// if anonymous and exported type should be processed
		// recursively current type should inherit properties of
		// anonymous one
		if name == "" {
			if shouldEmbed {
				r.reflectStructFields(st, definitions, f.Type)
			}
			return
		}
		// If a JSONSchemaAlias(prop string) method is defined, attempt
		// to use the provided object's type instead of the field's
		// type.
		var property *Schema
		if alias := customPropertyMethod(name); alias != nil {
			property = r.refOrReflectTypeToSchema(
				definitions,
				reflect.TypeOf(alias),
			)
		} else {
			property = r.refOrReflectTypeToSchema(definitions, f.Type)
		}
		property.fieldsFromTags(f, st, name)
		if property.Description == "" {
			property.Description = r.lookupComment(t, f.Name)
		}
		if getFieldDocString != nil {
			property.Description = getFieldDocString(f.Name)
		}
		if nullable {
			property = &Schema{
				OneOf: []*Schema{
					property,
					{
						Type: "null",
					},
				},
			}
		}
		st.Properties.Set(name, property)
		if required {
			st.Required = appendUniqueString(st.Required, name)
		}
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		handleField(f)
	}
	if r.AdditionalFields != nil {
		if af := r.AdditionalFields(t); af != nil {
			for _, sf := range af {
				handleField(sf)
			}
		}
	}
}

func (r *reflector) lookupComment(t reflect.Type, name string) string {
	if r.CommentMap == nil {
		return ""
	}
	n := fullyQualifiedTypeName(t)
	if name != "" {
		n = n + "." + name
	}
	return r.CommentMap[n]
}

// addDefinition will append the provided schema. If needed, an ID and anchor
// will also be added.
func (r *reflector) addDefinition(
	definitions schemaDefinitions,
	t reflect.Type,
	s *Schema,
) {
	name := r.typeName(t)
	if name == "" {
		return
	}
	definitions[name] = s
}

// refDefinition will provide a schema with a reference to an existing
// definition.
func (r *reflector) refDefinition(
	definitions schemaDefinitions,
	t reflect.Type,
) *Schema {
	if r.DoNotReference {
		return nil
	}
	name := r.typeName(t)
	if name == "" {
		return nil
	}
	if _, ok := definitions[name]; !ok {
		return nil
	}
	return &Schema{
		Ref: "#/$defs/" + name,
	}
}
func (r *reflector) lookupID(t reflect.Type) schemaID {
	if r.Lookup != nil {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		return r.Lookup(t)
	}
	return EmptyID
}
