package database

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

// AddFunds - Adds a sum to the wallet
func (c *Client) AddFunds(walletID uint64, sum float64) (err error) {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	var funds float64
	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select("funds").From("wallet").Where(sq.Eq{"id": walletID})
	err = query.QueryRow().Scan(&funds)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = sq.Update("wallet").Set("funds", funds+sum).Where(sq.Eq{"id": walletID}).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).Exec()

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// RemoveFunds - Removes a sum from the wallet
func (c *Client) RemoveFunds(walletID uint64, sum float64) (err error) {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	var funds float64

	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select("funds").From("wallet").Where(sq.Eq{"id": walletID})
	err = query.QueryRow().Scan(&funds)
	if err != nil {
		tx.Rollback()
		return err
	}

	if funds-sum < 0 {
		tx.Rollback()
		return fmt.Errorf("the balance of a wallet can not be negative, funds: %.2f, sum to be subtracted: %.2f", funds, sum)
	}

	_, err = sq.Update("wallet").Set("funds", funds).Where(sq.Eq{"id": walletID}).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).Exec()

	return tx.Commit()
}
