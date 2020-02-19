package handler

import (
	pbAuth "FullBottle/auth/proto/auth"
	"FullBottle/common"
	"FullBottle/models"
	"FullBottle/user/dao"
	pb "FullBottle/user/proto/user"
	"FullBottle/user/util"
	"context"
	"github.com/micro/go-micro/v2/util/log"
)

const JwtTokenExpire = int64(60 * 60 * 24)

type UserHandler struct{}

func (u *UserHandler) GetUserInfo(ctx context.Context, req *pb.GetUserRequest, resp *pb.GetUserResponse) error {
	uid := req.GetId()
	result, err := dao.GetUsersByQuery(models.Fields{"id": uid})
	if err != nil {
		log.Error(err)
		resp.Code, resp.Msg = common.DBConnError, "DB error"
		return nil
	}
	if len(result) == 0 {
		resp.Code, resp.Msg = common.UserNotFound, "User not found"
		return nil
	}

	user := result[0]

	var deleteTime int64
	if user.DeleteTime != nil {
		deleteTime = user.DeleteTime.Unix()
	}

	resp.Info = &pb.UserInfo{
		Id:         int64(user.ID),
		Username:   user.Username,
		Email:      user.Email,
		Role:       user.Role,
		AvatarUri:  user.AvatarUri,
		Status:     user.Status,
		CreateTime: user.CreateTime.Unix(),
		UpdateTime: user.UpdateTime.Unix(),
		DeleteTime: deleteTime,
	}

	resp.Code, resp.Msg = common.Success, "Success"

	return nil
}

func (u *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest, resp *pb.CreateUserResponse) error {
	info := req.GetVar()

	rows, err := dao.GetUsersByQuery(models.Fields{"email": info.Email})
	if err != nil {
		resp.Code, resp.Msg = common.DBConnError, "DB error"
		return nil
	} else if len(rows) > 0 {
		resp.Code, resp.Msg = common.EmailExisted, "Email existed"
		return nil
	}

	user := models.User{
		Email:     info.Email,
		Username:  info.Username,
		Password:  util.PasswordCrypto(info.Password),
		Role:      info.Role,
		AvatarUri: info.AvatarUri,
	}
	err = dao.CreateUser(&user)

	if err != nil {
		log.Error(err)
		resp.Code, resp.Msg = common.DBConnError, "DB error"
		return nil
	}

	resp.Code, resp.Msg = common.Success, "Success"
	return nil
}

func (u *UserHandler) ModifyUser(ctx context.Context, req *pb.ModifyUserRequest, resp *pb.ModifyUserResponse) error {
	info := req.GetVar()
	uid := req.GetId()

	rows, err := dao.GetUsersByQuery(models.Fields{"id": uid})
	if err != nil {
		log.Error(err)
		resp.Code, resp.Msg = common.DBConnError, "DB error"
		return nil
	}

	if len(rows) == 0 {
		resp.Code, resp.Msg = common.UserNotFound, "User not found"
		return nil
	}

	fields := models.Fields{
		"username":   info.Username,
		"password":   util.PasswordCrypto(info.Password),
		"avatar_uri": info.AvatarUri,
		"role":       info.Role,
	}

	basicUser := models.User{}
	basicUser.ID = uid
	err = dao.UpdateUser(&basicUser, fields)

	if err != nil {
		log.Error(err)
		resp.Code, resp.Msg = common.DBConnError, "DB error"
		return nil
	}

	resp.Code, resp.Msg = common.Success, "Success"
	return nil
}

func (u *UserHandler) UserLogin(ctx context.Context, req *pb.UserLoginRequest, resp *pb.UserLoginResponse) error {
	email := req.GetEmail()
	result, err := dao.GetUsersByQuery(models.Fields{"email": email})
	if err != nil {
		log.Error(err)
		resp.Code, resp.Msg = common.DBConnError, "DB error"
		return nil
	}
	if len(result) == 0 {
		resp.Code, resp.Msg = common.UserNotFound, "User not found"
		return nil
	}

	user := result[0]
	if pass := util.ComparePassword(user.Password, req.Password); !pass {
		resp.Code, resp.Msg = common.PasswordError, "Password error"
		return nil
	}

	authClient := common.GetAuthSrvClient()
	authResp, err := authClient.GenerateJwtToken(ctx, &pbAuth.GenerateJwtTokenRequest{
		UserId:user.ID,
		Expire:JwtTokenExpire,
	})
	if err != nil || authResp.Code != common.Success {
		resp.Code, resp.Msg = common.JwtError, "Cannot generate jwt token"
		return nil
	}

	resp.Code, resp.Msg = common.Success, "Success"
	resp.Token, resp.Expire = authResp.GetToken(), JwtTokenExpire
	return nil
}
