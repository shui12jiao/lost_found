package session

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
)

type Manager struct {
	cookieName string
	store      SessionStore
	lifetime   time.Duration
}

func NewSessionManager(cookieName string, store SessionStore, lifetime time.Duration) *Manager {
	manager := &Manager{
		cookieName: cookieName,
		store:      store,
		lifetime:   lifetime,
	}
	go manager.gc()
	return manager
}

func (manager *Manager) gc() {
	manager.store.GC(manager.lifetime)
}

func (manager *Manager) AddSession(value map[string]string) (string, error) {
	id := uuid.New()
	session := &Session{
		ID:       id,
		Value:    value,
		IssuedAt: time.Now(),
	}
	return id.String(), manager.store.Add(*session)
}

func (manager *Manager) ReadSession(sid string) (*Session, error) {
	session, err := manager.store.Read(sid)
	if session.IssuedAt.Add(manager.lifetime).Before(time.Now()) {
		session = nil
		err = ErrSessionExpired
	}
	return session, err
}

func (manager *Manager) DestorySession(sid string) error {
	return manager.store.Delete(sid)
}

type Session struct {
	ID       uuid.UUID
	Value    map[string]string
	IssuedAt time.Time
}
