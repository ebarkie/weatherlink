# Weatherlink

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](http://choosealicense.com/licenses/mit/)
[![Build Status](https://travis-ci.org/ebarkie/weatherlink.svg?branch=master)](https://travis-ci.org/ebarkie/weatherlink)

Go package for working with Davis Instruments Weatherlink protocol over a
Weatherlink IP, serial, or USB interface.

Current features include:
* Should work with any Davis station made after 2002.  Developed for a Vantage Pro
  2 Plus with all sensor types.
* Supports Weatherlink IP, serial, or USB (genuine or clone).
* Decodes DMP (archive), HILOWS, LPS 1 (loop 1), and LPS 2 (loop 2) events and sends
  them over a channel.
* Syncs console time.
* Includes a command broker that attempts to intelligently select what
  commands should be run but also accepts explicit commands.

Future features:
* Encode Weatherlink packets.  Useful for creating a virtual Weatherlink IP,
  even a multiplexed one.
* Support additional packet types like EERD (useful for getting lat/lon).

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

	// Goroutine to receive station events
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

	// Send an explicit command
	w.CmdQ <- weatherlink.CmdGetHiLows

	// Run for a period of time and then send a stop signal
	runTime := time.Duration(10 * time.Second)
	log.Printf("Receiving events for %s", runTime)
	time.Sleep(runTime)
	log.Println("Stopping command broker")
	w.Stop()
	// (program exits when goroutine is actually stopped)
}
```

## License

Copyright (c) 2016-2017 Eric Barkie. All rights reserved.  
Use of this source code is governed by the MIT license
that can be found in the [LICENSE](LICENSE) file.
