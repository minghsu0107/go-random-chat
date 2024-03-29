package uploader

type UploadedFilePresenter struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type UploadedFilesPresenter struct {
	UploadedFiles []UploadedFilePresenter `json:"uploaded_files"`
}

type GetPresignedUploadRequest struct {
	Extension string `form:"ext" binding:"required"`
}

type GetPresignedDownloadRequest struct {
	ObjectKeyBase64 string `form:"okb64" binding:"required"`
}

type PresignedUpload struct {
	ObjectKey string `json:"object_key"`
	Url       string `json:"url"`
}

type PresignedDownload struct {
	Url string `json:"url"`
}
