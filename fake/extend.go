package fake

func (f *Fake) Extend(key string) error {
	if f.FakeExtend != nil {
		return f.FakeExtend()
	}

	return nil
}
