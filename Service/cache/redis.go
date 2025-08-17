package cache

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type Redis struct {
	Redis redis.Client
}

func NewRedis(addr string) Redis {
	return Redis{Redis: *redis.NewClient(&redis.Options{
		Addr: addr,
	})}
}
func (r Redis) CacheTerms(terms []string) {
	r.Redis.Set("terms", terms, 10*time.Minute)

}
func (r Redis) GetCacheValue(key string) (string, error) {
	val, err := r.Redis.Get(key).Result()
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return "", err
	}
	fmt.Print(val)
	return val, nil
}
