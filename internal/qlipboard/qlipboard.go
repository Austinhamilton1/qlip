package qlipboard

import (
	"crypto/sha256"
	"time"

	"golang.design/x/clipboard"
)

type Format int

const Text Format = 1
const Image Format = 2

type Data struct {
	DataType Format
	Bytes    []byte
}

// This function initializes the clipboard and MUST be called
// before any subsequent calls to any of the other functions.
// It ensures that the correct libraries are installed to use the
// clipboard.
func Init() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
}

// Watches for changes in the data on the clipboard.
// Ensures that the data was changed from the last signalled
// change with a SHA256 hash.
func Watch(quit <-chan bool) <-chan Data {
	result := make(chan Data)

	// Read initial data from the clipboard
	data := clipboard.Read(clipboard.FmtText)
	if data == nil {
		data = clipboard.Read(clipboard.FmtImage)
	}
	hash := sha256.Sum256(data) // Hash ensures infinite recursion in network

	go func() {
		defer close(result)

		for {
			select {
			// quit signal channel gracefully closes Watcher
			case <-quit:
				return

			// Look for changes every half a second for practically live changes
			case <-time.After(100 * time.Millisecond):
				data = clipboard.Read(clipboard.FmtText)
				dataType := Text
				if data == nil {
					data = clipboard.Read(clipboard.FmtImage)
					dataType = Image
				}

				newHash := sha256.Sum256(data)

				if newHash != hash {
					hash = newHash
					result <- Data{dataType, data}
				}
			}
		}
	}()

	return result
}

// Synchronizes data to the clipboard on the local machine.
// Send data to the returned channel to add data to the clipboard.
func Sync(quit <-chan bool) chan<- Data {
	synchronizer := make(chan Data)

	go func() {
		for {
			select {
			// quit signal channel gracefully closes Synchronizer
			case <-quit:
				return

			// When changes are received, update the clipboard
			case signal := <-synchronizer:
				switch signal.DataType {
				case Text:
					clipboard.Write(clipboard.FmtText, signal.Bytes)
				case Image:
					clipboard.Write(clipboard.FmtImage, signal.Bytes)
				}
			}
		}
	}()

	return synchronizer
}
