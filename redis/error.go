package redis

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var executionFailedError = &tracer.Error{
	Kind: "executionFailedError",
	Desc: "The action failed unexpectedly due to some internal error.",
}

var lockAlreadyExistsError = &tracer.Error{
	Kind: "lockAlreadyExistsError",
	Desc: "The action expects the lock not to be taken. The lock was found to be taken. Therefore the action failed.",
}

func IsLockAlreadyExists(err error) bool {
	return errors.Is(err, lockAlreadyExistsError)
}

var lockNotFoundError = &tracer.Error{
	Kind: "lockNotFoundError",
	Desc: "The action expects the lock to exist. The lock was not found to exist. Therefore the action failed.",
}

func IsLockNotFound(err error) bool {
	return errors.Is(err, lockNotFoundError)
}

var lockKeyEmptyError = &tracer.Error{
	Kind: "lockKeyEmptyError",
	Desc: "The action expects the provided key not to be empty. The provided key was found to be empty. Therefore the action failed.",
}
