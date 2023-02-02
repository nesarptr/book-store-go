package utils

import (
	"fmt"
	"os"
)

func RemoveImage(imgUrl string) error {
	filename := fmt.Sprintf("./images/%s", imgUrl)
	if _, err := os.Stat(filename); err == nil {
		if err := os.Remove(filename); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}
