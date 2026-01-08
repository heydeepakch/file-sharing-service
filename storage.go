package main

import (
	"time"
	"encoding/json"
	"os"
	"sync"
)

type FileMeta struct {
	ID         string
	Name       string
	Path       string
	Size       int64
	UploadedAt time.Time
}

var mu sync.Mutex

func saveDB() {
	mu.Lock()
	defer mu.Unlock()

	data, _ := json.MarshalIndent(fileDB, "", "  ")
	os.WriteFile("db.json", data, 0644)
}

func loadDB() {
	data, err := os.ReadFile("db.json")
	if err != nil {
		return
	}
	json.Unmarshal(data, &fileDB)
}
