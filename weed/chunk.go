package weed

import (
	"bytes"
	"encoding/json"
	"github.com/vegchic/fullbottle/common"
	"io"
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
func (m *ChunkManifest) Reader() (io.Reader, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, common.NewWeedError(err)
	}
	return bytes.NewReader(b), nil
}
