package cimego

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type header struct {
	Authorization string
	ClientID      string
	ClientSecret  string
}

func get(client *http.Client, url string, header *header) (*APIResponseBody, error) {
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	addHeader(&r.Header, header)

	resp, err := client.Do(r)
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
		err = fmt.Errorf("http status code: %d, message: %s", data.Code, *data.Message)
		return nil, err
	}

	return &data, nil
}

func post(client *http.Client, url string, body any, header *header) (*APIResponseBody, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyBuffer := bytes.NewBuffer(bodyBytes)

	r, err := http.NewRequest(http.MethodPost, url, bodyBuffer)
	if err != nil {
		return nil, err
	}

	addHeader(&r.Header, header)

	resp, err := client.Do(r)
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
		err = fmt.Errorf("http Status Code: %d, message: %s", data.Code, *data.Message)
		return nil, err
	}
	return &data, nil
}

func addHeader(h *http.Header, data *header) {
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
