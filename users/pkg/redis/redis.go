package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	rd "github.com/go-redis/redis"
	"github.com/micro/micro/service/logger"

	"github.com/emadghaffari/kit-blog/users/config"
)

var (
	once sync.Once

	// DB variable
	DB rs = &redi{}
)

type rs interface {
	New()
	GetDB() *rd.Client
	Get(key string, dest interface{}) error
	Set(key string, value interface{}, duration time.Duration) error
	Del(key ...string) error
}

type redi struct {
	db *rd.Client
}

func (s *redi) New() {
	once.Do(func() {
		db, err := strconv.Atoi(config.Confs.Redis.DB)
		if err != nil {
			logger.Warn(err)
		}
		s.db = rd.NewClient(&rd.Options{
			Addr: config.Confs.Redis.Host,
			// Password: viper.GetString("redis.Password"), // no password set
			DB: db, // use default DB
		})

		if err := s.db.Ping(context.Background()).Err(); err != nil {
			logger.Warn(err)
		}

		fmt.Println("redis connected")
	})
}

func (s *redi) GetDB() *rd.Client {
	return s.db
}

// Set meth a new key,value
func (s *redi) Set(key string, value interface{}, duration time.Duration) error {
	p, err := json.Marshal(value)
	if err != nil {
		logger.Warn(err)
		return err
	}
	return s.db.Set(context.Background(), key, p, duration).Err()
}

// Get meth, get value with key
func (s *redi) Get(key string, dest interface{}) error {
	p, err := s.db.Get(context.Background(), key).Result()

	if p == "" {
		return fmt.Errorf("Value Not Found")
	}

	if err != nil {
		logger.Warn(err)
		return err
	}

	return json.Unmarshal([]byte(p), &dest)
}

func (s *redi) Del(key ...string) error {
	_, err := s.db.Del(context.Background(), key...).Result()
	if err != nil {
		return err
	}
	return nil
}
