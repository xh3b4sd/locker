package redis

import (
	"github.com/xh3b4sd/tracer"
)

func (l *Redis) Create(key string) (string, error) {
	var err error

	if key == "" {
		return "", tracer.Mask(lockKeyEmptyError)
	}

	var val string

	act := func() error {
		val, err = l.create(key)
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

	return val, nil
}

// TODO
func (l *Redis) create(key string) (string, error) {
	return "", nil
}
