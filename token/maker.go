package token

import "time"

type Maker interface {
	CreatToken(username string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
