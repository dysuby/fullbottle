package handler

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/v2/errors"
	pbbottle "github.com/vegchic/fullbottle/bottle/proto/bottle"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/kv"
	"github.com/vegchic/fullbottle/common/log"
	"github.com/vegchic/fullbottle/config"
	"github.com/vegchic/fullbottle/util"
	"github.com/vegchic/fullbottle/weed"
	"time"

	pb "github.com/vegchic/fullbottle/upload/proto/upload"
)

const UploadLockKey = "lock:token=%s"

type UploadHandler struct{}

func (*UploadHandler) GenerateUploadToken(ctx context.Context, req *pb.GenerateUploadTokenRequest, resp *pb.GenerateUploadTokenResponse) error {
	ownerId := req.GetOwnerId()
	filename := req.GetFilename()
	folderId := req.GetFolderId()
	hash := req.GetHash()
	size := req.GetSize()
	mime := req.GetMime()

	// create upload meta
	upload := weed.NewUploadMeta(ownerId, folderId, filename, hash, size, mime)
	var history weed.FileUploadMeta
	if err := kv.GetM(upload.Token, &history); err == nil {
		upload = &history
	} else {
		// store for 15 days
		if err := kv.SetM(upload.Token, upload, 15*24*time.Hour); err != nil {
			return err
		}
	}

	resp.NeedUpload = false

	// check file is already uploaded, then client only need to call UploadFile without file data
	bottleClient := common.BottleSrvClient()
	metaResp, err := bottleClient.GetFileMeta(ctx, &pbbottle.GetFileMetaRequest{Hash: hash})
	if err != nil {
		return err
	} else if metaResp.Id == 0 {
		resp.NeedUpload = true // meta not found
	} else {
		fileResp, err := bottleClient.GetFileByMeta(ctx, &pbbottle.GetFileByMetaRequest{Name: filename, FolderId: folderId, OwnerId: ownerId, MetaId: metaResp.Id})
		if err != nil {
			return err
		} else if fileResp.File.Id != 0 {
			return errors.New(config.UploadSrvName, "File has existed", common.ExistedError)
		}
	}

	resp.Uploaded = upload.UploadedChunks()
	resp.Token = upload.Token
	return nil
}

func (*UploadHandler) UploadFile(ctx context.Context, req *pb.UploadFileRequest, resp *pb.UploadFileResponse) error {
	token := req.GetToken()

	// lock for upload meta
	lock, err := kv.Obtain(fmt.Sprintf(UploadLockKey, token), 100*time.Millisecond)
	if err != nil {
		return err
	}
	defer lock.Release()

	// fetch upload meta
	upload := &weed.FileUploadMeta{}
	if err := kv.GetM(token, upload); err != nil {
		return err
	}

	if req.GetOwnerId() != upload.OwnerId {
		return errors.New(config.UploadSrvName, "Incorrect owner", common.NotFoundError)
	}

	defer func() {
		resp.Status = pb.UploadStatus(upload.Status)
		resp.Uploaded = upload.UploadedChunks()

		// refresh token
		redisErr := kv.RefreshMValue(token, upload)
		if err == nil {
			err = redisErr
		}
	}()

	// check file is already uploaded
	bottleClient := common.BottleSrvClient()
	metaResp, err := bottleClient.GetFileMeta(ctx, &pbbottle.GetFileMetaRequest{Hash: upload.Hash})
	if err != nil {
		return err
	} else if metaResp.Id != 0 {
		upload.SetStatus(weed.WeedDone)
	} else {
		// upload chunk
		offset := req.GetOffset()
		raw := req.GetRaw()
		chunkHash := req.GetChunkHash()
		if uploaded, err := upload.CheckChunkOffset(offset, int64(len(raw))); err != nil {
			return err
		} else if uploaded {
			return errors.New(config.UploadSrvName, "The chunk has been uploaded", common.ChunkUploadedError)
		}

		hash := util.BytesMd5(raw)
		if hash != chunkHash {
			return errors.New(config.UploadSrvName, "The chunk hash is incorrect", common.FileUploadingError)
		}

		err = upload.Upload(raw, offset, hash)
		if err != nil {
			return err
		}
	}

	// upload weed done, create meta and file
	if upload.Status == weed.WeedDone {
		if metaResp.Id == 0 {
			b, _ := upload.ChunkManifest.Json()
			createMeta, err := bottleClient.CreateFileMeta(ctx, &pbbottle.CreateFileMetaRequest{
				Fid:           upload.Fid,
				Size:          upload.Size,
				Hash:          upload.Hash,
				ChunkManifest: string(b),
			})
			if err != nil {
				return err
			}
			metaResp.Id = createMeta.Id
		}

		_, err := bottleClient.CreateFile(ctx, &pbbottle.CreateFileRequest{
			OwnerId:  upload.OwnerId,
			FolderId: upload.FolderId,
			Name:     upload.Name,
			MetaId:   metaResp.Id,
		})
		if err != nil {
			return err
		}
		upload.SetStatus(weed.DBDone)
	}

	return nil
}

func (*UploadHandler) GetFileUploadedChunks(ctx context.Context, req *pb.GetFileUploadedChunksRequest, resp *pb.GetFileUploadedChunksResponse) error {
	token := req.GetToken()

	// fetch upload meta
	upload := &weed.FileUploadMeta{}
	if err := kv.GetM(token, upload); err != nil {
		return err
	}

	if req.GetOwnerId() != upload.OwnerId {
		return errors.New(config.UploadSrvName, "Incorrect owner", common.NotFoundError)
	}

	resp.NeedUpload = true

	bottleClient := common.BottleSrvClient()
	metaResp, err := bottleClient.GetFileMeta(ctx, &pbbottle.GetFileMetaRequest{Hash: upload.Hash})
	if err != nil {
		return err
	} else if metaResp.Id != 0 {
		resp.NeedUpload = false
	}

	resp.Uploaded = upload.UploadedChunks()

	return nil
}

func (*UploadHandler) CancelFileUpload(ctx context.Context, req *pb.CancelFileUploadRequest, resp *pb.CancelFileUploadResponse) error {
	token := req.GetToken()

	// lock for upload meta
	lock, err := kv.Obtain(fmt.Sprintf(UploadLockKey, token), 100*time.Millisecond)
	if err != nil {
		return err
	}
	defer lock.Release()

	// fetch upload meta
	upload := &weed.FileUploadMeta{}
	if err := kv.GetM(token, upload); err != nil {
		return err
	}

	if req.GetOwnerId() != upload.OwnerId {
		return errors.New(config.UploadSrvName, "Incorrect owner", common.NotFoundError)
	}

	for _, c := range upload.Chunks {
		if err := weed.DeleteFile(c.Fid); err != nil {
			log.WithCtx(ctx).WithError(err).Errorf("Cancel file upload failed")
		}
	}

	if err := kv.Del(token); err != nil {
		log.WithCtx(ctx).WithError(err).Errorf("Delete upload meta failed")
	}

	return nil
}
