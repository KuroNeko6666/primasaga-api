package form

import "mime/multipart"

type FileForm struct {
	FileName   string
	FileBuffer multipart.File
	FileSize   int64
	FileType   string
}
