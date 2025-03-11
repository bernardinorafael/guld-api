package session

import (
	"time"

	"github.com/bernardinorafael/internal/_shared/util"
)

type session struct {
	id           string
	username     string
	refreshToken string
	agent        string
	ip           string
	revoked      bool
	expires      time.Time
	created      time.Time
	updated      time.Time
}

func NewFromDatabase(sess Entity) *session {
	return &session{
		id:           sess.ID,
		username:     sess.Username,
		refreshToken: sess.RefreshToken,
		agent:        sess.Agent,
		ip:           sess.IP,
		revoked:      sess.Revoked,
		expires:      sess.Expires,
		created:      sess.Created,
		updated:      sess.Updated,
	}
}

func New(username, refreshToken, agent, ip string) *session {
	return &session{
		id:           util.GenID("sess"),
		username:     username,
		refreshToken: refreshToken,
		agent:        agent,
		ip:           ip,
		revoked:      false,
		expires:      time.Now().Add(24 * time.Hour), // TODO: remove hard-coded expires
		created:      time.Now(),
		updated:      time.Now(),
	}
}

func (s *session) Revoke() {
	s.revoked = true
	s.updated = time.Now()
}

func (s *session) IsValidSession() bool {
	return !s.revoked && s.expires.After(time.Now())
}

func (s *session) Store() Entity {
	return Entity{
		ID:           s.ID(),
		Username:     s.Username(),
		RefreshToken: s.RefreshToken(),
		Agent:        s.Agent(),
		IP:           s.IP(),
		Revoked:      s.Revoked(),
		Expires:      s.Expires(),
		Created:      s.Created(),
		Updated:      s.Updated(),
	}
}

func (s *session) ID() string           { return s.id }
func (s *session) Username() string     { return s.username }
func (s *session) RefreshToken() string { return s.refreshToken }
func (s *session) Agent() string        { return s.agent }
func (s *session) IP() string           { return s.ip }
func (s *session) Revoked() bool        { return s.revoked }
func (s *session) Expires() time.Time   { return s.expires }
func (s *session) Created() time.Time   { return s.created }
func (s *session) Updated() time.Time   { return s.updated }
