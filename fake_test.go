package locker

import (
	"testing"
)

func Test_Locker_Interface_Fake(t *testing.T) {
	var _ Interface = &Fake{}
}
