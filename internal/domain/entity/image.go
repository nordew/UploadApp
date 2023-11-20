package entity

import (
	"io"
)

type Image struct {
	ID     string
	Name   string
	Size   int64
	Reader io.Reader
}
