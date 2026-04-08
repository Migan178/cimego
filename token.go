package cimego

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	ErrTokenNotFound = fmt.Errorf("해당하는 토큰을 찾을 수 없습니다")
	ErrTokenExpired  = fmt.Errorf("토큰이 만료되었습니다")
)

// RefreshToken 구조체는 Refresh Token의 정보를 담고 있는 구조체입니다.
type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

// RefreshTokenStorage는 ci.me에서 발급 받은 Token을 저장하기 위한 인터페이스입니다.
type RefreshTokenStorage interface {
	// SaveToken은 Token을 저장하는 메서드입니다.
	SaveToken(ctx context.Context, channelID string, newToken RefreshToken) error
	// GetToken은 Token을 가져오는 메서드입니다.
	GetToken(ctx context.Context, channelID string) (*RefreshToken, error)
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

// SaveToken은 Refresh Token을 저장하는 메서드입니다.
func (s *FileRefreshTokenStorage) SaveToken(ctx context.Context, channelID string, newToken RefreshToken) error {
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

// GetToken은 Refresh Token을 가져오는 메서드입니다.
func (s *FileRefreshTokenStorage) GetToken(ctx context.Context, channelID string) (*RefreshToken, error) {
	tokens, err := s.getTokens(ctx)
	if err != nil {
		return nil, err
	}

	if token, ok := tokens[channelID]; ok {
		return &token, nil
	}

	return nil, ErrTokenNotFound
}

func (s *FileRefreshTokenStorage) getTokens(ctx context.Context) (map[string]RefreshToken, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	rawBytes, err := os.ReadFile(s.filename)
	if err != nil {
		return nil, err
	}

	var tokens = make(map[string]RefreshToken)
	err = json.Unmarshal(rawBytes, &tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// AccessToken 구조체는 Access Token의 정보를 담고 있는 구조체입니다.
type AccessToken struct {
	AccessToken string
	ExpiresAt   time.Time
	TokenType   string
	Scope       string
}

// Expired는 해당 토큰이 만료된 상태인지 판별합니다.
func (t *AccessToken) Expired() bool {
	return time.Now().After(t.ExpiresAt)
}

// AccessTokenStorage는 ci.me의 Access Token들을 저장하기 위한 인터페이스입니다.
type AccessTokenStorage interface {
	// SaveToken은 Access Token을 저장합니다.
	SaveToken(ctx context.Context, channelID string, newToken AccessToken) error
	// GetToken은 RefreshToken을 저장합니다
	GetToken(ctx context.Context, channelID string) (*AccessToken, error)
}

// InMemoryAccessTokenStorage는 Access Token을 메모리에 저장하는 구조체입니다.
type InMemoryAccessTokenStorage struct {
	tokens map[string]AccessToken
	mu     *sync.RWMutex
}

// NewInMemoryAccessTokenStorage는 새로운 InMemoryAccessTokenStorage 구조체의 인스턴스를 생성합니다.
func NewInMemoryAccessTokenStorage() *InMemoryAccessTokenStorage {
	return &InMemoryAccessTokenStorage{
		tokens: make(map[string]AccessToken),
		mu:     &sync.RWMutex{},
	}
}

// SaveToken은 Access Token을 메모리에 저장합니다.
func (s *InMemoryAccessTokenStorage) SaveToken(ctx context.Context, channelID string, newToken AccessToken) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.tokens[channelID] = newToken

	return nil
}

// GetToken은 Access Token을 가져옵니다.
func (s *InMemoryAccessTokenStorage) GetToken(ctx context.Context, channelID string) (*AccessToken, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RLock()

	if token, ok := s.tokens[channelID]; ok {
		if token.Expired() {
			return nil, ErrTokenExpired
		}

		return &token, nil
	}

	return nil, ErrTokenNotFound
}
