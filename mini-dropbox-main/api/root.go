package api

import (
	"net/http"

	"github.com/manishlpu/assignment/utils"

	"github.com/gorilla/mux"
)

func New() (*mux.Router, error) {
	router := mux.NewRouter()

	dropboxRouter := router.PathPrefix("/api").Subrouter()
	dropboxRouter.Use(PanicRecoveryMiddleware)
	dropboxHandler(dropboxRouter)

	// Apply the CORS middleware to all routes
	router.Use(corsMiddleware)

	return router, nil
}

func PanicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				// Handle the panic
				utils.InfoLog("Panic recovered: ", r)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func DeleteInactiveRecords() error {
	ah := NewAPIHandler()

	records, err := ah.MetadataOps.FetchInactiveRecords()
	if err != nil {
		return err
	}

	bucketName := utils.GetEnvValue("S3_BUCKET", "assignment")
	for _, record := range records {
		s3Key := getS3KeyFromURI(record.S3ObjectKey)
		if err = ah.S3Ops.DeleteObject(bucketName, s3Key); err != nil {
			utils.ErrorLog("unable to remove the s3 object with following details: ", record)
			continue
		}
	}
	return nil
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		next.ServeHTTP(w, r)
	})
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
