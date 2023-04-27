package utils

import (
	"GC2-sheet/internal/configuration"
	"log"
)

func LogDebug(message string) {

	if configuration.GetOptionsDebug() {
		log.Println(message)
	}

}

func LogFatalDebug(message string) {

	if configuration.GetOptionsDebug() {
		log.Fatal(message)
	}

}
