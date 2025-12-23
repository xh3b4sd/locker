package locker

import (
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"
)

func (r *Redis) Delete(key string) error {
	var err error

	if key == "" {
		return tracer.Mask(lockKeyEmptyError,
			tracer.Context{Key: "method", Value: "Redis.Delete"},
			tracer.Context{Key: "lock key", Value: key},
		)
	}

	fnc := func() error {
		err := r.delete(key)
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

	return nil
}

func (r *Redis) delete(key string) error {
	var err error

	var con redis.Conn
	{
		con = r.poo.Get()
	}

	{
		defer con.Close() // nolint:errcheck
	}

	var arg []any
	{
		arg = append(arg,
			strings.Join([]string{r.pre, key}, r.del),
		)
	}

	{
		_, err = redis.Int64(con.Do("DEL", arg...))
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
