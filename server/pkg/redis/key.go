package redis

import "fmt"

const RedisKeyPrefix = "mewsoproxy:"

const (
	RedisKeyUserToken   = RedisKeyPrefix + "user:token:"
	RedisKeyRefresh     = RedisKeyPrefix + "refresh_token:"
	RedisKeyTokenBlk    = RedisKeyPrefix + "token_blacklist:"
	RedisKeyIdempotent  = RedisKeyPrefix + "idempotent:"
	RedisKeyRateLimit   = RedisKeyPrefix + "rate_limit:"
	RedisKeyUserRefresh = RedisKeyPrefix + "user:refresh:"
)

func UserTokenKey(userID uint) string {
	return fmt.Sprintf("%s%d", RedisKeyUserToken, userID)
}

func RefreshTokenKey(userID uint, tokenID string) string {
	return fmt.Sprintf("%s%d:%s", RedisKeyRefresh, userID, tokenID)
}

func TokenBlacklistKey(jti string) string {
	return fmt.Sprintf("%s%s", RedisKeyTokenBlk, jti)
}

func IdempotentKey(requestID string) string {
	return fmt.Sprintf("%s%s", RedisKeyIdempotent, requestID)
}

func RateLimitKey(bucket, identity string) string {
	return fmt.Sprintf("%s%s:%s", RedisKeyRateLimit, bucket, identity)
}

func UserRefreshKey(userID uint) string {
	return fmt.Sprintf("%s%d", RedisKeyUserRefresh, userID)
}

func RefreshIndexKey(refreshToken string) string {
	return fmt.Sprintf("%si:%s", RedisKeyRefresh, refreshToken)
}
