package pop

import "github.com/jmoiron/sqlx"

type dB struct {
	*sqlx.DB
}

func (db *dB) Transaction() (*tX, error) {
	return newTX(db)
}

func (db *dB) Rollback() error {
	return nil
}

func (db *dB) Commit() error {
	return nil
}
