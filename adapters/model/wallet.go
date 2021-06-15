package model

import "fmt"

// Wallet - wallet database model
// UUID - use UUID in order to make it safe from incremental data ids and data leaks
type Wallet struct {
	UUID           string  `db:"uuid" json:"-"`
	IDx            uint64  `db:"idx" json:"idx"`
	Funds          float64 `db:"funds" json:"funds"`
	OwnerAccountID uint64  `db:"owner_id" json:"owner_account_id"`
}

// Validate - Validates the wallet object
func (w *Wallet) Validate() error {
	if w.Funds < 0 {
		return fmt.Errorf("the balance of a wallet can not be negative")
	}

	return nil
}
