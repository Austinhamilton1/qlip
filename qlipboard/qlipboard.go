package qlipboard

import (
	"crypto/sha256"
	"time"
	"fmt"

	"golang.design/x/clipboard"
)

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
func Watch(quit <-chan bool) <-chan []byte {
	result := make(chan []byte)

	textData := clipboard.Read(clipboard.FmtText)	// Initial text data so existing clipboard is not copied
	textHash := sha256.Sum256(textData)					// Text hash ensures recursive clipping does not happen
	imgData := clipboard.Read(clipboard.FmtImage)
	imgHash := sha256.Sum256(imgData)

	go func() {
		defer close(result)

		for {
			select {
			// quit signal channel gracefully closes Watcher
			case <-quit:
				return

			// Look for changes every half a second for practically live changes
			case <-time.After(100 * time.Millisecond):
				textData = clipboard.Read(clipboard.FmtText)
				imgData = clipboard.Read(clipboard.FmtImage)
				newTextHash := sha256.Sum256(textData)
				newImgHash := sha256.Sum256(imgData)

				if len(textData) == 0 {
					textHash = newTextHash
				}
				if len(imgData) == 0 {
					imgHash = newImgHash
				}
				
				if newTextHash != textHash {
					textHash = newTextHash
					result <- textData
					fmt.Println("text")
				} else if newImgHash != imgHash {
					imgHash = newImgHash
					result <- imgData
					fmt.Println("image")
				}
			}
		}
	}()

	return result
}
