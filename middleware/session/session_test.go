package session

import (
	"lost_found/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSession(t *testing.T) {
	lifetime := time.Second * 3
	manager := NewSessionManager("test", NewMemorySessionStore(), lifetime)

	issuedAt := time.Now()
	openid := util.RandomString(28)
	session_key := util.RandomString(24)
	sessionValue := make(map[string]string, 2)
	sessionValue["openid"] = openid
	sessionValue["session_key"] = session_key

	sessionId, err := manager.AddSession(sessionValue)
	require.NoError(t, err)

	session, err := manager.ReadSession(sessionId)
	require.NoError(t, err)
	require.NotEmpty(t, session)

	require.Equal(t, sessionId, session.ID.String())
	require.Equal(t, openid, sessionValue["openid"])
	require.Equal(t, session_key, sessionValue["session_key"])
	require.WithinDuration(t, issuedAt, session.IssuedAt, time.Second)

	time.Sleep(lifetime)
	session, err = manager.ReadSession(sessionId)
	require.Error(t, err)
	require.EqualError(t, err, ErrSessionNotFound.Error())
	require.Empty(t, session)
}
