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

// mengambil seluruh data transaksi
func (h *handlerTransaction) FindTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	transactions, err := h.TransactionRepository.FindTransactions()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		res := dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertMultipleTransactionResponse(transactions),
	}
	json.NewEncoder(w).Encode(res)
}

// mengambil seluruh data transaksi milik user tertentu
func (h *handlerTransaction) FindTransactionsByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	idUser := int(userInfo["id"].(float64))

	transactions, err := h.TransactionRepository.FindTransactionsByUser(idUser)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		res := dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertMultipleTransactionResponse(transactions),
	}
	json.NewEncoder(w).Encode(res)
}

// mengambil 1 data transaksi
func (h *handlerTransaction) GetDetailTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	transaction, err := h.TransactionRepository.GetTransaction(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		res := dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertTransactionResponse(transaction),
	}
	json.NewEncoder(w).Encode(res)
}

// membuat transaksi baru
func (h *handlerTransaction) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request transactiondto.CreateTransactionRequest
	json.NewDecoder(r.Body).Decode(&request)

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	request.UserId = int(userInfo["id"].(float64))

	// memvalidasi inputan dari request body
	validation := validator.New()
	errValidation := validation.Struct(request)
	if errValidation != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: errValidation.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	newTransaction := models.Transaction{
		Id:        fmt.Sprintf("TRX-%d-%d", request.UserId, timeIn("Asia/Jakarta").UnixNano()),
		OrderDate: timeIn("Asia/Jakarta"),
		Total:     request.Total,
		Status:    "new",
		UserId:    request.UserId,
	}

	for _, order := range request.Books {
		newTransaction.Cart = append(newTransaction.Cart, models.CartResponse{
			Id:       order.Id,
			BookId:   order.BookId,
			OrderQty: order.OrderQty,
		})
	}

	transactionAdded, err := h.TransactionRepository.CreateTransaction(newTransaction)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	transaction, err := h.TransactionRepository.GetTransaction(transactionAdded.Id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		res := dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	// 1. Initiate Snap client
	var s = snap.Client{}
	s.New(os.Getenv("SERVER_KEY"), midtrans.Sandbox)

	// 2. Initiate Snap request
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transaction.Id,
			GrossAmt: int64(transaction.Total),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: transaction.User.Name,
			Phone: transaction.User.Phone,
			BillAddr: &midtrans.CustomerAddress{
				FName:   transaction.User.Name,
				Phone:   transaction.User.Phone,
				Address: transaction.User.Address,
				// Postcode: transaction.User.PostCode,
			},
			ShipAddr: &midtrans.CustomerAddress{
				FName:   transaction.User.Name,
				Phone:   transaction.User.Phone,
				Address: transaction.User.Address,
				// Postcode: transaction.User.PostCode,
			},
		},
	}

	// 3. Request create Snap transaction to Midtrans
	snapResp, _ := s.CreateTransactionToken(req)
	fmt.Println("Response :", snapResp)

	transaction, err = h.TransactionRepository.UpdateTokenTransaction(snapResp, transaction.Id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		res := dto.ErrorResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusCreated)
	res := dto.SuccessResult{
		Code: http.StatusOK, Data: convertTransactionResponse(transaction),
	}
	json.NewEncoder(w).Encode(res)
}

func (h *handlerTransaction) UpdateTransactionStatus(w http.ResponseWriter, r *http.Request) {
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

	if request.Status == "rejected" {
		SendTransactionMail("Rejected", transaction)
	} else if request.Status == "sent" {
		SendTransactionMail("Success, Product On Delivery", transaction)
	} else if request.Status == "done" {
		SendTransactionMail("Success, Product Received", transaction)
	}

	w.WriteHeader(http.StatusOK)
	res := dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertTransactionResponse(transaction),
	}
	json.NewEncoder(w).Encode(res)
}

func (h *handlerTransaction) Notification(w http.ResponseWriter, r *http.Request) {
	// 1. Initialize empty map
	var notificationPayload map[string]interface{}

	// 2. Parse JSON request body and use it to set json to payload
	err := json.NewDecoder(r.Body).Decode(&notificationPayload)
	if err != nil {
		// do something on error when decode
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	// 3. Get order-id from payload
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
				// TODO set transaction status on your database to 'challenge'
				// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
				h.TransactionRepository.UpdateTransaction("pending", transaction.Id)
			} else if fraudStatus == "accept" {
				// TODO set transaction status on your database to 'success'
				h.TransactionRepository.UpdateTransaction("success", transaction.Id)
				SendTransactionMail("Success", transaction)
			}
		} else if transactionStatus == "settlement" {
			// TODO set transaction status on your databaase to 'success'
			h.TransactionRepository.UpdateTransaction("success", transaction.Id)
			SendTransactionMail("success", transaction)
		} else if transactionStatus == "deny" {
			// TODO you can ignore 'deny', because most of the time it allows payment retries
			// and later can become success
			h.TransactionRepository.UpdateTransaction("failed", transaction.Id)
		} else if transactionStatus == "cancel" || transactionStatus == "expire" {
			// TODO set transaction status on your databaase to 'failure'
			h.TransactionRepository.UpdateTransaction("failed", transaction.Id)
			SendTransactionMail("Failed", transaction)
		} else if transactionStatus == "pending" {
			// TODO set transaction status on your databaase to 'pending' / waiting payment
			h.TransactionRepository.UpdateTransaction("pending", transaction.Id)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("ok"))
}

func convertMultipleTransactionResponse(transactions []models.Transaction) []transactiondto.TransactionResponse {
	var transactionsResponse []transactiondto.TransactionResponse

	for _, trx := range transactions {
		var trxResponse = transactiondto.TransactionResponse{
			Id:         trx.Id,
			MidtransId: trx.MidtransId,
			OrderDate:  trx.OrderDate.Format("Monday, 2 January 2006"),
			Total:      trx.Total,
			Status:     trx.Status,
			User:       trx.User,
		}

		for _, order := range trx.Cart {
			trxResponse.Book = append(trxResponse.Book, transactiondto.BookResponseForTransaction{
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

		transactionsResponse = append(transactionsResponse, trxResponse)
	}

	return transactionsResponse
}

func convertTransactionResponse(transaction models.Transaction) transactiondto.TransactionResponse {
	var transactionResponse = transactiondto.TransactionResponse{
		Id:         transaction.Id,
		MidtransId: transaction.MidtransId,
		OrderDate:  transaction.OrderDate.Format("Monday, 2 January 2006"),
		Total:      transaction.Total,
		Status:     transaction.Status,
		User:       transaction.User,
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

// fungsi untuk kirim email transaksi
func SendTransactionMail(status string, transaction models.Transaction) {
	var CONFIG_SMTP_HOST = "smtp.gmail.com"
	var CONFIG_SMTP_PORT = 587
	var CONFIG_SENDER_NAME = "WAYSBOOK <rafialfian770@gmail.com>"
	var CONFIG_AUTH_EMAIL = os.Getenv("SYSTEM_EMAIL")
	var CONFIG_AUTH_PASSWORD = os.Getenv("SYSTEM_PASSWORD")

	var tripName = transaction.User.Name
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
    </html>`, tripName, price, status, "Terima kasih"))

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

// func SendTransactionMail(status string, transaction models.Transaction) {

// 	var CONFIG_SMTP_HOST = os.Getenv("CONFIG_SMTP_HOST")
// 	var CONFIG_SMTP_PORT, _ = strconv.Atoi(os.Getenv("CONFIG_SMTP_PORT"))
// 	var CONFIG_SENDER_NAME = os.Getenv("CONFIG_SENDER_NAME")
// 	var CONFIG_AUTH_EMAIL = os.Getenv("CONFIG_AUTH_EMAIL")
// 	var CONFIG_AUTH_PASSWORD = os.Getenv("CONFIG_AUTH_PASSWORD")

// 	ac := accounting.Accounting{Symbol: "Rp", Precision: 2}

// 	var products []map[string]interface{}

// 	for _, order := range transaction.Cart {
// 		products = append(products, map[string]interface{}{
// 			"ProductName": order.Book.Title,
// 			"Price":       order.Book.Price,
// 			"Qty":         order.OrderQty,
// 			"SubTotal":    ac.FormatMoney(order.OrderQty * order.Book.Price),
// 		})
// 	}

// 	data := map[string]interface{}{
// 		"TransactionID":     transaction.Id,
// 		"TransactionStatus": status,
// 		"UserName":          transaction.User.Name,
// 		"OrderDate":         timeIn("Asia/Jakarta").Format("Monday, 2 January 2006"),
// 		"Total":             ac.FormatMoney(transaction.Total),
// 		"Products":          products,
// 	}

// 	// mengambil file template
// 	t, err := template.ParseFiles("view/notification_email.html")
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}

// 	bodyMail := new(bytes.Buffer)

// 	// mengeksekusi template, dan memparse "data" ke template
// 	t.Execute(bodyMail, data)

// 	// create new message
// 	trxMail := gomail.NewMessage()
// 	trxMail.SetHeader("From", CONFIG_SENDER_NAME)
// 	trxMail.SetHeader("To", transaction.User.Email)
// 	trxMail.SetHeader("Subject", "WAYSBOOK ORDER NOTIFICATION")
// 	trxMail.SetBody("text/html", bodyMail.String())

// 	trxDialer := gomail.NewDialer(
// 		CONFIG_SMTP_HOST,
// 		CONFIG_SMTP_PORT,
// 		CONFIG_AUTH_EMAIL,
// 		CONFIG_AUTH_PASSWORD,
// 	)

// 	err = trxDialer.DialAndSend(trxMail)
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	log.Println("Pesan terkirim!")
// }

// fungsi untuk mendapatkan waktu sesuai zona indonesia
func timeIn(name string) time.Time {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return time.Now().In(loc)
}
