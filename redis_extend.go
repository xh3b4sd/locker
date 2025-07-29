package locker

import (
	"errors"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"
)

func (r *Redis) Extend(key string, val string) error {
	var err error

	if key == "" {
		return tracer.Mask(lockKeyEmptyError,
			tracer.Context{Key: "method", Value: "Redis.Extend"},
			tracer.Context{Key: "lock key", Value: key},
			tracer.Context{Key: "lock value", Value: val},
		)
	}

	var ext bool

	fnc := func() error {
		ext, err = r.extend(key, val)
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	{
		err = r.bac.Backoff(fnc)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	{
		if !ext {
			return tracer.Mask(lockNotFoundError,
				tracer.Context{Key: "method", Value: "Redis.Extend"},
				tracer.Context{Key: "lock key", Value: key},
				tracer.Context{Key: "lock value", Value: val},
			)
		}
	}

	return nil
}

func (r *Redis) extend(key string, val string) (bool, error) {
	var err error

	var con redis.Conn
	{
		con = r.poo.Get()
		defer con.Close()
	}

	var arg []interface{}
	{
		arg = append(arg,
			strings.Join([]string{r.pre, key}, r.del),
			val,
			"XX", // only set the key if it already exists
			"EX",
			r.exp.Seconds(),
		)
	}

	var res string
	{
		res, err = redis.String(con.Do("SET", arg...))
		if errors.Is(err, redis.ErrNil) {
			return false, nil
		} else if err != nil {
			return false, tracer.Mask(err)
		}
	}

	{
		if res != "OK" {
			return false, tracer.Mask(lockAttemptFailedError,
				tracer.Context{Key: "method", Value: "Redis.Extend"},
				tracer.Context{Key: "lock key", Value: key},
				tracer.Context{Key: "lock value", Value: val},
			)
		}
	}

	return true, nil
}
