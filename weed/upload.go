package weed

import (
	"bytes"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/config"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
)

func UploadSingleFile(f io.Reader, name string, fid string, volumeUrl string, isManifest bool) (resp *http.Response, err error) {
	// read
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("file", name)
	if err != nil {
		return nil, common.NewWeedError(err)
	}
	if _, err = io.Copy(fw, f); err != nil {
		return nil, common.NewWeedError(err)
	}
	_ = w.Close()

	// upload
	base := &url.URL{
		Scheme: "http",
		Host:   volumeUrl,
		Path:   fid,
	}
	if isManifest {
		q := base.Query()
		q.Set("cm", "true")
		base.RawQuery = q.Encode()
	}
	req, err := http.NewRequest("POST", base.String(), &b)
	if err != nil {
		return nil, common.NewWeedError(err)
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err = client.Do(req)
	if err != nil {
		return nil, common.NewWeedError(err)
	}

	if !IsSuccessStatus(resp.StatusCode) {
		return nil, errors.New(config.WeedName, "weed return unexpected statuscode: "+strconv.Itoa(resp.StatusCode), common.WeedError)
	}

	return
}
