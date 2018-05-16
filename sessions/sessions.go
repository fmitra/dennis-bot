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

var codec cache.Codec

type Config struct {
	Host     string
	Port     int32
	Password string
	Db       int
}

func Init(config Config) {
	address := fmt.Sprintf("%s:%s", config.Host, strconv.Itoa(int(config.Port)))
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: config.Password,
		DB:       config.Db,
	})

	codec = cache.Codec{
		Redis: client,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
}

func Set(cacheKey string, v interface{}) {
	oneWeek := 25200 * time.Millisecond
	expireIn := time.Duration(oneWeek)
	codec.Set(&cache.Item{
		Key:        cacheKey,
		Object:     v,
		Expiration: expireIn,
	})
}

func Get(cacheKey string) (interface{}, error) {
	var v interface{}
	err := codec.Get(cacheKey, &v)
	if err != nil {
		return v, errors.New("No session found")
	}
	return v, nil
}
