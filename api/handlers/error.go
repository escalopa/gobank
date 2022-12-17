package handlers

import (
	"errors"
	"fmt"
)

var (
	// Account
	ErrNotAccountOwner = errors.New("account doesn't belong to authenticated user")

	// User
	ErrEmailSameAsOld = errors.New("new email is the same as old email")

	// Refresh token errors
	ErrMismatchedTokenOwner    = errors.New("refresh token doesn't belong to authenticated user")
	ErrInvalidRefreshToken     = errors.New("invalid refresh token")
	ErrBlockedRefreshToken     = errors.New("refresh token is blocked")
	ErrMismatchedRefreshTokens = errors.New("refresh token doesn't match with stored refresh token")
	ErrExpiredRefreshToken     = errors.New("refresh token has expired")

	ErrSameAccountTransfer = func(from, to int64) error {
		return fmt.Errorf(fmt.Sprintf("can't transfer to the same account, req.FromAccountId=%d, req.ToAccount=%d",
			from, to,
		))
	}

	ErrCurrencyMismatch = func(from, to string) error {
		return fmt.Errorf(fmt.Sprintf("currency mismatch account1.currency=%s, account2.currency=%s",
			from, to,
		))
	}
)
