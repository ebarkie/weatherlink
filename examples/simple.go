package main

import (
	"log"
	"os"
	"time"

	"github.com/ebarkie/weatherlink"
	"github.com/ebarkie/weatherlink/data"
)

func main() {
	// Enable weatherlink logging to standard output
	weatherlink.Error.SetOutput(os.Stdout)
	weatherlink.Warn.SetOutput(os.Stdout)
	weatherlink.Info.SetOutput(os.Stdout)

	// Open station
	w, err := weatherlink.Dial("192.168.1.254:22222")
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()
	log.Println("Opened station")

	// Goroutine to receive station events
	go func() {
		// Standard idler which reads loop packets and new archive
		// records when they're available.
		ec := w.Start(weatherlink.StdIdle)
		log.Println("Command broker started")

		// Keep retrieving events until the channel is closed
		for e := range ec {
			switch e.(type) {
			case data.Archive:
				log.Printf("Received archive record: %+v", e)
			case data.EEPROM:
				log.Printf("Received EEPROM configuration: %+v", e)
			case data.HiLows:
				log.Printf("Received record high and lows: %+v", e)
			case data.Loop:
				log.Printf("Received loop packet: %+v", e)
			default:
				log.Printf("Received unknown event of type: %T", e)
			}
		}
		log.Println("Command broker stopped")
	}()

	// Send an explicit command
	w.Q <- weatherlink.GetHiLows

	// Run for a period of time and then send a stop signal
	runTime := time.Duration(6 * time.Minute)
	log.Printf("Receiving events for %s", runTime)
	time.Sleep(runTime)
	log.Println("Stopping command broker")
	w.Stop()
}
