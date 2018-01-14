package main

import (
	"time"
	"strconv"
	"fmt"
	"errors"

	"github.com/vmihailenco/msgpack"
	"github.com/go-redis/redis"
	"github.com/go-redis/cache"
)

var client = redis.NewClient(&redis.Options{
	Addr: "127.0.0.1:6379",
	Password: "",
	DB:  0,
})

var codec = cache.Codec{
	Redis: client,
	Marshal: func(v interface{}) ([]byte, error) {
		return msgpack.Marshal(v)
	},
	Unmarshal: func(b []byte, v interface{}) error {
		return msgpack.Unmarshal(b, v)
	},
}

// Updates a session from redis using the users
// Telegram ID.
func updateSession(user User) {
	userId := strconv.Itoa(user.Id)
	oneWeek := 25200 * time.Millisecond
	expireIn := time.Duration(oneWeek)
	codec.Set(&cache.Item{
		Key: userId,
		Object: user,
		Expiration: expireIn,
	})
}

func getSession(userId int) (User, error) {
	var user User
	err := codec.Get(strconv.Itoa(userId), &user)
	if err != nil {
		return user, errors.New("No session found")
	}
	return user, nil
}

// Updates an intent session based on a keyword
func updateIntentSession(userId int, intentSession IntentSession) {
	userIdStr := strconv.Itoa(userId)
	sessionKey := fmt.Sprintf("%s_intent", userIdStr)
	fiveMinutes := 300 * time.Millisecond
	expireIn := time.Duration(fiveMinutes)
	codec.Set(&cache.Item{
		Key: sessionKey,
		Object: intentSession,
		Expiration: expireIn,
	})
}

func getIntentSession(userId int) (IntentSession, error) {
	var intentSession IntentSession
	sessionKey := fmt.Sprintf("%s_intent", strconv.Itoa(userId))
	err := codec.Get(sessionKey, &intentSession)
	if err != nil {
		return intentSession, errors.New("No session found")
	}
	return intentSession, nil
}

func clearIntentSession(userId int) {
	sessionKey := fmt.Sprintf("%s_intent", strconv.Itoa(userId))
	codec.Delete(sessionKey)
}
