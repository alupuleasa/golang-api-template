package model

// Transaction - a func or withdrawal event
type Transaction struct {
	WalletID  uint64
	Reference string
	Sum       float64
}
