package cimego

import (
	"encoding/json"
)

// ApiResponseBody는 ci.me의 API에서 반환되는 반환 구조체입니다.
type APIResponseBody struct {
	Code    int             `json:"code"`
	Message *string         `json:"message"`
	Content json.RawMessage `json:"content"`
}
