package weed

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/cache"
	"github.com/vegchic/fullbottle/config"
	"net/url"
	"time"
)

const (
	FilePath         = "/dir/assign"
	LookupVolumePath = "/dir/lookup"

	VolumeCacheKey   = "weed:volumeid=%s"
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

func (v *VolumeLookupInfo) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(*v)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (v *VolumeLookupInfo) Unmarshal(b []byte) error {
	buf := bytes.NewReader(b)
	dec := gob.NewDecoder(buf)

	err := dec.Decode(v)
	if err != nil {
		return err
	}

	return nil
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

	resp, err := client.Get(base.String())
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
	v := &VolumeLookupInfo{}
	key := fmt.Sprintf(VolumeCacheKey, volumeId)
	if err := cache.Get(key, v); err == nil {
		return v, nil
	}

	res, err := LookupVolumeNoCache(volumeId)
	if err != nil {
		return nil, err
	}

	if err := cache.Set(key, res, 24*time.Hour); err != nil {
		return nil, err
	}
	return res, nil
}

func LookupVolumeNoCache(volumeId string) (*VolumeLookupInfo, error) {
	base, err := MasterUrl()
	if err != nil {
		return nil, common.NewWeedError(err)
	}
	base.Path = LookupVolumePath
	q := base.Query()
	q.Set("volumeId", volumeId)
	base.RawQuery = q.Encode()

	resp, err := client.Get(base.String())
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
