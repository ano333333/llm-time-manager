package store

import "database/sql"

type Transaction interface {
	Commit() error
	Rollback() error
}

type TransactionStore interface {
	Begin() (Transaction, error)
}

type DefaultTransactionStore struct {
	DB *sql.DB
}

type DefaultTransaction struct {
	Tx *sql.Tx
}

func (s *DefaultTransactionStore) Begin() (Transaction, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	return DefaultTransaction{Tx: tx}, nil
}

func (t DefaultTransaction) Commit() error {
	return t.Tx.Commit()
}

func (t DefaultTransaction) Rollback() error {
	return t.Tx.Rollback()
}
