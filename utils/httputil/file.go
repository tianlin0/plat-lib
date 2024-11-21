package httputil

import (
	"encoding/base64"
	"mime/multipart"
	"path"
	"strings"
)

// fileInfo 返回文件信息
type fileInfo struct {
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
	ExtName     string `json:"ext_name"`
	Size        int64  `json:"size"`
	Data        string `json:"data"`
}

// GetUploadFileBase64 获取上传的文件为base64
func GetUploadFileBase64(f *multipart.FileHeader) (*fileInfo, error) {
	if f == nil {
		return nil, nil
	}
	fileExt := strings.ToLower(path.Ext(f.Filename))
	fileExt = strings.ReplaceAll(fileExt, ".", "")
	fileInfoTemp := &fileInfo{
		Name:        f.Filename,
		ContentType: f.Header.Get("Content-Type"),
		ExtName:     fileExt,
		Size:        f.Size,
		Data:        "",
	}

	if f.Size > 0 {
		file, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer func(file multipart.File) {
			_ = file.Close()
		}(file)

		sourceBuffer := make([]byte, f.Size)
		n, _ := file.Read(sourceBuffer)
		sourceString := base64.StdEncoding.EncodeToString(sourceBuffer[:n])
		fileInfoTemp.Data = sourceString
	}
	return fileInfoTemp, nil
}
