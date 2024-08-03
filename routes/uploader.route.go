package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ynoacamino/ynoa-uploader/db"
	"github.com/ynoacamino/ynoa-uploader/services"
)

func GetPublicFiles(w http.ResponseWriter, r *http.Request) {
	files, err := db.Query.GetPublicFiles(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(files)
}

func GetPrivateFiles(w http.ResponseWriter, r *http.Request) {
	var user db.File

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	files, err := db.Query.GetPrivateFiles(r.Context(), user.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(files)
}

func CreateFile(w http.ResponseWriter, r *http.Request) {
	newFile, err := newFileForm(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	defer file.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	fmt.Println(handler.Header.Get("Content-Type"))

	secureUrl, err := services.UploadFile(bytes.NewReader(buf.Bytes()), handler.Header.Get("Content-Type"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if newFile.FileName == "" {
		newFile.FileName = handler.Filename
	}

	createdFile, err := db.Query.CreateFile(r.Context(), db.CreateFileParams{
		FileName: newFile.FileName,
		Public:   newFile.Public,
		UserID:   newFile.UserID,
		FileUrl:  secureUrl,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(createdFile)
}

func UpdateFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	fileId, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var file db.File

	err = json.NewDecoder(r.Body).Decode(&file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	updatedFile, err := db.Query.UpdateFile(r.Context(), db.UpdateFileParams{
		FileName: file.FileName,
		Public:   file.Public,
		FileID:   int32(fileId),
	})
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		w.Write([]byte("File not found"))
		return
	}

	json.NewEncoder(w).Encode(updatedFile)
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	fileId, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	deletedFile, err := db.Query.DeleteFile(r.Context(), int32(fileId))
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		w.Write([]byte("File not found"))
		return
	}

	json.NewEncoder(w).Encode(deletedFile)
}

func newFileForm(r *http.Request) (db.File, error) {
	FileName := r.FormValue("FileName")
	Public, err := strconv.ParseBool(r.FormValue("Public"))
	if err != nil {
		return db.File{}, err
	}
	UserID := r.FormValue("User_ID")

	return db.File{
		FileName: FileName,
		Public:   Public,
		UserID:   UserID,
	}, nil
}

func SetUpUploaderRoutes(router *mux.Router) {
	router.HandleFunc("/", GetPrivateFiles).Methods("GET")
	router.HandleFunc("/public/", GetPublicFiles).Methods("GET")
	router.HandleFunc("/", CreateFile).Methods("POST")
	router.HandleFunc("/{id}/", UpdateFile).Methods("PUT")
	router.HandleFunc("/{id}/", DeleteFile).Methods("DELETE")
}
