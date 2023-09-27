package transaction

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/learn-go/web/module/daos"
	"github.com/learn-go/web/pkg/dbutils"
)

//go:generate mockgen -source=dao.go -destination=../mock/mock_dao/mock_dao.go -package=mock_dao

type Transaction interface {
	WithTransaction(txFunc func(Transaction) error) error
	GetHealth() error
	daos.Locker
}

// Transaction DAO for database access
type dao struct {
	db *sql.DB
	tx *sql.Tx
	daos.Locker
}

// NewTransaction creates a Transaction object
func NewTransaction(db *sql.DB, tx *sql.Tx) (*dao, error) {
	var querier dbutils.Querier
	if tx != nil {
		querier = tx
	} else if db != nil {
		querier = db
	} else {
		return nil, errors.New("failed to create new DAO")
	}
	daoLocker, _ := daos.NewDaoLocker(querier)
	return &dao{
		db:     db,
		tx:     tx,
		Locker: daoLocker,
	}, nil
}

func (d *dao) GetHealth() error {
	return d.db.Ping()
}

func (d *dao) WithTransaction(txFunc func(Transaction) error) error {
	var txTransaction *dao
	if d.tx != nil {
		// already in transaction
		return txFunc(d)
	}

	// create a new transaction
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
			_ = tx.Rollback()
			panic(r)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	// new DAO for this transaction
	txTransaction, err = NewTransaction(d.db, tx)
	if err != nil {
		return err
	}

	err = txFunc(txTransaction)
	return err
}
