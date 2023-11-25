package lock

import (
	"errors"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"
)

func (l *Lock) Exists(key string) (string, bool, error) {
	var err error

	if key == "" {
		return "", false, tracer.Maskf(lockKeyEmptyError, "Locker.Exists")
	}

	var val string
	var exi bool

	act := func() error {
		val, exi, err = l.exists(key)
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	{
		err = l.brk.Execute(act)
		if err != nil {
			return "", false, tracer.Mask(err)
		}
	}

	return val, exi, nil
}

func (l *Lock) exists(key string) (string, bool, error) {
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
