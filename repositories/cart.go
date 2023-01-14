package repositories

import (
	"waysbook/models"

	"gorm.io/gorm"
)

type CartRepository interface {
	CreateCart(cart models.Cart) (models.Cart, error)
	GetCart(ID int) (models.Cart, error)
	DeleteCart(cart models.Cart) (models.Cart, error)
	GetBookCart(ID int) (models.Book, error)
	GetTransactionID(ID int) (models.Transaction, error)
	GetCartByTransID(TransactionID int) ([]models.Cart, error)
	CreateTransaction(transaction models.Transaction) (models.Transaction, error)
}

func RepositoryCart(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateCart(cart models.Cart) (models.Cart, error) {
	err := r.db.Create(&cart).Error
	return cart, err
}

func (r *repository) GetCart(ID int) (models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("Book").First(&cart, ID).Error
	return cart, err
}

func (r *repository) DeleteCart(cart models.Cart) (models.Cart, error) {
	err := r.db.Delete(&cart).Error
	return cart, err
}

func (r *repository) GetBookCart(ID int) (models.Book, error) {
	var book models.Book
	err := r.db.First(&book, ID).Error
	return book, err
}

func (r *repository) GetTransactionID(userID int) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Find(&transaction, "status = ? AND user_id = ?", "waiting", userID).Error
	return transaction, err
}

func (r *repository) GetCartByTransID(TransactionID int) ([]models.Cart, error) {
	var cart []models.Cart
	err := r.db.Preload("Book").Find(&cart, "transaction_id = ?", TransactionID).Error
	return cart, err
}

func (r *repository) CreateTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Preload("User").Create(&transaction).Error
	return transaction, err
}
