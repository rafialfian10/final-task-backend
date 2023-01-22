package repositories

import (
	"waysbook/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindTransactions() ([]models.Transaction, error)
	FindTransactionsByUser(UserId int) ([]models.Transaction, error)
	GetTransaction(Id string) (models.Transaction, error)
	CreateTransaction(newTransaction models.Transaction) (models.Transaction, error)
	UpdateTransaction(status string, trxId string) (models.Transaction, error)
	UpdateTokenTransaction(token string, trxId string) (models.Transaction, error)
}

func RepositoryTransaction(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindTransactions() ([]models.Transaction, error) {
	var transaction []models.Transaction
	err := r.db.Preload("User").Preload("Cart").Preload("Cart.Book").Order("order_date desc").Find(&transaction).Error

	return transaction, err
}

func (r *repository) FindTransactionsByUser(UserId int) ([]models.Transaction, error) {
	var transaction []models.Transaction
	err := r.db.Preload("User").Preload("Cart").Preload("Cart.Book").Where("user_id = ?", UserId).Order("order_date desc").Find(&transaction).Error

	return transaction, err
}

func (r *repository) GetTransaction(Id string) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("User").Preload("Cart").Preload("Cart.Book").First(&transaction, "id = ?", Id).Error

	return transaction, err
}

func (r *repository) CreateTransaction(newTransaction models.Transaction) (models.Transaction, error) {
	err := r.db.Create(&newTransaction).Error

	return newTransaction, err
}

func (r *repository) UpdateTransaction(status string, trxId string) (models.Transaction, error) {
	var transaction models.Transaction
	r.db.Preload("User").Preload("Cart").Preload("Cart.Book").First(&transaction, "id = ?", trxId)

	// If is different & Status is "success" decrement available quota on data trip
	if status != transaction.Status && status == "success" {
		for _, ordr := range transaction.Cart {
			var book models.Book
			r.db.First(&book, ordr.Book.Id)
			book.Quota = book.Quota - ordr.OrderQty
			r.db.Model(&book).Updates(book)
		}
	}

	// If is different & Status is "reject" decrement available quota on data trip
	if status != transaction.Status && status == "rejected" {
		for _, ordr := range transaction.Cart {
			var book models.Book
			r.db.First(&book, ordr.BookId)
			book.Quota = book.Quota + ordr.OrderQty
			r.db.Model(&book).Updates(book)
		}
	}

	// change transaction status
	transaction.Status = status

	err := r.db.Model(&transaction).Updates(transaction).Error

	return transaction, err
}

func (r *repository) UpdateTokenTransaction(token string, trxId string) (models.Transaction, error) {
	var transaction models.Transaction
	r.db.Preload("User").Preload("Cart").Preload("Cart.Book").First(&transaction, "id = ?", trxId)

	// change transaction token
	transaction.MidtransId = token
	err := r.db.Model(&transaction).Updates(transaction).Error

	return transaction, err
}

func (r *repository) DeleteTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Delete(&transaction).Error

	return transaction, err
}
