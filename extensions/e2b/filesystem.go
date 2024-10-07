package e2b

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gorilla/websocket"
)

// Mkdir makes a directory in the sandbox file system.
func (s *Sandbox) Mkdir(path string) error {
	s.logger.Debug("Making directory", "path", path)
	s.msgCnt++
	jsVal, err := json.Marshal(Request{
		Params:  []any{path},
		JSONRPC: rpc,
		ID:      s.msgCnt,
		Method:  filesystemMakeDir,
	})
	if err != nil {
		return err
	}
	err = s.ws.WriteMessage(websocket.TextMessage, jsVal)
	if err != nil {
		return err
	}
	_, _, err = s.ws.ReadMessage()
	if err != nil {
		return err
	}
	return nil
}

// Ls lists the files and/or directories in the sandbox file system at
// the given path.
func (s *Sandbox) Ls(path string) ([]LsResult, error) {
	s.logger.Debug("Listing files and dirs", "path", path)
	s.msgCnt++
	jsVal, err := json.Marshal(Request{
		Params:  []any{path},
		JSONRPC: rpc,
		ID:      s.msgCnt,
		Method:  filesystemList,
	})
	if err != nil {
		return nil, err
	}
	err = s.ws.WriteMessage(websocket.TextMessage, jsVal)
	if err != nil {
		return nil, err
	}
	_, msr, err := s.ws.ReadMessage()
	if err != nil {
		return nil, err
	}
	var res LsResponse
	err = json.Unmarshal(msr, &res)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

// Read reads a file from the sandbox file system.
func (s *Sandbox) Read(
	path string,
) ([]byte, error) {
	s.logger.Debug("Reading from file", "path", path)
	s.msgCnt++
	jsnV, err := json.Marshal(Request{
		JSONRPC: rpc,
		Method:  filesystemRead,
		Params:  []any{path},
		ID:      s.msgCnt,
	})
	if err != nil {
		return nil, err
	}
	err = s.ws.WriteMessage(websocket.TextMessage, jsnV)
	if err != nil {
		return nil, err
	}
	_, message, err := s.ws.ReadMessage()
	if err != nil {
		return nil, err
	}
	var resp ReadResponse
	err = json.Unmarshal(message, &resp)
	if err != nil {
		return nil, err
	}
	return []byte(resp.Result), nil
}

// Write writes to a file to the sandbox file system.
func (s *Sandbox) Write(path string, data []byte) error {
	s.logger.Debug("Writing to file", "path", path)
	s.msgCnt++
	jsnV, err := json.Marshal(Request{
		JSONRPC: rpc,
		Method:  filesystemWrite,
		Params: []any{
			path,
			string(data),
		},
		ID: s.msgCnt,
	})
	if err != nil {
		return err
	}
	err = s.ws.WriteMessage(websocket.TextMessage, jsnV)
	if err != nil {
		return err
	}
	_, _, err = s.ws.ReadMessage()
	if err != nil {
		return err
	}
	return nil
}

// ReadBytes reads a file from the sandbox file system.
func (s *Sandbox) ReadBytes(path string) ([]byte, error) {
	s.logger.Debug("Reading Bytes", "path", path)
	s.msgCnt++
	jsnV, err := json.Marshal(Request{
		JSONRPC: rpc,
		Method:  filesystemReadBytes,
		Params: []any{
			path,
		},
		ID: s.msgCnt,
	})
	if err != nil {
		return nil, err
	}
	err = s.ws.WriteMessage(websocket.TextMessage, jsnV)
	if err != nil {
		return nil, err
	}
	sid, message, err := s.ws.ReadMessage()
	if err != nil {
		return nil, err
	}
	fmt.Println(string(message))
	fmt.Println(sid)
	return nil, nil
}

// Watch watches a directory in the sandbox file system.
//
// This is intended to be run in a goroutine as it will block until the
// connection is closed, an error occurs, or the context is canceled.
func (s *Sandbox) Watch(
	ctx context.Context,
	path string,
) (<-chan Event, error) {
	return nil, nil
}

// Upload uploads a file to the sandbox file system.
func (s *Sandbox) Upload(r io.Reader, path string) error {
	return nil
}

// Download downloads a file from the sandbox file system.
func (s *Sandbox) Download(path string) (io.ReadCloser, error) {
	return nil, nil
}
