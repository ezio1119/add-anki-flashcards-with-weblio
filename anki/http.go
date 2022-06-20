package anki

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	url       = "http://localhost:8765"
	version   = 6
	deckName  = "English"
	modelName = "Basic"
)

type requestBody[T params] struct {
	Action  string `json:"action"`
	Version int    `json:"version"`
	Params  *T     `json:"params,omitempty"`
}

type responseBody struct {
	Result json.RawMessage `json:"result"`
	Error  *string         `json:"error"`
}

func newRequestBody[T params](action string, params *T) *requestBody[T] {
	return &requestBody[T]{action, version, params}
}

func requestAPI[T params](ctx context.Context, action string, params *T) (*responseBody, error) {

	reqBody := newRequestBody(action, params)
	reqBodyByte, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	// fmt.Println(string(reqBodyByte))

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBodyByte))
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBoy := &responseBody{}

	if err := json.NewDecoder(res.Body).Decode(resBoy); err != nil {
		return nil, err
	}

	if resBoy.Error != nil {
		return nil, fmt.Errorf(*resBoy.Error)
	}

	return resBoy, nil
}
