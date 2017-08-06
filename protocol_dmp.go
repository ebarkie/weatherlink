// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

import (
	"encoding/hex"
	"time"
)

// getDmps retrieves all of the archive records *after* lastRecord and
// sends them to the archive channel ordered from oldest to newest. It
// also returns the timestamp of the last record it read.
//
// If lastRecord does not match an existing archive timestamp (which is the case if
// left uninitialized) then all records are returned.
func (w *Weatherlink) getDmps(ec chan interface{}, lastRecord time.Time) (newLastRecord time.Time, err error) {
	Debug.Printf("Retrieving archive records since %s", lastRecord)

	// If for some reason we return on error before any records are read
	// it's safer to at least return the original lastRecord instead of
	// a zeroed time.
	newLastRecord = lastRecord

	// Setup download.
	_, err = w.sendCommand([]byte("DMPAFT\n"), 0)
	if err != nil {
		Error.Printf("DMPAFT command error: %s, aborting", err.Error())
		return
	}
	var p Packet
	p, err = w.sendCommand(DmpAft(lastRecord).ToPacket(), 6)
	if err != nil {
		Error.Printf("Dmp metadata read error: %s, aborting", err.Error())
		return
	}
	// The response tells us the number of pages we need to download
	// and the offset of the first record we should look at within
	// the first page.
	dm := DmpMeta{}
	err = dm.FromPacket(p)
	if err != nil {
		// Most likely a CRC error so cancel gracefully.
		Error.Printf("Dmp metadata decode error: %s, aborting", err.Error())
		w.d.Write([]byte{esc})
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
	w.d.Write([]byte{ack})
	p = make(Packet, 267)
	for pageNum := 0; pageNum < dm.Pages; pageNum++ {
		_, err = w.d.ReadFull(p)
		if err != nil {
			// Page read failed before we got all of the expected pages.
			Error.Printf("Dmp download %d/%d interrupted: %s, aborting",
				pageNum, dm.Pages, err.Error())
			break
		}

		d := Dmp{}
		err = d.FromPacket(p)
		if err != nil {
			// Most likely a CRC error.  We could NAK it and retry the page
			// but in practice it's simpler and more reliable to just let
			// the next invocation retry.
			Error.Printf("Dmp page %d/%d decode error: %s, aborting",
				pageNum, dm.Pages, err.Error())
			w.d.Write([]byte{esc})
			break
		}

		// We have a valid decoded archive page
		Trace.Printf("Valid dmp packet (%d:%d/%d)",
			pageNum, int(p[0]), dm.Pages)
		Trace.Printf("Packet\n%s", hex.Dump(p))
		Trace.Printf("Decoded\n%+v", d)

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
				lastRecord != newLastRecord &&
				newLastRecord.After(d[recordNum].Timestamp) {
				break
			}

			newLastRecord = d[recordNum].Timestamp
			ec <- d[recordNum]
			Info.Printf("Retrieved archive record for %s", d[recordNum].Timestamp)
		}

		// ACK page as received OK so the next is sent.
		w.d.Write([]byte{ack})
	}

	return
}
