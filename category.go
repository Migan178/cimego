package cimego

import (
	"context"
	"encoding/json"
	"strconv"
)

// Category는 방송의 카테고리에 대한 정보를 담고 있는 구조체입니다.
type Category struct {
	CategoryID     string  `json:"categoryId"`
	CategoryType   string  `json:"categoryType"`
	CategoryValue  string  `json:"categoryValue"`
	PosterImageURL *string `json:"posterImageUrl"`
}

// Categories는 방송에 사용할 수 있는 카테고리들을 가져옵니다.
func (c *CIME) Categories(ctx context.Context, keyword string, size int) ([]Category, error) {
	resp, err := c.get(ctx, EndpointCategorySearch, &header{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
	}, map[string]string{
		"keyword": keyword,
		"size":    strconv.Itoa(size),
	})
	if err != nil {
		return nil, err
	}

	var content APIResponseContent[[]Category]
	err = json.Unmarshal(resp.Content, &content)
	if err != nil {
		return nil, err
	}

	return content.Data, nil
}
