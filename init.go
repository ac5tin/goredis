package redis

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

// Client contains redis pool
type Client struct {
	pool *redis.Pool
}

// NewRedisClient returns a new RedisClient
func NewRedisClient(addr string, db int, passwd string) *Client {
	pool := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr, redis.DialPassword(passwd), redis.DialDatabase(db))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	log.Printf("new redis pool at %s", addr)
	client := &Client{
		pool: pool,
	}
	return client
}
