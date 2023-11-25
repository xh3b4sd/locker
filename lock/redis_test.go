//go:build redis

package lock

import (
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/xh3b4sd/locker"
	"github.com/xh3b4sd/locker/pool"
	"github.com/xh3b4sd/tracer"
)

func Test_Redis_Delete(t *testing.T) {
	var err error

	var loc locker.Interface
	{
		loc = New(Config{
			Exp: 1 * time.Second,
			Poo: prgAll(pool.New()),
		})
	}

	var key string
	{
		key = "key"
	}

	// Deleting should not fail if there is no lock.
	{
		err = loc.Delete(key)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create the lock the first time.
	{
		_, err = loc.Create(key)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Delete the created lock from above.
	{
		err = loc.Delete(key)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create the lock the second time. This should not fail since the lock got
	// deleted above.
	{
		_, err = loc.Create(key)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func Test_Redis_Exists(t *testing.T) {
	var err error

	var loc locker.Interface
	{
		loc = New(Config{
			Exp: 1 * time.Second,
			Poo: prgAll(pool.New()),
		})
	}

	var key string
	{
		key = "key"
	}

	var one string
	var exi bool
	{
		one, exi, err = loc.Exists(key)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if one != "" {
			t.Fatal("expected", "empty string", "got", one)
		}
		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	var two string
	{
		two, err = loc.Create(key)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		one, exi, err = loc.Exists(key)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if one != two {
			t.Fatal("expected", two, "got", one)
		}
		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}
}

func Test_Redis_Expiry(t *testing.T) {
	var err error

	var loc locker.Interface
	{
		loc = New(Config{
			Exp: 1 * time.Second,
			Poo: prgAll(pool.New()),
		})
	}

	// Extending a lock that does not exist should fail.
	{
		err = loc.Extend("foo", "bar")
		if !IsLockNotFound(err) {
			t.Fatal("expected", "lockNotFoundError", "got", err)
		}
	}

	var key string
	{
		key = "key"
	}

	{
		_, err = loc.Create(key)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		time.Sleep(2 * time.Second)
	}

	// The first Create call should not hold the lock anymore due to expiry.
	var val string
	{
		val, err = loc.Create(key)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		time.Sleep(500 * time.Millisecond)
	}

	// We are repeating the same steps to verify whether the lock can be created
	// again after extending it.
	{
		err = loc.Extend(key, val)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		time.Sleep(500 * time.Millisecond)
	}

	// We are repeating the same steps to verify whether the lock can be created
	// again after extending it.
	{
		err = loc.Extend(key, val)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		time.Sleep(500 * time.Millisecond)
	}

	// We are repeating the same steps to verify whether the lock can be created
	// again after extending it.
	{
		err = loc.Extend(key, val)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create the lock the third time. This should fail since the lock got
	// extended above.
	{
		_, err = loc.Create(key)
		if !IsLockAlreadyExists(err) {
			t.Fatal("expected", "lockAlreadyExistsError", "got", err)
		}
	}
}

func Test_Redis_Locked(t *testing.T) {
	var err error

	var loc locker.Interface
	{
		loc = New(Config{
			Exp: 1 * time.Second,
			Poo: prgAll(pool.New()),
		})
	}

	var key string
	{
		key = "key"
	}

	// Create the lock the first time.
	{
		_, err = loc.Create(key)
		if err != nil {
			t.Fatal(err)
		}
	}

	var sta time.Time
	{
		sta = time.Now().UTC()
	}

	// Create the lock the second time. This should fail since the lock got
	// already created above.
	{
		_, err = loc.Create(key)
		if !IsLockAlreadyExists(err) {
			t.Fatal("expected", "lockAlreadyExistsError", "got", err)
		}
	}

	// Our locker is configured with the default breaker budget, which should be
	// 30 * 1 second. The breaker should retry on error, but if a lock is already
	// taken, then this error should cancel the breaker and overwrite the
	// remaining budget. That means, the test here should not take much longer to
	// execute than the lock expiry itself, given the breaker respects the
	// "already locked" error.
	{
		if time.Since(sta) > 5*time.Second {
			t.Fatal("test must pass within 5 seconds")
		}
	}
}

// prgAll is a convenience function for calling FLUSHALL. The provided pool is
// returned as is.
func prgAll(poo *redis.Pool) *redis.Pool {
	var con redis.Conn
	{
		con = poo.Get()
		defer con.Close()
	}

	{
		_, err := con.Do("FLUSHALL")
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}

	return poo
}
