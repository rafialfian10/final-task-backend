package repositories

import (
	"waysbook/models"

	"gorm.io/gorm"
)

type CartRepository interface {
	FindCarts(UserId int) ([]models.Cart, error)
	GetCart(Id int) (models.Cart, error)
	GetCartByBook(BookId int, UserId int) (models.Cart, error)
	CreateCart(newOrder models.Cart) (models.Cart, error)
	UpdateCart(cart models.Cart) (models.Cart, error)
	DeleteCart(cart models.Cart) (models.Cart, error)
}

func RepositoryCart(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindCarts(UserId int) ([]models.Cart, error) {
	var cart []models.Cart
	err := r.db.Preload("Book").Where("user_id = ?", UserId).Where("transaction_id IS NULL").Find(&cart).Error
	return cart, err
}

func (r *repository) GetCart(Id int) (models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("Book").Where("transaction_id IS NULL").First(&cart, "id = ?", Id).Error
	return cart, err
}

func (r *repository) GetCartByBook(bookId int, userId int) (models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("Book").Where("user_id = ?", userId).Where("transaction_id IS NULL").First(&cart, "product_id = ?", bookId).Error
	return cart, err
}

func (r *repository) CreateCart(cart models.Cart) (models.Cart, error) {
	err := r.db.Select("BookId", "OrderQty", "UserId").Create(&cart).Error
	return cart, err
}

func (r *repository) UpdateCart(cart models.Cart) (models.Cart, error) {
	err := r.db.Model(&cart).Updates(cart).Error
	return cart, err
}

func (r *repository) DeleteCart(cart models.Cart) (models.Cart, error) {
	err := r.db.Delete(&cart).Error
	return cart, err
}
