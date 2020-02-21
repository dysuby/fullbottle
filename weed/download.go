package weed

import (
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/config"
	"net/http"
	"net/url"
	"strconv"
)

func FetchFile(fid string, volumeUrl string) (resp *http.Response, err error) {
	base := url.URL{
		Scheme: "http",
		Host:   volumeUrl,
		Path:   fid,
	}

	resp, err = HttpClient().Get(base.String())
	if err != nil {
		return nil, common.NewWeedError(err)
	}
	if !IsSuccessStatus(resp.StatusCode) {
		return nil, errors.New(config.WeedName, "weed return unexpected statuscode: "+strconv.Itoa(resp.StatusCode), common.WeedError)
	}
	return
}
