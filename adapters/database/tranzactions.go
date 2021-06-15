package database

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/efimovalex/wallet/adapters/model"
)

// UpdateWalletFunds - Updates the wallet funds with a sum
func (c *Client) UpdateWalletFunds(walletID uint64, sum float64, ref string) (w *model.Wallet, transactionIDx uint64, err error) {
	tx, err := c.db.Begin()
	if err != nil {
		return nil, 0, err
	}
	w = &model.Wallet{}
	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select("idx,funds,owner_id,uuid").From("wallet").Where(sq.Eq{"idx": walletID}).RunWith(tx)
	err = query.QueryRow().Scan(&w.IDx, &w.Funds, &w.OwnerAccountID, &w.UUID)
	if err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	w.Funds += sum

	if err := w.Validate(); err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	_, err = sq.Update("wallet").Set("funds", w.Funds).Where(sq.Eq{"idx": walletID}).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).Exec()
	if err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	tq := sq.Insert("transactions").
		Columns("walled_id", "sum", "reference").
		Values(w.IDx, sum, ref).
		Suffix("RETURNING \"idx\"").
		RunWith(tx).
		PlaceholderFormat(sq.Dollar)

	err = tq.QueryRow().Scan(transactionIDx)
	if err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	if err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, 0, err
	}

	return w, transactionIDx, nil
}
