package database

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/efimovalex/wallet/adapters/model"
	"github.com/google/uuid"
)

// CreateWallet - creates a wallet in the database
func (c *Client) CreateWallet(ownerAccountID uint64) (w *model.Wallet, err error) {
	w = &model.Wallet{
		UUID:           uuid.NewString(),
		OwnerAccountID: ownerAccountID,
		Funds:          0.,
	}

	// validate wallet
	if w.Validate() != nil {
		return nil, err
	}

	tx, err := c.db.Begin()
	if err != nil {
		return nil, err
	}

	query := sq.Insert("wallet").
		Columns("uuid", "funds", "owner_id").
		Values(w.UUID, w.Funds, w.OwnerAccountID).
		Suffix("RETURNING \"idx\"").
		RunWith(tx).
		PlaceholderFormat(sq.Dollar)

	err = query.QueryRow().Scan(&w.IDx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return w, nil
}

// GetWallet - retrieves wallet by Id
func (c *Client) GetWallet(ID uint64) (wallet model.Wallet, err error) {
	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select("uuid, idx,funds,owner_id").From("wallet").Where(sq.Eq{"idx": ID}).
		RunWith(c.db).
		PlaceholderFormat(sq.Dollar)

	err = query.QueryRow().Scan(&wallet.UUID, &wallet.IDx, &wallet.Funds, &wallet.OwnerAccountID)
	if err != nil {
		return wallet, err
	}

	return wallet, nil
}

// GetWallets - retrieves all wallets
func (c *Client) GetWallets(limit, offset uint64) (wallets []model.Wallet, err error) {
	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select("uuid,idx,funds,owner_id").From("wallet").
		RunWith(c.db).
		PlaceholderFormat(sq.Dollar)

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var wallet model.Wallet
		err = rows.Scan(&wallet.UUID, &wallet.IDx, &wallet.Funds, &wallet.OwnerAccountID)
		if err != nil {
			return nil, err
		}

		wallets = append(wallets, wallet)
	}

	return wallets, nil
}
