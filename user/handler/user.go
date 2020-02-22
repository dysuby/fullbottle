package handler

import (
	"bytes"
	"context"
	"fmt"
	"github.com/micro/go-micro/v2/errors"
	"github.com/sirupsen/logrus"
	pbauth "github.com/vegchic/fullbottle/auth/proto/auth"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/db"
	"github.com/vegchic/fullbottle/common/log"
	"github.com/vegchic/fullbottle/config"
	"github.com/vegchic/fullbottle/user/dao"
	pb "github.com/vegchic/fullbottle/user/proto/user"
	"github.com/vegchic/fullbottle/user/util"
	"github.com/vegchic/fullbottle/weed"
	"io/ioutil"
	"time"
)

type UserHandler struct{}

func (u *UserHandler) GetUserInfo(ctx context.Context, req *pb.GetUserRequest, resp *pb.GetUserResponse) error {
	uid := req.GetUid()
	result, err := dao.GetUsersByQuery(db.Fields{"id": uid})
	if err != nil {
		return err
	} else if len(result) == 0 {
		return errors.New(config.UserSrvName, "User not found", common.UserNotFound)
	}

	user := result[0]

	resp.Uid = user.ID
	resp.Username = user.Username
	resp.Email = user.Email
	resp.Role = user.Role
	resp.AvatarFid = user.AvatarFid
	resp.Status = user.Status
	resp.CreateTime = user.CreateTime.Unix()
	resp.UpdateTime = user.UpdateTime.Unix()
	if user.DeleteTime != nil {
		resp.DeleteTime = user.DeleteTime.Unix()
	}

	return nil
}

func (u *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest, resp *pb.CreateUserResponse) error {
	rows, err := dao.GetUsersByQuery(db.Fields{"email": req.Email})
	if err != nil {
		return err
	} else if len(rows) > 0 {
		return errors.New(config.UserSrvName, "Email existed", common.EmailExisted)
	}

	user := dao.User{
		Email:    req.Email,
		Username: req.Username,
		Password: util.PasswordCrypto(req.Password),
	}

	if err = dao.CreateUser(&user); err != nil {
		return err
	}

	return nil
}

func (u *UserHandler) ModifyUser(ctx context.Context, req *pb.ModifyUserRequest, resp *pb.ModifyUserResponse) error {
	uid := req.GetUid()

	rows, err := dao.GetUsersByQuery(db.Fields{"id": uid})
	if err != nil {
		return err
	} else if len(rows) == 0 {
		return errors.New(config.UserSrvName, "User not found", common.UserNotFound)
	}

	fields := db.Fields{
		"username": req.Username,
		"password": util.PasswordCrypto(req.Password),
	}

	basicUser := dao.User{}
	basicUser.ID = uid

	if err = dao.UpdateUser(&basicUser, fields); err != nil {
		return err
	}

	return nil
}

func (u *UserHandler) UserLogin(ctx context.Context, req *pb.UserLoginRequest, resp *pb.UserLoginResponse) error {
	email := req.GetEmail()
	result, err := dao.GetUsersByQuery(db.Fields{"email": email})
	if err != nil {
		return err
	}
	if len(result) == 0 {
		return errors.New(config.UserSrvName, "User not found", common.UserNotFound)
	}

	user := result[0]
	if pass := util.ComparePassword(user.Password, req.Password); !pass {
		return errors.New(config.UserSrvName, "Password error", common.PasswordError)
	}

	authClient := common.AuthSrvClient()
	authResp, err := authClient.GenerateJwtToken(ctx, &pbauth.GenerateJwtTokenRequest{
		UserId: user.ID,
		Expire: config.JwtTokenExpire,
	})
	if err != nil {
		return err
	}

	resp.Token, resp.Expire = authResp.GetToken(), config.JwtTokenExpire
	return nil
}

func (u *UserHandler) GetUserAvatar(ctx context.Context, req *pb.GetUserAvatarRequest, resp *pb.GetUserAvatarResponse) error {
	rows, err := dao.GetUsersByQuery(db.Fields{"id": req.GetUid()})
	if err != nil {
		return err
	} else if len(rows) == 0 {
		return errors.New(config.UserSrvName, "User not found", common.UserNotFound)
	}

	user := rows[0]
	avatarFid := user.AvatarFid
	if avatarFid == "" {
		return errors.New(config.UserSrvName, "User has no avatar", common.EmptyAvatarError)
	}

	fid, err := weed.ParseFid(avatarFid)
	if err != nil {
		return errors.New(config.UserSrvName, "Invalid avatar fid", common.InternalError)
	}

	volume, err := weed.LookupVolume(fid.VolumeId)
	if err != nil {
		return err
	}
	avatarResp, err := weed.FetchFile(avatarFid, volume.Locations[0].Url)
	if err != nil {
		return err
	}

	body := avatarResp.Body
	defer body.Close()

	avatar, err := ioutil.ReadAll(body)
	if err != nil {
		return errors.New(config.UserSrvName, "Avatar lost due to: "+err.Error(), common.FileFetchError)
	}

	resp.Avatar = avatar
	resp.ContentType = avatarResp.Header.Get("Content-Type")
	return nil
}

func (u *UserHandler) UploadUserAvatar(ctx context.Context, req *pb.UploadUserAvatarRequest, resp *pb.UploadUserAvatarResponse) error {
	rows, err := dao.GetUsersByQuery(db.Fields{"id": req.GetUid()})
	if err != nil {
		return err
	} else if len(rows) == 0 {
		return errors.New(config.UserSrvName, "User not found", common.UserNotFound)
	}
	user := rows[0]

	var volume *weed.VolumeLookupInfo
	var volumeUrl string
	fid := user.AvatarFid
	if user.AvatarFid != "" {
		f, err := weed.ParseFid(user.AvatarFid)
		if err == nil {
			volume, err = weed.LookupVolume(f.VolumeId)
		}
		if err != nil {
			fid = "" // reset fid
			log.WithFields(logrus.Fields{
				"uid": user.ID,
			}).WithError(errors.New(config.UserSrvName, "Dirty avatar fid", common.InternalError))
		}
	}

	if fid != "" {
		volumeUrl = volume.Locations[0].Url
	} else {
		key, err := weed.AssignFileKey()
		if err != nil {
			return err
		}
		fid = key.Fid
		volumeUrl = key.Url
	}

	// write db first
	if fid != user.AvatarFid {
		err = dao.UpdateUser(&user, db.Fields{
			"avatar_fid": fid,
		})
		if err != nil {
			return err
		}
	}

	_, err = weed.UploadFile(bytes.NewReader(req.Avatar), fmt.Sprint(user.ID, "-", time.Now().Unix()), fid, volumeUrl)
	return err
}
