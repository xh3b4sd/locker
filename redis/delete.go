package redis

import (
	"github.com/xh3b4sd/tracer"
)

func (l *Redis) Delete(key string) error {
	var err error

	if key == "" {
		return tracer.Mask(lockKeyEmptyError)
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

// TODO
func (l *Redis) delete(key string) error {
	return nil
}
