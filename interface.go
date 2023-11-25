package locker

type Interface interface {
	// Create acquires a new distributed lock, if possible. If a new lock got
	// acquired, then the lock value will be returned. Note that creating a lock
	// for a key that got already acquired will result in an error.
	//
	//     @inp[0] the lock key to create
	//     @out[0] the lock value being created, if any
	//
	Create(string) (string, error)

	// Delete tries to release a distributed lock so that it can be acquired by
	// another process. If the lock could not be deleted, then it did not exist
	// and an error is returned.
	//
	//     @inp[0] the lock key to delete
	//
	Delete(string) error

	// Exists expresses whether a lock exists already for the given key. If the
	// lock does already exist, then true is returned, together with the lock
	// value.
	//
	//     @inp[0] the lock key to check
	//     @out[0] the lock value, if it exists
	//     @out[1] the bool expressing whether the lock exists
	//
	Exists(string) (string, bool, error)

	// Extend prevents a distributed lock from expiring.
	//
	//     @inp[0] the lock key to extend
	//
	Extend(string) error
}
