package utils

import "time"

func GetCurrentDate() string {

	currentTime := time.Now()

	return currentTime.Format(time.RFC1123Z)

}
