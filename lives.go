package cimego

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Live는 방송에 대한 정보들을 담고 있는 구조체입니다.
type Live struct {
	LiveID              string     `json:"liveId"`
	LiveTitle           string     `json:"liveTitle"`
	LiveThumbnailURL    *string    `json:"liveThumbnailImageUrl"`
	ConcurrentUserCount int        `json:"concurrentUserCount"`
	OpenedDate          *time.Time `json:"openedDate"`
	Adult               bool       `json:"adult"`
	Tags                []string   `json:"tags"`
	CategoryType        *string    `json:"categoryType"`
	LiveCategory        *string    `json:"liveCategory"`
	LiveCategoryValue   *string    `json:"liveCategoryValue"`
	ChannelID           string     `json:"channelId"`
	ChannelName         string     `json:"channelName"`
	ChannelHandle       string     `json:"channelHandle"`
	ChannelImageURL     *string    `json:"channelImageUrl"`
}

// Lives는 방송들의 목록을 LivesCursor로 반환합니다.
func (c *CIME) Lives(ctx context.Context, size int) (*LivesCursor, error) {
	lives, next, err := c.lives(ctx, size, "")
	if err != nil {
		return nil, err
	}

	return &LivesCursor{
		data: lives,
		cime: c,
		size: size,
		next: next,
	}, nil
}

func (c *CIME) lives(ctx context.Context, size int, next string) ([]Live, *string, error) {
	queryParams := map[string]string{"size": strconv.Itoa(size)}

	if next != "" {
		queryParams["next"] = next
	}

	resp, err := c.get(ctx, EndpointLives, &header{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
	}, queryParams)
	if err != nil {
		return nil, nil, err
	}

	var data APIResponseContent[[]Live]
	err = json.Unmarshal(resp.Content, &data)
	if err != nil {
		return nil, nil, err
	}

	return data.Data, data.Page.Next, nil
}

// LiveSettings는 해당 방송의 설정을 담고 있는 구조체입니다.
type LiveSettings struct {
	DefaultLiveTitle string    `json:"defaultLiveTitle"`
	Category         *Category `json:"category"`
	Tags             []string  `json:"tags"`
}

// LiveSettings는 해당 방송의 설정을 가져옵니다.
func (c *CIME) LiveSettings(ctx context.Context, channelID string) (*LiveSettings, error) {
	token, err := c.GetToken(ctx, channelID)
	if err != nil {
		return nil, err
	}

	resp, err := c.get(ctx, EndpointLivesSetting, &header{
		Authorization: token.AccessToken,
	}, nil)
	if err != nil {
		return nil, err
	}

	var content LiveSettings
	err = json.Unmarshal(resp.Content, &content)
	if err != nil {
		return nil, err
	}

	return &content, nil
}

// LiveSettingsUpdate는 라이브의 설정을 업데이트 하는 구조체입니다.
type LiveSettingsUpdate struct {
	DefaultLiveTitle *string   `json:"defaultLiveTitle,omitempty"`
	CategoryID       *string   `json:"categoryId,omitempty"`
	Tags             *[]string `json:"tags,omitempty"`
}

// UpdateLiveSettings은 해당 라이브의 설정을 일부 업데이트 합니다.
func (c *CIME) UpdateLiveSettings(ctx context.Context, channelID string, data *LiveSettingsUpdate) error {
	token, err := c.GetToken(ctx, channelID)
	if err != nil {
		return err
	}

	_, err = c.patch(ctx, EndpointLivesSetting, data, &header{
		Authorization: token.AccessToken,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}

// StreamKey는 해당 방송의 스트림키를 가져옵니다.
func (c *CIME) StreamKey(ctx context.Context, channelID string) (string, error) {
	token, err := c.GetToken(ctx, channelID)
	if err != nil {
		return "", err
	}

	resp, err := c.get(ctx, EndpointStreamKey, &header{
		Authorization: token.AccessToken,
	}, nil)
	if err != nil {
		return "", err
	}

	var content struct {
		StreamKey string `json:"streamKey"`
	}
	err = json.Unmarshal(resp.Content, &content)
	if err != nil {
		return "", err
	}

	return content.StreamKey, nil
}

type LiveStatus struct {
	IsLive   bool       `json:"isLive"`
	Title    *string    `json:"title"`
	OpenedAt *time.Time `json:"openedAt"`
}

// LiveStatus는 해당 채널의 방송 여부를 확인합니다.
// 해당 요청은 인증이 필요없습니다.
func (c *CIME) LiveStatus(ctx context.Context, channelID string) (*LiveStatus, error) {
	resp, err := c.get(ctx, fmt.Sprintf("%s/%s/%s/live-status", APIBaseURL, APIVersion, channelID), nil, nil)
	if err != nil {
		return nil, err
	}

	var content LiveStatus
	err = json.Unmarshal(resp.Content, &content)
	if err != nil {
		return nil, err
	}

	return &content, nil
}
