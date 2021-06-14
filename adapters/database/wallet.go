package database

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/efimovalex/wallet/adapters/model"
	"github.com/google/uuid"
)

// CreateWallet - creates a wallet in the database
func (c *Client) CreateWallet(OwnerAccountID uint64) (w *model.Wallet, err error) {
	w = &model.Wallet{OwnerAccountID: OwnerAccountID}
	tx, err := c.db.Begin()
	if err != nil {
		return nil, err
	}

	query := sq.Insert("wallet").
		Columns("uuid", "funds", "owner").
		Values(uuid.NewString(), 0., OwnerAccountID).
		Suffix("RETURNING \"id\"").
		RunWith(tx).
		PlaceholderFormat(sq.Dollar)

	err = query.QueryRow().Scan(&w.UUID)
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
func (c *Client) GetWallet(ID uint64) (wallets []model.Wallet, err error) {
	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select("uuid,id,funds,owner_id").From("wallet").Where(sq.Eq{"id": ID}).
		RunWith(c.db).
		PlaceholderFormat(sq.Dollar)

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	err = rows.Scan(&wallets)
	if err != nil {
		return nil, err
	}

	return wallets, nil
}

// GetWallets - retrieves all wallets
func (c *Client) GetWallets(limit, offset uint64) (wallets []model.Wallet, err error) {
	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select("uuid,id,funds,owner_id").From("wallet").
		RunWith(c.db).
		PlaceholderFormat(sq.Dollar)

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}

	err = rows.Scan(&wallets)
	if err != nil {
		return nil, err
	}

	return wallets, nil
}
