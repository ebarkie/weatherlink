package main

import (
	"log"
	"time"

	"github.com/ebarkie/weatherlink"
)

func main() {
	// Open station
	w, err := weatherlink.Dial("192.168.1.254:22222")
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()
	log.Println("Opened station")

	// Goroutine to handle station events
	go func() {
		ec := w.Start()
		log.Println("Command broker started")
		// Keep retrieving events until the channel is closed
		for e := range ec {
			switch e.(type) {
			case weatherlink.Archive:
				log.Printf("Received archive record: %+v", e)
			case weatherlink.HiLows:
				log.Printf("Received record high and lows: %+v", e)
			case weatherlink.Loop:
				log.Printf("Received loop packet: %+v", e)
			default:
				log.Printf("Received unknown event of type: %T", e)
			}
		}
		log.Println("Command broker stopped")
	}()

	// Throw in some extra commands to run
	w.CmdQ <- weatherlink.CmdGetHiLows

	runTime := time.Duration(10 * time.Second)
	log.Printf("Receiving events for %s", runTime)
	time.Sleep(runTime)
	log.Println("Stopping command broker")
	w.Stop()
}
