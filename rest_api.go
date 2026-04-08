package cimego

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var (
	ErrBadRequest    = fmt.Errorf("잘못된 요청")
	ErrUnauthorized  = fmt.Errorf("인증 실패")
	ErrNotFound      = fmt.Errorf("정보를 찾을 수 없음")
	ErrInternalError = fmt.Errorf("내부 오류")
)

// ApiResponseBody는 ci.me의 API에서 반환되는 반환 구조체입니다.
type APIResponseBody struct {
	Code    int             `json:"code"`
	Message *string         `json:"message"`
	Content json.RawMessage `json:"content"`
}

// APIResponseContent는 APIResponseBody의 Content 필드에 들어갈 타입(인증 관련 제외)에서 쓰이는 구조체입니다.
type APIResponseContent[T any] struct {
	Data T `json:"data"`
	Page *struct {
		Next *string `json:"next"`
	} `json:"page"`
}

type header struct {
	Authorization string
	ClientID      string
	ClientSecret  string
}

func (c *CIME) get(ctx context.Context, url string, header *header, queryParams map[string]string) (*APIResponseBody, error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	addHeader(r.Header, header)
	addQueryParams(r.URL.Query(), queryParams)

	resp, err := c.apiClient.Do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data APIResponseBody
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	if data.Code != 200 {
		return nil, returnErr(data)
	}

	return &data, nil
}

func (c *CIME) post(ctx context.Context, url string, body any, header *header, queryParams map[string]string) (*APIResponseBody, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyBuffer := bytes.NewBuffer(bodyBytes)

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyBuffer)
	if err != nil {
		return nil, err
	}

	addHeader(r.Header, header)
	addQueryParams(r.URL.Query(), queryParams)

	resp, err := c.apiClient.Do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data APIResponseBody
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return nil, err
	}

	if data.Code != 200 {
		return nil, returnErr(data)
	}

	return &data, nil
}

func returnErr(data APIResponseBody) error {
	switch data.Code {
	case 400:
		return fmt.Errorf("%w: %s", ErrBadRequest, *data.Message)
	case 401:
		return fmt.Errorf("%w: %s", ErrUnauthorized, *data.Message)
	case 404:
		return fmt.Errorf("%w: %s", ErrNotFound, *data.Message)
	case 500:
		return fmt.Errorf("%w: %s", ErrInternalError, *data.Message)

	default:
		return nil
	}
}

func addHeader(h http.Header, data *header) {
	h.Add("Content-Type", "application/json")

	if data != nil {
		if data.Authorization != "" {
			h.Add("Authorization", fmt.Sprintf("Bearer %s", data.Authorization))
		}

		if data.ClientID != "" {
			h.Add("Client-Id", data.ClientID)
		}

		if data.ClientSecret != "" {
			h.Add("Client-Secret", data.ClientSecret)
		}
	}
}

func addQueryParams(q url.Values, queryParams map[string]string) {
	if queryParams == nil {
		return
	}

	for k, v := range queryParams {
		q.Add(k, v)
	}
}
