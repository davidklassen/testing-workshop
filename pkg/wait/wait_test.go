package wait_test

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/davidklassen/testing-workshop/pkg/wait"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWait(t *testing.T) {
	t.Run("should wait for signal", func(t *testing.T) {
		var sig os.Signal
		var done = make(chan struct{})

		go func() {
			sig = wait.Wait(syscall.SIGHUP)
			done <- struct{}{}
		}()

		time.Sleep(time.Millisecond * 10)

		err := syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
		require.NoError(t, err)

		<-done
		assert.Equal(t, syscall.SIGHUP, sig)
	})

	t.Run("should handle multiple signals", func(t *testing.T) {
		var sig os.Signal
		var done = make(chan struct{})

		go func() {
			sig = wait.Wait(syscall.SIGHUP, syscall.SIGWINCH)
			done <- struct{}{}
		}()

		time.Sleep(time.Millisecond * 10)

		err := syscall.Kill(syscall.Getpid(), syscall.SIGWINCH)
		require.NoError(t, err)

		<-done
		assert.Equal(t, syscall.SIGWINCH, sig)
	})

	t.Run("should ignore unregistered signals", func(t *testing.T) {
		var sig os.Signal
		var done = make(chan struct{})

		go func() {
			sig = wait.Wait(syscall.SIGHUP)
			done <- struct{}{}
		}()

		time.Sleep(time.Millisecond * 10)

		err := syscall.Kill(syscall.Getpid(), syscall.SIGWINCH)
		require.NoError(t, err)

		select {
		case <-done:
			t.Fatalf("unexpected Wait success")
		case <-time.After(time.Millisecond * 10):
			assert.Equal(t, nil, sig)
		}
	})
}
