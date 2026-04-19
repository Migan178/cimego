// 해당 패키지는 대한민국의 방송 플랫폼인 씨미에 대한 비공식 Golang 래퍼입니다.
package cimego

import (
	"fmt"
	"net/http"
)

type OnChatEventFunc func(c *CIME, e *ChatEvent)
type OnDonationEventFunc func(c *CIME, e *DonationEvent)
type OnSubscriptionEventFunc func(c *CIME, e *SubscriptionEvent)

// CIME는 ci.me API에 접근하기 위한 구조체입니다.
type CIME struct {
	onChatEvent         OnChatEventFunc
	onDonationEvent     OnDonationEventFunc
	onSubscriptionEvent OnSubscriptionEventFunc

	RefreshTokens RefreshTokenStorage
	AccessTokens  AccessTokenStorage
	ClientID      string
	ClientSecret  string
	apiClient     *http.Client
	sessions      map[string]*session
}

// CIMEOptions는 CIME 구조체를 생성할 때 넘겨줄 설정입니다.
type CIMEOptions struct {
	OnChatEvent         OnChatEventFunc
	OnDonationEvent     OnDonationEventFunc
	OnSubscriptionEvent OnSubscriptionEventFunc

	RefreshTokenStorage RefreshTokenStorage
	AccessTokenStorage  AccessTokenStorage
	APIClient           *http.Client
}

// New는 새로운 CIME 구조체의 인스턴스를 생성합니다.
func New(clientID, secret string, opts *CIMEOptions) (*CIME, error) {
	if clientID == "" {
		return nil, fmt.Errorf("clientID 값은 필수로 있어야 합니다")
	}

	if secret == "" {
		return nil, fmt.Errorf("secret 값은 필수로 있어야 합니다")
	}

	var apiClient = &http.Client{}

	var (
		accessTokens  AccessTokenStorage  = NewInMemoryAccessTokenStorage()
		refreshTokens RefreshTokenStorage = NewFileRefreshTokenStorage("")
	)

	var (
		onChatEvent         OnChatEventFunc
		onDonationEvent     OnDonationEventFunc
		onSubscriptionEvent OnSubscriptionEventFunc
	)

	if opts != nil {
		if opts.OnChatEvent != nil {
			onChatEvent = opts.OnChatEvent
		}

		if opts.OnDonationEvent != nil {
			onDonationEvent = opts.OnDonationEvent
		}

		if opts.OnSubscriptionEvent != nil {
			onSubscriptionEvent = opts.OnSubscriptionEvent
		}

		if opts.APIClient != nil {
			apiClient = opts.APIClient
		}

		if opts.RefreshTokenStorage != nil {
			refreshTokens = opts.RefreshTokenStorage
		}

		if opts.AccessTokenStorage != nil {
			accessTokens = opts.AccessTokenStorage
		}
	}

	cime := &CIME{
		onChatEvent:         onChatEvent,
		onDonationEvent:     onDonationEvent,
		onSubscriptionEvent: onSubscriptionEvent,

		RefreshTokens: refreshTokens,
		AccessTokens:  accessTokens,
		ClientID:      clientID,
		ClientSecret:  secret,
		apiClient:     apiClient,
	}

	return cime, nil
}
