package session

import "time"

type SessionStore interface {
	//添加，读取，删除session
	Add(session Session) error
	Read(sid string) (*Session, error)
	Delete(sid string) error
	//回收过期会话
	GC(lifetime time.Duration)
}
