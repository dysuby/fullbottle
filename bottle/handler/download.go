package handler

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/bottle/dao"
	pb "github.com/vegchic/fullbottle/bottle/proto/bottle"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/kv"
	"github.com/vegchic/fullbottle/config"
	"github.com/vegchic/fullbottle/util"
	"github.com/vegchic/fullbottle/weed"
	"io"
	"time"
)

const DownloadTokenKey = "download:token=%s;user_id=%d"

type DownloadHandler struct{}

type DownloadURL struct {
	Url string
	Filename string
}

func (d *DownloadURL) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(*d)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (d *DownloadURL) Unmarshal(b []byte) error {
	buf := bytes.NewReader(b)
	dec := gob.NewDecoder(buf)

	err := dec.Decode(d)
	if err != nil {
		return err
	}

	return nil
}

func (*DownloadHandler) CreateDownloadUrl(ctx context.Context, req *pb.CreateDownloadUrlRequest, resp *pb.CreateDownloadUrlResponse) error {
	fileId := req.GetFileId()
	ownerId := req.GetOwnerId()
	userId := req.GetUserId()

	file, err := dao.GetFileById(ownerId, fileId)
	if err != nil {
		return err
	} else if file == nil {
		return errors.New(config.BottleSrvName, "File not found", common.NotFoundError)
	}

	fid := file.Metadata.Fid
	weedUrl, err := weed.GetDownloadUrl(fid)
	if err != nil {
		return err
	}

	downloadUrl := &DownloadURL{
		Url:      weedUrl.String(),
		Filename: file.Name,
	}

	token := util.GenToken(20)
	if err := kv.SetM(fmt.Sprintf(DownloadTokenKey, token, userId), downloadUrl, 5*time.Minute); err != nil {
		return err
	}

	resp.DownloadToken = token
	return nil
}

func (*DownloadHandler) GetWeedDownloadUrl(ctx context.Context, req *pb.GetWeedDownloadUrlRequest, resp *pb.GetWeedDownloadUrlResponse) error {
	token := req.GetDownloadToken()
	userId := req.GetUserId()

	downloadUrl := &DownloadURL{}
	err := kv.GetM(fmt.Sprintf(DownloadTokenKey, token, userId), downloadUrl)
	if err != nil {
		return err
	}

	resp.WeedUrl = downloadUrl.Url
	resp.Filename = downloadUrl.Filename
	return nil
}

func (*DownloadHandler) GetImageThumbnail(ctx context.Context, req *pb.GetImageThumbnailRequest, resp *pb.GetImageThumbnailResponse) error {
	fileId := req.GetFileId()
	ownerId := req.GetOwnerId()
	height := int(req.GetHeight())
	if height == 0 {
		height = 500
	}
	width := int(req.GetWidth())
	if width == 0 {
		width = 500
	}

	file, err := dao.GetFileById(ownerId, fileId)
	if err != nil {
		return err
	} else if file == nil {
		return errors.New(config.BottleSrvName, "File not found", common.NotFoundError)
	}
	if file.Size > config.PreviewSizeLimit {
		return errors.New(config.BottleSrvName, "File is too large", common.BadArgError)
	}

	cm := weed.ChunkManifest{}
	if err := json.Unmarshal([]byte(file.Metadata.ChunkManifest), &cm); err != nil {
		return errors.New(config.BottleSrvName, "File meta error", common.InternalError)
	}

	// comment here due to some dirty data
	//if !strings.HasPrefix(cm.Mime, "image") {
	//	return errors.New(config.BottleSrvName, "File isn't image", common.BadArgError)
	//}

	fid := file.Metadata.Fid

	imgResp, err := weed.FetchFile(fid)
	if err != nil {
		return err
	}

	body := imgResp.Body
	defer body.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, body); err != nil {
		return errors.New(config.BottleSrvName, "Avatar lost due to: "+err.Error(), common.FileFetchError)
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		return errors.New(config.BottleSrvName, "Cannot decode image", common.FileFetchError)
	}

	bounds := img.Bounds()
	var factor float32
	if bounds.Max.X < width {
		width = bounds.Max.X
	}
	if bounds.Max.Y < height {
		height = bounds.Max.Y
	}
	if bounds.Max.X < bounds.Max.Y {
		factor = float32(height) / float32(bounds.Max.Y)
	} else {
		factor = float32(width) / float32(bounds.Max.X)
	}
	th := imaging.Thumbnail(img, int(float32(bounds.Max.X)*factor), int(float32(bounds.Max.Y)*factor), imaging.Lanczos)

	reader := bytes.NewBuffer(nil)
	err = imaging.Encode(reader, th, imaging.JPEG)
	if err != nil {
		return errors.New(config.BottleSrvName, "Cannot encode image", common.FileFetchError)
	}

	resp.Thumbnail = reader.Bytes()
	return nil
}
