package weed

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
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
	c := cache.Client()
	key := "weed_volumeid:" + volumeId
	if r, err := c.Get(key).Bytes(); err == nil {
		v := &VolumeLookupInfo{}
		err := v.Unmarshal(r)
		if err != nil {
			return nil, common.NewWeedError(err)
		}
		return v, nil
	}

	res, err := LookupVolumeNoCache(volumeId)
	if err != nil {
		return nil, err
	}
	b, err := res.Marshal()
	if err != nil {
		return nil, common.NewWeedError(err)
	}

	if err := c.Set(key, b, 24 * time.Hour).Err(); err != nil {
		return nil, common.NewRedisError(err)
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
