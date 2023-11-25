package lock

import (
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"
)

func (l *Lock) Delete(key string) error {
	var err error

	if key == "" {
		return tracer.Maskf(lockKeyEmptyError, "Locker.Delete")
	}

	act := func() error {
		err := l.delete(key)
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

	return nil
}

func (l *Lock) delete(key string) error {
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

	{
		_, err = redis.Int64(con.Do("DEL", arg...))
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
