package handlers

import db "github.com/escalopa/gobank/db/sqlc"

func fromUserToUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		CreatedAt:         user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
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
