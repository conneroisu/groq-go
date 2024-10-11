package e2b

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
)

func (s *Sandbox) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Accept", "application/json")
	// Check whether Content-Type is already set, Upload Files API requires
	// Content-Type == multipart/form-data
	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < http.StatusOK ||
		res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("request to create sandbox failed: %s\nbody: %s", res.Status, getBody(res))
	}
	return decodeResponse(res.Body, v)
}

func (s *Sandbox) readWSResponse(v interface{}) (err error) {
	_, resp, err := s.ws.ReadMessage()
	if err != nil {
		return err
	}
	s.logger.Debug("read", "resp", string(resp))
	return decodeBytes(resp, v)
}

func decodeString(body io.Reader, output *string) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	*output = string(b)
	return nil
}
func decodeResponse(body io.Reader, v any) error {
	if v == nil {
		return nil
	}
	switch o := v.(type) {
	case *string:
		return decodeString(body, o)
	default:
		return json.NewDecoder(body).Decode(v)
	}
}
func getBody(resp *http.Response) string {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(b)
}
func decodeBytes(body []byte, v any) error {
	switch o := v.(type) {
	case *[]byte:
		*o = body
	default:
		return json.Unmarshal(body, v)
	}
	return nil
}

func (s *Sandbox) writeRequest(req Request) (err error) {
	s.msgCnt++
	req.ID = s.msgCnt
	jsVal, err := json.Marshal(req)
	if err != nil {
		return err
	}
	s.logger.Debug("write", "method", req.Method, "id", req.ID, "params", req.Params)
	err = s.ws.WriteMessage(websocket.TextMessage, jsVal)
	if err != nil {
		return fmt.Errorf("failed to write %s request (%d): %w", req.Method, req.ID, err)
	}
	return nil
}
