package handlers

import db "github.com/escalopa/gobank/db/sqlc"

func mapUserToResponse(user *db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		CreatedAt:         user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
	}
}

func mapAccountToResponse(account db.Account) *accountResponse {
	return &accountResponse{
		ID:        account.ID,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
	}
}

func mapTransferToResponse(transfer db.Transfer) *transferResponse {
	return &transferResponse{
		ID: transfer.ID,
		// FromAccount: transfer.FromAccountID,
		ToAccountID: transfer.ToAccountID,
		// FromEntry:   transfer.FromEntryID,
		Amount:    transfer.Amount,
		CreatedAt: transfer.CreatedAt,
	}
}

func fromTransferTxToTransferResponse(result db.TransferTxResult) transferResponse {
	return transferResponse{
		ID:          result.Transfer.ID,
		FromAccount: result.FromAccount,
		ToAccountID: result.ToAccount.ID,
		FromEntry:   result.FromEntry,
		Amount:      result.Transfer.Amount,
	}
}
