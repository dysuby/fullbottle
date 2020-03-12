package service

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/kv"
	"github.com/vegchic/fullbottle/config"
	"github.com/vegchic/fullbottle/util"
	"time"
)

type AccessToken struct {
	Id       int64
	SharerId int64
	ViewerId int64
	Token    string
}

func (at *AccessToken) init() {
	raw := fmt.Sprintf("%d:%d:%d:%d", at.Id, at.SharerId, at.ViewerId, time.Now().Unix())
	at.Token = util.Md5(raw)
}

func (at *AccessToken) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(*at)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (at *AccessToken) Unmarshal(b []byte) error {
	buf := bytes.NewReader(b)
	dec := gob.NewDecoder(buf)

	err := dec.Decode(at)
	if err != nil {
		return err
	}

	return nil
}

func NewAccessToken(id, sharerId, viewerId int64) *AccessToken {
	at := &AccessToken{
		Id:       id,
		SharerId: sharerId,
		ViewerId: viewerId,
	}
	at.init()
	return at
}

func ValidateAccessToken(accessToken string, viewerId int64) (AccessToken, error) {
	var at AccessToken
	if err := kv.GetM(accessToken, &at); err != nil {
		return at, err
	}

	if viewerId != at.ViewerId {
		return at, errors.New(config.ShareSrvName, "Invalid token", common.ConflictError)
	}
	return at, nil
}
