package locker

import (
	"github.com/xh3b4sd/locker/fake"
)

func Fake() Interface {
	return &fake.Fake{}
}
