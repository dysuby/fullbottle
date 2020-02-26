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

const FolderMaxLevel = 10 // 最大层数
const FolderMaxSub = 255  // 文件夹下同一层最多文件+文件夹数

const DefaultCapacity = 1 << 30 // 1GB

const (
	MaxMsgSendSize = 64 << 20 // client
	MaxMsgRecvSize = 64 << 20 // client
	MaxMsgSize     = 64 << 20 // server
)

const (
	ReadAction = "read"
	WriteAction = "write"
)