package handler

import (
	"bytes"
	"context"
	"fmt"
	"github.com/micro/go-micro/v2/errors"
	"github.com/sirupsen/logrus"
	pbauth "github.com/vegchic/fullbottle/auth/proto/auth"
	pbbottle "github.com/vegchic/fullbottle/bottle/proto/bottle"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/kv"
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

const UserInfoKey = "info:usr_id=%d"

type UserHandler struct{}

func (u *UserHandler) GetUserInfo(ctx context.Context, req *pb.GetUserRequest, resp *pb.GetUserResponse) error {
	uid := req.GetUid()
	user := &dao.User{}
	key := fmt.Sprintf(UserInfoKey, uid)
	if err := kv.Get(key, user); err != nil {
		user, err := dao.GetUsersById(uid)
		if err != nil {
			return err
		} else if user == nil {
			return errors.New(config.UserSrvName, "User not found", common.NotFoundError)
		}
		_ = kv.Set(key, user, 24*time.Hour)
	}

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
	if u, err := dao.GetUsersByEmail(req.GetEmail()); err != nil {
		return err
	} else if u != nil {
		return errors.New(config.UserSrvName, "Email existed", common.ExistedError)
	}

	user := &dao.User{
		Email:    req.Email,
		Username: req.Username,
		Password: util.PasswordCrypto(req.Password),
	}

	err := dao.CreateUser(user)
	if err != nil {
		return err
	}

	bottleClient := common.BottleSrvClient()
	_, err = bottleClient.InitBottle(ctx, &pbbottle.InitBottleRequest{Uid: user.ID, Capacity: config.DefaultCapacity})
	if err != nil {
		log.WithError(err).Errorf("Cannot init user bottle")
	}

	_ = kv.Get(fmt.Sprintf(UserInfoKey, user.ID), user)
	return nil
}

func (u *UserHandler) ModifyUser(ctx context.Context, req *pb.ModifyUserRequest, resp *pb.ModifyUserResponse) error {
	uid := req.GetUid()

	user, err := dao.GetUsersById(uid)
	if err != nil {
		return err
	} else if user == nil {
		return errors.New(config.UserSrvName, "User not found", common.NotFoundError)
	}

	fields := db.Fields{
		"username": req.Username,
		"password": util.PasswordCrypto(req.Password),
	}

	if err = dao.UpdateUser(user, fields); err != nil {
		return err
	}
	_ = kv.Del(fmt.Sprintf(UserInfoKey, uid))
	return nil
}

func (u *UserHandler) UserLogin(ctx context.Context, req *pb.UserLoginRequest, resp *pb.UserLoginResponse) error {
	email := req.GetEmail()
	user, err := dao.GetUsersByEmail(email)
	if err != nil {
		return err
	} else if user == nil {
		return errors.New(config.UserSrvName, "User not found", common.NotFoundError)
	}

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
	uid := req.GetUid()
	user := &dao.User{}
	key := fmt.Sprintf(UserInfoKey, uid)
	if err := kv.Get(key, user); err != nil {
		user, err := dao.GetUsersById(uid)
		if err != nil {
			return err
		} else if user == nil {
			return errors.New(config.UserSrvName, "User not found", common.NotFoundError)
		}
		_ = kv.Set(key, user, 24*time.Hour)
	}

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
	user, err := dao.GetUsersById(req.GetUid())
	if err != nil {
		return err
	} else if user == nil {
		return errors.New(config.UserSrvName, "User not found", common.NotFoundError)
	}

	var volume *weed.VolumeLookupInfo
	var volumeUrl string
	fid := user.AvatarFid
	// if already exist, then try to rewrite
	if user.AvatarFid != "" {
		f, err := weed.ParseFid(user.AvatarFid)
		if err == nil {
			volume, err = weed.LookupVolume(f.VolumeId)
		}
		if err != nil {
			fid = "" // reset fid
			log.WithFields(logrus.Fields{
				"userId": user.ID,
			}).WithError(errors.New(config.UserSrvName, "Dirty avatar fid", common.InternalError))
		}
	}

	// if failed, get new file key
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
		err = dao.UpdateUser(user, db.Fields{
			"avatar_fid": fid,
		})
		if err != nil {
			return err
		}
	}

	_, err = weed.UploadSingleFile(bytes.NewReader(req.Avatar), fmt.Sprint(user.ID, "-", time.Now().Unix()),
		fid, volumeUrl, false)
	_ = kv.Del(fmt.Sprintf(UserInfoKey, req.GetUid()))
	return err
}
