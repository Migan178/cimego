package cimego

import (
	"context"
	"encoding/json"
)

// ChatAllowedGroup은 채팅을 칠수 있는 유저의 종류를 나타내는 타입입니다.
type ChatAllowedGroup string

const (
	ChatAllowedGroupAll      ChatAllowedGroup = "ALL"
	ChatAllowedGroupFollower ChatAllowedGroup = "FOLLOWER"
	ChatAllowedGroupManager  ChatAllowedGroup = "MANAGER"
)

// ChatSettings는 채팅의 설정을 담고 있는 구조체입니다.
type ChatSettings struct {
	ChatAllowedGroup            ChatAllowedGroup `json:"chatAllowedGroup"`
	MinFollowerMinute           int              `json:"minFollowerMinute"`
	FollowerSubscriberChatAllow *bool            `json:"followerSubscriberChatAllow"`
}

// ChatSettings는 해당 방송의 채팅의 설정을 가져옵니다
func (c *CIME) ChatSettings(ctx context.Context, channelID string) (*ChatSettings, error) {
	token, err := c.GetToken(ctx, channelID)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, EndpointChatSettings, &header{
		Authorization: token.AccessToken,
	}, nil)
	if err != nil {
		return nil, err
	}

	var content ChatSettings
	err = json.Unmarshal(resp.Content, &content)
	if err != nil {
		return nil, err
	}

	return &content, nil
}

// ChatSettingsUpdate는 채팅의 설정을 업데이트 하는 구조체입니다.
type ChatSettingsUpdate struct {
	ChatEmojiMode               *bool             `json:"chatEmojiMode"`
	ChatSlowModeSec             *int              `json:"chatSlowModeSec"`
	ChatAllowedGroup            *ChatAllowedGroup `json:"chatAllowedGroup"`
	MinFollowerMinute           *int              `json:"minFollowerMinute"`
	FollowerSubscriberChatAllow *bool             `json:"followerSubscriberChatAllow"`
}

// UpdateChatSettings는 해당 채널의 채팅 설정을 업데이트합니다.
func (c *CIME) UpdateChatSettings(ctx context.Context, channelID string, data *ChatSettingsUpdate) error {
	token, err := c.GetToken(ctx, channelID)
	if err != nil {
		return err
	}

	_, err = c.put(ctx, EndpointChatSettings, data, &header{
		Authorization: token.AccessToken,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}

// ChatSenderType은 채팅을 보낼 때 발신자의 종류를 정하는 타입입니다.
type ChatSenderType string

const (
	ChatSenderTypeApp  ChatSenderType = "APP"
	ChatSenderTypeUser ChatSenderType = "USER"
)

// SendChat은 해당 채널에 채팅을 보냅니다.
func (c *CIME) SendChat(ctx context.Context, channelID string, senderType ChatSenderType, message string) (string, error) {
	token, err := c.GetToken(ctx, channelID)
	if err != nil {
		return "", err
	}

	resp, err := c.post(ctx, EndpointChatSend, map[string]any{
		"message":    message,
		"senderType": senderType,
	}, &header{
		Authorization: token.AccessToken,
	}, nil)
	if err != nil {
		return "", err
	}

	var content struct {
		ID string `json:"messageId"`
	}
	err = json.Unmarshal(resp.Content, &content)
	if err != nil {
		return "", err
	}

	return content.ID, nil
}

// SetChatNotice는 채팅에 공지사항을 등록합니다.
// message 또는 messageID 둘 중 하나를 반드시 제공해야 하며,
// 둘 다 없을 시, 400 에러가 반환됩니다.
func (c *CIME) SetChatNotice(ctx context.Context, channelID, message, messageID string) error {
	token, err := c.GetToken(ctx, channelID)
	if err != nil {
		return err
	}

	body := make(map[string]string)

	if message != "" {
		body["message"] = message
	}

	if messageID != "" {
		body["messageId"] = messageID
	}

	_, err = c.post(ctx, EndpointChatNotice, body, &header{
		Authorization: token.AccessToken,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}
