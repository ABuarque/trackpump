package service

import "io"

// Storage defines how storage services should work
type Storage interface {
	Put(fileName string, data io.Reader) (string, error)
}
