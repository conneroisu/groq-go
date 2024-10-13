package e2b

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (s *Sandbox) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Accept", "application/json")
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
	if v == nil {
		return nil
	}
	switch o := v.(type) {
	case *string:
		return decodeString(res.Body, o)
	default:
		return json.NewDecoder(res.Body).Decode(v)
	}
}
func decodeString(body io.Reader, output *string) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	*output = string(b)
	return nil
}
func getBody(resp *http.Response) string {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(b)
}
