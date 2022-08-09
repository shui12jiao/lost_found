package session

import (
	"container/list"
	"sync"
	"time"
)

type MemorySessionStore struct {
	lock     sync.Mutex
	sessions map[string]*list.Element
	list     *list.List
}

func NewMemorySessionStore() *MemorySessionStore {
	return &MemorySessionStore{
		sessions: make(map[string]*list.Element),
		list:     list.New(),
	}
}

func (store *MemorySessionStore) Add(session Session) error {
	store.lock.Lock()
	defer store.lock.Unlock()
	element := store.list.PushBack(session)
	store.sessions[session.ID.String()] = element
	return nil
}

func (store *MemorySessionStore) Read(sid string) (*Session, error) {
	if element, ok := store.sessions[sid]; ok {
		session := element.Value.(Session)
		return &session, nil
	}
	return nil, ErrSessionNotFound
}

func (store *MemorySessionStore) Delete(sid string) error {
	store.lock.Lock()
	defer store.lock.Unlock()
	if element, ok := store.sessions[sid]; ok {
		store.list.Remove(element)
		delete(store.sessions, sid)
		return nil
	}
	return ErrSessionNotFound
}

func (store *MemorySessionStore) GC(lifetime time.Duration) {
	for {
		time.Sleep(time.Second * 2)

		element := store.list.Back()
		if element == nil {
			time.Sleep(lifetime)
		}

		if element.Value.(Session).IssuedAt.Add(lifetime).After(time.Now()) {
			store.Delete(element.Value.(Session).ID.String())
		}
	}
}
