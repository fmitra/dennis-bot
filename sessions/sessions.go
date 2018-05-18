package sessions

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
)

type Config struct {
	Host     string
	Port     int32
	Password string
	Db       int
}

type Session struct {
	codec cache.Codec
}

func New(config Config) *Session {
	address := fmt.Sprintf("%s:%s", config.Host, strconv.Itoa(int(config.Port)))
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: config.Password,
		DB:       config.Db,
	})

	codec := cache.Codec{
		Redis: client,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}

	return &Session{codec}
}

func (s *Session) Set(cacheKey string, v interface{}) {
	oneWeek := 25200 * time.Millisecond
	expireIn := time.Duration(oneWeek)
	fmt.Printf("setting %v", v)
	s.codec.Set(&cache.Item{
		Key:        cacheKey,
		Object:     v,
		Expiration: expireIn,
	})
}

func (s *Session) Get(cacheKey string, v interface{}) error {
	err := s.codec.Get(cacheKey, &v)
	if err != nil {
		return errors.New("No session found")
	}
	return nil
}
