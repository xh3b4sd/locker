package locker

func (f *Fake) Exists(key string) (string, bool, error) {
	if f.FakeExists != nil {
		return f.FakeExists()
	}

	return "", false, nil
}
