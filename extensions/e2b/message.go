package e2b

const (
	rpc = "2.0"

	filesystemWrite      Method = "filesystem_write"
	filesystemRead       Method = "filesystem_read"
	filesystemList       Method = "filesystem_list"
	filesystemRemove     Method = "filesystem_remove"
	filesystemMakeDir    Method = "filesystem_makeDir"
	filesystemReadBytes  Method = "filesystem_readBase64"
	filesystemWriteBytes Method = "filesystem_writeBase64"
	// TODO: Check this one.
	filesystemSubscribe = "filesystem_subscribe"
)

type (
	// Method is a JSON-RPC method.
	Method string
	// Request is a JSON-RPC request.
	Request struct {
		// JSONRPC is the JSON-RPC version of the message.
		JSONRPC string `json:"jsonrpc"`
		// Method is the method of the message.
		Method Method `json:"method"`
		// ID is the ID of the message.
		ID int `json:"id"`
		// Params is the params of the message.
		Params []any `json:"params"`
	}
	// LsResponse is a JSON-RPC response.
	LsResponse struct {
		// JSONRPC is the JSON-RPC version of the message.
		JSONRPC string `json:"jsonrpc"`
		// Method is the method of the message.
		Method Method `json:"method"`
		// ID is the ID of the message.
		ID     int        `json:"id"`
		Result []LsResult `json:"result"`
	}
	// LsResult is a result of the list request.
	LsResult struct {
		Name  string `json:"name"`
		IsDir bool   `json:"isDir"`
	}
)
