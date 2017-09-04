package main

import (
	"log"
	"os"
	"time"

	"github.com/ebarkie/weatherlink"
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
		ec := w.Start()
		log.Println("Command broker started")
		// Keep retrieving events until the channel is closed
		for e := range ec {
			switch e.(type) {
			case weatherlink.Archive:
				log.Printf("Received archive record: %+v", e)
			case weatherlink.EEPROM:
				log.Printf("Received EEPROM configuration: %+v", e)
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

	// Send an explicit command
	w.CmdQ <- weatherlink.CmdGetHiLows

	// Run for a period of time and then send a stop signal
	runTime := time.Duration(30 * time.Second)
	log.Printf("Receiving events for %s", runTime)
	time.Sleep(runTime)
	log.Println("Stopping command broker")
	w.Stop()
}
