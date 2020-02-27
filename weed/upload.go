package weed

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/config"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
)

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

	Status int	// shouldn't modify directly

	Hash string
	Mime string
	Fid  string

	// chunk info
	ChunkManifest
}

func (f *FileUploadMeta) init() {
	raw := fmt.Sprintf("%d:%s:%d:%s", f.OwnerId, f.Name, f.FolderId, f.Hash)
	// todo sign
	f.Token = raw
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
func (f *FileUploadMeta) Upload(raw []byte, offset int64) error {
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
		f.Chunks = append(f.Chunks, &ChunkInfo{
			Fid:    key.Fid,
			Offset: offset,
			Size:   reader.Size(),
		})

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
		// upload
		key, err := AssignFileKey()
		if err != nil {
			return err
		}
		reader, err := f.Reader()
		if err != nil {
			return err
		}
		_, err = UploadSingleFile(reader, f.Name, key.Fid, key.Url, true)
		if err != nil {
			return err
		}

		// update info
		f.SetStatus(WeedDone)
		f.Fid = key.Fid
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
		ChunkManifest: ChunkManifest{Size: size, Name: filename},
	}
	meta.init()
	return meta
}

func UploadSingleFile(f io.Reader, name string, fid string, volumeUrl string, isManifest bool) (resp *http.Response, err error) {
	// read
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("file", name)
	if err != nil {
		return nil, common.NewWeedError(err)
	}
	if _, err = io.Copy(fw, f); err != nil {
		return nil, common.NewWeedError(err)
	}
	_ = w.Close()

	// upload
	base := &url.URL{
		Scheme: "http",
		Host:   volumeUrl,
		Path:   fid,
	}
	if isManifest {
		q := base.Query()
		q.Set("cm", "true")
		base.RawQuery = q.Encode()
	}
	req, err := http.NewRequest("POST", base.String(), &b)
	if err != nil {
		return nil, common.NewWeedError(err)
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err = client.Do(req)
	if err != nil {
		return nil, common.NewWeedError(err)
	}

	if !IsSuccessStatus(resp.StatusCode) {
		return nil, errors.New(config.WeedName, "weed return unexpected statuscode: "+strconv.Itoa(resp.StatusCode), common.WeedError)
	}

	return
}
