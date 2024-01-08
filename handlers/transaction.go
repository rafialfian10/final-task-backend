package handlers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	_ "time/tzdata"
	dto "waysbook/dto/result"
	transactiondto "waysbook/dto/transaction"
	"waysbook/models"
	"waysbook/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gopkg.in/gomail.v2"
)

type handlerTransaction struct {
	TransactionRepository repositories.TransactionRepository
}

func HandlerTransaction(transactionRepository repositories.TransactionRepository) *handlerTransaction {
	return &handlerTransaction{transactionRepository}
}

func (h *handlerTransaction) FindTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	transactions, err := h.TransactionRepository.FindTransactions()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		res := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(res)
		return
	}

	// for i, p := range transactions {
	// 	fileImage := os.Getenv("PATH_FILE") + p.Image
	// 	transactions[i].Image = fileImage
	// }

	w.WriteHeader(http.StatusOK)
	res := dto.SuccessResult{Code: http.StatusOK, Data: ConvertMultipleTransactionResponse(transactions)}
	json.NewEncoder(w).Encode(res)
}

func (h *handlerTransaction) FindTransactionsByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	idUser := int(userInfo["id"].(float64))

	transactions, err := h.TransactionRepository.FindTransactionsByUser(idUser)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		res := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(res)
		return
	}

	// for i, p := range transactions {
	// 	fileImage := os.Getenv("PATH_FILE") + p.Image
	// 	transactions[i].Image = fileImage
	// }

	w.WriteHeader(http.StatusOK)
	res := dto.SuccessResult{Code: http.StatusOK, Data: ConvertMultipleTransactionResponse(transactions)}
	json.NewEncoder(w).Encode(res)
}

func (h *handlerTransaction) GetTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	transaction, err := h.TransactionRepository.GetTransaction(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		res := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(res)
		return
	}

	// transaction.Image = os.Getenv("PATH_FILE") + transaction.Image

	w.WriteHeader(http.StatusOK)
	res := dto.SuccessResult{Code: http.StatusOK, Data: ConvertTransactionResponse(transaction)}
	json.NewEncoder(w).Encode(res)
}

func (h *handlerTransaction) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request transactiondto.CreateTransactionRequest
	json.NewDecoder(r.Body).Decode(&request)

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	request.UserId = int(userInfo["id"].(float64))

	// middleware
	// dataContex := r.Context().Value("dataFileTrans")
	// fileImage := dataContex.(string)

	// validasi semua input
	validation := validator.New()
	errValidation := validation.Struct(request)
	if errValidation != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: errValidation.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// cloudinary
	// var ctx = context.Background()
	// var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	// var API_KEY = os.Getenv("API_KEY")
	// var API_SECRET = os.Getenv("API_SECRET")

	// cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	// // Upload file to Cloudinary
	// resp, err := cld.Upload.Upload(ctx, fileImage, uploader.UploadParams{Folder: "waysbook"})
	// fmt.Println(resp.SecureURL)

	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	transaction := models.Transaction{
		Id:        fmt.Sprintf("%d-%d", request.UserId, TimeIn("Asia/Jakarta").UnixNano()),
		OrderDate: TimeIn("Asia/Jakarta"),
		Total:     request.Total,
		Status:    "pending",
		UserId:    request.UserId,
		// Image:     resp.SecureURL,
	}

	for _, order := range request.Books {
		transaction.Cart = append(transaction.Cart, models.CartResponse{
			Id:       order.Id,
			BookId:   order.BookId,
			OrderQty: order.OrderQty,
		})
	}
	fmt.Println("transaction bro", transaction)

	addTransaction, err := h.TransactionRepository.CreateTransaction(transaction)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	transactionAdded, err := h.TransactionRepository.GetTransaction(addTransaction.Id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// 1. Initiate Snap client
	var s = snap.Client{}
	s.New(os.Getenv("SERVER_KEY"), midtrans.Sandbox)

	// 2. Initiate Snap request
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transactionAdded.Id,
			GrossAmt: int64(transactionAdded.Total),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: transactionAdded.User.Name,
			Phone: transactionAdded.User.Phone,
			BillAddr: &midtrans.CustomerAddress{
				FName:   transactionAdded.User.Name,
				Phone:   transactionAdded.User.Phone,
				Address: transactionAdded.User.Address,
			},
			ShipAddr: &midtrans.CustomerAddress{
				FName:   transactionAdded.User.Name,
				Phone:   transactionAdded.User.Phone,
				Address: transactionAdded.User.Address,
			},
		},
	}

	// 3. Request create Snap transaction to Midtrans
	snapResp, _ := s.CreateTransactionToken(req)
	fmt.Println("Midtrans Id :", snapResp)

	updateTransaction, err := h.TransactionRepository.UpdateTokenTransaction(snapResp, transactionAdded.Id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		res := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(res)
		return
	}

	// mengambil data transaction yang baru diupdate
	transactionUpdated, _ := h.TransactionRepository.GetTransaction(updateTransaction.Id)

	w.WriteHeader(http.StatusCreated)
	res := dto.SuccessResult{Code: http.StatusOK, Data: ConvertTransactionResponse(transactionUpdated)}
	json.NewEncoder(w).Encode(res)
}

func (h *handlerTransaction) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	var request transactiondto.UpdateTransactionRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	fmt.Println("ID : ", id)
	fmt.Println("Status : ", request.Status)

	// memeriksa transaksi yang ingin diupdate
	_, err = h.TransactionRepository.GetTransaction(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		res := dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	// mengupdate status transaksi
	transaction, err := h.TransactionRepository.UpdateTransaction(request.Status, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	// mengambil data transaksi yang sudah diupdate
	transaction, err = h.TransactionRepository.GetTransaction(transaction.Id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		res := dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	if request.Status == "reject" {
		SendTransactionMail("Rejected", transaction)
	} else if request.Status == "sent" {
		SendTransactionMail("Success, Product On Delivery", transaction)
	} else if request.Status == "done" {
		SendTransactionMail("Success, Product Received", transaction)
	}

	w.WriteHeader(http.StatusOK)
	res := dto.SuccessResult{Code: http.StatusOK, Data: ConvertTransactionResponse(transaction)}
	json.NewEncoder(w).Encode(res)
}

func (h *handlerTransaction) UpdateTransactionByAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	status := r.FormValue("status")
	request := transactiondto.UpdateTransactionByAdminRequest{
		Status: status,
	}

	// mengambil data dari request form
	json.NewDecoder(r.Body).Decode(&request)
	fmt.Println("status", request.Status)

	// get data yang ingin diupdate berdasarkan id yang didapatkan dari url
	_, err := h.TransactionRepository.GetTransaction(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := dto.ErrorResult{Code: http.StatusNotFound, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// fmt.Println("ID after", id)

	// mengirim data transaction yang sudah diupdate ke database
	transactionUpdated, err := h.TransactionRepository.UpdateTransaction(request.Status, id)
	fmt.Println("Transaction Updated", transactionUpdated)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	getTransactionUpdated, err := h.TransactionRepository.GetTransaction(transactionUpdated.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.SuccessResult{Code: http.StatusOK, Data: ConvertTransactionResponse(getTransactionUpdated)}

	// mengirim response
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) Notification(w http.ResponseWriter, r *http.Request) {
	// 1. Initialize empty map
	var notificationPayload map[string]interface{}

	// 2. Parse JSON request body and use it to set json to payload
	err := json.NewDecoder(r.Body).Decode(&notificationPayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	// 3. Get order_id from payload
	orderId, exists := notificationPayload["order_id"].(string)
	if !exists {
		// do something when key `order_id` not found
		return
	}

	// 4. Check transaction to Midtrans with param orderId
	transaction, err := h.TransactionRepository.GetTransaction(orderId)
	// jika transaksi di database tidak ditemukan, atau sudah dihapus, maka hentikan fungsi notification (menghindari app crash)
	if err != nil {
		fmt.Println("Transaction not found")
		return
	}

	transactionStatus := notificationPayload["transaction_status"].(string)
	fraudStatus := notificationPayload["fraud_status"].(string)

	if transactionStatus != "" {
		// 5. Do set transaction status based on response from check transaction status
		if transactionStatus == "capture" {
			if fraudStatus == "challenge" {
				h.TransactionRepository.UpdateTransaction("pending", transaction.Id)
			} else if fraudStatus == "accept" {
				h.TransactionRepository.UpdateTransaction("success", transaction.Id)
				transaction.Status = "success"
				SendTransactionMail("Success", transaction)
			}
		} else if transactionStatus == "settlement" {
			h.TransactionRepository.UpdateTransaction("success", transaction.Id)
			transaction.Status = "success"
			SendTransactionMail("success", transaction)
		} else if transactionStatus == "deny" {
			h.TransactionRepository.UpdateTransaction("failed", transaction.Id)
			transaction.Status = "failed"
		} else if transactionStatus == "cancel" || transactionStatus == "expire" {
			h.TransactionRepository.UpdateTransaction("failed", transaction.Id)
			transaction.Status = "failed"
			SendTransactionMail("Failed", transaction)
		} else if transactionStatus == "pending" {
			transaction.Status = "pending"
			h.TransactionRepository.UpdateTransaction("pending", transaction.Id)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("ok"))
}

func SendTransactionMail(status string, transaction models.Transaction) {
	var CONFIG_SMTP_HOST = "smtp.gmail.com"
	var CONFIG_SMTP_PORT = 587
	var CONFIG_SENDER_NAME = "WAYSBOOK <rafialfian770@gmail.com>"
	var CONFIG_AUTH_EMAIL = os.Getenv("SYSTEM_EMAIL")
	var CONFIG_AUTH_PASSWORD = os.Getenv("SYSTEM_PASSWORD")

	var bookName = transaction.User.Name
	var price = strconv.Itoa(transaction.Total)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", transaction.User.Email)
	mailer.SetHeader("Subject", "Status Transaction")
	mailer.SetBody("text/html", fmt.Sprintf(`<!DOCTYPE html>
    <html lang="en">
      <head>
      <meta charset="UTF-8" />
      <meta http-equiv="X-UA-Compatible" content="IE=edge" />
      <meta name="viewport" content="width=device-width, initial-scale=1.0" />
      <title>Document</title>
      <style>
        h1 {
        color: brown;
        }
      </style>
      </head>
      <body>
      <h2>Product payment :</h2>
      <ul style="list-style-type:none;">
        <li>Name : %s</li>
        <li>Total payment: Rp.%s</li>
        <li>Status : %s</li>
		<li>Iklan : %s</li>
      </ul>
      </body>
    </html>`, bookName, price, status, "Terima kasih"))

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ConvertTransactionResponse(transaction models.Transaction) transactiondto.TransactionResponse {
	var transactionResponse = transactiondto.TransactionResponse{
		Id:         transaction.Id,
		MidtransId: transaction.MidtransId,
		OrderDate:  transaction.OrderDate.Format("Monday, 2 January 2006"),
		Total:      transaction.Total,
		Status:     transaction.Status,
		User:       transaction.User,
		// Image:      transaction.Image,
	}

	for _, order := range transaction.Cart {
		transactionResponse.Book = append(transactionResponse.Book, transactiondto.BookResponseForTransaction{
			Id:                 order.Id,
			Title:              order.Book.Title,
			PublicationDate:    order.Book.PublicationDate.Format("2 january 2006"),
			ISBN:               order.Book.ISBN,
			Pages:              order.Book.Pages,
			Author:             order.Book.Author,
			Price:              order.Book.Price,
			IsPromo:            order.Book.IsPromo,
			Discount:           order.Book.Discount,
			PriceAfterDiscount: order.Book.PriceAfterDiscount,
			Description:        order.Book.Description,
			Book:               order.Book.Book,
			Thumbnail:          order.Book.Thumbnail,
			Quota:              order.Book.Quota,
			OrderQty:           order.OrderQty,
		})
	}
	return transactionResponse
}

func ConvertMultipleTransactionResponse(transactions []models.Transaction) []transactiondto.TransactionResponse {
	var transactionsResponse []transactiondto.TransactionResponse

	for _, t := range transactions {
		var tResponse = transactiondto.TransactionResponse{
			Id:         t.Id,
			MidtransId: t.MidtransId,
			OrderDate:  t.OrderDate.Format("Monday, 2 January 2006"),
			Total:      t.Total,
			Status:     t.Status,
			User:       t.User,
			// Image:      t.Image,
		}

		for _, order := range t.Cart {
			tResponse.Book = append(tResponse.Book, transactiondto.BookResponseForTransaction{
				Id:                 order.Id,
				Title:              order.Book.Title,
				PublicationDate:    order.Book.PublicationDate.Format("2 january 2006"),
				ISBN:               order.Book.ISBN,
				Pages:              order.Book.Pages,
				Author:             order.Book.Author,
				Price:              order.Book.Price,
				IsPromo:            order.Book.IsPromo,
				Discount:           order.Book.Discount,
				PriceAfterDiscount: order.Book.PriceAfterDiscount,
				Description:        order.Book.Description,
				Book:               order.Book.Book,
				Thumbnail:          order.Book.Thumbnail,
				Quota:              order.Book.Quota,
				OrderQty:           order.OrderQty,
			})
		}
		transactionsResponse = append(transactionsResponse, tResponse)
	}
	return transactionsResponse
}

func TimeIn(name string) time.Time {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return time.Now().In(loc)
}
