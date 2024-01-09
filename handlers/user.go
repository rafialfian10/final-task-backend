package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	dto "waysbook/dto/result"
	usersdto "waysbook/dto/user"
	"waysbook/models"
	"waysbook/repositories"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type handlerUser struct {
	UserRepository repositories.UserRepository
}

func HandlerUser(UserRepository repositories.UserRepository) *handlerUser {
	return &handlerUser{UserRepository}
}

func (h *handlerUser) FindUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := h.UserRepository.FindUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: users}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerUser) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	user, err := h.UserRepository.GetUser(userId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseUser(user)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerUser) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	dataPhoto := r.Context().Value("dataPhoto")
	if dataPhoto == nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: "dataPhoto is nil"}
		json.NewEncoder(w).Encode(response)
		return
	}

	filePhoto, ok := dataPhoto.(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: "dataPhoto is not a string"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// cloudinary
	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)
	respPhoto, _ := cld.Upload.Upload(ctx, filePhoto, uploader.UploadParams{Folder: "waysbook"})

	var photoSecureURL string
	if respPhoto != nil && respPhoto.SecureURL != "" {
		photoSecureURL = respPhoto.SecureURL
	}

	request := usersdto.UpdateUserRequest{
		Name:    r.FormValue("name"),
		Email:   r.FormValue("email"),
		Gender:  r.FormValue("gender"),
		Phone:   r.FormValue("phone"),
		Address: r.FormValue("address"),
		Photo:   photoSecureURL,
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	user, _ := h.UserRepository.GetUser(id)

	// name
	if r.FormValue("name") != "" {
		user.Name = r.FormValue("name")
	}

	// email
	if r.FormValue("email") != "" {
		user.Email = r.FormValue("email")
	}

	// gender
	if r.FormValue("gender") != "" {
		user.Gender = r.FormValue("gender")
	}

	// phone
	if r.FormValue("phone") != "" {
		user.Phone = r.FormValue("phone")
	}

	// address
	if r.FormValue("address") != "" {
		user.Address = r.FormValue("address")
	}

	// photo
	if request.Photo != "" {
		user.Photo = request.Photo
	}

	newUser, err := h.UserRepository.UpdateUser(user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	newUserResponse, err := h.UserRepository.GetUser(newUser.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// jika tidak ada error maka SuccessResult
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseUser(newUserResponse)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerUser) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	user, err := h.UserRepository.GetUser(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := h.UserRepository.DeleteUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseUser(data)}
	json.NewEncoder(w).Encode(response)
}

func convertResponseUser(u models.User) usersdto.UserResponse {
	return usersdto.UserResponse{
		Id:      u.Id,
		Name:    u.Name,
		Email:   u.Email,
		Gender:  u.Gender,
		Phone:   u.Phone,
		Address: u.Address,
		Photo:   u.Photo,
		Role:    u.Role,
	}
}
