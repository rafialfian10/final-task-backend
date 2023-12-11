package middleware

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	dto "waysbook/dto/result"
)

func UploadFilePdf(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("book")

		if err != nil && r.Method == "PATCH" {
			ctx := context.WithValue(r.Context(), "dataPdf", "")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "no such file document"}
			json.NewEncoder(w).Encode(response)
			return
		}
		defer file.Close()

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		filetype := http.DetectContentType(buff)
		if filetype != "application/pdf" {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "The provided file format is not allowed. Please upload pdf"}
			json.NewEncoder(w).Encode(response)
			return
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		const MAX_UPLOAD_SIZE = 50 << 20
		r.ParseMultipartForm(MAX_UPLOAD_SIZE)
		if r.ContentLength > MAX_UPLOAD_SIZE {
			w.WriteHeader(http.StatusBadRequest)
			response := Result{Code: http.StatusBadRequest, Message: "Max size is 50mb"}
			json.NewEncoder(w).Encode(response)
			return
		}

		// Temporary file handling
		tempFile, err := ioutil.TempFile("uploads/document", "document-*.pdf")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}
		defer tempFile.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		tempFile.Write(fileBytes)

		data := tempFile.Name()
		// filename := data[8:] // split uploads/

		ctx := context.WithValue(r.Context(), "dataPdf", data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
