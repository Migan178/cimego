package cimego

import (
	"context"
	"encoding/json"
	"strconv"
	"time"
)

// AddRestrictedChannel은 해당 채널을 차단/추방 합니다.
func (c *CIME) AddRestrictedChannel(ctx context.Context, channelID, targetChannelID string) error {
	token, err := c.GetToken(ctx, channelID)
	if err != nil {
		return err
	}

	_, err = c.post(ctx, EndpointRestrictChannels, map[string]string{
		"targetChannelId": targetChannelID,
	}, &header{
		Authorization: token.AccessToken,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}

// RestrictedChannel은 해당 채널에서 차단/추방 당한 채널의 정보를 담고 있는 구조체입니다.
type RestrictedChannel struct {
	RestrictedChanelID      string     `json:"restrictedChannelId"`
	RestrictedChannelName   string     `json:"restrictedChannelName"`
	RestrictedChannelHandle string     `json:"restrictedChannelHandle"`
	CreatedAt               time.Time  `json:"createdDate"`
	ReleasesAt              *time.Time `json:"releaseDate"`
}

// RestrictedChannels는 차단/추방 된 채널들을 가져옵니다.
func (c *CIME) RestrictedChannels(ctx context.Context, channelID string, size int) (*RestrictedChannelsCursor, error) {
	restrictedChannels, next, err := c.restrictedChannels(ctx, channelID, size, "")
	if err != nil {
		return nil, err
	}

	return &RestrictedChannelsCursor{
		data:      restrictedChannels,
		channelID: channelID,
		cime:      c,
		size:      size,
		next:      next,
	}, nil
}

func (c *CIME) restrictedChannels(ctx context.Context, channelID string, size int, next string) ([]RestrictedChannel, *string, error) {
	token, err := c.GetToken(ctx, channelID)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.get(ctx, EndpointRestrictChannels, &header{
		Authorization: token.AccessToken,
	}, map[string]string{
		"size": strconv.Itoa(size),
		"next": next,
	})
	if err != nil {
		return nil, nil, err
	}

	var content APIResponseContent[[]RestrictedChannel]
	err = json.Unmarshal(resp.Content, &content)
	if err != nil {
		return nil, nil, err
	}

	return content.Data, content.Page.Next, nil
}

// DeleteRestrictedChannel은 차단/추방 된 채널의 차단/추방을 해제합니다.
func (c *CIME) DeleteRestrictedChannel(ctx context.Context, channelID, targetChannelID string) error {
	token, err := c.GetToken(ctx, channelID)
	if err != nil {
		return err
	}

	// 이거 그래서 왜 DELETE 요청 날릴 때 바디 실어서 보내야함??????????
	_, err = c.delete(ctx, EndpointRestrictChannels, map[string]string{
		"targetChannelId": targetChannelID,
	}, &header{
		Authorization: token.AccessToken,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}
