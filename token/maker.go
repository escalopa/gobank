package token

import "time"

const (
	AccessTokenExpiration  = time.Hour
	RefreshTokenExpiration = time.Hour * 24 * 7
)

type Maker interface {
	CreateToken(username string) (string, *Payload, error)
	CreateRefreshToken(username string) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
