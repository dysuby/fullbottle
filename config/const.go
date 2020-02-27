package config

const (
	ApiName       = "fullbottle.api.v1"
	UserSrvName   = "fullbottle.srv.user."
	AuthSrvName   = "fullbottle.srv.auth"
	BottleSrvName = "fullbottle.srv.bottle"

	WeedName  = "fullbottle.weed"
	DBName    = "fullbottle.mysql"
	RedisName = "fullbottle.redis"
)

const AppIss = "github.com/vegchic/fullbottle"

const JwtTokenExpire = int64(60 * 60 * 24)

const AvatarMaxSize = 1 << 20 // 1mb

const DefaultCapacity = 1 << 30 // 1GB
const NeedChunkLimit = 4 << 20

const (
	MaxMsgSendSize = 64 << 20 // client
	MaxMsgRecvSize = 64 << 20 // client
	MaxMsgSize     = 64 << 20 // server
)
