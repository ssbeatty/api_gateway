package utils

import (
	"os"
	"time"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func TimeStandardFormat(time time.Time, preciseMode bool) string {
	var layout = StandardFormat
	if preciseMode {
		layout = PreciseFormat
	}
	return time.Format(layout)
}
