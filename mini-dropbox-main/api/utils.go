package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/manishlpu/assignment/utils"
)

var (
	SUCCESS_MSG = map[string]interface{}{
		"status": "success",
	}
	FAILURE_MSG = map[string]interface{}{
		"status": "failure",
		"error":  nil,
	}
)

func getSuccessMessage() []byte {
	data, err := json.Marshal(SUCCESS_MSG)
	if err != nil {
		return nil
	}
	return data
}

func getFailureMessage(err error) []byte {
	// Append error to failure message
	FAILURE_MSG["error"] = err.Error()

	data, err := json.Marshal(FAILURE_MSG)
	if err != nil {
		return nil
	}
	return data
}

func getCustomMessage(msg map[string]interface{}) ([]byte, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func getS3KeyFromURI(uri string) string {
	bucketName := utils.GetEnvValue("S3_BUCKET", "bucket")

	prefix := fmt.Sprintf("https://%s.s3.amazonaws.com/", bucketName)
	return strings.TrimPrefix(uri, prefix)
}
