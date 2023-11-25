package redis

import (
	"github.com/xh3b4sd/tracer"
)

func (l *Redis) Exists(key string) (string, bool, error) {
	var err error

	if key == "" {
		return "", false, tracer.Mask(lockKeyEmptyError)
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

// TODO
func (l *Redis) exists(key string) (string, bool, error) {
	return "", false, nil
}
