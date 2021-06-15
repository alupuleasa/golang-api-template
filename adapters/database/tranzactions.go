package database

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/efimovalex/wallet/adapters/model"
)

// UpdateWalletFunds - Updates the wallet funds with a sum
func (c *Client) UpdateWalletFunds(walletID uint64, sum float64) (w *model.Wallet, err error) {
	tx, err := c.db.Begin()
	if err != nil {
		return nil, err
	}
	w = &model.Wallet{}
	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select("idx,funds,owner_id,uuid").From("wallet").Where(sq.Eq{"idx": walletID}).RunWith(tx)
	err = query.QueryRow().Scan(&w.IDx, &w.Funds, &w.OwnerAccountID, &w.UUID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	w.Funds += sum

	if err := w.Validate(); err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = sq.Update("wallet").Set("funds", w.Funds).Where(sq.Eq{"idx": walletID}).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).Exec()
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
