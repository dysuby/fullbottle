package utils

import (
	"FullBottle/config"
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

const (
	WeedFileKey = "dir/assign"
)

type FileInfo struct {
	Size      int
}

type FileKeyInfo struct {
	Count     int      `json:"count"`
	Fid       string   `json:"fid"`
	Url       string   `json:"url"`
	PublicUrl string   `json:"public_url"`
	FileInfo  FileInfo `json:"-"`
}

func JoinUrl(url string, path string) string {
	if strings.HasPrefix(url, "http") {
		return strings.Join([]string{url, "/", path}, "")
	}
	return strings.Join([]string{"http://", url, "/", path}, "")

}

func getFileKey() (f FileKeyInfo, err error) {
	url := JoinUrl(config.GetConfig().Weed.Master, WeedFileKey)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&f)
	return
}

func UploadFile(file io.Reader, filename string) (info FileKeyInfo, err error) {
	client := http.DefaultClient

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

	info, err = getFileKey()
	if err != nil {
		return
	}

	url := JoinUrl(info.Url, info.Fid)
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}
	err = json.NewDecoder(resp.Body).Decode(&info.FileInfo)
	return
}
