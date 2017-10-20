# Weatherlink

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](http://choosealicense.com/licenses/mit/)
[![Build Status](https://travis-ci.org/ebarkie/weatherlink.svg?branch=master)](https://travis-ci.org/ebarkie/weatherlink)

Go package for working with Davis Instruments Weatherlink protocol over a
Weatherlink IP, serial, or USB interface.

Features:
* Should work with any Davis station made after 2002.  Developed using a Vantage Pro
  2 Plus with all sensor types.
* Device support for Weatherlink IP, serial, USB (genuine or clone).
* Device simulator.
* Decodes DMP (archive), EEPROM (configuration), HILOWS, LPS 1 (loop 1), and
  LPS 2 (loop 2) events and writes them to a channel.
* Partial encoding (work in progress).
* Sync console time.
* Command broker that coordinates commands.  Use the standard idler or define a
  custom one.

## Installation

```
$ go get github.com/ebarkie/weatherlink
```

## Usage

See [USAGE](USAGE.md).

## Example

```go
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
		// Standard idler which reads loop packets and new archive
		// records when they're available.
		//ec := w.Start(weatherlink.StdIdle)

		// Custom idler which only reads loop packets and ignores archive
		// records.
		ec := w.Start(func(c *weatherlink.Conn, ec chan<- interface{}) error {
			return c.GetLoops(ec)
		})
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
	w.CmdQ <- weatherlink.GetHiLows

	// Run for a period of time and then send a stop signal
	runTime := time.Duration(6 * time.Minute)
	log.Printf("Receiving events for %s", runTime)
	time.Sleep(runTime)
	log.Println("Stopping command broker")
	w.Stop()
}
```

## License

Copyright (c) 2016-2017 Eric Barkie. All rights reserved.  
Use of this source code is governed by the MIT license
that can be found in the [LICENSE](LICENSE) file.
