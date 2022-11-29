package api

import (
	"errors"
	"fmt"
)

var (
	ErrNotAccountOwner = errors.New("account doesn't belong to authenticated user")

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
