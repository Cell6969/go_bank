package db

import (
	"context"
	"database/sql"
)

// VerifyEmailTxParams contains input parameters of transfer transaction
type VerifyEmailTxParams struct {
	EmailId    int64
	SecretCode string
}

// VerifyEmailTxResult contains result of VerifyEmailTx
type VerifyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

// TransfersTx performs a money transfer from one account into another account
// it creates a transfer record, add account entries, and update accounts balance within a single database transaction
func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		//  Update verify email table
		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}

		// Update user table
		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			Username: result.VerifyEmail.Username,
			IsEmailVerified: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		})

		return err
	})

	return result, err
}
