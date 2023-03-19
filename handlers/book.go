package handlers

import (
	"context"
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

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type handlerBook struct {
	BookRepository repositories.BookRepository
}

func HanlderBook(BookRepository repositories.BookRepository) *handlerBook {
	return &handlerBook{BookRepository}
}

// function get all books
func (h *handlerBook) FindBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	books, err := h.BookRepository.FindBooks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	for i, p := range books {
		imagePath := os.Getenv("PATH_FILE") + p.Thumbnail
		books[i].Thumbnail = imagePath
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: books}
	json.NewEncoder(w).Encode(response)
}

// function get all books promo
func (h *handlerBook) FindBooksPromo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	books, err := h.BookRepository.FindBookPromo()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: books}
	json.NewEncoder(w).Encode(response)
}

// function get detail book
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

	book.Thumbnail = os.Getenv("PATH_FILE") + book.Thumbnail

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: book}
	json.NewEncoder(w).Encode(response)
}

// function get book by id
func (h *handlerBook) CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get pdf name for book attachment
	dataPDF := r.Context().Value("dataPDF")
	filePDF := dataPDF.(string)

	// get image name for thumbnail
	dataImage := r.Context().Value("dataFile")
	fileImage := dataImage.(string)

	//parse data
	isbn, _ := strconv.Atoi(r.FormValue("isbn"))
	pages, _ := strconv.Atoi(r.FormValue("pages"))
	price, _ := strconv.Atoi(r.FormValue("price"))
	quota, _ := strconv.Atoi(r.FormValue("quota"))

	request := booksdto.CreateBookRequest{
		Title:           r.FormValue("title"),
		PublicationDate: r.FormValue("publication_date"),
		ISBN:            isbn,
		Pages:           pages,
		Author:          r.FormValue("author"),
		Quota:           quota,
		Description:     r.FormValue("description"),
		Price:           price,
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// cloudinary
	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	// Add your Cloudinary credentials ...
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)
	cld1, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	resp, err := cld.Upload.Upload(ctx, fileImage, uploader.UploadParams{Folder: "waysbook"})
	resp1, err1 := cld1.Upload.Upload(ctx, filePDF, uploader.UploadParams{Folder: "waysbook"})

	if err != nil {
		fmt.Println(err.Error())
	}

	if err1 != nil {
		fmt.Println(err.Error())
	}

	publicationDate, _ := time.Parse("2006-01-02", r.FormValue("publication_date"))

	book := models.Book{
		Title:              request.Title,
		PublicationDate:    publicationDate,
		ISBN:               int(isbn),
		Pages:              int(pages),
		Author:             request.Author,
		Quota:              int(quota),
		Description:        request.Description,
		Price:              int(price),
		IsPromo:            false,
		Discount:           0,
		PriceAfterDiscount: 0,
		Book:               resp1.SecureURL,
		Thumbnail:          resp.SecureURL,
	}

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
	response := dto.SuccessResult{Code: http.StatusOK, Data: bookResponse}
	json.NewEncoder(w).Encode(response)
}

// function update book
func (h *handlerBook) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	book, err := h.BookRepository.GetBook(int(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// get pdf name for book attachment
	dataPDF := r.Context().Value("dataPDF")
	filePDF := dataPDF.(string)

	// get image name for thumbnail
	dataImage := r.Context().Value("dataFile")
	fileImage := dataImage.(string)

	// cloudinary
	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	// Add your Cloudinary credentials ...
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)
	cld1, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	resp, err := cld.Upload.Upload(ctx, fileImage, uploader.UploadParams{Folder: "waysbook"})
	resp1, err1 := cld1.Upload.Upload(ctx, filePDF, uploader.UploadParams{Folder: "waysbook"})

	if err != nil {
		fmt.Println(err.Error())
	}

	if err1 != nil {
		fmt.Println(err.Error())
	}

	// title
	if r.FormValue("title") != "" {
		book.Title = r.FormValue("title")
	}

	// parse time
	date, _ := time.Parse("2006-01-02", r.FormValue("publication_date"))
	time := time.Now()
	if date != time {
		book.PublicationDate = date
	}

	// isbn
	isbn, _ := strconv.Atoi(r.FormValue("isbn"))
	if isbn != 0 {
		book.ISBN = isbn
	}

	// pages
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

	// quota
	quota, _ := strconv.Atoi(r.FormValue("quota"))
	if quota != 0 {
		book.Quota = quota
	}

	// description
	if r.FormValue("description") != "" {
		book.Description = r.FormValue("description")
	}

	// image
	if resp.SecureURL != "" {
		book.Thumbnail = resp.SecureURL
	}

	// book
	if resp1.SecureURL != "" {
		book.Book = resp1.SecureURL
	}

	newBook, _ := h.BookRepository.UpdateBook(book)

	// jika ada error maka tampilkan ErrorResult
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// panggil function getTrip agar setelah data di create data id akan keluar response
	newBookResponse, err := h.BookRepository.GetBook(newBook.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// jika tidak ada error maka SuccessResult
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseBook(newBookResponse)}
	json.NewEncoder(w).Encode(response)
}

// function update book promo
func (h *handlerBook) UpdateBookPromo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	book, err := h.BookRepository.GetBook(int(id))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// discount
	discount, _ := strconv.Atoi(r.FormValue("discount"))
	if discount != 0 {
		book.Discount = discount
	}

	newBook, err := h.BookRepository.UpdateBookPromo(book.Id, discount)

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

// function convert response book
func convertResponseBook(u models.Book) booksdto.BookResponse {
	return booksdto.BookResponse{
		Id:              u.Id,
		Title:           u.Title,
		PublicationDate: u.PublicationDate.Format("2 January 2006"),
		ISBN:            u.ISBN,
		Pages:           u.Pages,
		Author:          u.Author,
		Quota:           u.Quota,
		Price:           u.Price,
		Description:     u.Description,
		Book:            u.Book,
		Thumbnail:       u.Thumbnail,
	}
}
