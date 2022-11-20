package uploader

import (
	"context"
	"io"
	"net/http"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
	"github.com/minghsu0107/go-random-chat/pkg/common"
)

// @Summary Upload a file
// @Description Upload a file to S3 bucket
// @Tags uploader
// @Accept mpfd
// @param file formData file true "file to upload"
// @Produce json
// @param Authorization header string true "channel authorization"
// @Success 201 {object} gin.H
// @Failure 400 {object} common.ErrResponse
// @Failure 401 {object} common.ErrResponse
// @Failure 503 {object} common.ErrResponse
// @Router /uploader/file [post]
func (r *HttpServer) UploadFile(c *gin.Context) {
	channelID, ok := c.Request.Context().Value(common.ChannelKey).(uint64)
	if !ok {
		response(c, http.StatusUnauthorized, common.ErrUnauthorized)
		return
	}
	c.Request.ParseMultipartForm(r.maxMemory)
	fileHeader, err := c.FormFile("file")
	if err != nil {
		r.logger.Error(err)
		response(c, http.StatusBadRequest, ErrReceiveFile)
		return
	}

	f, err := fileHeader.Open()
	if err != nil {
		r.logger.Error(err)
		response(c, http.StatusBadRequest, ErrOpenFile)
		return
	}

	extension := filepath.Ext(fileHeader.Filename)
	newFileName := newObjectKey(channelID, extension)
	if err := r.putFileToS3(c.Request.Context(), r.s3Bucket, newFileName, f); err != nil {
		r.logger.Error(err)
		response(c, http.StatusInternalServerError, ErrUploadFile)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"file_name": fileHeader.Filename,
		"file_url":  joinStrs(r.s3Endpoint, "/", r.s3Bucket, "/", newFileName),
	})
}

func (r *HttpServer) putFileToS3(ctx context.Context, bucket, fileName string, f io.Reader) error {
	_, err := r.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		ACL:    types.ObjectCannedACLPublicRead,
		Body:   f,
	})
	if err != nil {
		return err
	}
	return nil
}
