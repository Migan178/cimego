package cimego

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

// Channel은 채널에 대한 정보를 담고 있는 구조체입니다.
type Channel struct {
	ChannelID          string  `json:"channelId"`
	ChannelName        string  `json:"channelName"`
	ChannelHandle      string  `json:"channelHandle"`
	ChannelImageURL    *string `json:"channelImageUrl"`
	ChannelDescription string  `json:"channelDescription"`
	FollowerCount      int     `json:"followerCount"`
}

// Channels은 채널들에 대한 정보를 가져옵니다.
// 이는 ClientID와 ClientSecrets만 필요하며, Access Token은 사용되지 않습니다.
func (c *CIME) Channels(ctx context.Context, channelIDs []string) ([]Channel, error) {
	var channelIDsStr string

	for idx, channelID := range channelIDs {
		if idx == 0 {
			channelIDsStr = channelID
		}

		channelIDsStr += "," + channelID
	}

	resp, err := c.get(ctx, EndpointChannels, &header{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
	}, map[string]string{"channelIds": channelIDsStr})
	if err != nil {
		return nil, err
	}

	var channels []Channel
	err = json.Unmarshal(resp.Content, &channels)
	if err != nil {
		return nil, err
	}

	return channels, nil
}

// ChannelFollower는 해당 채널의 팔로워의 정보를 담고 있는 구조체입니다.
type ChannelFollower struct {
	ChannelID     string    `json:"channelId"`
	ChannelName   string    `json:"channelName"`
	ChannelHandle string    `json:"channelHandle"`
	CreatedDate   time.Time `json:"createdDate"`
}

// ChannelFollowers는 해당 채널의 팔로워 목록을 가져옵니다.
// 이는 Access Token을 사용하며, 해당 Access Token은 READ:CHANNEL 스코프가 필요합니다.
func (c *CIME) ChannelFollowers(ctx context.Context, channelID string, page, size int) ([]ChannelFollower, error) {
	token, err := c.AccessTokens.GetToken(ctx, channelID)
	if err != nil {
		if errors.Is(err, ErrTokenNotFound) || errors.Is(err, ErrTokenExpired) {
			token, err = c.Refresh(ctx, channelID)
			if err != nil {
				return nil, err
			}
		}

		return nil, err
	}

	resp, err := c.get(ctx, EndpointChannelFollowers, &header{
		Authorization: token.AccessToken,
	}, map[string]string{
		"page": strconv.Itoa(page),
		"size": strconv.Itoa(size),
	})
	if err != nil {
		return nil, err
	}

	var followers []ChannelFollower
	err = json.Unmarshal(resp.Content, &followers)
	if err != nil {
		return nil, err
	}

	return followers, nil
}

// ChannelSubscriber는 채널의 구독자를 가져올 때 정렬 방식을 선택하는 타입입니다.
type ChannelSubscriberSort string

const (
	ChannelSubscriberSortRecent ChannelSubscriberSort = "RECENT"
	ChannelSubscriberSortLonger ChannelSubscriberSort = "LONGER"
)

// ChannelSubscriber는 채널의 구독자에 대한 정보를 담고 있는 구조체입니다.
type ChannelSubscriber struct {
	ChannelID     string    `json:"channelId"`
	ChannelName   string    `json:"channelName"`
	ChannelHandle string    `json:"channelHandle"`
	Month         int       `json:"month"`
	TierNo        int       `json:"tierNo"`
	CreatedDate   time.Time `json:"createdDate"`
}

// ChannelSubscribers는 채널의 구독자 목록을 가져옵니다.
// 이는 Access Token을 사용하며, 해당 Access Token은 READ:SUBSCRIPTION 스코프가 필요합니다.
func (c *CIME) ChannelSubscribers(ctx context.Context, channelID string, page, size int, sort ChannelSubscriberSort) ([]ChannelSubscriber, error) {
	token, err := c.AccessTokens.GetToken(ctx, channelID)
	if err != nil {
		if errors.Is(err, ErrTokenNotFound) || errors.Is(err, ErrTokenExpired) {
			token, err = c.Refresh(ctx, channelID)
			if err != nil {
				return nil, err
			}
		}

		return nil, err
	}

	resp, err := c.get(ctx, EndpointChannelSubscribers, &header{
		Authorization: token.AccessToken,
	}, map[string]string{
		"page": strconv.Itoa(page),
		"size": strconv.Itoa(size),
		"sort": string(sort),
	})
	if err != nil {
		return nil, err
	}

	var subscribers []ChannelSubscriber
	err = json.Unmarshal(resp.Content, &subscribers)
	if err != nil {
		return nil, err
	}

	return subscribers, nil
}

// ManagerRole은 채널의 관리자가 어떠한 역할인지에 대한 타입입니다.
type ManagerRole string

const (
	ManagerRoleStreamingChannelOwner      ManagerRole = "STREAMING_CHANNEL_OWNER"
	ManagerRoleStreamingChannelManager    ManagerRole = "STREAMING_CHANNEL_MANAGER"
	ManagerRoleStreamingChatManager       ManagerRole = "STREAMING_CHAT_MANAGER"
	ManagerRoleStreamingSettlementManager ManagerRole = "STREAMING_SETTLEMENT_MANAGER"
)

// ChannelManager는 채널의 관리자에 대한 정보를 담고 있는 구조체입니다.
type ChannelManager struct {
	ManagerChannelID     string      `json:"managerChannelId"`
	ManagerChannelName   string      `json:"managerChannelName"`
	ManagerChannelHandle string      `json:"managerChannelHandle"`
	ManagerRole          ManagerRole `json:"userRole"`
	CreatedDate          time.Time   `json:"createdDate"`
}

// ChannelManagers는 채널의 관리자 목록을 가져옵니다.
// 이는 Access Token을 사용하며, 해당 Access Token은 READ:CHANNEL 스코프가 필요합니다.
func (c *CIME) ChannelManagers(ctx context.Context, channelID string) ([]ChannelManager, error) {
	token, err := c.AccessTokens.GetToken(ctx, channelID)
	if err != nil {
		if errors.Is(err, ErrTokenNotFound) || errors.Is(err, ErrTokenExpired) {
			token, err = c.Refresh(ctx, channelID)
			if err != nil {
				return nil, err
			}
		}

		return nil, err
	}

	resp, err := c.get(ctx, EndpointChannelManagers, &header{
		Authorization: token.AccessToken,
	}, nil)
	if err != nil {
		return nil, err
	}

	var managers []ChannelManager
	err = json.Unmarshal(resp.Content, &managers)
	if err != nil {
		return nil, err
	}

	return managers, nil
}
