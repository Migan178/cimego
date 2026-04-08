package cimego

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

var (
	ErrTokenNotFound = fmt.Errorf("해당하는 토큰을 찾을 수 없습니다")
)

// Token 구조체는 Refresh Token의 정보를 담고 있는 구조체입니다.
type Token struct {
	AccessToken  string `json:"-"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

// RefreshTokenStorage는 ci.me에서 발급 받은 Token을 저장하기 위한 저장소입니다.
type RefreshTokenStorage interface {
	// SaveToken은 Refresh Token을 저장하는 메서드입니다.
	SaveToken(ctx context.Context, userID string, newToken Token) error
	// GetToken은 Refresh Token을 가져오는 메서드입니다.
	GetToken(ctx context.Context, userID string) (*Token, error)
}

// FileRefreshTokenStorage는 json 형식의 파일에 ci.me API의 Refresh Token을 저장하는 구조체입니다.
type FileRefreshTokenStorage struct {
	filename string
}

// NewFileRefreshTokenStorage는 새로운 FileRefreshTokenStorage 구조체의 인스턴스를 생성합니다.
// filename 매개변수는 선택 값입니다.
func NewFileRefreshTokenStorage(filename string) *FileRefreshTokenStorage {
	if filename == "" {
		filename = "cime_token.json"
	}

	return &FileRefreshTokenStorage{
		filename: filename,
	}
}

// SaveToken은 Token을 저장하는 메서드입니다.
func (s *FileRefreshTokenStorage) SaveToken(ctx context.Context, channelID string, newToken Token) error {
	tokens, err := s.getTokens(ctx)
	if err != nil {
		return err
	}

	tokens[channelID] = newToken

	rawBytes, err := json.Marshal(tokens)
	if err != nil {
		return err
	}

	return os.WriteFile(s.filename, []byte(rawBytes), 0777)
}

// GetToken은 Token을 가져오는 메서드입니다.
func (s *FileRefreshTokenStorage) GetToken(ctx context.Context, channelID string) (*Token, error) {
	tokens, err := s.getTokens(ctx)
	if err != nil {
		return nil, err
	}

	if token, ok := tokens[channelID]; ok {
		return &token, nil
	}

	return nil, ErrTokenNotFound
}

func (s *FileRefreshTokenStorage) getTokens(_ context.Context) (map[string]Token, error) {
	rawBytes, err := os.ReadFile(s.filename)
	if err != nil {
		return nil, err
	}

	var tokens = make(map[string]Token)
	err = json.Unmarshal(rawBytes, &tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
