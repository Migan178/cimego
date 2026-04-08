package cimego

import (
	"context"
	"encoding/json"
	"time"
)

// GrantType은 토큰을 얻을 방식을 선택합니다.
type GrantType string

const (
	GrantTypeAuthorizationCode GrantType = "authorization_code"
	GrantTypeRefreshToken      GrantType = "refresh_token"
)

// AuthorizationPayload는 인증용 토큰을 얻기 위해 보내야하는 구조체입니다.
type AuthorizationPayload struct {
	GrantType    GrantType `json:"grantType"`
	ClientID     string    `json:"clientId"`
	ClientSecret string    `json:"clientSecret"`
	Code         string    `json:"code,omitempty"`
	RefreshToken string    `json:"refreshToken,omitempty"`
}

type token struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int    `json:"expiresIn"`
	Scope        string `json:"scope"`
}

// Authorize는 토큰을 가져옵니다.
func (c *CIME) Authorize(ctx context.Context, authorizeCode string) error {
	payload := AuthorizationPayload{
		GrantType:    GrantTypeAuthorizationCode,
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Code:         authorizeCode,
	}

	resp, err := c.post(EndpointAuthorization, payload, nil, nil)
	if err != nil {
		return err
	}

	var data token
	err = json.Unmarshal(resp.Content, &data)
	if err != nil {
		return err
	}

	me, err := c.Me(ctx, data.AccessToken)
	if err != nil {
		return err
	}

	err = c.AccessTokens.SaveToken(ctx, me.ChannelID, AccessToken{
		AccessToken: data.AccessToken,
		ExpiresAt:   time.Now().Add(time.Duration(data.ExpiresIn)*time.Second - 5*time.Minute),
		TokenType:   data.TokenType,
		Scope:       data.Scope,
	})
	if err != nil {
		return err
	}

	return c.RefreshTokens.SaveToken(ctx, me.ChannelID, RefreshToken{
		RefreshToken: data.RefreshToken,
		TokenType:    data.TokenType,
		Scope:        data.Scope,
	})
}

// Refresh는 channelID에 연결된 Access Token이 만료되었을 때 해당 토큰을 새로 발급 받습니다.
func (c *CIME) Refresh(ctx context.Context, channelID string) error {
	oldToken, err := c.RefreshTokens.GetToken(ctx, channelID)
	if err != nil {
		return err
	}

	payload := AuthorizationPayload{
		GrantType:    GrantTypeRefreshToken,
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RefreshToken: oldToken.RefreshToken,
	}

	resp, err := c.post(EndpointAuthorization, payload, nil, nil)
	if err != nil {
		return err
	}

	var data token
	err = json.Unmarshal(resp.Content, &data)
	if err != nil {
		return err
	}

	err = c.AccessTokens.SaveToken(ctx, channelID, AccessToken{
		AccessToken: data.AccessToken,
		ExpiresAt:   time.Now().Add(time.Duration(data.ExpiresIn)*time.Second - 5*time.Minute),
		TokenType:   data.TokenType,
		Scope:       data.Scope,
	})
	if err != nil {
		return err
	}

	return c.RefreshTokens.SaveToken(ctx, channelID, RefreshToken{
		RefreshToken: data.RefreshToken,
		TokenType:    data.TokenType,
		Scope:        data.Scope,
	})
}

// User는 해당 사용자 계정과 연결된 채널의 정보를 담는 구조체입니다.
type User struct {
	ChannelID     string `json:"channelId"`
	ChannelName   string `json:"channelName"`
	ChannelHandle string `json:"channelHandle"`
}

// Me는 AccessToken에 연결된 사용자의 채널 정보를 가져옵니다.
func (c *CIME) Me(ctx context.Context, accessToken string) (*User, error) {
	resp, err := c.get(EndpointMe, &header{Authorization: accessToken}, nil)
	if err != nil {
		return nil, err
	}

	var data User
	err = json.Unmarshal(resp.Content, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
