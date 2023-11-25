package redis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/breakr"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	// Brk is the budget implementation used to retry redis connections on
	// failure.
	Brk breakr.Interface
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
	brk breakr.Interface
	exp time.Duration
	poo *redis.Pool
	pre string
}

func New(c Config) *Redis {
	if c.Brk == nil {
		c.Brk = breakr.Default()
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

	var l *Redis
	{
		l = &Redis{
			brk: c.Brk,
			exp: c.Exp,
			poo: c.Poo,
			pre: c.Pre,
		}
	}

	return l
}
