package locker

import "testing"

func Test_Locker_Interface_Redis(t *testing.T) {
	var _ Interface = &Redis{}
}
