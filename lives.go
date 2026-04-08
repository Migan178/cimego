package cimego

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

var (
	ErrLiveCursorIsEnded = fmt.Errorf("더 이상 다음으로 넘어갈 수 없음")
)

// LivesCursor는 커서 기반으로 방송 목록을 조작하는 구조체입니다.
type LivesCursor struct {
	data []Live
	cime *CIME
	size int
	next *string
}

// Next는 다음 커서로 이동합니다.
func (d *LivesCursor) Next(ctx context.Context) (*LivesCursor, error) {
	if d.next == nil {
		return nil, ErrLiveCursorIsEnded
	}

	lives, next, err := d.cime.lives(ctx, d.size, *d.next)
	if err != nil {
		return nil, err
	}

	d.data = lives
	d.next = next
	return d, nil
}

// Data는 해당 커서의 데이터들을 반환합니다.
func (d *LivesCursor) Data() []Live {
	dataCopy := make([]Live, len(d.data))
	copy(dataCopy, d.data)
	return dataCopy
}

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

	resp, err := c.get(ctx, ENdpointLives, &header{
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
