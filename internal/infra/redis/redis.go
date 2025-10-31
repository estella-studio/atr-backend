package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/estella-studio/atr-backend/internal/infra/env"
	"github.com/redis/go-redis/v9"
)

type RedisItf interface {
	Set(key string, value string)
}

type Redis struct {
	Client     *redis.Client
	address    string
	port       uint
	username   string
	password   string
	db         int
	expiration int
}

func NewRedis(env *env.Env) (*redis.Client, *Redis) {
	Redis := Redis{
		address:    env.RedisAddress,
		port:       env.RedisPort,
		username:   env.RedisUsername,
		password:   env.RedisPassword,
		db:         env.RedisDatabase,
		expiration: env.RedisExpiration,
	}

	redis := New(&Redis)

	Redis.Client = redis

	Test(&Redis)

	return redis,
		&Redis
}

func New(r *Redis) *redis.Client {
	url := fmt.Sprintf(
		"redis://%s:%s@%s:%d/%d",
		r.username,
		r.password,
		r.address,
		r.port,
		r.db,
	)

	opts, err := redis.ParseURL(url)
	if err != nil {
		log.Panic(err)
	}

	return redis.NewClient(opts)
}

func Test(r *Redis) {
	ctx := context.Background()
	key := "test"
	value := "test"

	log.Println("testing redis connection")

	r.Client.Set(ctx, key, value, time.Duration(r.expiration))

	result, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		log.Panic(err)
	}

	if result != value {
		log.Panic("redis testing failed")
	}

	r.Client.Del(ctx, key)

	log.Println("redis testing success")
}

func (r *Redis) Set(key string, value string) {
	ctx := context.Background()

	r.Client.Set(ctx, key, value, time.Duration(r.expiration))
}
