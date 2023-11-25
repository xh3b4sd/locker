package pool

import (
	"net"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

func New() *redis.Pool {
	var add string
	{
		add = envAdd()
	}

	var opt []redis.DialOption
	{
		opt = []redis.DialOption{
			redis.DialConnectTimeout(time.Second),
			redis.DialReadTimeout(time.Second),
			redis.DialWriteTimeout(time.Second),
		}
	}

	var p *redis.Pool
	{
		p = &redis.Pool{
			MaxIdle:     100,
			MaxActive:   100,
			IdleTimeout: 5 * time.Minute,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", add, opt...)
				if err != nil {
					return nil, err
				}

				return c, nil
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				// For 1 minute, connections are neither considered stale nor broken.
				// During that time connections may break and the pool will not discard
				// the broken connection until 1 minute has passed since the current
				// connection was found to be working. If a connection breaks, for
				// instance because redis restarted, and a minute passed, the pool will
				// establish a new connection and the client will work just fine again.
				// If the time configured here is too long, it should be considered to
				// lower that threshold.
				if time.Since(t) < time.Minute {
					return nil
				}

				_, err := c.Do("PING")
				if err != nil {
					// By returning an error the pool realizes the connection is broken
					// and will then establish another one to work with.
					return err
				}

				return nil
			},
		}
	}

	return p
}

func envAdd() string {
	var hos string
	{
		hos = os.Getenv("REDIS_HOST")
		if hos == "" {
			hos = "127.0.0.1"
		}
	}

	var por string
	{
		por = os.Getenv("REDIS_PORT")
		if por == "" {
			por = "6379"
		}
	}

	return net.JoinHostPort(hos, por)
}
