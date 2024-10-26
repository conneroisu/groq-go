package tools

const (
	ToolTypeFunction ToolType = "function" // ToolTypeFunction is the function tool type.
)

type (
	// Tool represents the tool.
	Tool struct {
		Type     ToolType           `json:"type"`               // Type is the type of the tool.
		Function FunctionDefinition `json:"function,omitempty"` // Function is the tool's functional definition.
	}
	// ToolType is the tool type.
	//
	// string
	ToolType string
	// ToolChoice represents the tool choice.
	ToolChoice struct {
		Type     ToolType     `json:"type"`               // Type is the type of the tool choice.
		Function ToolFunction `json:"function,omitempty"` // Function is the function of the tool choice.
	}
	// ToolFunction represents the tool function.
	ToolFunction struct {
		Name string `json:"name"` // Name is the name of the tool function.
	}
	// FunctionDefinition represents the function definition.
	FunctionDefinition struct {
		Name        string             `json:"name"`
		Description string             `json:"description"`
		Parameters  FunctionParameters `json:"parameters"`
	}
	// FunctionParameters represents the function parameters of a tool.
	FunctionParameters struct {
		Type                 string                        `json:"type"`
		Properties           map[string]PropertyDefinition `json:"properties"`
		Required             []string                      `json:"required"`
		AdditionalProperties bool                          `json:"additionalProperties,omitempty"`
	}
	// PropertyDefinition represents the property definition.
	PropertyDefinition struct {
		Type        string `json:"type"`
		Description string `json:"description"`
	}
	// ToolCall represents a tool call.
	ToolCall struct {
		// Index is not nil only in chat completion chunk object
		Index    *int         `json:"index,omitempty"` // Index is the index of the tool call.
		ID       string       `json:"id"`              // ID is the id of the tool call.
		Type     string       `json:"type"`            // Type is the type of the tool call.
		Function FunctionCall `json:"function"`        // Function is the function of the tool call.
	}
	// FunctionCall represents a function call.
	FunctionCall struct {
		Name      string `json:"name,omitempty"`      // Name is the name of the function call.
		Arguments string `json:"arguments,omitempty"` // Arguments is the arguments of the function call in JSON format.
	}
)
