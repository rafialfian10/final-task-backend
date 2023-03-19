package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	cartdto "waysbook/dto/cart"
	dto "waysbook/dto/result"
	"waysbook/models"
	"waysbook/repositories"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type handlerCart struct {
	CartRepository repositories.CartRepository
}

func HandlerCart(orderRepository repositories.CartRepository) *handlerCart {
	return &handlerCart{orderRepository}
}

// function get all carts
func (h *handlerCart) FindCarts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	cart, err := h.CartRepository.FindCarts(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := dto.SuccessResult{Code: http.StatusOK, Data: convertMultipleCartResponse(cart)}
	json.NewEncoder(w).Encode(res)
}

// function get detail cart
func (h *handlerCart) GetCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	cart, err := h.CartRepository.GetCart(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := dto.SuccessResult{Code: http.StatusOK, Data: convertCartResponse(cart)}
	json.NewEncoder(w).Encode(res)
}

// function add cart
func (h *handlerCart) CreateCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request cartdto.CreateCartRequest
	json.NewDecoder(r.Body).Decode(&request)

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	// periksa order dengan product id yang sama
	cart, err := h.CartRepository.GetCartByBook(request.BookId, userId)
	if err != nil {
		// jika belum ada, maka buat baru
		newCart := models.Cart{
			UserId:   userId,
			BookId:   request.BookId,
			OrderQty: 1,
		}

		addCart, err := h.CartRepository.CreateCart(newCart)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			res := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
			json.NewEncoder(w).Encode(res)
			return
		}

		cartAdded, err := h.CartRepository.GetCart(addCart.Id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			res := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
			json.NewEncoder(w).Encode(res)
			return
		}

		w.WriteHeader(http.StatusOK)
		res := dto.SuccessResult{Code: http.StatusOK, Data: convertCartResponse(cartAdded)}
		json.NewEncoder(w).Encode(res)
		return
	}

	// bila sudah ada, maka cukup tambahkan order qty
	cart.OrderQty = cart.OrderQty + 1

	cartUpdated, err := h.CartRepository.UpdateCart(cart)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(res)
		return
	}

	cart, err = h.CartRepository.GetCart(cartUpdated.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := dto.SuccessResult{Code: http.StatusOK, Data: convertCartResponse(cart)}
	json.NewEncoder(w).Encode(res)
}

func (h *handlerCart) UpdateCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// mengabil event dari request body
	var request cartdto.UpdateCartRequest
	json.NewDecoder(r.Body).Decode(&request)

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	cart, err := h.CartRepository.GetCart(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := dto.ErrorResult{
			Code: http.StatusBadRequest, Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	if request.Event == "add" {
		cart.OrderQty = cart.OrderQty + 1
	} else if request.Event == "less" {
		cart.OrderQty = cart.OrderQty - 1
	}

	if request.OrderQty != 0 {
		cart.OrderQty = request.OrderQty
	}

	cartUpdated, err := h.CartRepository.UpdateCart(cart)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := dto.ErrorResult{
			Code: http.StatusBadRequest, Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	cart, err = h.CartRepository.GetCart(cartUpdated.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := dto.ErrorResult{
			Code: http.StatusBadRequest, Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertCartResponse(cart),
	}
	json.NewEncoder(w).Encode(res)
}

func (h *handlerCart) DeleteCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	cart, err := h.CartRepository.GetCart(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := dto.ErrorResult{
			Code: http.StatusBadRequest, Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	cartDeleted, err := h.CartRepository.DeleteCart(cart)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res := dto.ErrorResult{
			Code: http.StatusBadRequest, Message: err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertCartResponse(cartDeleted),
	}
	json.NewEncoder(w).Encode(res)
}

func convertMultipleCartResponse(carts []models.Cart) []cartdto.CartResponse {
	var cartResponse []cartdto.CartResponse

	for _, cart := range carts {
		cartResponse = append(cartResponse, cartdto.CartResponse{
			Id:            cart.Id,
			BookId:        cart.BookId,
			Book:          cart.Book,
			UserId:        cart.UserId,
			BookTitle:     cart.Book.Title,
			BookThumbnail: cart.Book.Thumbnail,
			Author:        cart.Book.Author,
			OrderQty:      cart.OrderQty,
		})
	}

	return cartResponse
}

func convertCartResponse(cart models.Cart) cartdto.CartResponse {
	return cartdto.CartResponse{
		Id:            cart.Id,
		BookId:        cart.BookId,
		Book:          cart.Book,
		UserId:        cart.UserId,
		BookTitle:     cart.Book.Title,
		BookThumbnail: cart.Book.Thumbnail,
		Author:        cart.Book.Author,
		OrderQty:      cart.OrderQty,
	}
}
