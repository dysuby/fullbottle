package weed

import (
	"encoding/json"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/config"
	"net/url"
)

const (
	FilePath         = "/dir/assign"
	LookupVolumePath = "/dir/lookup"
)

type FileKeyInfo struct {
	Count     int    `json:"count"`
	Fid       string `json:"fid"`
	Url       string `json:"url"`
	PublicUrl string `json:"publicUrl"`
}

type VolumeLookupInfo struct {
	VolumeId  string `json:"volumeId"`
	Locations []struct {
		PublicUrl string `json:"publicUrl"`
		Url       string `json:"url"`
	} `json:"locations"`
}

func MasterUrl() (u *url.URL, err error) {
	u, err = url.Parse(config.C().Weed.Master)
	return
}

func AssignFileKey() (*FileKeyInfo, error) {
	base, err := MasterUrl()
	if err != nil {
		return nil, common.NewWeedError(err)
	}
	base.Path = FilePath

	resp, err := HttpClient().Get(base.String())
	if err != nil {
		return nil, common.NewWeedError(err)
	}
	defer resp.Body.Close()

	var key *FileKeyInfo
	err = json.NewDecoder(resp.Body).Decode(&key)
	if err != nil {
		return nil, common.NewWeedError(err)
	}

	return key, nil
}

func LookupVolume(volumeId string) (*VolumeLookupInfo, error) {
	base, err := MasterUrl()
	if err != nil {
		return nil, common.NewWeedError(err)
	}
	base.Path = LookupVolumePath
	q := base.Query()
	q.Set("volumeId", volumeId)
	base.RawQuery = q.Encode()

	resp, err := HttpClient().Get(base.String())
	if err != nil {
		return nil, common.NewWeedError(err)
	}

	defer resp.Body.Close()

	var info *VolumeLookupInfo
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return nil, common.NewWeedError(err)
	}

	if len(info.Locations) == 0 {
		return nil, errors.New(config.WeedName, "Volume lost", common.WeedError)
	}

	return info, nil
}
