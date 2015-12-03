/*
Package Redis
*/
package Redis

import (
	"github.com/Rabbit-Go/logger"
	"github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// redis connection pool. Usage 'Redis.Pool.Get()'.
// won't return any error.
var Pool *redis.Pool

var log = logger.Log.WithFields(logrus.Fields{
	"tag": "redis",
})

// Create Redis connection pool.
//     redis://localhost:6379
func NewPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(server)
			if err != nil {
				log.WithError(err).Error("Dial Redis Failed")
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					log.WithError(err).Error("Auth Redis Failed")
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				log.WithError(err).Error("PingPong Missing")
			}
			return err
		},
	}
}
