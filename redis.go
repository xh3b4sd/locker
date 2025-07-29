package locker

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/choreo/backoff"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	// Bac is the backoff implementation used to retry redis connections on
	// failure.
	Bac backoff.Interface

	// Exp is the lock's expiry, so that locks can expire after a certain amount
	// of time of inactivity. Defaults to 30 seconds. Disabled with -1.
	Exp time.Duration

	// Poo is the redis connection pool to select client connections from.
	Poo *redis.Pool

	// Pre is the prefix of the underlying redis key used to coordinate the
	// distributed lock.
	Pre string
}

type Redis struct {
	bac backoff.Interface
	exp time.Duration
	del string
	poo *redis.Pool
	pre string
}

func New(c Config) *Redis {
	if c.Bac == nil {
		c.Bac = backoff.New(backoff.Config{
			Bac: []time.Duration{
				1 * time.Second,
				1 * time.Second,
				1 * time.Second,
			},
		})
	}
	if c.Exp == 0 {
		c.Exp = 30 * time.Second
	}
	if c.Poo == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Poo must not be empty", c)))
	}
	if c.Pre == "" {
		c.Pre = "loc"
	}

	return &Redis{
		bac: c.Bac,
		exp: c.Exp,
		del: ":",
		poo: c.Poo,
		pre: c.Pre,
	}
}
