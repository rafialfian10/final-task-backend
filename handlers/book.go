package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"
	booksdto "waysbook/dto/book"
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

func (h *handlerBook) FindBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	books, err := h.BookRepository.FindBooks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	// for i, p := range books {
	// 	fileImage := os.Getenv("PATH_FILE") + p.Thumbnail
	// 	books[i].Thumbnail = fileImage
	// }

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: ConvertMultipleBookResponse(books)}
	json.NewEncoder(w).Encode(response)
}

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

	// book.Thumbnail = os.Getenv("PATH_FILE") + book.Thumbnail

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: ConvertBookResponse(book)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerBook) CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get pdf name for book attachment
	dataPdf := r.Context().Value("dataPdf")
	filePdf := dataPdf.(string)

	// get image name for thumbnail
	dataImage := r.Context().Value("dataImage")
	fileImage := dataImage.(string)

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

	respImage, _ := cld.Upload.Upload(ctx, fileImage, uploader.UploadParams{Folder: "waysbook"})
	respDocument, _ := cld1.Upload.Upload(ctx, filePdf, uploader.UploadParams{Folder: "waysbook"})

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
		Book:               respDocument.SecureURL,
		Thumbnail:          respImage.SecureURL,
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
	response := dto.SuccessResult{Code: http.StatusOK, Data: ConvertBookResponse(bookResponse)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerBook) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// get pdf name for book attachment
	dataPdf := r.Context().Value("dataPdf")
	if dataPdf == nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: "dataPdf is nil"}
		json.NewEncoder(w).Encode(response)
		return
	}

	filePdf, ok := dataPdf.(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: "dataPdf is not a string"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// get image name for thumbnail
	dataImage := r.Context().Value("dataImage")
	if dataImage == nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: "dataImage is nil"}
		json.NewEncoder(w).Encode(response)
		return
	}

	fileImage, ok := dataImage.(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: "dataImage is not a string"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// cloudinary
	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)
	cld1, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	respImage, _ := cld.Upload.Upload(ctx, fileImage, uploader.UploadParams{Folder: "waysbook"})
	respDocument, _ := cld1.Upload.Upload(ctx, filePdf, uploader.UploadParams{Folder: "waysbook"})

	var imageSecureURL string
	if respImage != nil && respImage.SecureURL != "" {
		imageSecureURL = respImage.SecureURL
	}

	var documentSecureURL string
	if respDocument != nil && respDocument.SecureURL != "" {
		documentSecureURL = respDocument.SecureURL
	}

	isbn, _ := strconv.Atoi(r.FormValue("isbn"))
	pages, _ := strconv.Atoi(r.FormValue("pages"))
	price, _ := strconv.Atoi(r.FormValue("price"))
	quota, _ := strconv.Atoi(r.FormValue("quota"))

	request := booksdto.UpdateBookRequest{
		Title:           r.FormValue("title"),
		PublicationDate: r.FormValue("publication_date"),
		ISBN:            isbn,
		Pages:           pages,
		Author:          r.FormValue("author"),
		Price:           price,
		Description:     r.FormValue("description"),
		Quota:           quota,
		Book:            documentSecureURL,
		Thumbnail:       imageSecureURL,
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	book, _ := h.BookRepository.GetBook(id)

	if request.Title != "" {
		book.Title = request.Title
	}

	// parse time
	date, _ := time.Parse("2006-01-02", r.FormValue("publication_date"))
	time := time.Now()
	if date != time {
		book.PublicationDate = date
	}

	if isbn != 0 {
		book.ISBN = isbn
	}

	if pages != 0 {
		book.Pages = pages
	}

	if request.Author != "" {
		book.Author = request.Author
	}

	if price != 0 {
		book.Pages = price
	}

	if request.Description != "" {
		book.Description = request.Description
	}

	if quota != 0 {
		book.Pages = quota
	}

	if request.Thumbnail != "" {
		book.Thumbnail = request.Thumbnail
	}

	if request.Book != "" {
		book.Book = request.Book
	}

	newBook, _ := h.BookRepository.UpdateBook(book)

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
	response := dto.SuccessResult{Code: http.StatusOK, Data: ConvertBookResponse(newBookResponse)}
	json.NewEncoder(w).Encode(response)
}

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

func (h *handlerBook) DeleteBookThumbnail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := h.BookRepository.DeleteBookThumbnail(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := h.BookRepository.GetBook(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: ConvertBookResponse(data)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerBook) DeleteBookDocument(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := h.BookRepository.DeleteBookDocument(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := h.BookRepository.GetBook(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: ConvertBookResponse(data)}
	json.NewEncoder(w).Encode(response)
}

func ConvertBookResponse(book models.Book) booksdto.BookResponse {
	return booksdto.BookResponse{
		Id:                 book.Id,
		Title:              book.Title,
		PublicationDate:    book.PublicationDate.Format("2 January 2006"),
		ISBN:               book.ISBN,
		Pages:              book.Pages,
		Author:             book.Author,
		Price:              book.Price,
		IsPromo:            book.IsPromo,
		Discount:           book.Discount,
		PriceAfterDiscount: book.PriceAfterDiscount,
		Description:        book.Description,
		Quota:              book.Quota,
		Book:               book.Book,
		Thumbnail:          book.Thumbnail,
	}
}

func ConvertMultipleBookResponse(books []models.Book) []booksdto.BookResponse {
	var bookResponse []booksdto.BookResponse

	for _, book := range books {
		bookResponse = append(bookResponse, booksdto.BookResponse{
			Id:                 book.Id,
			Title:              book.Title,
			PublicationDate:    book.PublicationDate.Format("Monday, 2 January 2006"),
			ISBN:               book.ISBN,
			Pages:              book.Pages,
			Author:             book.Author,
			Price:              book.Price,
			IsPromo:            book.IsPromo,
			Discount:           book.Discount,
			PriceAfterDiscount: book.Discount,
			Description:        book.Description,
			Quota:              book.Quota,
			Book:               book.Book,
			Thumbnail:          book.Thumbnail,
		})
	}
	return bookResponse
}
