package main

import (
	"time"
)

type FileMeta struct {
	ID         string
	Name       string
	Path       string
	Size       int64
	UploadedAt time.Time
}
