package locker

import (
	"errors"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"
)

func (r *Redis) Exists(key string) (string, bool, error) {
	var err error

	if key == "" {
		return "", false, tracer.Mask(lockKeyEmptyError,
			tracer.Context{Key: "method", Value: "Redis.Exists"},
			tracer.Context{Key: "lock key", Value: key},
		)
	}

	var val string
	var exi bool

	fnc := func() error {
		val, exi, err = r.exists(key)
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	{
		err = r.bac.Backoff(fnc)
		if err != nil {
			return "", false, tracer.Mask(err)
		}
	}

	return val, exi, nil
}

func (r *Redis) exists(key string) (string, bool, error) {
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
		)
	}

	var res string
	{
		res, err = redis.String(con.Do("GET", arg...))
		if errors.Is(err, redis.ErrNil) {
			return "", false, nil
		} else if err != nil {
			return "", false, tracer.Mask(err)
		}
	}

	return res, res != "", nil
}
