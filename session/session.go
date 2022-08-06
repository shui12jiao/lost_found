package session

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

type Manager struct {
	cookieName string
	store      SessionStore
	lifetime   time.Duration
}

func NewManager(cookieName string, store SessionStore, lifetime time.Duration) *Manager {
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

func (manager *Manager) AddSession(value ...any) (uuid.UUID, error) {
	id := uuid.New()
	session := &Session{
		ID:       id,
		Value:    value,
		IssuedAt: time.Now(),
	}
	return id, manager.store.Add(*session)
}

func (manager *Manager) ReadSession(id uuid.UUID) (*Session, error) {
	return manager.store.Read(id.String())
}

func (manager *Manager) DestorySession(id uuid.UUID) error {
	return manager.store.Delete(id.String())
}

type Session struct {
	ID       uuid.UUID
	Value    []any
	IssuedAt time.Time
}
