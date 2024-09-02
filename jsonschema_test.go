package groq_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/conneroisu/groq-go"
)

func TestDefinition_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		def  groq.Definition
		want string
	}{
		{
			name: "Test with empty Definition",
			def:  groq.Definition{},
			want: `{"properties":{}}`,
		},
		{
			name: "Test with Definition properties set",
			def: groq.Definition{
				Type:        groq.String,
				Description: "A string type",
				Properties: map[string]groq.Definition{
					"name": {
						Type: groq.String,
					},
				},
			},
			want: `{
   "type":"string",
   "description":"A string type",
   "properties":{
      "name":{
         "type":"string",
         "properties":{}
      }
   }
}`,
		},
		{
			name: "Test with nested Definition properties",
			def: groq.Definition{
				Type: groq.Object,
				Properties: map[string]groq.Definition{
					"user": {
						Type: groq.Object,
						Properties: map[string]groq.Definition{
							"name": {
								Type: groq.String,
							},
							"age": {
								Type: groq.Integer,
							},
						},
					},
				},
			},
			want: `{
   "type":"object",
   "properties":{
      "user":{
         "type":"object",
         "properties":{
            "name":{
               "type":"string",
               "properties":{}
            },
            "age":{
               "type":"integer",
               "properties":{}
            }
         }
      }
   }
}`,
		},
		{
			name: "Test with complex nested Definition",
			def: groq.Definition{
				Type: groq.Object,
				Properties: map[string]groq.Definition{
					"user": {
						Type: groq.Object,
						Properties: map[string]groq.Definition{
							"name": {
								Type: groq.String,
							},
							"age": {
								Type: groq.Integer,
							},
							"address": {
								Type: groq.Object,
								Properties: map[string]groq.Definition{
									"city": {
										Type: groq.String,
									},
									"country": {
										Type: groq.String,
									},
								},
							},
						},
					},
				},
			},
			want: `{
   "type":"object",
   "properties":{
      "user":{
         "type":"object",
         "properties":{
            "name":{
               "type":"string",
               "properties":{}
            },
            "age":{
               "type":"integer",
               "properties":{}
            },
            "address":{
               "type":"object",
               "properties":{
                  "city":{
                     "type":"string",
                     "properties":{}
                  },
                  "country":{
                     "type":"string",
                     "properties":{}
                  }
               }
            }
         }
      }
   }
}`,
		},
		{
			name: "Test with Array type Definition",
			def: groq.Definition{
				Type: groq.Array,
				Items: &groq.Definition{
					Type: groq.String,
				},
				Properties: map[string]groq.Definition{
					"name": {
						Type: groq.String,
					},
				},
			},
			want: `{
   "type":"array",
   "items":{
      "type":"string",
      "properties":{
         
      }
   },
   "properties":{
      "name":{
         "type":"string",
         "properties":{}
      }
   }
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantBytes := []byte(tt.want)
			var want map[string]interface{}
			err := json.Unmarshal(wantBytes, &want)
			if err != nil {
				t.Errorf("Failed to Unmarshal JSON: error = %v", err)
				return
			}

			got := structToMap(t, tt.def)
			gotPtr := structToMap(t, &tt.def)

			if !reflect.DeepEqual(got, want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, want)
			}
			if !reflect.DeepEqual(gotPtr, want) {
				t.Errorf("MarshalJSON() gotPtr = %v, want %v", gotPtr, want)
			}
		})
	}
}

func structToMap(t *testing.T, v any) map[string]any {
	t.Helper()
	gotBytes, err := json.Marshal(v)
	if err != nil {
		t.Errorf("Failed to Marshal JSON: error = %v", err)
		return nil
	}

	var got map[string]interface{}
	err = json.Unmarshal(gotBytes, &got)
	if err != nil {
		t.Errorf("Failed to Unmarshal JSON: error =  %v", err)
		return nil
	}
	return got
}
