package repositories

import (
	"waysbook/models"

	"gorm.io/gorm"
)

type BookRepository interface {
	FindBooks() ([]models.Book, error)
	GetBook(ID int) (models.Book, error)
	FindBookPromo() ([]models.Book, error)
	CreateBook(book models.Book) (models.Book, error)
	UpdateBook(book models.Book) (models.Book, error)
	UpdateBookPromo(ID int, discount int) (models.Book, error)
	GetBooksByPromo() ([]models.Book, error)
	DeleteBook(book models.Book) (models.Book, error)
	DeleteBookThumbnail(ID int) (models.Book, error)
	DeleteBookDocument(ID int) (models.Book, error)
}

func RepositoryBook(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindBooks() ([]models.Book, error) {
	var books []models.Book
	err := r.db.Debug().Find(&books).Error

	return books, err
}

func (r *repository) FindBookPromo() ([]models.Book, error) {
	var books []models.Book
	err := r.db.Where("is_promo= ?", true).Find(&books).Error
	return books, err
}

func (r *repository) GetBook(Id int) (models.Book, error) {
	var book models.Book
	err := r.db.Debug().First(&book, "id=?", Id).Error

	return book, err
}

func (r *repository) CreateBook(book models.Book) (models.Book, error) {
	err := r.db.Debug().Create(&book).Error
	return book, err
}

func (r *repository) UpdateBook(book models.Book) (models.Book, error) {
	err := r.db.Debug().Model(&book).Updates(book).Error

	return book, err
}

func (r *repository) UpdateBookPromo(ID int, discount int) (models.Book, error) {
	var book models.Book
	r.db.First(&book, "id=?", ID)

	book.IsPromo = true
	book.Discount = discount

	// Calculate Price After Discount
	book.PriceAfterDiscount = book.Price - (book.Price * discount / 100)
	err := r.db.Model(&book).Updates(book).Error

	return book, err
}

func (r *repository) GetBooksByPromo() ([]models.Book, error) {
	var books []models.Book
	err := r.db.Preload("User").Find(&books, "is_promo=?", true).Error

	return books, err
}

func (r *repository) DeleteBook(book models.Book) (models.Book, error) {
	err := r.db.Debug().Delete(&book).Error

	return book, err
}

func (r *repository) DeleteBookThumbnail(ID int) (models.Book, error) {
	return models.Book{}, r.db.Model(&models.Book{}).Where("id = ?", ID).UpdateColumn("thumbnail", gorm.Expr("NULL")).Error
}

func (r *repository) DeleteBookDocument(ID int) (models.Book, error) {
	return models.Book{}, r.db.Model(&models.Book{}).Where("id = ?", ID).UpdateColumn("book", gorm.Expr("NULL")).Error
}
