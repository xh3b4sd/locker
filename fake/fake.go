package fake

type Fake struct {
	FakeCreate func() (string, error)
	FakeDelete func() error
	FakeExists func() (string, bool, error)
	FakeExtend func() error
}
