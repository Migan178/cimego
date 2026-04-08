package cimego

import (
	"context"
	"encoding/json"
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

	resp, err := c.get(EndpointChannels, &header{
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
