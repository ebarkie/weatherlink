// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

import (
	"encoding/hex"
	"time"

	"github.com/ebarkie/weatherlink/data"
)

// Archive interval.
const archInt = 5 * time.Minute // XXX Read from EEPROM.

// GetDmps downloads all archive records *after* lastRec and sends
// them to the event channel ordered from oldest to newest. It
// returns the time of the last record it read.
//
// If lastRec does not match an existing archive timestamp (which is the case if
// left uninitialized) then all records in memory are returned.
func (c Conn) GetDmps(ec chan<- interface{}, lastRec time.Time) (newLastRec time.Time, err error) {
	const (
		nak = 0x15 // Not acknowledge
		esc = 0x1b // Escape
	)

	Debug.Printf("Retrieving archive records since %s", lastRec)

	// If for some reason we return on error before any records are read
	// it's safer to at least return the original lastRec instead of
	// a zeroed time.
	newLastRec = lastRec

	// Setup download.
	_, err = c.writeCmd([]byte("DMPAFT\n"), []byte{ack}, 0)
	if err != nil {
		Error.Printf("DMPAFT command error: %s, aborting", err.Error())
		return
	}
	var p []byte
	p, err = data.DmpAft(lastRec).MarshalBinary()
	if err != nil {
		Error.Printf("DmpAft marshal error: %s, aborting", err.Error())
		return
	}
	p, err = c.writeCmd(p, []byte{ack}, 6)
	if err != nil {
		Error.Printf("Dmp metadata read error: %s, aborting", err.Error())
		return
	}

	// The response tells us the number of pages we need to download
	// and the offset of the first record we should look at within
	// the first page.
	dm := data.DmpMeta{}
	err = dm.UnmarshalBinary(p)
	if err != nil {
		// Most likely a CRC error so cancel gracefully.
		Error.Printf("Dmp metadata decode error: %s, aborting", err.Error())
		c.d.Write([]byte{esc})
		return
	}
	// If numPages is 0 then it means there's nothing newer than what
	// we have so we're done.
	if dm.Pages == 0 {
		Debug.Println("No newer archive records")
		return
	}

	// Start download.
	// ACK to begin and then loop through all pages we were told are
	// available.  There are 5 records per page.
	Debug.Printf("Starting %d page dmp download", dm.Pages)
	c.d.Write([]byte{ack})
	p = make([]byte, 267)
	for pageNum := 0; pageNum < dm.Pages; pageNum++ {
		_, err = c.d.ReadFull(p)
		if err != nil {
			// Page read failed before we got all of the expected pages.
			Error.Printf("Dmp download %d/%d interrupted: %s, aborting",
				pageNum, dm.Pages, err.Error())
			break
		}
		Trace.Printf("Packet\n%s", hex.Dump(p))

		d := data.Dmp{}
		err = d.UnmarshalBinary(p)
		if err == data.ErrBadCRC {
			// NAK and retry the page.
			Error.Printf("Dmp page %d/%d error: %s, retrying",
				pageNum, dm.Pages, err.Error())
			c.d.Write([]byte{nak})
			pageNum--
			continue
		} else if err != nil {
			Error.Printf("Dmp page %d/%d error: %s, aborting",
				pageNum, dm.Pages, err.Error())
			break
		}

		// We have a valid decoded archive page
		Debug.Printf("Valid dmp page (%d:%d/%d)", pageNum, int(p[0]), dm.Pages-1)
		Trace.Printf("Decoded\n%s", Sdump(d))

		for recordNum := 0; recordNum < len(d); recordNum++ {
			// On the first page skip anything before the offset
			// given during the download setup.
			//
			// On the last page, after reading at least one
			// record bail out as soon as we hit one where the
			// date is older than the previous record.
			if pageNum == 0 && recordNum < dm.FirstPageOffset {
				continue
			} else if pageNum == dm.Pages-1 &&
				lastRec != newLastRec &&
				newLastRec.After(d[recordNum].Timestamp) {
				break
			}

			newLastRec = d[recordNum].Timestamp
			ec <- d[recordNum]
			Info.Printf("Retrieved archive record for %s", d[recordNum].Timestamp)
		}

		// ACK page as received OK so the next is sent.
		c.d.Write([]byte{ack})
	}

	return
}
