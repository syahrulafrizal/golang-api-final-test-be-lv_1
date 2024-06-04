package main

import (
	httpHandler "app/app/delivery/http"
	mongorepo "app/app/repository/mongo"
	"app/app/usecase"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/Yureka-Teknologi-Cipta/yureka/services/mongodb"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	mongo := mongodb.Connect(timeoutContext, os.Getenv("MONGO_URI"), "")

	// init repo
	mongorepo := mongorepo.NewMongodbRepo(mongo)

	// init usecase
	uc := usecase.NewAppUsecase(mongorepo, timeoutContext)

	// init gin
	ginEngine := gin.New()
	ginEngine.Use(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.Error("Panic Recover : ", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Something went wrong"))
			}
		}()
		c.Next()
	})

	// init route
	httpHandler.NewRouteHandler(ginEngine, uc)

	port := os.Getenv("PORT")

	logrus.Infof("Service running on port %s", port)
	ginEngine.Run(":" + port)
}
