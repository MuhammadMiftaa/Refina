package repository

import "gorm.io/gorm"

type Transaction interface {
	Commit() error
	Rollback() error
}

type TxManager interface {
	Begin() (Transaction, error)
}

type GormTx struct {
	db *gorm.DB
}

type GormTxManager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) TxManager {
	return &GormTxManager{db: db}
}

func (txm *GormTx) Commit() error {
	return txm.db.Commit().Error
}

func (txm *GormTx) Rollback() error {
	return txm.db.Rollback().Error
}

func (txm *GormTxManager) Begin() (Transaction, error) {
	tx := txm.db.Begin()
	return &GormTx{db: tx}, tx.Error
}
