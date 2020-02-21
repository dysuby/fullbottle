package config

const (
	ApiName     = "fullbottle.api"
	UserSrvName = "fullbottle.srv.user."
	AuthSrvName = "fullbottle.srv.auth"
	WeedName    = "fullbottle.weed"
	DBName      = "fullbottle.mysql"
)

const AppIss = "github.com/vegchic/fullbottle"

const JwtTokenExpire = int64(60 * 60 * 24)

const AvatarMaxSize = 1 << 20 // 1mb
