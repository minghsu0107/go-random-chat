package uploader

import (
	"context"
	"io"
	"net/http"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/minghsu0107/go-random-chat/pkg/common"
	log "github.com/sirupsen/logrus"
)

func (r *HttpServer) UploadFile(c *gin.Context) {
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
	if err := r.putFileToS3(c.Request.Context(), r.s3Bucket, newFileName, f); err != nil {
		log.Error(err)
		response(c, http.StatusServiceUnavailable, ErrUploadFile)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"file_name": fileHeader.Filename,
		"file_url":  joinStrs(r.s3Endpoint, "/", r.s3Bucket, "/", newFileName),
	})
}

func (r *HttpServer) putFileToS3(ctx context.Context, bucket, fileName string, f io.Reader) error {
	_, err := r.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
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
