package weed

import (
	"errors"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/log"
	"net/http"
	"net/url"
)

func DeleteFile(fid string) (err error) {
	defer func() {
		log.Infof("Weed: delete fid: %s, err = %v", fid, err)
	}()

	f, err := ParseFid(fid)
	if err != nil {
		return common.NewWeedError(errors.New("invalid fid"))
	}

	volume, err := LookupVolume(f.VolumeId)
	if err != nil {
		return err
	}

	base := url.URL{
		Scheme: "http",
		Host:   volume.Locations[0].Url,
		Path:   fid,
	}

	req, err := http.NewRequest("DELETE", base.String(), nil)
	if err != nil {
		return common.NewWeedError(err)
	}

	_, err = client.Do(req)
	if err != nil {
		return common.NewWeedError(err)
	}

	return nil
}
