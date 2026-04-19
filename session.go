package cimego

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type sessionEventType string

const (
	sessionEventTypeChat         sessionEventType = "CHAT"
	sessionEventTypeDonation     sessionEventType = "DONATION"
	sessionEventTypeSubscription sessionEventType = "SUBSCRIPTION"
)

type sessionCreateBody struct {
	URL string `json:"url"`
}

type sessionPayload struct {
	Type sessionEventType `json:"type"`
	Data json.RawMessage  `json:"data"`
}

type session struct {
	conn        *websocket.Conn
	sessionType string
	channelID   string
	cime        *CIME
	key         string
	sessionKey  string
}

func (c *CIME) createSessionURL(ctx context.Context, sessionType string, channelID string) (string, string, error) {
	var h *header
	var endpoint string

	if sessionType == "user" {
		token, err := c.GetToken(ctx, channelID)
		if err != nil {
			return "", "", err
		}

		h = &header{
			Authorization: token.AccessToken,
		}

		endpoint = EndpointSessionsAuth
	}

	if sessionType == "client" {
		h = &header{
			ClientID:     c.ClientID,
			ClientSecret: c.ClientSecret,
		}

		endpoint = EndpointSessionAuthClient
	}

	resp, err := c.get(ctx, endpoint, h, nil)
	if err != nil {
		return "", "", err
	}

	var content sessionCreateBody
	err = json.Unmarshal(resp.Content, &content)
	if err != nil {
		return "", "", err
	}

	query, err := url.ParseQuery(content.URL)
	if err != nil {
		return "", "", err
	}

	sessionKey := query.Get("sessionKey")

	return content.URL, sessionKey, nil
}

func (c *CIME) newSession(ctx context.Context, sessionType string, channelID string) (*session, error) {
	var key string

	if sessionType == "user" {
		key = channelID
	}

	if sessionType == "client" {
		key = c.ClientID
	}

	if session, ok := c.sessions[key]; ok {
		return session, nil
	}

	sessionURL, sessionKey, err := c.createSessionURL(ctx, sessionType, channelID)
	if err != nil {
		return nil, err
	}

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, sessionURL, nil)
	if err != nil {
		return nil, err
	}

	session := &session{
		conn:        conn,
		cime:        c,
		sessionType: sessionType,
		channelID:   channelID,
		key:         key,
		sessionKey:  sessionKey,
	}

	go session.heartbeat()
	go session.listen()

	return session, nil
}

func (s *session) heartbeat() {
	ticker := time.NewTicker(1 * time.Minute)

	for range ticker.C {
		err := s.conn.WriteJSON(map[string]string{"type": "PING"})
		if err != nil {
			s.reconnect()
		}
	}
}

func (s *session) listen() {
	for {
		var payload sessionPayload

		err := s.conn.ReadJSON(&payload)
		if err != nil {
			var payload map[string]string

			err = s.conn.ReadJSON(&payload)
			if err != nil {
				s.reconnect()
			}

			if payload["action"] == "PONG" {
				continue
			}
		}

		s.handleEvent(&payload)
	}
}

func (s *session) reconnect() {
	err := s.close()
	if err != nil {
		return
	}

	session, err := s.cime.newSession(context.Background(), s.sessionType, s.channelID)
	if err != nil {
		return
	}

	s.cime.sessions[s.key] = session
}

func (s *session) close() error {
	err := s.conn.Close()
	if err != nil {
		return err
	}

	delete(s.cime.sessions, s.key)

	return nil
}
