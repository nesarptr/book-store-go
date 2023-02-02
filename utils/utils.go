package utils

import (
	"fmt"
	"os"
	"time"
)

func RemoveImage(filename string, maxRetries int, done chan error) {
	imgUrl := fmt.Sprintf("./images/%s", filename)
	if _, err := os.Stat(imgUrl); err == nil {
		for i := 0; i < maxRetries; i++ {
			if err := os.Remove(imgUrl); err == nil {
				done <- nil
			}
			time.Sleep(time.Second)
		}
		done <- fmt.Errorf("failed to remove file after %d retries", maxRetries)
	} else {
		done <- err
	}
}
