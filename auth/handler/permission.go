package handler

import (
	"context"
	pb "github.com/vegchic/fullbottle/auth/proto/auth"
	pbbottle "github.com/vegchic/fullbottle/bottle/proto/bottle"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/config"
)

type PermissionHandler struct {}

func (*PermissionHandler) CheckFolderAccess(ctx context.Context, req *pb.CheckFolderAccessRequest, resp *pb.CheckFolderAccessResponse) error {
	userId := req.GetUserId()
	folderId := req.GetFolderId()

	bottleClient := common.BottleSrvClient()
	ownerResp, err := bottleClient.GetEntryOwner(ctx, &pbbottle.GetEntryOwnerRequest{EntryId:folderId, IsFolder:true})
	if err != nil {
		return err
	}

	// check logic
	if ownerResp.OwnerId != userId {
		resp.Actions = []string{}
	} else {
		resp.Actions = []string{config.ReadAction, config.WriteAction}
	}

	return nil
}
