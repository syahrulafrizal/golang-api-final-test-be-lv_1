package middleware

import (
	"app/helpers"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type appMiddleware struct {
	secret string
	cache  CacheConfig
}

type CacheConfig struct {
	enabled     bool
	headerKeys  []string
	store       *redis.Client
	storeTTL    time.Duration
	cachePrefix string
}

func NewMiddleware(redis *redis.Client) Middleware {
	ttl, _ := time.ParseDuration(os.Getenv("REDIS_TTL"))
	// default ttl redis
	if ttl == 0 {
		ttl = 1 * time.Minute
	}

	useRedis, _ := strconv.ParseBool(os.Getenv("USE_REDIS"))

	return &appMiddleware{
		secret: helpers.GetJWTSecretKey(),
		cache: CacheConfig{
			enabled:     useRedis,
			store:       redis,
			cachePrefix: "gin:cache:",
			storeTTL:    ttl,
			headerKeys: []string{
				"User-Agent",
				"Accept",
				"Accept-Encoding",
				"Accept-Language",
				"Cookie",
				"Authorization",
			},
		},
	}
}

type Middleware interface {
	Auth() gin.HandlerFunc
	Cache(expiry ...time.Duration) gin.HandlerFunc
}
