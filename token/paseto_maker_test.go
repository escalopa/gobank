package token

import (
	"testing"
	"time"

	"github.com/escalopa/gobank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	pasetoMaker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	issuedAt := time.Now()

	token, payload, err := pasetoMaker.CreateToken(username)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = pasetoMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)

	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
}

// func TestPasetoMakerExpired(t *testing.T) {
// 	JWTMaker, err := NewPasetoMaker(util.RandomString(32))
// 	require.NoError(t, err)

// 	username := util.RandomOwner()

// 	token, payload, err := JWTMaker.CreateToken(username)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, payload)

// 	payload, err = JWTMaker.VerifyToken(token)
// 	require.EqualError(t, err, ErrTokenExpired.Error())
// 	require.Nil(t, payload)
// }
