package e2b

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/gorilla/websocket"
)

type (
	// defaultWSHandler is a handler for websockets.
	defaultWSHandler struct {
		logger *slog.Logger
		ws     *websocket.Conn
		Map    sync.Map
		msgCnt int
	}
	decResp struct {
		Method string `json:"method"`
		ID     int    `json:"id"`
		Params struct {
			Subscription string `json:"subscription"`
		}
	}
	// Different things we have to identify
	// 1. ID of response (matches request ID)
	// 2. Subscription ID (matches response.result from subscribe)

)

func newWSHandler(ctx context.Context, url string, logger *slog.Logger) (*defaultWSHandler, error) {
	h := defaultWSHandler{
		// logger: slog.New(slog.NewJSONHandler(io.Discard, nil)),
		logger: logger,
		msgCnt: 1,
	}
	var err error
	h.ws, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	go func() {
		err := h.read(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}()
	return &h, nil
}

// Write writes a request to the websocket.
func (h *defaultWSHandler) Write(req Request) (err error) {
	req.ID = h.msgCnt
	h.logger.Debug("writing request", "method", req.Method, "id", req.ID, "params", req.Params)
	h.Map.Store(req.ID, req.ResponseCh)
	jsVal, err := json.Marshal(req)
	if err != nil {
		return err
	}
	err = h.ws.WriteMessage(websocket.TextMessage, jsVal)
	if err != nil {
		return fmt.Errorf("failed to write %s request (%d): %w", req.Method, req.ID, err)
	}
	h.msgCnt++
	return nil
}

// Read reads a response from the websocket.
//
// If the context is cancelled, the websocket will be closed.
func (h *defaultWSHandler) read(ctx context.Context) (err error) {
	defer func() {
		err = h.ws.Close()
	}()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			_, resp, err := h.ws.ReadMessage()
			if err != nil {
				return err
			}
			var decResp decResp
			err = json.Unmarshal(resp, &decResp)
			if err != nil {
				return err
			}
			h.logger.Debug("reading response", "method", decResp.Method, "id", decResp.ID, "body", string(resp))
			if decResp.Params.Subscription != "" {
				toR, ok := h.Map.Load(decResp.Params.Subscription)
				if !ok {
					h.logger.Debug("subscription not found", "id", decResp.Params.Subscription)
				}
				toRCh, ok := toR.(chan []byte)
				if !ok {
					h.logger.Debug("subscription not found", "id", decResp.Params.Subscription)
				}
				toRCh <- resp
			}
			if decResp.ID != 0 {
				toR, ok := h.Map.Load(decResp.ID)
				if !ok {
					h.logger.Debug("response not found", "id", decResp.ID, "values", func() []string {
						vals := make([]string, 0)
						h.Map.Range(func(key, value any) bool {
							vals = append(vals, fmt.Sprintf("%s: %v", key, value))
							return true
						})
						return vals
					}())
				}
				toRCh, ok := toR.(chan []byte)
				if !ok {
					h.logger.Debug("responsech not found", "id", decResp.ID)
				}
				toRCh <- resp
			}
		}
	}
}
