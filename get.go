package redis

import (
	"github.com/gomodule/redigo/redis"
)

// Get returns a redis connection from the pool
func (r *Client) Get() redis.Conn {
	return r.pool.Get()
}
