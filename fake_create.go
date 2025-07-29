package locker

func (f *Fake) Create(key string) (string, error) {
	if f.FakeCreate != nil {
		return f.FakeCreate()
	}

	return "", nil
}
