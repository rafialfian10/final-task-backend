package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	dto "waysbook/dto/result"
)

// type contextImage string

// const dataImageKey contextImage = "dataThumbnail"

func UploadFileThumbnail(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Upload file
		// FormFile returns the first file for the given key `myFile`
		// it also returns the FileHeader so we can get the Filename,
		// the Header and the size of the file
		file, _, err := r.FormFile("thumbnail")

		if err != nil && r.Method == "PATCH" {
			ctx := context.WithValue(r.Context(), "dataThumbnail", "")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "no such file thumbnail"}
			json.NewEncoder(w).Encode(response)
			return
		}
		defer file.Close()

		// setup file type filtering
		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/jpg" {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "The provided file format is not allowed. Please upload a JPG, JPEG or PNG image"}
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

		// setup max-upload
		const MAX_UPLOAD_SIZE = 10 << 20
		r.ParseMultipartForm(MAX_UPLOAD_SIZE)
		if r.ContentLength > MAX_UPLOAD_SIZE {
			w.WriteHeader(http.StatusBadRequest)
			response := Result{Code: http.StatusBadRequest, Message: "Max size in 10mb"}
			json.NewEncoder(w).Encode(response)
			return
		}

		// Create a temporary file within our temp-images directory that follows
		// a particular naming pattern
		tempFile, err := ioutil.TempFile("uploads/thumbnail", "thumbnail-*.png")
		if err != nil {
			fmt.Println(err)
			fmt.Println("path upload error")
			json.NewEncoder(w).Encode(err)
			return
		}
		defer tempFile.Close()

		// read all of the contents of our uploaded file into a
		// byte array
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}

		// write this byte array to our temporary file
		tempFile.Write(fileBytes)

		data := tempFile.Name()
		// filename := data[8:] // split uploads/

		// add filename to ctx
		ctx := context.WithValue(r.Context(), "dataThumbnail", data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// func UploadFileTransaction(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		file, _, err := r.FormFile("image")

// 		if err != nil && r.Method == "PATCH" {
// 			ctx := context.WithValue(r.Context(), "dataFileTrans", "false")
// 			next.ServeHTTP(w, r.WithContext(ctx))
// 			return
// 		}

// 		if err != nil {
// 			fmt.Println(err)
// 			json.NewEncoder(w).Encode("Error Retrieving the File")
// 			return
// 		}
// 		defer file.Close()

// 		const MAX_UPLOAD_SIZE = 10 << 20 // masksimal file upload 10mb

// 		// var MAX_UPLOAD_SIZE akan diparse
// 		r.ParseMultipartForm(MAX_UPLOAD_SIZE)

// 		// if contentLength lebih besar dari file yang diupload maka panggil ErrorResult
// 		if r.ContentLength > MAX_UPLOAD_SIZE {
// 			w.WriteHeader(http.StatusBadRequest)
// 			response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "Max size in 1mb"}
// 			json.NewEncoder(w).Encode(response)
// 			return
// 		}

// 		// jika ukuran file sudah dibawah maksimal upload file maka file masuk ke folder upload
// 		tempFile, err := ioutil.TempFile("uploads", "image-*.png")
// 		if err != nil {
// 			fmt.Println(err)
// 			fmt.Println("path upload error")
// 			json.NewEncoder(w).Encode(err)
// 			return
// 		}
// 		defer tempFile.Close()

// 		// baca semua isi file yang kita upload, jika ada error maka tampilkan err
// 		fileBytes, err := ioutil.ReadAll(file)
// 		if err != nil {
// 			fmt.Println(err)
// 		}

// 		// write this byte array to our temporary file
// 		tempFile.Write(fileBytes)

// 		data := tempFile.Name()
// 		// filepath := data[8:] // split uploads(huruf paling 8 depan akan diambil)

// 		// filename akan ditambahkan kedalam variable ctx. dan r.Context akan di panggil jika ingin upload file
// 		ctx := context.WithValue(r.Context(), "dataFileTrans", data)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }
