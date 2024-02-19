package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/manishlpu/assignment/models"
	"github.com/manishlpu/assignment/utils"
	"github.com/gorilla/mux"
)

// Uploads the file to blob storage (s3 here).
func (ah *APIHandler) uploadFile(w http.ResponseWriter, r *http.Request) {
	utils.DebugLog("inside uploadFile")

	// Parse the multipart form data
	err := r.ParseMultipartForm(10 << 9) // Setting the max limit to 100MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve the description and uploaded file
	desc := r.FormValue("description")
	file, header, err := r.FormFile("upload_file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Specify the S3 bucket and object key where you want to upload the file
	bucketName := utils.GetEnvValue("S3_BUCKET", "dropbox_files")
	s3ObjectKey := header.Filename + "_" + fmt.Sprint(time.Now().UnixNano())

	idChan := make(chan int64)
	go func(fh *multipart.FileHeader, uri, description string) {
		ext := fh.Filename[strings.LastIndexByte(fh.Filename, '.'):]
		mimeType := mime.TypeByExtension(ext)

		record := models.Metadata{
			Filename:    fh.Filename,
			SizeInBytes: fh.Size,
			S3ObjectKey: uri,
			MimeType:    mimeType,
			Description: description,
			Status:      1,
		}
		// Insert the metadata into RDBMS using goroutine
		id, err := ah.MetadataOps.SaveRecord(record)
		if err != nil {
			utils.ErrorLog("Error saving metadata for upload: ", err)
			return
		}
		idChan <- id
	}(header, fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, s3ObjectKey), desc)

	utils.DebugLog("File size in bytes: ", header.Size)
	if header.Size <= 5*1024*1024 {
		err = ah.S3Ops.UploadObject(bucketName, s3ObjectKey, file)
		if err != nil {
			utils.ErrorLog("Error uploading metadata for upload: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(getFailureMessage(errors.New("unable to upload object")))
			return
		}
	} else {
		err = ah.S3Ops.UploadObjectParts(bucketName, s3ObjectKey, file)
		if err != nil {
			utils.ErrorLog("Error uploading metadata for upload: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(getFailureMessage(errors.New("unable to upload object")))
			return
		}
	}

	id := <-idChan
	close(idChan)

	jsonBytes, err := getCustomMessage(map[string]interface{}{
		"id": id,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

// Fetch the file metadata from persistent storage (s3 here).
func (ah *APIHandler) getFile(w http.ResponseWriter, r *http.Request) {
	utils.DebugLog("inside getFile")

	vars := mux.Vars(r)
	id := vars["fileID"]

	w.Header().Add("Content-Type", "application/json")
	if utils.IsEmptyString(id) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getFailureMessage(errors.New("unique id of file is required")))
		return
	}

	fileID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getFailureMessage(err))
		return
	}

	data, err := ah.MetadataOps.GetRecord(fileID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(getFailureMessage(err))
		return
	}
	if data == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(getFailureMessage(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// Upload the new file to blob storage and update metadata with new url.
func (ah *APIHandler) updateFile(w http.ResponseWriter, r *http.Request) {
	utils.DebugLog("inside updateFile")

	vars := mux.Vars(r)
	id := vars["fileID"]

	w.Header().Add("Content-Type", "application/json")
	if utils.IsEmptyString(id) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getFailureMessage(errors.New("unique id of file is required")))
		return
	}

	fileID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getFailureMessage(err))
		return
	}

	// Parse the multipart form data
	err = r.ParseMultipartForm(10 << 9) // Setting the max limit to 100MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Fetch the record with given ID to verify if it exists
	record, err := ah.MetadataOps.GetRecord(fileID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(getFailureMessage(err))
		return
	}

	if record == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getFailureMessage(errors.New("no file exists with given id")))
		return
	}

	// Retrieve the description and uploaded file
	desc := r.FormValue("description")
	file, header, err := r.FormFile("upload_file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Specify the S3 bucket and object key where you want to upload the file
	bucketName := utils.GetEnvValue("S3_BUCKET", "dropbox_files")
	s3ObjectKey := header.Filename + "_" + fmt.Sprint(time.Now().UnixNano())
	newS3Key := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, s3ObjectKey)

	boolChan := make(chan bool)
	go func(fh *multipart.FileHeader, uri, description string) {
		ext := fh.Filename[strings.LastIndexByte(fh.Filename, '.'):]
		mimeType := mime.TypeByExtension(ext)

		newRecord := models.Metadata{
			Filename:    fh.Filename,
			SizeInBytes: fh.Size,
			S3ObjectKey: uri,
			MimeType:    mimeType,
			Description: description,
			Status:      1,
		}
		// Insert the metadata into RDBMS using goroutine
		if err := ah.MetadataOps.UpdateRecord(fileID, newRecord); err != nil {
			utils.ErrorLog("Error saving metadata for upload: ", err)
			return
		}
		boolChan <- true
	}(header, newS3Key, desc)

	utils.DebugLog("File size in bytes: ", header.Size)
	if header.Size <= 5*1024*1024 {
		err = ah.S3Ops.UploadObject(bucketName, s3ObjectKey, file)
		if err != nil {
			utils.ErrorLog("Error uploading metadata for upload: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(getFailureMessage(errors.New("unable to upload object")))
			return
		}
	} else {
		err = ah.S3Ops.UploadObjectParts(bucketName, s3ObjectKey, file)
		if err != nil {
			utils.ErrorLog("Error uploading metadata for upload: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(getFailureMessage(errors.New("unable to upload object")))
			return
		}
	}

	<-boolChan
	close(boolChan)

	// Remove the previous uploaded object from blob store
	go func(bucket, s3Key, newKey string) {
		if !strings.EqualFold(s3Key, newKey) {
			if err = ah.S3Ops.DeleteObject(bucket, getS3KeyFromURI(s3Key)); err != nil {
				utils.ErrorLog("error deleting object: ", err)
				return
			}
		}
	}(bucketName, record.S3ObjectKey, newS3Key)

	w.Header().Set("Content-Type", "application/json")
	w.Write(getSuccessMessage())
}

// Soft deletes the file from blob storage.
func (ah *APIHandler) deleteFile(w http.ResponseWriter, r *http.Request) {
	utils.DebugLog("inside deleteFile")

	vars := mux.Vars(r)
	id := vars["fileID"]

	w.Header().Add("Content-Type", "application/json")
	if utils.IsEmptyString(id) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getFailureMessage(errors.New("unique id of file is required")))
		return
	}

	fileID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getFailureMessage(err))
		return
	}

	// Validate if record with the given id exists
	exists, err := ah.MetadataOps.Exists(fileID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(getFailureMessage(err))
		return
	}

	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getFailureMessage(errors.New("no such record with given id exists")))
		return
	}

	if err = ah.MetadataOps.DeactivateRecord(fileID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(getFailureMessage(err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(getSuccessMessage())
}

func (ah *APIHandler) listFiles(w http.ResponseWriter, r *http.Request) {
	utils.DebugLog("inside listFiles")

	w.Header().Add("Content-Type", "application/json")

	data, err := ah.MetadataOps.FetchRecords()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(getFailureMessage(err))
		return
	}

	if data == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(getFailureMessage(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
