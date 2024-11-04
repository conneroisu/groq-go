package tools

const (
	// ToolTypeFunction is the function tool type.
	ToolTypeFunction ToolType = "function"
)

type (
	// Tool represents the tool.
	Tool struct {
		// Type is the type of the tool.
		Type ToolType `json:"type"`
		// Function is the tool's functional definition.
		Function FunctionDefinition `json:"function,omitempty"`
	}
	// ToolType is the tool type.
	//
	// string
	ToolType string
	// ToolChoice represents the tool choice.
	ToolChoice struct {
		// Type is the type of the tool choice.
		Type ToolType `json:"type"`
		// Function is the function of the tool choice.
		Function ToolFunction `json:"function,omitempty"`
	}
	// ToolFunction represents the tool function.
	ToolFunction struct {
		// Name is the name of the tool function.
		Name string `json:"name"`
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
		Index *int `json:"index,omitempty"`
		// ID is the id of the tool call.
		ID string `json:"id"`
		// Type is the type of the tool call.
		Type string `json:"type"`
		// Function is the function of the tool call.
		Function FunctionCall `json:"function"`
	}
	// FunctionCall represents a function call.
	FunctionCall struct {
		// Name is the name of the function call.
		Name string `json:"name,omitempty"`
		// Arguments is the arguments of the function call in JSON format.
		Arguments string `json:"arguments,omitempty"`
	}
)
