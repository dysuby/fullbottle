package handler

import (
	"context"
	pbAuth "github.com/vegchic/fullbottle/auth/proto/auth"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/log"
	"github.com/vegchic/fullbottle/config"
	"github.com/vegchic/fullbottle/models"
	"github.com/vegchic/fullbottle/user/dao"
	pb "github.com/vegchic/fullbottle/user/proto/user"
	"github.com/vegchic/fullbottle/user/util"
	"github.com/micro/go-micro/v2/errors"
)

type UserHandler struct{}

func (u *UserHandler) GetUserInfo(ctx context.Context, req *pb.GetUserRequest, resp *pb.GetUserResponse) error {
	uid := req.GetId()
	result, err := dao.GetUsersByQuery(models.Fields{"id": uid})
	if err != nil {
		log.Errorf(err, "DB error")
		return common.NewDBError(config.UserSrvName, err)
	}
	if len(result) == 0 {
		return errors.New(config.UserSrvName, "User not found", common.UserNotFound)
	}

	user := result[0]

	var deleteTime int64
	if user.DeleteTime != nil {
		deleteTime = user.DeleteTime.Unix()
	}

	resp.Id = int64(user.ID)
	resp.Username = user.Username
	resp.Email = user.Email
	resp.Role = user.Role
	resp.AvatarUrl = user.AvatarUrl
	resp.Status = user.Status
	resp.CreateTime = user.CreateTime.Unix()
	resp.UpdateTime = user.UpdateTime.Unix()
	resp.DeleteTime = deleteTime

	return nil
}

func (u *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest, resp *pb.CreateUserResponse) error {
	rows, err := dao.GetUsersByQuery(models.Fields{"email": req.Email})
	if err != nil {
		log.Errorf(err, "DB error")
		return common.NewDBError(config.UserSrvName, err)
	} else if len(rows) > 0 {
		return errors.New(config.UserSrvName, "Email existed", common.EmailExisted)
	}

	user := models.User{
		Email:    req.Email,
		Username: req.Username,
		Password: util.PasswordCrypto(req.Password),
	}

	err = dao.CreateUser(&user)

	if err != nil {
		log.Errorf(err, "DB error")
		return common.NewDBError(config.UserSrvName, err)
	}

	return nil
}

func (u *UserHandler) ModifyUser(ctx context.Context, req *pb.ModifyUserRequest, resp *pb.ModifyUserResponse) error {
	uid := req.GetId()

	rows, err := dao.GetUsersByQuery(models.Fields{"id": uid})
	if err != nil {
		log.Errorf(err, "DB error")
		return common.NewDBError(config.UserSrvName, err)
	}

	if len(rows) == 0 {
		return errors.New(config.UserSrvName, "User not found", common.UserNotFound)
	}

	fields := models.Fields{
		"username":   req.Username,
		"password":   util.PasswordCrypto(req.Password),
		"avatar_url": req.AvatarUrl,
	}

	basicUser := models.User{}
	basicUser.ID = uid
	err = dao.UpdateUser(&basicUser, fields)

	if err != nil {
		log.Errorf(err, "DB error")
		return common.NewDBError(config.UserSrvName, err)
	}

	return nil
}

func (u *UserHandler) UserLogin(ctx context.Context, req *pb.UserLoginRequest, resp *pb.UserLoginResponse) error {
	email := req.GetEmail()
	result, err := dao.GetUsersByQuery(models.Fields{"email": email})
	if err != nil {
		log.Errorf(err, "DB error")
		return common.NewDBError(config.UserSrvName, err)
	}
	if len(result) == 0 {
		return errors.New(config.UserSrvName, "User not found", common.UserNotFound)
	}

	user := result[0]
	if pass := util.ComparePassword(user.Password, req.Password); !pass {
		return errors.New(config.UserSrvName, "Password error", common.PasswordError)
	}

	authClient := common.GetAuthSrvClient()
	authResp, err := authClient.GenerateJwtToken(ctx, &pbAuth.GenerateJwtTokenRequest{
		UserId: user.ID,
		Expire: config.JwtTokenExpire,
	})
	if err != nil {
		return err
	}

	resp.Token, resp.Expire = authResp.GetToken(), config.JwtTokenExpire
	return nil
}
