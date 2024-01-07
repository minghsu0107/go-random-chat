package uploader

import (
	"context"
	b64 "encoding/base64"
	"io"
	"net/http"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
	"github.com/minghsu0107/go-random-chat/pkg/common"
)

// @Summary Upload files (deprecated)
// @Description Upload files to S3 bucket (deprecated; use presigned urls instead)
// @Tags uploader
// @Accept mpfd
// @param files formData []file true "files to upload" collectionFormat(multi)
// @Produce json
// @param Authorization header string true "channel authorization"
// @Success 201 {object} UploadedFilesPresenter
// @Failure 400 {object} common.ErrResponse
// @Failure 401 {object} common.ErrResponse
// @Failure 500 {object} common.ErrResponse
// @Router /uploader/upload/files [post]
func (r *HttpServer) UploadFiles(c *gin.Context) {
	channelID, ok := c.Request.Context().Value(common.ChannelKey).(uint64)
	if !ok {
		response(c, http.StatusUnauthorized, common.ErrUnauthorized)
		return
	}
	if err := c.Request.ParseMultipartForm(r.maxMemory); err != nil {
		r.logger.Error("error parsing multipart form into memory: " + err.Error())
		response(c, http.StatusInternalServerError, common.ErrServer)
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		r.logger.Error("parse multipart form error: " + err.Error())
		response(c, http.StatusBadRequest, ErrReceiveFile)
		return
	}
	fileHeaders := form.File["files"]

	var uploadedFiles []UploadedFilePresenter

	for _, fileHeader := range fileHeaders {
		f, err := fileHeader.Open()
		if err != nil {
			r.logger.Error("error opening multipart file header: " + err.Error())
			response(c, http.StatusBadRequest, ErrOpenFile)
			return
		}

		extension := filepath.Ext(fileHeader.Filename)
		newFileName := newObjectKey(channelID, extension)
		if err := r.putFileToS3(c.Request.Context(), r.s3Bucket, newFileName, f); err != nil {
			r.logger.Error("error putting file to S3: ", err.Error())
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

// @Summary Get presigned upload url
// @Description Get presigned url for uploading a file to S3
// @Tags uploader
// @Produce json
// @Param ext query string true "file extension"
// @param Authorization header string true "channel authorization"
// @Success 200 {object} PresignedUpload
// @Failure 400 {object} common.ErrResponse
// @Failure 401 {object} common.ErrResponse
// @Failure 500 {object} common.ErrResponse
// @Router /uploader/upload/presigned [get]
func (r *HttpServer) GetPresignedUpload(c *gin.Context) {
	channelID, ok := c.Request.Context().Value(common.ChannelKey).(uint64)
	if !ok {
		response(c, http.StatusUnauthorized, common.ErrUnauthorized)
		return
	}
	var req GetPresignedUploadRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	objectKey := newObjectKey(channelID, common.Join(".", req.Extension))
	res, err := r.presigner.PutObject(c.Request.Context(), r.s3Bucket, objectKey)
	if err != nil {
		r.logger.Error("get presigned upload url failed: " + err.Error())
		response(c, http.StatusInternalServerError, common.ErrServer)
		return
	}

	c.JSON(http.StatusOK, &PresignedUpload{
		ObjectKey: objectKey,
		Url:       res.URL,
	})
}

// @Summary Get presigned download url
// @Description Get presigned url for downloading a file from S3
// @Tags uploader
// @Produce json
// @Param okb64 query string true "base64-encoded object key"
// @param Authorization header string true "channel authorization"
// @Success 200 {object} PresignedDownload
// @Failure 400 {object} common.ErrResponse
// @Failure 401 {object} common.ErrResponse
// @Failure 500 {object} common.ErrResponse
// @Router /uploader/download/presigned [get]
func (r *HttpServer) GetPresignedDownload(c *gin.Context) {
	channelID, ok := c.Request.Context().Value(common.ChannelKey).(uint64)
	if !ok {
		response(c, http.StatusUnauthorized, common.ErrUnauthorized)
		return
	}
	var req GetPresignedDownloadRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	objectKeyByte, err := b64.URLEncoding.DecodeString(req.ObjectKeyBase64)
	if err != nil {
		response(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	objectKey := byteSlice2String(objectKeyByte)
	targetChannelID, err := getChannelIDFromObjectKey(objectKey)
	if err != nil {
		response(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	if channelID != targetChannelID {
		response(c, http.StatusUnauthorized, common.ErrUnauthorized)
		return
	}

	res, err := r.presigner.GetObject(c.Request.Context(), r.s3Bucket, objectKey)
	if err != nil {
		r.logger.Error("get presigned download url failed: " + err.Error())
		response(c, http.StatusInternalServerError, common.ErrServer)
		return
	}

	c.JSON(http.StatusOK, &PresignedDownload{res.URL})
}
