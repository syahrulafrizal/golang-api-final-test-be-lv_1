package main

import (
	httpHandler "app/app/delivery/http"
	"app/app/delivery/http/middleware"
	mongorepo "app/app/repository/mongo"
	"app/app/usecase"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	yureka_mongodb "github.com/Yureka-Teknologi-Cipta/yureka/services/mongodb"
	yureka_redis "github.com/Yureka-Teknologi-Cipta/yureka/services/redis"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	timeoutStr := os.Getenv("TIMEOUT")
	if timeoutStr == "" {
		timeoutStr = "5"
	}
	timeout, _ := strconv.Atoi(timeoutStr)
	timeoutContext := time.Duration(timeout) * time.Second

	// logger
	writers := make([]io.Writer, 0)
	if logSTDOUT, _ := strconv.ParseBool(os.Getenv("LOG_TO_STDOUT")); logSTDOUT {
		writers = append(writers, os.Stdout)
	}

	if logFILE, _ := strconv.ParseBool(os.Getenv("LOG_TO_FILE")); logFILE {
		logMaxSize, _ := strconv.Atoi(os.Getenv("LOG_MAX_SIZE"))
		if logMaxSize == 0 {
			logMaxSize = 50 //default 50 megabytes
		}

		logFilename := os.Getenv("LOG_FILENAME")
		if logFilename == "" {
			logFilename = "server.log"
		}

		lg := &lumberjack.Logger{
			Filename:   logFilename,
			MaxSize:    logMaxSize,
			MaxBackups: 1,
			LocalTime:  true,
		}

		writers = append(writers, lg)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(io.MultiWriter(writers...))

	// set gin writer to logrus
	gin.DefaultWriter = logrus.StandardLogger().Writer()

	// init mongo database
	mongo := yureka_mongodb.Connect(timeoutContext, os.Getenv("MONGO_URL"), "")

	// init redis database
	var redisClient *redis.Client
	if useRedis, err := strconv.ParseBool(os.Getenv("USE_REDIS")); err == nil && useRedis {
		redisClient = yureka_redis.Connect(timeoutContext, os.Getenv("REDIS_URL"))
	}

	// init repo
	mongorepo := mongorepo.NewMongodbRepo(mongo)

	// init usecase
	uc := usecase.NewAppUsecase(mongorepo, timeoutContext)

	// init middleware
	mdl := middleware.NewMiddleware(redisClient)

	// init gin
	ginEngine := gin.New()

	// add exception handler
	ginEngine.Use(mdl.Recovery())

	// add logger
	ginEngine.Use(mdl.Logger(io.MultiWriter(writers...)))

	// cors
	ginEngine.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           12 * time.Hour,
	}))

	// default route
	ginEngine.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]any{
			"message": "It works",
		})
	})

	// init route
	httpHandler.NewRouteHandler(ginEngine.Group(""), mdl, uc)

	port := os.Getenv("PORT")

	logrus.Infof("Service running on port %s", port)
	ginEngine.Run(":" + port)
}
