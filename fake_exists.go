package locker

func (f *Fake) Exists(key string) (string, bool, error) {
	if f.FakeExtend != nil {
		return f.FakeExists()
	}

	return "", false, nil
}
