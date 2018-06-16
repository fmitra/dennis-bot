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

type Session interface {
	Set(cacheKey string, v interface{}, timeInSeconds int)
	Get(cacheKey string, v interface{}) error
	Delete(cacheKey string) error
}

type Client struct {
	codec cache.Codec
}

func NewClient(config Config) *Client {
	address := fmt.Sprintf("%s:%s", config.Host, strconv.Itoa(int(config.Port)))
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: config.Password,
		DB:       config.Db,
	})

	codec := cache.Codec{
		Redis: redisClient,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}

	return &Client{codec}
}

func (c *Client) Delete(cacheKey string) error {
	err := c.codec.Delete(cacheKey)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Set(cacheKey string, v interface{}, timeInSeconds int) {
	// One hour default duration
	duration := time.Duration(3600) * time.Second
	if timeInSeconds != 0 {
		duration = time.Duration(timeInSeconds) * time.Second
	}
	expireIn := duration
	c.codec.Set(&cache.Item{
		Key:        cacheKey,
		Object:     v,
		Expiration: expireIn,
	})
}

func (c *Client) Get(cacheKey string, v interface{}) error {
	err := c.codec.Get(cacheKey, &v)
	if err != nil {
		return errors.New("No session found")
	}
	return nil
}
