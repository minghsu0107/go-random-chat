package uploader

import "errors"

var (
	ErrOpenFile    = errors.New("fail to open file")
	ErrReceiveFile = errors.New("no file is received")
	ErrUploadFile  = errors.New("fail to upload file")
)
