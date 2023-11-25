package fake

func (f *Fake) Delete(key string) error {
	if f.FakeDelete != nil {
		return f.FakeDelete()
	}

	return nil
}
