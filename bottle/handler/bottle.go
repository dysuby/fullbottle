package handler

import (
	"context"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/bottle/dao"
	pb "github.com/vegchic/fullbottle/bottle/proto/bottle"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/db"
	"github.com/vegchic/fullbottle/config"
)

type BottleHandler struct{}

func (b *BottleHandler) InitBottle(ctx context.Context, req *pb.InitBottleRequest, resp *pb.InitBottleResponse) error {
	uid := req.GetUid()

	if bottle, err := dao.GetBottlesByUserId(uid); err != nil {
		return err
	} else if bottle != nil {
		return errors.New(config.BottleSrvName, "Bottle existed", common.InternalError)
	}

	bottle := dao.BottleMeta{
		UserID:   uid,
		Capacity: req.GetCapacity(),
		Remain:   req.GetCapacity(),
	}
	return dao.InitBottle(&bottle)
}

func (b *BottleHandler) GetBottleMetadata(ctx context.Context, req *pb.GetBottleMetadataRequest, resp *pb.GetBottleMetadataResponse) error {
	uid := req.GetUid()
	bottle, err := dao.GetBottlesByUserId(uid)
	if err != nil {
		return err
	} else if bottle == nil {
		bottle = &dao.BottleMeta{
			UserID:   uid,
			Capacity: config.DefaultCapacity,
			Remain:   config.DefaultCapacity,
		}
		err := dao.InitBottle(bottle)
		if err != nil {
			return err
		}
	}

	resp.Bid = bottle.ID
	resp.Capacity = bottle.Capacity
	resp.Remain = bottle.Remain
	resp.RootId = bottle.RootID

	return nil
}

func (b *BottleHandler) UpdateBottle(ctx context.Context, req *pb.UpdateBottleRequest, resp *pb.UpdateBottleResponse) error {
	bid := req.GetBid()

	bottle, err := dao.GetBottlesById(bid)
	if err != nil {
		return err
	} else if bottle != nil {
		return errors.New(config.BottleSrvName, "Bottle existed", common.InternalError)
	}

	return dao.UpdateBottle(bottle, db.Fields{
		"capacity": req.GetCapacity(),
	})
}

func (b *BottleHandler) GetEntryOwner(ctx context.Context, req *pb.GetEntryOwnerRequest, resp *pb.GetEntryOwnerResponse) error {
	// ugly code
	eid := req.GetEntryId()

	if req.IsFolder {
		if f, err := dao.GetFolderOwner(eid); err != nil {
			return err
		} else if f == nil {
			return errors.New(config.BottleSrvName, "Folder not found", common.NotFoundError)
		} else {
			resp.OwnerId = f.OwnerId
			return nil
		}
	} else {
		if f, err := dao.GetFileOwner(eid); err != nil {
			return err
		} else if f == nil {
			return errors.New(config.BottleSrvName, "Folder not found", common.NotFoundError)
		} else {
			resp.OwnerId = f.OwnerId
			return nil
		}
	}
}
