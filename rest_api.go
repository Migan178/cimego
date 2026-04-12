package cimego

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

type APIErrorResponseBody struct {
	Message    json.RawMessage `json:"message"`
	Error      string          `json:"error"`
	StatusCode int             `json:"statusCode"`
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

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, returnErr(respBody)
	}

	var data APIResponseBody
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return nil, err
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

	if resp.StatusCode != 200 {
		return nil, returnErr(respBody)
	}

	var data APIResponseBody
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *CIME) patch(ctx context.Context, url string, body any, header *header, queryParams map[string]string) (*APIResponseBody, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyBuffer := bytes.NewBuffer(bodyBytes)

	r, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bodyBuffer)
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

	if resp.StatusCode != 200 {
		return nil, returnErr(respBody)
	}

	var data APIResponseBody
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *CIME) put(ctx context.Context, url string, body any, header *header, queryParams map[string]string) (*APIResponseBody, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyBuffer := bytes.NewBuffer(bodyBytes)

	r, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bodyBuffer)
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

	if resp.StatusCode != 200 {
		return nil, returnErr(respBody)
	}

	var data APIResponseBody
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func returnErr(body []byte) error {
	var data APIErrorResponseBody
	err := json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	var message string
	err = json.Unmarshal(data.Message, &message)
	if err != nil {
		if _, ok := errors.AsType[*json.UnmarshalTypeError](err); ok {
			messages := []string{}
			err = json.Unmarshal(data.Message, &data)
			if err != nil {
				return err
			}

			for idx, msg := range messages {
				if idx == 0 {
					message = msg
				}

				message += "," + msg
			}
		}
	}

	switch data.StatusCode {
	case 400:
		return fmt.Errorf("%w: %s", ErrBadRequest, message)
	case 401:
		return fmt.Errorf("%w: %s", ErrUnauthorized, message)
	case 404:
		return fmt.Errorf("%w: %s", ErrNotFound, message)
	case 500:
		return fmt.Errorf("%w: %s", ErrInternalError, message)

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
