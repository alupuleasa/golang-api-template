package model

import "fmt"

// Transaction - a func or withdrawal event
type Transaction struct {
	IDx       uint64
	WalletID  uint64
	Reference string
	Sum       float64
}

func (t *Transaction) Validate() error {
	if len(t.Reference) == 0 {
		return fmt.Errorf("Transaction reference cannot be empty")
	}
	if t.WalletID == 0 {
		return fmt.Errorf("Transaction wallet id reference cannot be empty")
	}
	if t.Sum == 0 {
		return fmt.Errorf("Transaction sum cannot be 0")
	}

	return nil
}
