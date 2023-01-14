package repositories

import (
	"waysbook/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	UpdateTransaction(transaction models.Transaction) (models.Transaction, error)
	GetTransaction(userID int) (models.Transaction, error)
	GetTransactionAdmin(ID int) (models.Transaction, error)
	UpdateTrans(status string, ID int) error
	GetOneTrans(ID string) (models.Transaction, error)
	GetOrderByID(userID int) ([]models.Transaction, error)
	FindTransaction() ([]models.Transaction, error)
}

func RepositoryTransaction(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) UpdateTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Preload("Cart").Preload("Cart.Book").Preload("User").Save(&transaction).Error
	return transaction, err
}

func (r *repository) GetTransaction(userID int) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("Cart").Preload("Cart.Book").Preload("User").Where("user_id = ? AND status =?", userID, "waiting").First(&transaction).Error
	return transaction, err
}

func (r *repository) GetTransactionAdmin(ID int) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("User").First(&transaction, ID).Error
	return transaction, err
}

func (r *repository) UpdateTrans(status string, ID int) error {
	var transaction models.Transaction
	r.db.Preload("User").Preload("Cart.Book").First(&transaction, ID)
	transaction.Status = status
	err := r.db.Debug().Save(&transaction).Error

	return err
}

func (r *repository) GetOneTrans(ID string) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("User").Preload("Cart").Preload("Cart.Book").First(&transaction, ID).Error
	return transaction, err
}

func (r *repository) FindTransaction() ([]models.Transaction, error) {
	var transaction []models.Transaction
	err := r.db.Preload("User").Preload("Cart").Preload("Cart.Book").Not("status = ?", "waiting").Find(&transaction).Error
	return transaction, err
}

func (r *repository) GetOrderByID(userID int) ([]models.Transaction, error) {
	var transaction []models.Transaction
	err := r.db.Preload("User").Preload("Cart").Preload("Cart.Book").Not("status = ?", "waiting").Where("user_id = ? AND status=?", userID, "success").Find(&transaction).Error
	return transaction, err
}
