package repositories

import (
	"waysbook/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindTransactions() ([]models.Transaction, error)
	FindTransactionsByUser(Id int) ([]models.Transaction, error)
	GetTransaction(Id string) (models.Transaction, error)
	CreateTransaction(newTransaction models.Transaction) (models.Transaction, error)
	UpdateTransaction(status string, Id string) (models.Transaction, error)
	UpdateTokenTransaction(token string, Id string) (models.Transaction, error)
}

func RepositoryTransaction(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindTransactions() ([]models.Transaction, error) {
	var transaction []models.Transaction
	err := r.db.Preload("User").Preload("Cart").Preload("Cart.Book").Order("order_date desc").Find(&transaction).Error

	return transaction, err
}

func (r *repository) FindTransactionsByUser(Id int) ([]models.Transaction, error) {
	var transaction []models.Transaction
	err := r.db.Preload("User").Preload("Cart").Preload("Cart.Book").Where("user_id = ?", Id).Order("order_date desc").Find(&transaction).Error

	return transaction, err
}

func (r *repository) GetTransaction(Id string) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("User").Preload("Cart").Preload("Cart.Book").First(&transaction, "id = ?", Id).Error

	return transaction, err
}

func (r *repository) CreateTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Debug().Create(&transaction).Error

	return transaction, err
}

func (r *repository) UpdateTransaction(status string, Id string) (models.Transaction, error) {
	var transaction models.Transaction
	r.db.Preload("User").Preload("Cart").Preload("Cart.Book").First(&transaction, "id = ?", Id)

	// jika status dan transaksi status berbeda & Status adalah "reject" maka quota trip akan dikurangi
	if status != transaction.Status && status == "success" {
		for _, ordr := range transaction.Cart {
			var book models.Book
			r.db.First(&book, ordr.Book.Id)
			book.Quota = book.Quota - ordr.OrderQty
			r.db.Model(&book).Updates(book)
		}
	}

	// // jika status dan transaksi status berbeda & Status adalah "reject" maka quota book akan di tambahkan kembali
	if status != transaction.Status && status == "reject" {
		for _, ordr := range transaction.Cart {
			var book models.Book
			r.db.First(&book, ordr.BookId)
			book.Quota = book.Quota + ordr.OrderQty
			r.db.Model(&book).Updates(book)
		}
	}

	// ubah status transaksi
	transaction.Status = status

	err := r.db.Model(&transaction).Updates(transaction).Error

	return transaction, err
}

func (r *repository) UpdateTokenTransaction(token string, Id string) (models.Transaction, error) {
	var transaction models.Transaction
	r.db.Preload("User").Preload("Cart").Preload("Cart.Book").First(&transaction, "id = ?", Id)

	// change transaction token
	transaction.MidtransId = token
	err := r.db.Model(&transaction).Updates(transaction).Error

	return transaction, err
}

func (r *repository) DeleteTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Delete(&transaction).Error

	return transaction, err
}
