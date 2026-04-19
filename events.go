package cimego

import (
	"context"
	"encoding/json"
	"time"
)

// ChatEvent는 채팅 이벤트가 수신 되었을 때, 채팅 내용을 담고 있는 구조체입니다.
type ChatEvent struct {
	ChannelID       string `json:"channelID"`
	SenderChannelID string `json:"senderChannelId"`
	Profile         struct {
		Nickname string `json:"nickname"`
	} `json:"profile"`
	Content   string            `json:"content"`
	CreatedAt time.Time         `json:"messageTime"`
	Emojis    map[string]string `json:"emojis"`
}

// SubscribeChatEvent는 채팅 이벤트를 구독합니다.
// channelID가 공란일 시, client 세션을 사용합니다.
func (c *CIME) SubscribeChatEvent(ctx context.Context, channelID string) error {
	return c.subscribeEvent(ctx, sessionEventTypeChat, channelID)
}

// DonationType은 도네이션의 형식입니다.
type DonationType string

const (
	DonationTypeChat  DonationType = "CHAT"
	DonationTypeVideo DonationType = "VIDEO"
)

// DonationEvent는 후원 이벤트가 수신 되었을 때, 후원 내용이 담겨 있는 구조체입니다.
type DonationEvent struct {
	Type             DonationType      `json:"donationType"`
	ChannelID        string            `json:"channelID"`
	DonatorChannelID string            `json:"donatorChannelId"`
	DonatorNickname  string            `json:"donatorNickname"`
	PayAmount        string            `json:"payAmount"`
	DonationText     string            `json:"donationText"`
	Emojis           map[string]string `json:"emojis"`
}

// SubscribeDonationEvent는 후원 이벤트를 구독합니다.
// channelID가 공란일 시, client 세션을 사용합니다.
func (c *CIME) SubscribeDonationEvent(ctx context.Context, channelID string) error {
	return c.subscribeEvent(ctx, sessionEventTypeDonation, channelID)
}

// SubscriptionEvent는 구독 이벤트가 수신 되었을 때, 구독 내용이 담겨 있는 구조체입니다.
type SubscriptionEvent struct {
	ChannelID             string            `json:"channelId"`
	SubscriberChannelID   string            `json:"subscriberChannelId"`
	SubscriberChannelName string            `json:"subscriberChannelName"`
	Month                 int               `json:"month"`
	TierNo                int               `json:"tierNo"`
	SubscriptionMessage   string            `json:"subscriptionMessage"`
	Emojis                map[string]string `json:"emojis"`
}

// SubscribeSubscriptionEvent는 구독 이벤트를 구독합니다.
// channelID가 공란일 시, client 세션을 사용합니다.
func (c *CIME) SubscribeSubscriptionEvent(ctx context.Context, channelID string) error {
	return c.subscribeEvent(ctx, sessionEventTypeSubscription, channelID)
}

func (c *CIME) subscribeEvent(ctx context.Context, eventType sessionEventType, channelID string) error {
	var ok bool
	var err error
	var key string
	var endpoint string
	var session *session
	var sessionType string

	if channelID != "" {
		key = channelID
		sessionType = "user"
	} else {
		key = c.ClientID
		sessionType = "client"
	}

	switch eventType {
	case sessionEventTypeChat:
		endpoint = EndpointSessionEventsChat
	case sessionEventTypeDonation:
		endpoint = EndpointSessionEventsDonation
	case sessionEventTypeSubscription:
		endpoint = EndpointSessionEventsSubscription
	}

	session, ok = c.sessions[key]
	if !ok {
		session, err = c.newSession(ctx, sessionType, channelID)
		if err != nil {
			return err
		}

		c.sessions[key] = session
	}

	_, err = c.post(ctx, endpoint, nil, nil, map[string]string{
		"sessionKey": session.sessionKey,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *session) handleEvent(payload *sessionPayload) {
	switch payload.Type {
	case sessionEventTypeChat:
		if s.cime.onChatEvent == nil {
			return
		}

		var data ChatEvent
		err := json.Unmarshal(payload.Data, &data)
		if err != nil {
			return
		}

		s.cime.onChatEvent(s.cime, &data)
	case sessionEventTypeDonation:
		if s.cime.onDonationEvent == nil {
			return
		}

		var data DonationEvent
		err := json.Unmarshal(payload.Data, &data)
		if err != nil {
			return
		}

		s.cime.onDonationEvent(s.cime, &data)
	case sessionEventTypeSubscription:
		if s.cime.onSubscriptionEvent == nil {
			return
		}

		var data SubscriptionEvent
		err := json.Unmarshal(payload.Data, &data)
		if err != nil {
			return
		}

		s.cime.onSubscriptionEvent(s.cime, &data)
	}
}
