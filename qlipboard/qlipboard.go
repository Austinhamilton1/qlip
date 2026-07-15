package qlipboard

import (
	"crypto/sha256"
	"time"

	"golang.design/x/clipboard"
)

func Init() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
}

func Watch(quit <-chan bool) <-chan []byte {
	result := make(chan []byte)
	data := clipboard.Read(clipboard.FmtText)
	hash := sha256.Sum256(data)

	go func() {
		defer close(result)

		for {
			select {
			case <-quit:
				return

			case <-time.After(500 * time.Millisecond):
				data = clipboard.Read(clipboard.FmtText)
				newHash := sha256.Sum256(data)
				if newHash != hash {
					hash = newHash
					result <- data
				}
			}
		}
	}()

	return result
}
