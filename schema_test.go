//go:build !test
// +build !test

package groq

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestID(t *testing.T) {
	base := "https://github.com/conneroisu/groq-go/schema"
	id := schemaID(base)

	assert.Equal(t, base, id.String())

	id = id.Add("user")
	assert.EqualValues(t, base+"/user", id)

	id = id.Anchor("Name")
	assert.EqualValues(t, base+"/user#Name", id)

	id = id.Anchor("Title")
	assert.EqualValues(t, base+"/user#Title", id)

	id = id.Def("Name")
	assert.EqualValues(t, base+"/user#/$defs/Name", id)
}

func TestIDValidation(t *testing.T) {
	id := schemaID("https://invopop.com/schema/user")
	assert.NoError(t, id.Validate())

	id = "https://encoding/json"
	if assert.Error(t, id.Validate()) {
		assert.Contains(
			t,
			id.Validate().Error(),
			"hostname does not look valid",
		)
	}

	id = "time"
	if assert.Error(t, id.Validate()) {
		assert.Contains(t, id.Validate().Error(), "hostname")
	}

	id = "http://invopop.com"
	if assert.Error(t, id.Validate()) {
		assert.Contains(t, id.Validate().Error(), "path")
	}

	id = "foor://invopop.com/schema/user"
	if assert.Error(t, id.Validate()) {
		assert.Contains(t, id.Validate().Error(), "schema")
	}

	id = "invopop.com\n/test"
	if assert.Error(t, id.Validate()) {
		assert.Contains(t, id.Validate().Error(), "invalid URL")
	}
}

// Pet defines the user's fury friend.
type Pet struct {
	// Name of the animal.
	Name string `json:"name" jsonschema:"title=Name"`
}

// Pets is a collection of Pet objects.
type Pets []*Pet

// NamedPets is a map of animal names to pets.
type NamedPets map[string]*Pet

type (
	// Plant represents the plants the user might have and serves as a test
	// of structs inside a `type` set.
	Plant struct {
		Variant string `json:"variant"                 jsonschema:"title=Variant"` // This comment will be used
		// Multicellular is true if the plant is multicellular
		Multicellular bool `json:"multicellular,omitempty" jsonschema:"title=Multicellular"` // This comment will be ignored
	}
)

// User is used as a base to provide tests for comments.
// Don't forget to checkout the nested path.
type User struct {
	// Unique sequential identifier.
	ID      int            `json:"id"                jsonschema:"required"`
	Name    string         `json:"name"              jsonschema:"required,minLength=1,maxLength=20,pattern=.*,description=this is a property,title=the name,example=joe,example=lucy,default=alex"`
	Friends []int          `json:"friends,omitempty"                                                                                                                                               jsonschema_description:"list of IDs, omitted when empty"`
	Tags    map[string]any `json:"tags,omitempty"`

	// An array of pets the user cares for.
	Pets Pets `json:"pets"`

	// Set of animal names to pets
	NamedPets NamedPets `json:"named_pets"`

	// Set of plants that the user likes
	Plants []*Plant `json:"plants" jsonschema:"title=Plants"`
}

var updateFixtures = flag.Bool("update", false, "set to update fixtures")

var compareFixtures = flag.Bool(
	"compare",
	false,
	"output failed fixtures with .out.json",
)

type GrandfatherType struct {
	FamilyName string `json:"family_name" jsonschema:"required"`
}

type SomeBaseType struct {
	ID               string `json:"id"` // to test composition override
	SomeBaseProperty int    `json:"some_base_property"`
	// The jsonschema required tag is nonsensical for private and ignored properties.
	// Their presence here tests that the fields *will not* be required in the output
	// schema, even if they are tagged required.
	SomeIgnoredBaseProperty   string          `json:"-"                  jsonschema:"required"`
	SomeSchemaIgnoredProperty string          `                          jsonschema:"-,required"`
	Grandfather               GrandfatherType `json:"grand"`

	SomeUntaggedBaseProperty bool `jsonschema:"required"`
}

type MapType map[string]any

type ArrayType []string

type nonExported struct {
	PublicNonExported int
}

type ProtoEnum int32

func (ProtoEnum) EnumDescriptor() ([]byte, []int) { return []byte(nil), []int{0} }

const (
	Unset ProtoEnum = iota
	Great
)

type TestUser struct {
	SomeBaseType
	nonExported
	MapType

	ID       int               `json:"id"                jsonschema:"required,minimum=bad,maximum=bad,exclusiveMinimum=bad,exclusiveMaximum=bad,default=bad"`
	Name     string            `json:"name"              jsonschema:"required,minLength=1,maxLength=20,pattern=.*,description=this is a property,title=the name,example=joe,example=lucy,default=alex,readOnly=true"`
	Password string            `json:"password"          jsonschema:"writeOnly=true"`
	Friends  []int             `json:"friends,omitempty"                                                                                                                                                             jsonschema_description:"list of IDs, omitted when empty"`
	Tags     map[string]string `json:"tags,omitempty"`
	Options  map[string]any    `json:"options,omitempty"`

	TestFlag       bool
	TestFlagFalse  bool `json:",omitempty" jsonschema:"default=false"`
	TestFlagTrue   bool `json:",omitempty" jsonschema:"default=true"`
	IgnoredCounter int  `json:"-"`

	// Tests for RFC draft-wright-json-schema-validation-00, section 7.3
	BirthDate time.Time `json:"birth_date,omitempty"`
	Website   url.URL   `json:"website,omitempty"`
	IPAddress net.IP    `json:"network_address,omitempty"`

	// Tests for RFC draft-wright-json-schema-hyperschema-00, section 4
	Photo  []byte `json:"photo,omitempty"  jsonschema:"required"`
	Photo2 Bytes  `json:"photo2,omitempty" jsonschema:"required"`

	// Tests for jsonpb enum support
	Feeling ProtoEnum `json:"feeling,omitempty"`

	Age   int    `json:"age"   jsonschema:"minimum=18,maximum=120,exclusiveMaximum=121,exclusiveMinimum=17"`
	Email string `json:"email" jsonschema:"format=email"`
	UUID  string `json:"uuid"  jsonschema:"format=uuid"`

	// Test for "extras" support
	Baz       string `jsonschema_extras:"foo=bar,hello=world,foo=bar1"`
	BoolExtra string `jsonschema_extras:"isTrue=true,isFalse=false"    json:"bool_extra,omitempty"`

	// Tests for simple enum tags
	Color      string  `json:"color"          jsonschema:"enum=red,enum=green,enum=blue"`
	Rank       int     `json:"rank,omitempty" jsonschema:"enum=1,enum=2,enum=3"`
	Multiplier float64 `json:"mult,omitempty" jsonschema:"enum=1.0,enum=1.5,enum=2.0"`

	// Tests for enum tags on slices
	Roles      []string  `json:"roles"                jsonschema:"enum=admin,enum=moderator,enum=user"`
	Priorities []int     `json:"priorities,omitempty" jsonschema:"enum=-1,enum=0,enum=1,enun=2"`
	Offsets    []float64 `json:"offsets,omitempty"    jsonschema:"enum=1.570796,enum=3.141592,enum=6.283185"`

	// Test for raw JSON
	Anything any             `json:"anything,omitempty"`
	Raw      json.RawMessage `json:"raw"`
}

type CustomTime time.Time

type CustomTypeField struct {
	CreatedAt CustomTime
}

type CustomTimeWithInterface time.Time

type CustomTypeFieldWithInterface struct {
	CreatedAt CustomTimeWithInterface
}

func (CustomTimeWithInterface) JSONSchema() *schema {
	return &schema{
		Type:   "string",
		Format: "date-time",
	}
}

type RootOneOf struct {
	Field1 string     `json:"field1" jsonschema:"oneof_required=group1"`
	Field2 string     `json:"field2" jsonschema:"oneof_required=group2"`
	Field3 any        `json:"field3" jsonschema:"oneof_type=string;array"`
	Field4 string     `json:"field4" jsonschema:"oneof_required=group1"`
	Field5 ChildOneOf `json:"child"`
	Field6 any        `json:"field6" jsonschema:"oneof_ref=Outer;OuterNamed;OuterPtr"`
}

type ChildOneOf struct {
	Child1 string `json:"child1" jsonschema:"oneof_required=group1"`
	Child2 string `json:"child2" jsonschema:"oneof_required=group2"`
	Child3 any    `json:"child3" jsonschema:"oneof_required=group2,oneof_type=string;array"`
	Child4 string `json:"child4" jsonschema:"oneof_required=group1"`
}

type RootAnyOf struct {
	Field1 string     `json:"field1" jsonschema:"anyof_required=group1"`
	Field2 string     `json:"field2" jsonschema:"anyof_required=group2"`
	Field3 any        `json:"field3" jsonschema:"anyof_type=string;array"`
	Field4 string     `json:"field4" jsonschema:"anyof_required=group1"`
	Field5 ChildAnyOf `json:"child"`
}

type ChildAnyOf struct {
	Child1 string `json:"child1" jsonschema:"anyof_required=group1"`
	Child2 string `json:"child2" jsonschema:"anyof_required=group2"`
	Child3 any    `json:"child3" jsonschema:"anyof_required=group2,oneof_type=string;array"`
	Child4 string `json:"child4" jsonschema:"anyof_required=group1"`
}

type Text string

type TextNamed string

type Outer struct {
	TextNamed
	Text `json:",omitempty"`
	Inner
}

type OuterNamed struct {
	Text  `json:"text,omitempty"`
	Inner `json:"inner"`
}

type OuterInlined struct {
	Text  `json:"text,omitempty"`
	Inner `json:",inline"`
}

type OuterPtr struct {
	*Inner
	Text `json:",omitempty"`
}

type Inner struct {
	Foo string `yaml:"foo"`
}

type MinValue struct {
	Value int `json:"value4" jsonschema_extras:"minimum=0"`
}
type Bytes []byte

type TestNullable struct {
	Child1 string `json:"child1" jsonschema:"nullable"`
}

type CompactDate struct {
	Year  int
	Month int
}

type UserWithAnchor struct {
	Name string `json:"name" jsonschema:"anchor=Name"`
}

func (CompactDate) JSONSchema() *schema {
	return &schema{
		Type:        "string",
		Title:       "Compact Date",
		Description: "Short date that only includes year and month",
		Pattern:     "^[0-9]{4}-[0-1][0-9]$",
	}
}

type TestDescriptionOverride struct {
	FirstName  string `json:"FirstName"`
	LastName   string `json:"LastName"`
	Age        uint   `json:"age"`
	MiddleName string `json:"middle_name,omitempty"`
}

func (TestDescriptionOverride) GetFieldDocString(fieldName string) string {
	switch fieldName {
	case "FirstName":
		return "test2"
	case "LastName":
		return "test3"
	case "Age":
		return "test4"
	case "MiddleName":
		return "test5"
	default:
		return ""
	}
}

type LookupName struct {
	Given   string `json:"first"`
	Surname string `json:"surname"`
}

type LookupUser struct {
	Name  *LookupName `json:"name"`
	Alias string      `json:"alias,omitempty"`
}

type CustomSliceOuter struct {
	Slice CustomSliceType `json:"slice"`
}

type CustomSliceType []string

func (CustomSliceType) JSONSchema() *schema {
	return &schema{
		OneOf: []*schema{{
			Type: "string",
		}, {
			Type: "array",
			Items: &schema{
				Type: "string",
			},
		}},
	}
}

type CustomMapType map[string]string

func (CustomMapType) JSONSchema() *schema {
	properties := newProperties()
	properties.Set("key", &schema{
		Type: "string",
	})
	properties.Set("value", &schema{
		Type: "string",
	})
	return &schema{
		Type: "array",
		Items: &schema{
			Type:       "object",
			Properties: properties,
			Required:   []string{"key", "value"},
		},
	}
}

type CustomMapOuter struct {
	MyMap CustomMapType `json:"my_map"`
}

type PatternTest struct {
	WithPattern string `json:"with_pattern" jsonschema:"minLength=1,pattern=[0-9]{1\\,4},maxLength=50"`
}

type RecursiveExample struct {
	Text  string              `json:"text"`
	Child []*RecursiveExample `json:"children,omitempty"`
}

type KeyNamedNested struct {
	NestedNotRenamedProperty string
	NotRenamed               string
}

type KeyNamed struct {
	ThisWasLeftAsIs      string
	NotComingFromJSON    bool           `json:"coming_from_json_tag_not_renamed"`
	NestedNotRenamed     KeyNamedNested `json:"nested_not_renamed"`
	UnicodeShenanigans   string
	RenamedByComputation int `                                        jsonschema_description:"Description was preserved"`
}

type SchemaExtendTestBase struct {
	FirstName  string `json:"FirstName"`
	LastName   string `json:"LastName"`
	Age        uint   `json:"age"`
	MiddleName string `json:"middle_name,omitempty"`
}

type SchemaExtendTest struct {
	SchemaExtendTestBase `json:",inline"`
}

func (SchemaExtendTest) JSONSchemaExtend(base *schema) {
	base.Properties.Delete("FirstName")
	base.Properties.Delete("age")
	val, _ := base.Properties.Get("LastName")
	val.Description = "some extra words"
	base.Required = []string{"LastName"}
}

type Expression struct {
	Value int `json:"value" jsonschema_extras:"foo=bar=='baz'"`
}

type PatternEqualsTest struct {
	WithEquals          string `jsonschema:"pattern=foo=bar"`
	WithEqualsAndCommas string `jsonschema:"pattern=foo,=bar"`
}

func TestReflector(t *testing.T) {
	r := new(reflector)
	s := "http://example.com/schema"
	r.SetBaseSchemaID(s)
	assert.EqualValues(t, s, r.BaseSchemaID)
}

func TestReflectFromType(t *testing.T) {
	r := new(reflector)
	tu := new(TestUser)
	typ := reflect.TypeOf(tu)

	s := r.ReflectFromType(typ)
	assert.EqualValues(
		t,
		"https://github.com/conneroisu/groq-go/test-user",
		s.ID,
	)

	x := struct {
		Test string
	}{
		Test: "foo",
	}
	typ = reflect.TypeOf(x)
	s = r.Reflect(typ)
	assert.Empty(t, s.ID)
}

func TestSchemaGeneration(t *testing.T) {
	tests := []struct {
		typ       any
		reflector *reflector
		fixture   string
	}{
		{&TestUser{}, &reflector{}, "testdata/test_user.json"},
		{&UserWithAnchor{}, &reflector{}, "testdata/user_with_anchor.json"},
		{
			&TestUser{},
			&reflector{AssignAnchor: true},
			"testdata/test_user_assign_anchor.json",
		},
		{
			&TestUser{},
			&reflector{AllowAdditionalProperties: true},
			"testdata/allow_additional_props.json",
		},
		{
			&TestUser{},
			&reflector{RequiredFromJSONSchemaTags: true},
			"testdata/required_from_jsontags.json",
		},
		{
			&TestUser{},
			&reflector{ExpandedStruct: true},
			"testdata/defaults_expanded_toplevel.json",
		},
		{
			&TestUser{},
			&reflector{IgnoredTypes: []any{GrandfatherType{}}},
			"testdata/ignore_type.json",
		},
		{
			&TestUser{},
			&reflector{DoNotReference: true},
			"testdata/no_reference.json",
		},
		{
			&TestUser{},
			&reflector{DoNotReference: true, AssignAnchor: true},
			"testdata/no_reference_anchor.json",
		},
		{
			&RootOneOf{},
			&reflector{RequiredFromJSONSchemaTags: true},
			"testdata/oneof.json",
		},
		{
			&RootAnyOf{},
			&reflector{RequiredFromJSONSchemaTags: true},
			"testdata/anyof.json",
		},
		{&CustomTypeField{}, &reflector{
			Mapper: func(i reflect.Type) *schema {
				if i == reflect.TypeOf(CustomTime{}) {
					return &schema{
						Type:   "string",
						Format: "date-time",
					}
				}
				return nil
			},
		}, "testdata/custom_type.json"},
		{
			LookupUser{},
			&reflector{BaseSchemaID: "https://example.com/schemas"},
			"testdata/base_schema_id.json",
		},
		{LookupUser{}, &reflector{
			Lookup: func(i reflect.Type) schemaID {
				switch i {
				case reflect.TypeOf(LookupUser{}):
					return schemaID("https://example.com/schemas/lookup-user")
				case reflect.TypeOf(LookupName{}):
					return schemaID("https://example.com/schemas/lookup-name")
				}
				return EmptyID
			},
		}, "testdata/lookup.json"},
		{&LookupUser{}, &reflector{
			BaseSchemaID:   "https://example.com/schemas",
			ExpandedStruct: true,
			AssignAnchor:   true,
			Lookup: func(i reflect.Type) schemaID {
				switch i {
				case reflect.TypeOf(LookupUser{}):
					return schemaID("https://example.com/schemas/lookup-user")
				case reflect.TypeOf(LookupName{}):
					return schemaID("https://example.com/schemas/lookup-name")
				}
				return EmptyID
			},
		}, "testdata/lookup_expanded.json"},
		{
			&Outer{},
			&reflector{ExpandedStruct: true},
			"testdata/inlining_inheritance.json",
		},
		{
			&OuterNamed{},
			&reflector{ExpandedStruct: true},
			"testdata/inlining_embedded.json",
		},
		{
			&OuterNamed{},
			&reflector{ExpandedStruct: true, AssignAnchor: true},
			"testdata/inlining_embedded_anchored.json",
		},
		{
			&OuterInlined{},
			&reflector{ExpandedStruct: true},
			"testdata/inlining_tag.json",
		},
		{
			&OuterPtr{},
			&reflector{ExpandedStruct: true},
			"testdata/inlining_ptr.json",
		},
		{&MinValue{}, &reflector{}, "testdata/schema_with_minimum.json"},
		{&TestNullable{}, &reflector{}, "testdata/nullable.json"},
		{&GrandfatherType{}, &reflector{
			AdditionalFields: func(_ reflect.Type) []reflect.StructField {
				return []reflect.StructField{
					{
						Name:      "Addr",
						Type:      reflect.TypeOf((*net.IP)(nil)).Elem(),
						Tag:       "json:\"ip_addr\"",
						Anonymous: false,
					},
				}
			},
		}, "testdata/custom_additional.json"},
		{
			&TestDescriptionOverride{},
			&reflector{},
			"testdata/test_description_override.json",
		},
		{&CompactDate{}, &reflector{}, "testdata/compact_date.json"},
		{&CustomSliceOuter{}, &reflector{}, "testdata/custom_slice_type.json"},
		{&CustomMapOuter{}, &reflector{}, "testdata/custom_map_type.json"},
		{
			&CustomTypeFieldWithInterface{},
			&reflector{},
			"testdata/custom_type_with_interface.json",
		},
		{&RecursiveExample{}, &reflector{}, "testdata/recursive.json"},
		{&KeyNamed{}, &reflector{
			KeyNamer: func(s string) string {
				switch s {
				case "ThisWasLeftAsIs":
					fallthrough
				case "NotRenamed":
					fallthrough
				case "nested_not_renamed":
					return s
				case "coming_from_json_tag_not_renamed":
					return "coming_from_json_tag"
				case "NestedNotRenamed":
					return "nested-renamed"
				case "NestedNotRenamedProperty":
					return "nested-renamed-property"
				case "UnicodeShenanigans":
					return "✨unicode✨  s̸̥͝h̷̳͒e̴̜̽n̸̡̿a̷̘̔n̷̘͐i̶̫̐ǵ̶̯a̵̘͒n̷̮̾s̸̟̓"
				case "RenamedByComputation":
					return fmt.Sprintf("%.2f", float64(len(s))+1/137.0)
				}
				return "unknown case"
			},
		}, "testdata/keynamed.json"},
		{MapType{}, &reflector{}, "testdata/map_type.json"},
		{ArrayType{}, &reflector{}, "testdata/array_type.json"},
		{SchemaExtendTest{}, &reflector{}, "testdata/custom_type_extend.json"},
		{Expression{}, &reflector{}, "testdata/schema_with_expression.json"},
		{&PatternTest{}, &reflector{}, "testdata/commas_in_pattern.json"},
	}

	for _, tt := range tests {
		name := strings.TrimSuffix(filepath.Base(tt.fixture), ".json")
		t.Run(name, func(t *testing.T) {
			compareSchemaOutput(t,
				tt.fixture, tt.reflector, tt.typ,
			)
		})
	}
}

func TestBaselineUnmarshal(t *testing.T) {
	r := &reflector{}
	compareSchemaOutput(t, "testdata/test_user.json", r, &TestUser{})
}

func compareSchemaOutput(t *testing.T, f string, r *reflector, obj any) {
	t.Helper()
	expectedJSON, err := os.ReadFile(f)
	require.NoError(t, err)

	actualSchema := r.Reflect(obj)
	actualJSON, _ := json.MarshalIndent(
		actualSchema,
		"",
		"  ",
	) //nolint:errchkjson

	if *updateFixtures {
		_ = os.WriteFile(f, actualJSON, 0600)
	}

	if !assert.JSONEq(t, string(expectedJSON), string(actualJSON)) {
		if *compareFixtures {
			_ = os.WriteFile(
				strings.TrimSuffix(f, ".json")+".out.json",
				actualJSON,
				0600,
			)
		}
	}
}

func fixtureContains(t *testing.T, f, s string) {
	t.Helper()
	b, err := os.ReadFile(f)
	require.NoError(t, err)
	assert.Contains(t, string(b), s)
}

func TestSplitOnUnescapedCommas(t *testing.T) {
	tests := []struct {
		strToSplit string
		expected   []string
	}{
		{
			`Hello,this,is\,a\,string,haha`,
			[]string{`Hello`, `this`, `is,a,string`, `haha`},
		},
		{`hello,no\\,split`, []string{`hello`, `no\,split`}},
		{`string without commas`, []string{`string without commas`}},
		{`ünicode,𐂄,Ж\,П,ᠳ`, []string{`ünicode`, `𐂄`, `Ж,П`, `ᠳ`}},
		{`empty,,tag`, []string{`empty`, ``, `tag`}},
	}

	for _, test := range tests {
		actual := splitOnUnescapedCommas(test.strToSplit)
		require.Equal(t, test.expected, actual)
	}
}

func TestArrayExtraTags(t *testing.T) {
	type URIArray struct {
		TestURIs []string `jsonschema:"type=array,format=uri,pattern=^https://.*"`
	}

	r := new(reflector)
	schema := r.Reflect(&URIArray{})
	d := schema.Definitions["URIArray"]
	require.NotNil(t, d)
	props := d.Properties
	require.NotNil(t, props)
	p, found := props.Get("TestURIs")
	require.True(t, found)

	pt := p.Items.Format
	require.Equal(t, pt, "uri")
	pt = p.Items.Pattern
	require.Equal(t, pt, "^https://.*")
}

func TestFieldNameTag(t *testing.T) {
	type Config struct {
		Name  string `yaml:"name"`
		Count int    `yaml:"count"`
	}

	r := reflector{
		FieldNameTag: "yaml",
	}
	compareSchemaOutput(t, "testdata/test_config.json", &r, &Config{})
}

func TestFieldOneOfRef(t *testing.T) {
	type Server struct {
		IPAddress      any   `json:"ip_address,omitempty"       jsonschema:"oneof_ref=#/$defs/ipv4;#/$defs/ipv6"`
		IPAddresses    []any `json:"ip_addresses,omitempty"     jsonschema:"oneof_ref=#/$defs/ipv4;#/$defs/ipv6"`
		IPAddressAny   any   `json:"ip_address_any,omitempty"   jsonschema:"anyof_ref=#/$defs/ipv4;#/$defs/ipv6"`
		IPAddressesAny []any `json:"ip_addresses_any,omitempty" jsonschema:"anyof_ref=#/$defs/ipv4;#/$defs/ipv6"`
	}

	r := &reflector{}
	compareSchemaOutput(t, "testdata/oneof_ref.json", r, &Server{})
}

func TestNumberHandling(t *testing.T) {
	type NumberHandler struct {
		Int64   int64   `json:"int64"   jsonschema:"default=12"`
		Float32 float32 `json:"float32" jsonschema:"default=12.5"`
	}

	r := &reflector{}
	compareSchemaOutput(
		t,
		"testdata/number_handling.json",
		r,
		&NumberHandler{},
	)
	fixtureContains(t, "testdata/number_handling.json", `"default": 12`)
	fixtureContains(t, "testdata/number_handling.json", `"default": 12.5`)
}

func TestArrayHandling(t *testing.T) {
	type ArrayHandler struct {
		MinLen []string  `json:"min_len" jsonschema:"minLength=2,default=qwerty"`
		MinVal []float64 `json:"min_val" jsonschema:"minimum=2.5"`
	}

	r := &reflector{}
	compareSchemaOutput(t, "testdata/array_handling.json", r, &ArrayHandler{})
	fixtureContains(t, "testdata/array_handling.json", `"minLength": 2`)
	fixtureContains(t, "testdata/array_handling.json", `"minimum": 2.5`)
}

func TestUnsignedIntHandling(t *testing.T) {
	type UnsignedIntHandler struct {
		MinLen   []string `json:"min_len"   jsonschema:"minLength=0"`
		MaxLen   []string `json:"max_len"   jsonschema:"maxLength=0"`
		MinItems []string `json:"min_items" jsonschema:"minItems=0"`
		MaxItems []string `json:"max_items" jsonschema:"maxItems=0"`
	}

	r := &reflector{}
	compareSchemaOutput(
		t,
		"testdata/unsigned_int_handling.json",
		r,
		&UnsignedIntHandler{},
	)
	fixtureContains(t, "testdata/unsigned_int_handling.json", `"minLength": 0`)
	fixtureContains(t, "testdata/unsigned_int_handling.json", `"maxLength": 0`)
	fixtureContains(t, "testdata/unsigned_int_handling.json", `"minItems": 0`)
	fixtureContains(t, "testdata/unsigned_int_handling.json", `"maxItems": 0`)
}

func TestJSONSchemaFormat(t *testing.T) {
	type WithCustomFormat struct {
		Dates []string `json:"dates" jsonschema:"format=date"`
		Odds  []string `json:"odds"  jsonschema:"format=odd"`
	}

	r := &reflector{}
	compareSchemaOutput(
		t,
		"testdata/with_custom_format.json",
		r,
		&WithCustomFormat{},
	)
	fixtureContains(t, "testdata/with_custom_format.json", `"format": "date"`)
	fixtureContains(t, "testdata/with_custom_format.json", `"format": "odd"`)
}

type AliasObjectA struct {
	PropA string `json:"prop_a"`
}
type AliasObjectB struct {
	PropB string `json:"prop_b"`
}
type AliasObjectC struct {
	ObjB *AliasObjectB `json:"obj_b"`
}
type AliasPropertyObjectBase struct {
	Object any `json:"object"`
}

func (AliasPropertyObjectBase) JSONSchemaProperty(prop string) any {
	if prop == "object" {
		return &AliasObjectA{}
	}
	return nil
}

func (AliasObjectB) JSONSchemaAlias() any {
	return AliasObjectA{}
}

func TestJSONSchemaProperty(t *testing.T) {
	r := &reflector{}
	compareSchemaOutput(
		t,
		"testdata/schema_property_alias.json",
		r,
		&AliasPropertyObjectBase{},
	)
}

func TestJSONSchemaAlias(t *testing.T) {
	r := &reflector{}
	compareSchemaOutput(t, "testdata/schema_alias.json", r, &AliasObjectB{})
	compareSchemaOutput(t, "testdata/schema_alias_2.json", r, &AliasObjectC{})
}
