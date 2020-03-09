package weed

import (
	"errors"
	"github.com/vegchic/fullbottle/common"
	"net/http"
	"net/url"
	"strconv"
)

func FetchFile(fid string) (resp *http.Response, err error) {
	f, err := ParseFid(fid)
	if err != nil {
		return nil, common.NewWeedError(errors.New("invalid avatar fid"))
	}

	volume, err := LookupVolume(f.VolumeId)
	if err != nil {
		return nil, err
	}

	base := url.URL{
		Scheme: "http",
		Host:   volume.Locations[0].Url,
		Path:   fid,
	}

	resp, err = client.Get(base.String())
	if err != nil {
		return nil, common.NewWeedError(err)
	}
	if !IsSuccessStatus(resp.StatusCode) {
		return nil, common.NewWeedError(errors.New("weed return unexpected statuscode: "+strconv.Itoa(resp.StatusCode)))
	}
	return
}

func GetDownloadUrl(fid string) (*url.URL, error) {
	parsed, err := ParseFid(fid)
	if err != nil {
		return nil, err
	}

	volume, err := LookupVolume(parsed.VolumeId)
	if err != nil {
		return nil, err
	}

	base := &url.URL{
		Scheme: "http",
		Host:   volume.Locations[0].Url,
		Path:   fid,
	}
	return base, nil
}
