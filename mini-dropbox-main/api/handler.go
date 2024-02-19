package api

import (
	"net/http"

	"github.com/manishlpu/assignment/utils"
	"github.com/gorilla/mux"
)

type APIHandler struct {
	utils.MetadataOps
	utils.S3Ops
}

func NewAPIHandler() *APIHandler {
	persistenceDB, err := utils.NewPersistenceDBLayer()
	if err != nil {
		panic(err)
	}

	s3Client, err := utils.NewS3Client()
	if err != nil {
		panic(err)
	}

	return &APIHandler{
		persistenceDB,
		s3Client,
	}
}

func dropboxHandler(r *mux.Router) {
	dh := NewAPIHandler()

	r.HandleFunc("/files/upload", dh.uploadFile).Methods("POST")
	r.HandleFunc("/files/{fileID}", dh.getFile).Methods("GET")
	r.HandleFunc("/files/{fileID}", dh.updateFile).Methods("PUT")
	r.HandleFunc("/files/{fileID}", dh.deleteFile).Methods("DELETE")
	r.HandleFunc("/files/{fileID}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	}).Methods("OPTIONS")
	r.HandleFunc("/files", dh.listFiles).Methods("GET")

}
