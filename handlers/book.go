package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	bookdto "waysbook/dto/books"
	booksdto "waysbook/dto/books"
	dto "waysbook/dto/result"
	"waysbook/models"
	"waysbook/repositories"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// var path_file = "http://localhost:4000/uploads/"

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

	for i, p := range books {
		imagePath := os.Getenv("PATH_FILE") + p.Thumbnail
		books[i].Thumbnail = imagePath
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: books}
	json.NewEncoder(w).Encode(response)
}

// function get all book promo
func (h *handlerBook) FindBookPromo(w http.ResponseWriter, r *http.Request) {
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
	dataImageContex := r.Context().Value("dataFile")
	filepath := dataImageContex.(string)

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

	resp, err := cld.Upload.Upload(ctx, filepath, uploader.UploadParams{Folder: "waysbook"})
	resp1, err := cld1.Upload.Upload(ctx, filePDF, uploader.UploadParams{Folder: "waysbook"})

	if err != nil {
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
		BookAttachment:     resp1.SecureURL,
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
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseBook(bookResponse)}
	json.NewEncoder(w).Encode(response)
}

// function update book
func (h *handlerBook) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// get pdf name for book attachment
	dataPDF := r.Context().Value("dataPDF")
	filePDF := dataPDF.(string)

	// get image name for thumbnail
	dataImageContex := r.Context().Value("dataFile")
	filepath := dataImageContex.(string)

	// cloudinary
	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	// Add your Cloudinary credentials ...
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)
	cld1, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	resp, err := cld.Upload.Upload(ctx, filepath, uploader.UploadParams{Folder: "waysbook"})
	resp1, err := cld1.Upload.Upload(ctx, filePDF, uploader.UploadParams{Folder: "waysbook"})

	if err != nil {
		fmt.Println(err.Error())
	}

	//parse data
	isbn, _ := strconv.Atoi(r.FormValue("isbn"))
	pages, _ := strconv.Atoi(r.FormValue("pages"))
	price, _ := strconv.Atoi(r.FormValue("price"))
	quota, _ := strconv.Atoi(r.FormValue("quota"))

	request := bookdto.UpdateBookRequest{
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
	err = validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	book, _ := h.BookRepository.GetBook(id)

	if request.Title != "" {
		book.Title = request.Title
	}

	publicDate, _ := time.Parse("2006-01-02", request.PublicationDate)
	if request.PublicationDate != "" {
		book.PublicationDate = publicDate
	}

	if isbn != 0 {
		book.ISBN = isbn
	}

	if request.Pages != 0 {
		book.Pages = request.Pages
	}

	if request.Author != "" {
		book.Author = request.Author
	}

	if request.Quota != 0 {
		book.Quota = request.Quota
	}

	if request.Description != "" {
		book.Description = request.Description
	}

	if request.Price != 0 {
		book.Price = request.Price
	}

	if filepath != "false" {
		book.Thumbnail = resp.SecureURL
	}

	if filePDF != "false" {
		book.BookAttachment = resp1.SecureURL
	}

	book.UpdateAt = time.Now()

	book, _ = h.BookRepository.UpdateBook(book)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: book}
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
		BookAttachment:  u.BookAttachment,
		Thumbnail:       u.Thumbnail,
	}
}
