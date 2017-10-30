package ratelimiter

import (
	"strconv"

	"github.com/go-redis/redis"
)

// implement Counter interface
var _ Counter = &RedisCounter{}

// RedisCounter limits calls to the api based on IP
type RedisCounter struct {
	c *redis.Client
}

// NewRedisCounter creates a new RedisCounter with default values
func NewRedisCounter(c *redis.Client) *RedisCounter {
	return &RedisCounter{
		c: c,
	}
}

// Increment increments the counter for the IP
func (r *RedisCounter) Increment(ip string) (count int, err error) {
	res := r.c.Incr("ratelimit:" + ip)
	if err := res.Err(); err != nil {
		return 0, err
	}
	return int(res.Val()), nil
}

// Decrement decrements the counter for the IP
func (r *RedisCounter) Decrement(ip string) (count int, err error) {
	res := r.c.Decr("ratelimit:" + ip)
	if err := res.Err(); err != nil {
		return 0, err
	}
	return int(res.Val()), nil
}

// Count returns the current count
func (r *RedisCounter) Count(ip string) (count int, err error) {
	res := r.c.Get("ratelimit:" + ip)
	i, err := strconv.Atoi(res.Val())
	if err != nil {
		return 0, err
	}
	return i, res.Err()
}

// Set sets the ip count
func (r *RedisCounter) Set(ip string, count int) error {
	res := r.c.Set("ratelimit:"+ip, count, 0)
	return res.Err()
}
