package cimego

import (
	"context"
	"fmt"
)

var (
	ErrCursorIsEnded = fmt.Errorf("더 이상 다음으로 넘어갈 수 없음")
)

// LivesCursor는 커서 기반으로 방송 목록을 조작하는 구조체입니다.
type LivesCursor struct {
	data []Live
	cime *CIME
	size int
	next *string
}

// Next는 다음 커서로 이동합니다.
func (c *LivesCursor) Next(ctx context.Context) (*LivesCursor, error) {
	if c.next == nil {
		return nil, ErrCursorIsEnded
	}

	lives, next, err := c.cime.lives(ctx, c.size, *c.next)
	if err != nil {
		return nil, err
	}

	c.data = lives
	c.next = next
	return c, nil
}

// Data는 해당 커서의 데이터들을 반환합니다.
func (c *LivesCursor) Data() []Live {
	dataCopy := make([]Live, len(c.data))
	copy(dataCopy, c.data)
	return dataCopy
}

// RestrictedChannelsCursor는 커서 기반으로 차단/추방 된 유저 목록을 조작하는 구조체입니다.
type RestrictedChannelsCursor struct {
	data      []RestrictedChannel
	channelID string
	cime      *CIME
	size      int
	next      *string
}

// Next는 다음 커서로 넘어갑니다.
func (c *RestrictedChannelsCursor) Next(ctx context.Context) (*RestrictedChannelsCursor, error) {
	if c.next == nil {
		return nil, ErrCursorIsEnded
	}

	restrictedChannels, next, err := c.cime.restrictedChannels(ctx, c.channelID, c.size, *c.next)
	if err != nil {
		return nil, err
	}

	c.data = restrictedChannels
	c.next = next
	return c, nil
}

// Data는 해당 커서의 데이터들을 반환합니다.
func (c *RestrictedChannelsCursor) Data() []RestrictedChannel {
	dataCopy := make([]RestrictedChannel, len(c.data))
	copy(dataCopy, c.data)
	return dataCopy
}
