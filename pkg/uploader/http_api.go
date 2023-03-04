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

// @Summary Upload files
// @Description Upload files to S3 bucket
// @Tags uploader
// @Accept mpfd
// @param files formData []file true "files to upload" collectionFormat(multi)
// @Produce json
// @param Authorization header string true "channel authorization"
// @Success 201 {object} gin.H
// @Failure 400 {object} common.ErrResponse
// @Failure 401 {object} common.ErrResponse
// @Failure 503 {object} common.ErrResponse
// @Router /uploader/files [post]
func (r *HttpServer) UploadFiles(c *gin.Context) {
	channelID, ok := c.Request.Context().Value(common.ChannelKey).(uint64)
	if !ok {
		response(c, http.StatusUnauthorized, common.ErrUnauthorized)
		return
	}
	if err := c.Request.ParseMultipartForm(r.maxMemory); err != nil {
		r.logger.Error(err)
		response(c, http.StatusInternalServerError, common.ErrServer)
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		r.logger.Error(err)
		response(c, http.StatusBadRequest, ErrReceiveFile)
		return
	}
	fileHeaders := form.File["files[]"]

	var uploadedFiles []UploadedFilePresenter

	for _, fileHeader := range fileHeaders {
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
		uploadedFiles = append(uploadedFiles, UploadedFilePresenter{
			Name: fileHeader.Filename,
			Url:  joinStrs(r.s3Endpoint, "/", r.s3Bucket, "/", newFileName),
		})
	}

	c.JSON(http.StatusCreated, &UploadedFilesPresenter{
		UploadedFiles: uploadedFiles,
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
