package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	cartdto "waysbook/dto/cart"
	dto "waysbook/dto/result"
	"waysbook/models"
	"waysbook/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type handlerCart struct {
	CartRepository repositories.CartRepository
}

func HandlerCart(CartRepository repositories.CartRepository) *handlerCart {
	return &handlerCart{CartRepository}
}

func (h *handlerCart) CreateCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	request := new(cartdto.CreateCartRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "cek dto"}
		json.NewEncoder(w).Encode(response)
		return
	}

	validation := validator.New()
	err := validation.Struct(request)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "error validation"}
		json.NewEncoder(w).Encode(response)
		return
	}

	transaction, err := h.CartRepository.GetTransactionID(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "ID Transaction Not Found!"}
		json.NewEncoder(w).Encode(response)
		return
	}

	book, err := h.CartRepository.GetBookCart(request.BookId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "Book Not Found!"}
		json.NewEncoder(w).Encode(response)
		return
	}

	if transaction.Id == 0 {

		transId := int(time.Now().Unix())

		transaction := models.Transaction{
			Id:     transId,
			UserId: userId,
			Status: "waiting",
		}

		createTrans, err := h.CartRepository.CreateTransaction(transaction)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "Cart Failed!"}
			json.NewEncoder(w).Encode(response)
			return
		}

		dataCart := models.Cart{
			BookId:        request.BookId,
			TransactionId: int(createTrans.Id),
			Total:         book.Price,
		}

		cart, err := h.CartRepository.CreateCart(dataCart)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "Cart Failed!"}
			json.NewEncoder(w).Encode(response)
			return
		}

		res, _ := h.CartRepository.GetCart(cart.Id)

		w.WriteHeader(http.StatusOK)
		response := dto.SuccessResult{Code: http.StatusOK, Data: res}
		json.NewEncoder(w).Encode(response)
	} else {
		dataCart := models.Cart{
			BookId:        request.BookId,
			TransactionId: int(transaction.Id),
			Total:         book.Price,
		}

		cart, err := h.CartRepository.CreateCart(dataCart)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "Cart Failed!"}
			json.NewEncoder(w).Encode(response)
			return
		}

		res, _ := h.CartRepository.GetCart(cart.Id)

		w.WriteHeader(http.StatusOK)
		response := dto.SuccessResult{Code: http.StatusOK, Data: res}
		json.NewEncoder(w).Encode(response)
	}
}

func (h *handlerCart) DeleteCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	cart, err := h.CartRepository.GetCart(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := h.CartRepository.DeleteCart(cart)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	dataDetele := data.Id
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: dataDetele}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerCart) GetCartByTransID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int64(userInfo["id"].(float64))

	transaction, _ := h.CartRepository.GetTransactionID(int(userID))

	cart, err := h.CartRepository.GetCartByTransID(int(transaction.Id))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: cart}
	json.NewEncoder(w).Encode(response)
}
