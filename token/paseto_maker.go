package token

import (
	"fmt"

	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeyLen {
		return nil, fmt.Errorf("secretKet len is less than the min value %d", minSecretKeyLen)
	}

	return &PasetoMaker{paseto.NewV2(), []byte(secretKey)}, nil
}

func (pasetoMaker *PasetoMaker) CreateToken(username string) (string, *Payload, error) {
	payload, err := NewPayload(username, AccessTokenExpiration)
	if err != nil {
		return "", payload, err
	}

	token, err := pasetoMaker.paseto.Encrypt(pasetoMaker.symmetricKey, payload, nil)
	return token, payload, err
}

func (pasetoMaker *PasetoMaker) CreateRefreshToken(username string) (string, *Payload, error) {
	payload, err := NewPayload(username, RefreshTokenExpiration)
	if err != nil {
		return "", payload, err
	}

	token, err := pasetoMaker.paseto.Encrypt(pasetoMaker.symmetricKey, payload, nil)
	return token, payload, err
}

func (pasetoMaker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := pasetoMaker.paseto.Decrypt(token, pasetoMaker.symmetricKey, payload, nil)
	if err != nil {
		return nil, err
	}

	if err = payload.Valid(); err != nil {
		return nil, err
	}

	return payload, nil
}
