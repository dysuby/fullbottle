package weed

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/kv"
	"github.com/vegchic/fullbottle/config"
	"github.com/vegchic/fullbottle/util"
	"sort"
	"strings"
	"time"
)

const ChunkHashKey = "chunk:token=%s;offset=%d"

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

const (
	Inited = iota
	Uploading
	Manifest
	WeedDone
	DBDone
)

type FileUploadMeta struct {
	Token    string
	OwnerId  int64
	FolderId int64

	Status int // shouldn't modify directly

	Hash string
	Mime string
	Fid  string

	// chunk info
	ChunkManifest
	ChunkSize int64
}

func (f *FileUploadMeta) init() {
	raw := fmt.Sprintf("%d:%s:%d:%s", f.OwnerId, f.Name, f.FolderId, f.Hash)
	f.Token = util.Sha256(raw, config.C().App.Upload.Secret)
}

func (f *FileUploadMeta) SetStatus(s int) {
	// cannot rollback
	if s < f.Status {
		return
	}
	f.Status = s
}

func (f *FileUploadMeta) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(*f)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (f *FileUploadMeta) Unmarshal(b []byte) error {
	buf := bytes.NewReader(b)
	dec := gob.NewDecoder(buf)

	err := dec.Decode(f)
	if err != nil {
		return err
	}

	return nil
}

// always upload in chunk, just make it simple
func (f *FileUploadMeta) Upload(raw []byte, offset int64, hash string) error {
	// if done, return
	if f.Status == WeedDone || f.Status == DBDone {
		return nil
	}
	// if in upload chunk step
	if f.Status == Inited || f.Status == Uploading {
		// reset the status
		f.Status = Uploading

		// upload
		reader := bytes.NewReader(raw)
		key, err := AssignFileKey()
		if err != nil {
			return err
		}

		_, err = UploadSingleFile(reader, f.Name, key.Fid, key.Url, false)
		if err != nil {
			return err
		}

		// add chunk info
		chunk := &ChunkInfo{
			Fid:    key.Fid,
			Offset: offset,
			Size:   reader.Size(),
		}

		if err := f.StoreChunkHash(chunk, hash); err != nil {
			// TODO delete chunk from weed
			return err
		}
		f.Chunks = append(f.Chunks, chunk)

		// calculate if upload step is finish
		uploaded := int64(0)
		for _, c := range f.Chunks {
			uploaded += c.Size
		}
		if uploaded == f.Size {
			f.SetStatus(Manifest)
		}
	}

	// if in manifest step
	if f.Status == Manifest {
		// if file hash incorrect, clear all chunks
		if err := f.CheckFileHash(); err != nil {
			f.SetStatus(Uploading)
			for _, c := range f.Chunks {
				_ = DeleteFile(c.Fid)
			}
			f.Chunks = make([]*ChunkInfo, 0)
			return err
		}

		// upload
		key, err := AssignFileKey()
		if err != nil {
			return err
		}
		b, err := f.ChunkManifest.Json()
		if err != nil {
			return err
		}
		_, err = UploadSingleFile(bytes.NewBuffer(b), f.Name, key.Fid, key.Url, true)
		if err != nil {
			return err
		}

		// update info
		f.SetStatus(WeedDone)
		f.Fid = key.Fid
	}
	return nil
}

func (f *FileUploadMeta) UploadedChunks() []int64 {
	var ranges []int64
	for _, c := range f.Chunks {
		ranges = append(ranges, c.Offset)
	}
	return ranges
}

func (f *FileUploadMeta) CheckChunkOffset(offset int64, size int64) (uploaded bool, err error) {
	if offset >= f.Size {
		return false, common.NewWeedError(errors.New("invalid offset"))
	}

	for _, c := range f.Chunks {
		if c.Offset == offset {
			return true, nil
		}
	}

	if size != f.ChunkSize && offset+size < f.Size {
		return false, common.NewWeedError(errors.New("invalid chunk size"))
	}
	return false, nil
}

func (f *FileUploadMeta) StoreChunkHash(chunk *ChunkInfo, hash string) error {
	client := kv.C()
	if err := client.Set(fmt.Sprintf(ChunkHashKey, f.Token, chunk.Offset), hash, 15*24*time.Hour).Err(); err != nil {
		return common.NewRedisError(err)
	}
	return nil
}

func (f *FileUploadMeta) GetChunkHashFromCache(chunk *ChunkInfo) (string, error) {
	client := kv.C()
	if val, err := client.Get(fmt.Sprintf(ChunkHashKey, f.Token, chunk.Offset)).Result(); err != nil {
		return "", common.NewRedisError(err)
	} else {
		return val, nil
	}
}

func (f *FileUploadMeta) CheckFileHash() error {
	var hash strings.Builder
	hash.Grow(32 * len(f.Chunks)) // md5 hash

	// make a copy and sort by offset
	cs := make([]*ChunkInfo, len(f.Chunks))
	copy(cs, f.Chunks)
	sort.Slice(cs, func(i, j int) bool {
		return cs[i].Offset < cs[j].Offset
	})

	for _, c := range cs {
		h, err := f.GetChunkHashFromCache(c)
		if err != nil {
			return err
		}
		hash.WriteString(h)
	}

	total := util.Md5(hash.String())
	if total != f.Hash {
		return common.NewWeedError(errors.New("check total file hash failed"))
	}
	return nil
}

func NewUploadMeta(ownerId int64, folderId int64, filename string, hash string, size int64, mime string) *FileUploadMeta {
	meta := &FileUploadMeta{
		OwnerId:       ownerId,
		FolderId:      folderId,
		Hash:          hash,
		Mime:          mime,
		Status:        Inited,
		ChunkManifest: ChunkManifest{Size: size, Name: filename, Mime: mime},
		ChunkSize:     config.DefaultChunkSize,
	}
	meta.init()
	return meta
}
