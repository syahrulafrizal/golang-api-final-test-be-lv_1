package main

import (
	http_admin "app/app/delivery/http/admin"
	"app/app/delivery/http/middleware"
	http_public "app/app/delivery/http/public"
	mongorepo "app/app/repository/mongo"
	usecase_admin "app/app/usecase/admin"
	usecase_public "app/app/usecase/public"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	yureka_mongodb "github.com/Yureka-Teknologi-Cipta/yureka/services/mongodb"
	yureka_redis "github.com/Yureka-Teknologi-Cipta/yureka/services/redis"
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
	ucAdmin := usecase_admin.NewAppUsecase(usecase_admin.RepoInjection{
		MongoDBRepo: mongorepo,
	}, timeoutContext)
	ucPublic := usecase_public.NewAppUsecase(usecase_public.RepoInjection{
		MongoDBRepo: mongorepo,
	}, timeoutContext)

	// init middleware
	mdl := middleware.NewMiddleware(redisClient)

	// gin mode realease when go env is production
	if os.Getenv("GO_ENV") == "production" || os.Getenv("GO_ENV") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// init gin
	ginEngine := gin.New()

	// add exception handler
	// ginEngine.Use(mdl.Recovery())

	// add logger
	ginEngine.Use(mdl.Logger(io.MultiWriter(writers...)))

	// cors
	ginEngine.Use(mdl.Cors())

	// default route
	ginEngine.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]any{
			"message": "It works",
		})
	})
	// Serve static files from the /media directory
	ginEngine.Static("/media", "./media")

	// init route
	http_admin.NewRouteHandler(ginEngine.Group("admin"), mdl, ucAdmin)
	http_public.NewRouteHandler(ginEngine.Group("public"), mdl, ucPublic)

	port := os.Getenv("PORT")

	logrus.Infof("Service running on port %s", port)
	ginEngine.Run(":" + port)
}
