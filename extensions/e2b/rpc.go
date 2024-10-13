package e2b

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type (
	// WSHandler is a handler for websockets.
	WSHandler struct {
		ws     *websocket.Conn
		idMap  sync.Map
		subMap sync.Map
	}
	decResp struct {
		ID     int `json:"id"`
		Params struct {
			Subscription string `json:"subscription"`
		}
	}
	// Different things we have to identify
	// 1. ID of response (matches request ID)
	// 2. Subscription ID (matches response.result from subscribe)

)

func newWSHandler(ctx context.Context, url string) (*WSHandler, error) {
	h := &WSHandler{}
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
	return h, nil
}

// Sub subscribes to a subscription id
func (h *WSHandler) Sub(id string, resCh chan []byte) {
	h.subMap.Store(id, resCh)
}

// UnSub subscribes to a subscription id
func (h *WSHandler) UnSub(id string) {
	h.subMap.Delete(id)
}

// Write writes a request to the websocket.
func (h *WSHandler) Write(req Request) (err error) {
	h.idMap.Store(req.ID, req.ResponseCh)
	jsVal, err := json.Marshal(req)
	if err != nil {
		return err
	}
	err = h.ws.WriteMessage(websocket.TextMessage, jsVal)
	if err != nil {
		return fmt.Errorf("failed to write %s request (%d): %w", req.Method, req.ID, err)
	}
	return nil
}

// Read reads a response from the websocket.
//
// If the context is cancelled, the websocket will be closed.
func (h *WSHandler) read(ctx context.Context) (err error) {
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
			if decResp.Params.Subscription != "" {
				toR, ok := h.subMap.Load(decResp.Params.Subscription)
				if !ok {
					return fmt.Errorf("subscription not found: %s", decResp.Params.Subscription)
				}
				toR.(chan []byte) <- resp
			}
			if decResp.ID == 0 {
				toR, ok := h.idMap.Load(decResp.ID)
				if !ok {
					return fmt.Errorf("response not found: %d", decResp.ID)
				}
				toRCh := toR.(chan []byte)
				toRCh <- resp
				// close(toRCh)
			}
		}
	}
}
