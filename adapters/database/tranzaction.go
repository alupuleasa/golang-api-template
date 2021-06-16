package database

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/efimovalex/wallet/adapters/model"
)

// UpdateWalletFunds - Updates the wallet funds with a sum
func (c *Client) UpdateWalletFunds(walletID uint64, sum float64, ref string) (w *model.Wallet, t *model.Transaction, err error) {
	tx, err := c.db.Begin()
	if err != nil {
		return nil, nil, err
	}
	w = &model.Wallet{}
	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select("idx,funds,owner_id,uuid").From("wallet").Where(sq.Eq{"idx": walletID}).RunWith(tx)
	err = query.QueryRow().Scan(&w.IDx, &w.Funds, &w.OwnerAccountID, &w.UUID)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	w.Funds += sum

	if err := w.Validate(); err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	_, err = sq.Update("wallet").Set("funds", w.Funds).Where(sq.Eq{"idx": walletID}).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).Exec()
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	t = &model.Transaction{WalletID: w.IDx, Sum: sum, Reference: ref}
	if err = t.Validate(); err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	tq := sq.Insert("transaction").
		Columns("wallet_id", "sum", "reference").
		Values(t.WalletID, t.Sum, t.Reference).
		Suffix("RETURNING \"idx\"").
		RunWith(tx).
		PlaceholderFormat(sq.Dollar)

	err = tq.QueryRow().Scan(&t.IDx)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, nil, err
	}

	return w, t, nil
}

// UpdateTransaction - Updates the Transaction with the payment references to paypal
func (c *Client) UpdateTransaction(tID uint64, ref string) (t *model.Transaction, err error) {
	tx, err := c.db.Begin()
	if err != nil {
		return nil, err
	}

	t = &model.Transaction{}
	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select("idx,wallet_id,sum,reference").From("transaction").Where(sq.Eq{"idx": tID}).RunWith(tx)
	err = query.QueryRow().Scan(&t.IDx, &t.WalletID, &t.Sum, &t.Reference)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	t.Reference = ref

	if err = t.Validate(); err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = sq.Update("transaction").Set("reference", ref).Where(sq.Eq{"idx": tID}).
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

	return t, nil
}
