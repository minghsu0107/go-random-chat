package uploader

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Presigner struct {
	presignClient  *s3.PresignClient
	lifetimeSecond int64
}

// GetObject makes a presigned request that can be used to get an object from a bucket.
// The presigned request is valid for the specified number of seconds.
func (presigner *Presigner) GetObject(ctx context.Context, bucketName string, objectKey string) (*v4.PresignedHTTPRequest, error) {
	request, err := presigner.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(presigner.lifetimeSecond * int64(time.Second))
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't get a presigned request to get %v:%v, reason: %v", bucketName, objectKey, err)
	}
	return request, nil
}

// PutObject makes a presigned request that can be used to put an object in a bucket.
// The presigned request is valid for the specified number of seconds.
func (presigner *Presigner) PutObject(ctx context.Context, bucketName string, objectKey string) (*v4.PresignedHTTPRequest, error) {
	request, err := presigner.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(presigner.lifetimeSecond * int64(time.Second))
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't get a presigned request to put %v:%v, reason: %v", bucketName, objectKey, err)
	}
	return request, nil
}
