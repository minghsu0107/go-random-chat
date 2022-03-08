package upload

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minghsu0107/go-random-chat/pkg/common"
	log "github.com/sirupsen/logrus"
)

var (
	uploader *s3manager.Uploader

	s3Endpoint = os.Getenv("S3_ENDPOINT")
	s3Region   = os.Getenv("S3_REGION")
	s3Bucket   = os.Getenv("S3_BUCKET")
	accessKey  = os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey  = os.Getenv("AWS_SECRET_KEY")
)

func init() {
	creds := credentials.NewStaticCredentials(accessKey, secretKey, "")

	config := &aws.Config{
		Credentials:      creds,
		Endpoint:         aws.String(s3Endpoint),
		Region:           aws.String(s3Region),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
		MaxRetries:       aws.Int(3),
	}

	sess := session.Must(session.NewSession(config))
	uploader = s3manager.NewUploader(sess)
}

func (r *Router) UploadFile(c *gin.Context) {
	channelID, ok := c.Request.Context().Value(common.ChannelKey).(uint64)
	if !ok {
		response(c, http.StatusUnauthorized, common.ErrUnauthorized)
		return
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Error(err)
		response(c, http.StatusBadRequest, ErrReceiveFile)
		return
	}

	f, err := fileHeader.Open()
	if err != nil {
		log.Error(err)
		response(c, http.StatusBadRequest, ErrOpenFile)
		return
	}

	extension := filepath.Ext(fileHeader.Filename)
	newFileName := newObjectKey(channelID, extension)
	if err := putFileToS3(c.Request.Context(), s3Bucket, newFileName, f); err != nil {
		log.Error(err)
		response(c, http.StatusServiceUnavailable, ErrUploadFile)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"file_name": fileHeader.Filename,
		"file_url":  joinStrs(s3Endpoint, "/", s3Bucket, "/", newFileName),
	})
}

func putFileToS3(ctx context.Context, bucket, fileName string, f io.Reader) error {
	_, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		ACL:    aws.String("public-read"),
		Body:   f,
	})
	if err != nil {
		return err
	}
	return nil
}

func newObjectKey(channelID uint64, extension string) string {
	return joinStrs(strconv.FormatUint(channelID, 10), "/", uuid.New().String(), extension)
}

func joinStrs(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}
