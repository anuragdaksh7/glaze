package config

import "github.com/hibiken/asynq"

var RedisClient *asynq.Client

func InitRedisClient() {
	config, err := LoadConfig(".")
	if err != nil {
		panic(err)
	}

	RedisClient = asynq.NewClient(asynq.RedisClientOpt{Addr: config.RedisURL})
}
