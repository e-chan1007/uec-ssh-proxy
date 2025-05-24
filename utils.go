package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

func ShuffleSlice[T any](slice []T) []T {
	shuffledSlice := make([]T, len(slice))
	copy(shuffledSlice, slice)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	r.Shuffle(len(shuffledSlice), func(i, j int) {
		shuffledSlice[i], shuffledSlice[j] = shuffledSlice[j], shuffledSlice[i]
	})

	return shuffledSlice
}

func CheckPortAvailability(addr string, interval, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, interval)
		if err == nil {
			conn.Close()
			return nil
		}

		if opErr, ok := err.(*net.OpError); ok && opErr.Op == "dial" {
			time.Sleep(interval)
			continue
		}
		return fmt.Errorf("unexpected error while polling port: %w", err)
	}
	return fmt.Errorf("port %s did not become available within %s", addr, timeout)
}
