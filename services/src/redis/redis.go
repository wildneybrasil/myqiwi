// random
package redis

import (
	"fmt"
	"time"

	"gopkg.in/redis.v3"
)

var (
	ADDR = "redis:6379"
)

func Del(key string) error {
	client := redis.NewClient(&redis.Options{
		Addr:     ADDR,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	fmt.Println("KEY:" + key)

	client.Del(key).Result()

	client.Close()

	return nil
}
func Get(key string) (*string, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     ADDR,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	fmt.Println("KEY:" + key)

	val, err := client.Get(key).Result()
	if err != nil {
		return nil, err
	}
	client.Close()

	return &val, nil
}
func Set(key string, value string, t time.Duration) error {
	client := redis.NewClient(&redis.Options{
		Addr:     ADDR,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	fmt.Println("KEY:" + key)
	fmt.Println("VALUE:" + value)

	err := client.Set(key, value, t).Err()
	if err != nil {
		panic(err)
	}
	client.Close()
	return err
}
