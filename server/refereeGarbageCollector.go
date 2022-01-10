package main

import (
	"fmt"
	"time"
)

func RefereeGarbageCollector() {
	for {
		tts := time.Hour
		fmt.Println("Next Garbage Collector running in 60 min at : ", time.Now().Local().Add(tts), " !")
		time.Sleep(tts)
		for _, refereeID := range refereeToRemove {
			fmt.Println("removing : ", refereeID)
			delete(referees, refereeID)
		}
	}
}
