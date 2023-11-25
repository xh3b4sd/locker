package redis

import (
	"errors"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"
)

func (l *Redis) Extend(key string, val string) error {
	var err error

	if key == "" {
		return tracer.Maskf(lockKeyEmptyError, "Locker.Extend")
	}

	var ext bool

	act := func() error {
		ext, err = l.extend(key, val)
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	{
		err = l.brk.Execute(act)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	{
		if !ext {
			return tracer.Mask(lockNotFoundError)
		}
	}

	return nil
}

func (l *Redis) extend(key string, val string) (bool, error) {
	var err error

	var con redis.Conn
	{
		con = l.poo.Get()
		defer con.Close()
	}

	var arg []interface{}
	{
		arg = append(arg,
			strings.Join([]string{l.pre, key}, l.del),
			val,
			"XX", // only set the key if it already exists
			"EX",
			l.exp.Seconds(),
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
			return false, tracer.Mask(executionFailedError)
		}
	}

	return true, nil
}
