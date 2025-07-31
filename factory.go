package locker

func Default() Interface {
	return New(Config{
		Poo: Pool(),
	})
}
