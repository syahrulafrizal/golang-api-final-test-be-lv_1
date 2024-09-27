package redisrepo

import (
	"app/domain"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisRepo struct {
	Conn       *redis.Client
	Prefix     string
	UseRedis   bool
	DefaultTTL time.Duration
}

func NewRedisRepo(Conn *redis.Client) domain.CacheRepo {
	ttl, _ := time.ParseDuration(os.Getenv("REDIS_TTL"))
	// default ttl redis
	if ttl == 0 {
		ttl = 1 * time.Minute
	}

	useRedis, _ := strconv.ParseBool(os.Getenv("USE_REDIS"))
	redisKeyPrefix := os.Getenv("REDIS_KEY_PREFIX")

	return &redisRepo{
		Conn:       Conn,
		Prefix:     redisKeyPrefix + "appc:",
		UseRedis:   useRedis,
		DefaultTTL: ttl,
	}
}
