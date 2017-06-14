# Weatherlink

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](http://choosealicense.com/licenses/mit/)
[![Build Status](https://travis-ci.org/ebarkie/weatherlink.svg?branch=master)](https://travis-ci.org/ebarkie/weatherlink)

Go package for working with Davis Instruments Weatherlink protocol over a
Weatherlink IP, serial, or USB interface.

Current features include:
* Should work with any Davis station made after 2002.  Developed for a Vantage Pro
  2 Plus with all sensor types.
* Supports Weatherlink IP, serial, or USB (genuine or clone).
* Decodes DMP (archive), LOOP1, and LOOP2 packets and sends data over
  Go channels.
* Syncs console time.
* Includes a command broker that attempts to intelligently select what
  commands should be run but also accepts commands via a channel.

Future features:
* Encode Weatherlink packets.  Useful for creating a virtual Weatherlink IP,
  even a multiplexed one.
* Encode/decode additional packet types like HILOW.

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

	"github.com/ebarkie/weatherlink"
)

func main() {
	archive := make(chan weatherlink.Archive)
	loops := make(chan weatherlink.Loop)

	// Start retrieving data from weather station
	go func() {
		wl, err := weatherlink.Dial("192.168.1.254:22222")
		if err != nil {
			log.Fatal(err)
		}
		defer wl.Close()
		wl.Archive = archive
		wl.Loops = loops

		err = wl.Start() // Block forever
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Get incoming records
	go func() {
		for {
			select {
			case a := <-archive:
				log.Printf("Received archive record: %+v", a)

				// Do something
			case l := <-loops:
				log.Printf("Received loop record: %+v", l)

				// Do something
			}
		}
	}()
}
```

## License

Copyright (c) 2016-2017 Eric Barkie. All rights reserved.  
Use of this source code is governed by the MIT license
that can be found in the [LICENSE](LICENSE) file.
