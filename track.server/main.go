package main

import (
	"log"
	"track/config"
	cacheinfra "track/infrastructure/cache"
	"track/internal/user"
	"track/logger"
	"track/router"
)

var _config config.Config

func init() {
	var err error
	_config, err = config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	logger.InitLogger(_config)
	logger.Logger.Info("Logger initialized")
	config.ConnectDB()
	logger.Logger.Info("DB connection established")
	config.SyncDB()
	logger.Logger.Info("DB sync completed")
	defer logger.Logger.Sync()
}

func main() {
	_ = cacheinfra.NewInMemoryCache(50)
	_ = cacheinfra.NewInMemoryCache(50)

	userSvc := user.NewService()

	userHandler := user.NewHandler(userSvc)

	router.InitRouter(
		userHandler,
	)
	log.Fatal(router.Start("0.0.0.0:" + _config.PORT))
}
