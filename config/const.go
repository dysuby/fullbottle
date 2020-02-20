package config

const (
	ApiName     = "fullbottle.api"
	UserSrvName = "fullbottle.srv.user."
	AuthSrvName = "fullbottle.srv.auth"
)

const AppIss = "github.com/vegchic/fullbottle"

const JwtTokenExpire = int64(60 * 60 * 24)

const AvatarMaxSize = 2 * (1 << 20) // 2mb
