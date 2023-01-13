package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	booksdto "waysbook/dto/books"
	dto "waysbook/dto/result"
	"waysbook/models"
	"waysbook/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

var path_file_books = "http://localhost:5000/uploads/"

type handlerBook struct {
	BookRepository repositories.BookRepository
}

func HanlderBook(BookRepository repositories.BookRepository) *handlerBook {
	return &handlerBook{BookRepository}
}

// function get all book
func (h *handlerBook) FindBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	books, err := h.BookRepository.FindBooks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	for i, data := range books {
		books[i].Thumbnail = path_file_books + data.Thumbnail
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: books}
	json.NewEncoder(w).Encode(response)
}

// function get book by id
func (h *handlerBook) GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	book, err := h.BookRepository.GetBook(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	book.Thumbnail = path_file_books + book.Thumbnail

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: book}
	json.NewEncoder(w).Encode(response)
}

// function get book by id
func (h *handlerBook) CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get pdf name for book attachment
	dataPdfContext := r.Context().Value("dataPdf")
	pdfName := dataPdfContext.(string)

	// get image name for thumbnail
	dataImageContext := r.Context().Value("dataFile")
	imageName := dataImageContext.(string)

	//parse data
	isbn, _ := strconv.Atoi(r.FormValue("isbn"))
	pages, _ := strconv.Atoi(r.FormValue("pages"))
	price, _ := strconv.Atoi(r.FormValue("price"))

	// struct createTripRequest (dto) untuk menampung data
	request := booksdto.CreateBookRequest{
		Title:           r.FormValue("title"),
		PublicationDate: r.FormValue("publication_date"),
		ISBN:            isbn,
		Pages:           pages,
		Author:          r.FormValue("author"),
		Description:     r.FormValue("description"),
		Price:           price,
		Thumbnail:       imageName,
	}

	fmt.Println("request :", request)

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	publicationDate, _ := time.Parse("2006-01-02", r.FormValue("publication_date"))

	filePath := os.Getenv("FILE_PATH")

	// Get user id from token
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	book := models.Book{
		Title:              request.Title,
		PublicationDate:    publicationDate,
		ISBN:               int(isbn),
		Pages:              int(pages),
		Author:             request.Author,
		Description:        request.Description,
		Price:              int(price),
		IsPromo:            false,
		Discount:           0,
		PriceAfterDiscount: 0,
		BookAttachment:     filePath + pdfName,
		Thumbnail:          filePath + imageName,
		UserId:             userId,
	}

	fmt.Println("book :", book)

	dataBook, err := h.BookRepository.CreateBook(book)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	bookResponse, err := h.BookRepository.GetBook(dataBook.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseBook(bookResponse)}
	json.NewEncoder(w).Encode(response)
}

// function update book
func (h *handlerBook) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get pdf name for book attachment
	dataPdfContext := r.Context().Value("dataPdf")
	pdfName := dataPdfContext.(string)

	// get image name for thumbnail
	dataImageContext := r.Context().Value("dataFile")
	imageName := dataImageContext.(string)

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	book, err := h.BookRepository.GetBook(int(id))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	request := booksdto.UpdateBookRequest{
		BookAttachment: pdfName,
		Thumbnail:      imageName,
	}

	// title
	if r.FormValue("title") != "" {
		book.Title = r.FormValue("title")
	}

	// publication date
	date, _ := time.Parse("2006-01-02", r.FormValue("publication_date"))
	time := time.Now()
	if date != time {
		book.PublicationDate = date
	}

	// ISBN
	isbn, _ := strconv.Atoi(r.FormValue("ISBN"))
	if isbn != 0 {
		book.ISBN = isbn
	}

	// page
	pages, _ := strconv.Atoi(r.FormValue("pages"))
	if pages != 0 {
		book.Pages = pages
	}

	// author
	if r.FormValue("author") != "" {
		book.Author = r.FormValue("author")
	}

	// price
	price, _ := strconv.Atoi(r.FormValue("price"))
	if price != 0 {
		book.Price = price
	}

	// description
	if r.FormValue("description") != "" {
		book.Description = r.FormValue("description")
	}

	// discount
	discount, _ := strconv.Atoi(r.FormValue("discount"))
	if discount != 0 {
		book.Discount = discount
	}

	// book attachment
	if request.BookAttachment != "" {
		book.BookAttachment = request.BookAttachment
	}

	// thumbnail image
	if request.Thumbnail != "" {
		book.Thumbnail = request.Thumbnail
	}

	newBook, err := h.BookRepository.UpdateBook(book.Id, discount)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	newBookResponse, err := h.BookRepository.GetBook(newBook.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: newBookResponse}
	json.NewEncoder(w).Encode(response)
}

// function delete book
func (h *handlerBook) DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	book, err := h.BookRepository.GetBook(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := h.BookRepository.DeleteBook(book)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: data}
	json.NewEncoder(w).Encode(response)
}

// func update book promo
// func (h *handlerBook) UpdateBookPromo(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	request := booksdto.UpdateBookPromoRequest{
// 		Id:       r.FormValue("book_id"),
// 		Discount: r.FormValue("discount"),
// 	}

// 	bookId, _ := strconv.Atoi(request.Id)
// 	discount, _ := strconv.Atoi(request.Discount)

// 	book, err := h.BookRepository.UpdateBookPromo(bookId, discount)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
// 		json.NewEncoder(w).Encode(response)
// 		return
// 	}

// 	filePath := os.Getenv("FILE_PATH")

// 	bookResponse := booksdto.BookResponse{
// 		Id:                 book.Id,
// 		Title:              book.Title,
// 		PublicationDate:    book.PublicationDate,
// 		Pages:              book.Pages,
// 		ISBN:               book.ISBN,
// 		Price:              book.Price,
// 		IsPromo:            book.IsPromo,
// 		Discount:           book.Discount,
// 		PriceAfterDiscount: book.PriceAfterDiscount,
// 		Description:        book.Description,
// 		BookAttachment:     filePath + book.BookAttachment,
// 		Thumbnail:          filePath + book.Thumbnail,
// 		Author:             book.Author,
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	response := dto.SuccessResult{Code: http.StatusOK, Data: bookResponse}
// 	json.NewEncoder(w).Encode(response)
// }

// func (h *handlerBook) GetBooksByPromo(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	books, err := h.BookRepository.GetBooksByPromo()
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
// 		json.NewEncoder(w).Encode(response)
// 		return
// 	}

// 	filePath := os.Getenv("FILE_PATH")

// 	bookResponse := make([]booksdto.BookResponse, 0)
// 	for _, book := range books {
// 		bookResponse = append(bookResponse, booksdto.BookResponse{
// 			Id:                 book.Id,
// 			Title:              book.Title,
// 			PublicationDate:    book.PublicationDate,
// 			Pages:              book.Pages,
// 			ISBN:               book.ISBN,
// 			Price:              book.Price,
// 			IsPromo:            book.IsPromo,
// 			Discount:           book.Discount,
// 			PriceAfterDiscount: book.PriceAfterDiscount,
// 			Description:        book.Description,
// 			BookAttachment:     filePath + book.BookAttachment,
// 			Thumbnail:          filePath + book.Thumbnail,
// 			Author:             book.Author,
// 		})
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	response := dto.SuccessResult{Code: http.StatusOK, Data: bookResponse}
// 	json.NewEncoder(w).Encode(response)
// }

func convertResponseBook(u models.Book) booksdto.BookResponse {
	return booksdto.BookResponse{
		Id:              u.Id,
		Title:           u.Title,
		PublicationDate: u.PublicationDate.Format("2 January 2006"),
		ISBN:            u.ISBN,
		Pages:           u.Pages,
		Author:          u.Author,
		Price:           u.Price,
		Description:     u.Description,
		BookAttachment:  u.BookAttachment,
		Thumbnail:       u.Thumbnail,
	}
}
