package token

import (
	"testing"
	"time"

	"github.com/escalopa/gobank/util"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	JWTMaker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := JWTMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = JWTMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)

	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, payload.ExpireAt, expiredAt, time.Second)

}

func TestJWTMakerExpired(t *testing.T) {
	JWTMaker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()

	token, payload, err := JWTMaker.CreateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	payload, err = JWTMaker.VerifyToken(token)
	require.EqualError(t, err, ErrTokenExpired.Error())
	require.Nil(t, payload)
}

func TestJWTMakerInvalid(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	JWTMaker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = JWTMaker.VerifyToken(token)
	require.EqualError(t, err, ErrTokenInvalid.Error())
	require.Nil(t, payload)
}
