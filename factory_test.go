package locker

import (
	"testing"

	"github.com/xh3b4sd/locker/fake"
	"github.com/xh3b4sd/locker/lock"
)

func Test_Factory_Interface_Default(t *testing.T) {
	var _ Interface = &lock.Lock{}
}

func Test_Factory_Interface_Fake(t *testing.T) {
	var _ Interface = &fake.Fake{}
}
