package repository

import (
	"database/sql"
)

func useTX(conn *sql.DB, fn func(tx *sql.Tx) error) error {
	var tx *sql.Tx
	var err error

	tx, err = conn.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			err = tx.Commit()
		} else {
			err = tx.Rollback()
		}
	}()

	err = fn(tx)
	return err
}
