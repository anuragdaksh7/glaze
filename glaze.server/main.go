package main

import (
	"encoding/base64"
	"glaze/config"
	"glaze/global"
	cacheinfra "glaze/infrastructure/cache"
	"glaze/internal/tasks"
	"glaze/internal/user"
	"glaze/internal/webhooks"
	"glaze/internal/worker"
	"glaze/internal/workspace"
	"glaze/logger"
	"glaze/models"
	"glaze/router"
	"log"

	"github.com/hibiken/asynq"
)

func init() {
	var err error
	global.GlobalConf, err = config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	key, err := base64.StdEncoding.DecodeString(global.GlobalConf.EncryptKey)
	if err != nil {
		panic(err)
	}

	models.SetEncryptionKey(key)
	logger.InitLogger(global.GlobalConf)
	logger.Logger.Info("Logger initialized")
	config.ConnectDB()
	logger.Logger.Info("DB connection established")
	config.SyncDB()
	logger.Logger.Info("DB sync completed")
	config.InitRedisClient()
	logger.Logger.Info("Redis initialized")
	defer logger.Logger.Sync()
}

func main() {
	_ = cacheinfra.NewInMemoryCache(50)
	_ = cacheinfra.NewInMemoryCache(50)

	buildWorker, err := worker.NewBuildWorker(config.DB)
	if err != nil {
		log.Fatalf("Failed to initialize worker: %v", err)
	}

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "localhost:6379"},
		asynq.Config{
			Concurrency: 1, // Stick to 1 for your laptop to avoid CPU melt
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeBuildDeployment, buildWorker.ProcessBuildTask)

	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatalf("Asynq server error: %v", err)
		}
	}()
	
	userSvc := user.NewService()
	workspaceSvc := workspace.NewService()
	webhookSvc := webhooks.NewService(config.RedisClient)

	userHandler := user.NewHandler(userSvc)
	workspacehandler := workspace.NewHandler(workspaceSvc)
	webhookHandler := webhooks.NewHandler(webhookSvc)

	router.InitRouter(
		userHandler,
		workspacehandler,
		webhookHandler,
	)
	log.Fatal(router.Start("0.0.0.0:" + global.GlobalConf.PORT))
}
