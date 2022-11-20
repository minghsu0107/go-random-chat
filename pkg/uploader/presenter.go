package uploader

type UploadedFilePresenter struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type UploadedFilesPresenter struct {
	UploadedFiles []UploadedFilePresenter `json:"uploaded_files"`
}
