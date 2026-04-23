package main

import (
	"glaze/config"
	"glaze/global"
	cacheinfra "glaze/infrastructure/cache"
	"glaze/internal/user"
	"glaze/internal/workspace"
	"glaze/logger"
	"glaze/models"
	"glaze/router"
	"log"
)

func init() {
	var err error
	global.GlobalConf, err = config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	models.SetEncryptionKey(global.GlobalConf.EncryptKey)
	logger.InitLogger(global.GlobalConf)
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
	workspaceSvc := workspace.NewService()

	userHandler := user.NewHandler(userSvc)
	workspacehandler := workspace.NewHandler(workspaceSvc)

	router.InitRouter(
		userHandler,
		workspacehandler,
	)
	log.Fatal(router.Start("0.0.0.0:" + global.GlobalConf.PORT))
}
