package utils

import (
	"FullBottle/config"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func GetWeedFilerUrl() string {
	return config.GetSingleConfig("weed", "filer")
}

func GenFilePath(path string, filename string) string {
	return fmt.Sprintf("/%s/%s", path, filename)
}

func UploadFile(file multipart.File, filename string, path string) (resp *http.Response, err error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		return
	}
	if _, err = io.Copy(fw, file); err != nil {
		return
	}

	_ = w.Close()

	url := GetWeedFilerUrl() + GenFilePath(path, filename)
	req, err := http.NewRequest("POST",url, &b)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	client := http.DefaultClient
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	return
}
