package redis

import "github.com/xh3b4sd/tracer"

var executionFailedError = &tracer.Error{
	Kind: "executionFailedError",
	Desc: "The action failed unexpectedly due to some internal error.",
}

var lockKeyEmptyError = &tracer.Error{
	Kind: "lockKeyEmptyError",
	Desc: "The action expects the provided key not to be empty. The provided key was found to be empty. Therefore the action failed.",
}
