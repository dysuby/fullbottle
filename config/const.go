package config

const (
	ApiName       = "fullbottle.api.v1"
	UserSrvName   = "fullbottle.srv.user."
	AuthSrvName   = "fullbottle.srv.auth"
	BottleSrvName = "fullbottle.srv.bottle"
	ShareSrvName  = "fullbottle.srv.share"
	UploadSrvName = "fullbottle.srv.upload"

	WeedName  = "fullbottle.weed"
	DBName    = "fullbottle.mysql"
	RedisName = "fullbottle.redis"
)

const AppIss = "github.com/vegchic/fullbottle"

const JwtTokenExpire = int64(60 * 60 * 24)

const AvatarMaxSize = 1 << 20 // 1mb

const DefaultCapacity = 1 << 30 // 1GB
const DefaultChunkSize = 1 << 20
const PreviewSizeLimit = 5 << 20 // 5mb

const (
	MaxMsgSendSize = 64 << 20 // client
	MaxMsgRecvSize = 64 << 20 // client
	MaxMsgSize     = 64 << 20 // server
)
