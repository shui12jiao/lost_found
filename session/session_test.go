package session

import (
	"lost_found/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSession(t *testing.T) {
	lifetime := time.Second * 3
	manager := NewManager("test", NewMemorySessionStore(), lifetime)

	issuedAt := time.Now()
	openid := util.RandomString(28)
	session_key := util.RandomString(24)

	sessionId, err := manager.AddSession(openid, session_key)
	require.NoError(t, err)

	session, err := manager.ReadSession(sessionId)
	require.NoError(t, err)
	require.NotEmpty(t, session)

	require.Equal(t, sessionId, session.ID)
	require.Equal(t, openid, session.Value[0])
	require.Equal(t, session_key, session.Value[1])
	require.WithinDuration(t, issuedAt, session.IssuedAt, time.Second)

	time.Sleep(lifetime)
	session, err = manager.ReadSession(sessionId)
	require.Error(t, err)
	require.EqualError(t, err, ErrSessionNotFound.Error())
	require.Empty(t, session)
}
