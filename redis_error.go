package locker

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var lockAlreadyExistsError = &tracer.Error{
	Description: "The locker expects the lock not to be taken. The lock was found to be taken already. Therefore the locking attempt failed.",
}

func IsLockAlreadyExists(err error) bool {
	return errors.Is(err, lockAlreadyExistsError)
}

var lockAttemptFailedError = &tracer.Error{
	Description: "The locker tried to acquire a lock. The lock could not be acquired. Therefore the locking attempt failed.",
}

func IsLockAttemptFailed(err error) bool {
	return errors.Is(err, lockAttemptFailedError)
}

var lockNotFoundError = &tracer.Error{
	Description: "The locker expects the lock to exist. The lock was not found to exist. Therefore the locking attempt failed.",
}

func IsLockNotFound(err error) bool {
	return errors.Is(err, lockNotFoundError)
}

var lockKeyEmptyError = &tracer.Error{
	Description: "The locker expects the provided key not to be empty. The provided key was found to be empty. Therefore the locking attempt failed.",
}
