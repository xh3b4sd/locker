# locker

Distributed redis lock.

```go
import (
	"github.com/xh3b4sd/locker/lock"
	"github.com/xh3b4sd/locker/pool"
)

func main() {
	var err error

	// Initialize a new locker instance that can be reused for different locks.
	// Every locker requires a redis connection pool.
	var loc locker.Interface
	{
		loc = lock.New(lock.Config{
			Poo: pool.New(),
		})
	}

	// Keys define a lock scope. They are used to distinguish different locks.
	var key string
	{
		key = "key"
	}

	// Creating a lock returns a unique lock value that can be used to extend an
	// existing lock as well as to verify whether a specific lock exists.
	var val string
	{
		val, err = loc.Create(key)
		if err != nil {
			panic(err)
		}
	}

	// Extending the lock can be done using its key and value.
	{
		err = loc.Extend(key, val)
		if err != nil {
			panic(err)
		}
	}

	// Locks can be deleted in order to free the lock key for further locking.
	{
		err = loc.Delete(key)
		if err != nil {
			panic(err)
		}
	}
}
```

### Conformance Tests

```
docker run --rm --name redis-stack -p 6379:6379 -p 8001:8001 redis/redis-stack:latest
```

```
go test ./... -race -tags integration
```

### Redis Port

```
export REDIS_PORT=6381
```

```
docker run --rm --name redis-stack-rescue -p 6381:6379 -p 8003:8001 redis/redis-stack:latest
```
