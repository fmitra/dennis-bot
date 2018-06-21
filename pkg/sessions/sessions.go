// Package sessions provides a client to interact with the cache layer.
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

// Config provides settings to connect to a Redis cache.
type Config struct {
	Host     string
	Port     int32
	Password string
	Db       int
}

// Session is an interface to interact with the cache layer.
type Session interface {
	Set(cacheKey string, v interface{}, timeInSeconds int)
	Get(cacheKey string, v interface{}) error
	Delete(cacheKey string) error
}

// Client provides methods to interact with the cache layer.
type Client struct {
	codec cache.Codec
}

// NewClient returns a Client connected to the cache layer.
func NewClient(c Config) *Client {
	address := fmt.Sprintf("%s:%s", c.Host, strconv.Itoa(int(c.Port)))
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: c.Password,
		DB:       c.Db,
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

// Delete removes an item from the cache.
func (c *Client) Delete(cacheKey string) error {
	err := c.codec.Delete(cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// Set adds an item to the cahce. Cache timeout is provided in seconds
// and defaults to one hour if a value of 0 is provided for the timeout.
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

// Get retrieves an item from the cache.
func (c *Client) Get(cacheKey string, v interface{}) error {
	err := c.codec.Get(cacheKey, &v)
	if err != nil {
		return errors.New("no session found")
	}
	return nil
}
