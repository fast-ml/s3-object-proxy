package permission

import (
	"errors"
	"log"

	"gopkg.in/redis.v4"
)

func NewRedisLookup(addr string) *RedisLookup {
	return &RedisLookup{
		Client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}

}

type RedisLookup struct {
	Client *redis.Client
}

func (l *RedisLookup) Name() string {
	return "Redis"
}

func (l *RedisLookup) CanView(key string) (*LookupResponse, error) {
	resp := l.Client.Get(key)
	log.Println("response received")
	if resp == nil {
		log.Println("response nil")
		return nil, errors.New("redis response was nil")
	}
	err := resp.Err()
	if err != nil {
		if err.Error() == "nil" {
			return &LookupResponse{Miss: true}, nil
		}
		// else {
		// 	l.Client.Set(key, PUBLIC, 0)
		// }
		return nil, err
	}
	log.Println("getting public/private status")
	val, err := resp.Int64()
	if err != nil {
		return nil, err
	}
	return &LookupResponse{Public: val == PUBLIC}, nil
}
