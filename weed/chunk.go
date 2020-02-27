package weed

import (
	"encoding/json"
	"github.com/vegchic/fullbottle/common"
)

type ChunkInfo struct {
	Fid    string `json:"fid"`
	Offset int64  `json:"offset"`
	Size   int64  `json:"size"`
}

type ChunkList []*ChunkInfo

type ChunkManifest struct {
	Name   string    `json:"name,omitempty"`
	Mime   string    `json:"mime,omitempty"`
	Size   int64     `json:"size,omitempty"`
	Chunks ChunkList `json:"chunks,omitempty"`
}

// generate manifest reader
func (m *ChunkManifest) Json() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, common.NewWeedError(err)
	}
	return b, nil
}
