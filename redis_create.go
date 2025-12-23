package locker

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"
)

func (r *Redis) Create(key string) (string, error) {
	var err error

	if key == "" {
		return "", tracer.Mask(lockKeyEmptyError,
			tracer.Context{Key: "method", Value: "Redis.Create"},
			tracer.Context{Key: "lock key", Value: key},
		)
	}

	var val string
	var cre bool

	fnc := func() error {
		val, cre, err = r.create(key)
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	{
		err = r.bac.Backoff(fnc)
		if err != nil {
			return "", tracer.Mask(err)
		}
	}

	{
		if !cre {
			return "", tracer.Mask(lockAlreadyExistsError,
				tracer.Context{Key: "method", Value: "Redis.Create"},
				tracer.Context{Key: "lock key", Value: key},
				tracer.Context{Key: "lock value", Value: val},
			)
		}
	}

	return val, nil
}

func (r *Redis) create(key string) (string, bool, error) {
	var err error

	var con redis.Conn
	{
		con = r.poo.Get()
	}

	{
		defer con.Close() // nolint:errcheck
	}

	var val string
	{
		val = strconv.FormatInt(time.Now().UTC().Unix(), 10)
	}

	var arg []any
	{
		arg = append(arg,
			strings.Join([]string{r.pre, key}, r.del),
			val,
			"NX", // only set the key if it does not already exist
			"EX",
			r.exp.Seconds(),
		)
	}

	var res string
	{
		res, err = redis.String(con.Do("SET", arg...))
		if errors.Is(err, redis.ErrNil) {
			return "", false, nil
		} else if err != nil {
			return "", false, tracer.Mask(err)
		}
	}

	{
		if res != "OK" {
			return "", false, tracer.Mask(lockAttemptFailedError,
				tracer.Context{Key: "method", Value: "Redis.Create"},
				tracer.Context{Key: "lock key", Value: key},
				tracer.Context{Key: "lock value", Value: val},
			)
		}
	}

	return val, true, nil
}
