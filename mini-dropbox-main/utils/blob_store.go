package utils

import (
	"errors"
	"io"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type blobStore struct {
	client *s3.S3
}

func NewS3Client() (S3Ops, error) {
	// Create an AWS session
	s3Region := GetEnvValue("S3_REGION", "ap-south-1")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(s3Region),
		// Using server based access control for security purposes
	})
	if err != nil {
		return nil, errors.New("error creating AWS session")
	}

	// Create an S3 client, and return it
	return &blobStore{
		s3.New(awsSession),
	}, nil
}

type S3Ops interface {
	DeleteObject(bucket, key string) error
	UploadObject(bucket, key string, file io.Reader) error
	UploadObjectParts(bucket, key string, file io.Reader) error
}

func (bs *blobStore) DeleteObject(bucket, key string) error {
	// Create input for the DeleteObject operation
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	// Delete the S3 object
	_, err := bs.client.DeleteObject(input)
	DebugLog("S3 object is deleted successfully. Key:", key)

	return err
}

func (bs *blobStore) UploadObject(bucket, key string, file io.Reader) error {
	// Create an uploader with the S3 client and specify the bucket and object key
	uploader := s3manager.NewUploaderWithClient(bs.client)

	// Upload input parameters
	upParams := &s3manager.UploadInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   file,
	}
	_, err := uploader.Upload(upParams)
	return err
}

func (bs *blobStore) UploadObjectParts(bucket, key string, file io.Reader) error {
	// Create an uploader with the S3 client and specify the bucket and object key
	uploader := s3manager.NewUploaderWithClient(bs.client)

	// Upload input parameters
	upParams := &s3manager.UploadInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   file,
	}

	// Perform upload with options different than the those in the Uploader.
	_, err := uploader.Upload(upParams, func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024 // 5MB part size
		u.LeavePartsOnError = true   // Don't delete the parts if the upload fails.
	})

	return err
}
