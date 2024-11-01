package jigsawstack

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/conneroisu/groq-go/pkg/builders"
)

const (
	uploadEndpoint Endpoint = "/v1/store/file"
	kvEndpoint     Endpoint = "/v1/store/kv"
)

type (
	// StorageResponse represents a response structure for file API.
	StorageResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message,omitempty"`
		URL     string `json:"url"`
		Key     string `json:"key"`
		Value   string `json:"value,omitempty"`
	}
)

// https://docs.jigsawstack.com/api-reference/store/file/get
// Upload Retrieve Delete

// FileAdd uploads a file to the Jigsaw Stack file store.
//
// https://docs.jigsawstack.com/api-reference/store/file/add
//
// POST https://api.jigsawstack.com/v1/store/file
func (j *JigsawStack) FileAdd(
	ctx context.Context,
	key string,
	contentType string,
	body io.Reader,
) (string, error) {
	// TODO: may need to santize the key
	url := j.baseURL + string(uploadEndpoint) + "?key=" + key
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodPost,
		url,
		builders.WithBody(body),
		builders.WithContentType(contentType),
	)
	if err != nil {
		return "", err
	}
	var resp StorageResponse
	err = j.sendRequest(req, &resp)
	if err != nil {
		return "", err
	}
	if !resp.Success {
		return "", fmt.Errorf("failed to upload file: %s", resp.Message)
	}
	return "", nil
}

// FileGet retrieves a file from the Jigsaw Stack file store.
//
// https://docs.jigsawstack.com/api-reference/store/file/get
//
// GET https://api.jigsawstack.com/v1/store/file/{fileName}
func (j *JigsawStack) FileGet(ctx context.Context, fileName string) (string, error) {
	// TODO: may need to santize the fileName
	url := j.baseURL + string(uploadEndpoint) + "/" + fileName
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodGet,
		url,
	)
	if err != nil {
		return "", err
	}
	var resp StorageResponse
	err = j.sendRequest(req, &resp)
	if err != nil {
		return "", err
	}
	if !resp.Success {
		return "", fmt.Errorf("failed to retrieve file: %s", resp.Message)
	}
	// TODO: may need to return the file content from url
	return resp.Message, nil
}

// FileDelete deletes a file from the Jigsaw Stack file store.
//
// https://docs.jigsawstack.com/api-reference/store/file/delete
//
// DELETE https://api.jigsawstack.com/v1/store/file/{fileName}
func (j *JigsawStack) FileDelete(fileName string) error {
	// TODO: may need to santize the fileName
	url := j.baseURL + string(uploadEndpoint) + "/" + fileName
	req, err := builders.NewRequest(
		context.Background(),
		j.header,
		http.MethodDelete,
		url,
	)
	if err != nil {
		return err
	}
	var resp StorageResponse
	err = j.sendRequest(req, &resp)
	if err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("failed to delete file: %s", resp.Message)
	}
	return nil
}

// KVAdd adds a key value pair to the Jigsaw Stack key-value store.
//
// https://docs.jigsawstack.com/api-reference/store/kv/add
//
// POST https://api.jigsawstack.com/v1/store/kv
func (j *JigsawStack) KVAdd(
	ctx context.Context,
	key string,
	value string,
) error {
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodPost,
		j.baseURL+string(kvEndpoint),
		builders.WithBody(map[string]string{
			"key":   key,
			"value": value,
		}),
	)
	if err != nil {
		return err
	}
	var resp StorageResponse
	err = j.sendRequest(req, &resp)
	if err != nil {
		return err
	}
	return nil
}

// KVGet retrieves a key value pair from the Jigsaw Stack key-value store.
//
// https://docs.jigsawstack.com/api-reference/store/kv/get
//
// GET https://api.jigsawstack.com/v1/store/kv/{key}
func (j *JigsawStack) KVGet(
	ctx context.Context,
	key string,
) (string, error) {
	url := j.baseURL + string(kvEndpoint) + "/" + key
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodGet,
		url,
	)
	if err != nil {
		return "", err
	}
	var resp StorageResponse
	err = j.sendRequest(req, &resp)
	if err != nil {
		return "", err
	}
	return "", nil
}

// KVDelete deletes a key value pair from the Jigsaw Stack key-value store.
//
// https://docs.jigsawstack.com/api-reference/store/kv/delete
//
// DELETE https://api.jigsawstack.com/v1/store/kv/{key}
func (j *JigsawStack) KVDelete(
	ctx context.Context,
	key string,
) (string, error) {
	url := j.baseURL + string(kvEndpoint) + "/" + key
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodDelete,
		url,
	)
	if err != nil {
		return "", err
	}
	var resp StorageResponse
	err = j.sendRequest(req, &resp)
	if err != nil {
		return "", err
	}
	return resp.Message, nil
}
