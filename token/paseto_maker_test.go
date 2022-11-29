package token

import (
	"testing"
	"time"

	"github.com/escalopa/go-bank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	pasetoMaker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := pasetoMaker.CreatToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := pasetoMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)

	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, payload.ExpireAt, expiredAt, time.Second)
}

func TestPasetoMakerExpired(t *testing.T) {
	JWTMaker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()

	token, err := JWTMaker.CreatToken(username, -time.Minute)
	require.NoError(t, err)

	payload, err := JWTMaker.VerifyToken(token)
	require.EqualError(t, err, ErrTokenExpired.Error())
	require.Nil(t, payload)
}
