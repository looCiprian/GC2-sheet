package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func GetUniqueHostnameName() string {

	currentTime := time.Now()

	currentTimeS := currentTime.Format("02-01-2006")

	unixString := strconv.FormatInt(currentTime.Unix(), 10)

	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Sprintf("%s-%s", currentTime, unixString[len(unixString)-5:])
	}

	return fmt.Sprintf("%s-%s-%s", currentTimeS, hostname, unixString[len(unixString)-5:])

}
