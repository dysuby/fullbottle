package weed

import (
	"errors"
	"strings"
)

type Fid struct {
	VolumeId   string
	FileId     string
	FileCookie string
}

func ParseFid(fid string) (Fid, error) {
	parts := strings.Split(fid, ",")
	if len(parts) == 2 && len(parts[1]) > 2 {
		return Fid{
			VolumeId:   parts[0],
			FileId:     parts[1][:2],
			FileCookie: parts[1][2:],
		}, nil
	}
	return Fid{}, errors.New("invalid fid format")
}

func IsSuccessStatus(statuscode int) bool {
	return statuscode >= 200 && statuscode < 300
}
