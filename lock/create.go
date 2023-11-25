package lock

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/tracer"
)

func (l *Lock) Create(key string) (string, error) {
	var err error

	if key == "" {
		return "", tracer.Maskf(lockKeyEmptyError, "Locker.Create")
	}

	var val string
	var cre bool

	act := func() error {
		val, cre, err = l.create(key)
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	{
		err = l.brk.Execute(act)
		if err != nil {
			return "", tracer.Mask(err)
		}
	}

	{
		if !cre {
			return "", tracer.Mask(lockAlreadyExistsError)
		}
	}

	return val, nil
}

func (l *Lock) create(key string) (string, bool, error) {
	var err error

	var con redis.Conn
	{
		con = l.poo.Get()
		defer con.Close()
	}

	var val string
	{
		val = strconv.FormatInt(time.Now().UTC().Unix(), 10)
	}

	var arg []interface{}
	{
		arg = append(arg,
			strings.Join([]string{l.pre, key}, l.del),
			val,
			"NX", // only set the key if it does not already exist
			"EX",
			l.exp.Seconds(),
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
			return "", false, tracer.Mask(executionFailedError)
		}
	}

	return val, true, nil
}
