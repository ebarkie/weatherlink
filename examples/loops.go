package main

import (
	"log"

	"gitlab.com/ebarkie/weatherlink"
	"gitlab.com/ebarkie/weatherlink/data"
)

func main() {
	// Open station
	w, err := weatherlink.Dial("192.168.1.254:22222")
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	// Custom idler which only reads loop packets and ignores archive
	// records.
	ec := w.Start(func(c *weatherlink.Conn, ec chan<- interface{}) error {
		return c.GetLoops(ec)
	})

	// Keep retrieving loop events forever.
	for e := range ec {
		log.Printf("Received loop packet: %+v", e.(data.Loop))
	}
}
