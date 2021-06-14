package model

// Wallet - wallet database model
// UUID - use UUID in order to make it safe from incremental data ids and data leaks
type Wallet struct {
	UUID           string
	IDx            uint64
	Funds          float64
	OwnerAccountID uint64
}
