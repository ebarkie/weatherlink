// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

import (
	"encoding/hex"
	"strconv"

	"gitlab.com/ebarkie/weatherlink/data"
)

// GetLoops starts a stream of loop packets and sends them to the
// event channel. It exits when numLoops is hit, an archive record
// was written, or a command is pending.
func (c *Conn) GetLoops(ec chan<- interface{}) (err error) {
	// The preferred exit condition is sensing a new archive record so
	// try to get 30 seconds beyond that.
	numLoops := (int(archInt.Seconds()) + 30) / 2

	Info.Printf("Retrieving %d loop packets", numLoops)

	// Start a stream of LOOP1&2 packets, loop through, decode, and
	// send each one to the loops channel.
	_, err = c.writeCmd([]byte("LPS 3 "+strconv.Itoa(numLoops)+"\n"), []byte{ack}, 0)
	if err != nil {
		Error.Printf("LPS command error: %s, aborting", err.Error())
		return
	}

	p := make([]byte, 99)
	var l data.Loop
	nextArcRec := -1
	for loopNum := 0; loopNum < numLoops; loopNum++ {
		_, err = c.d.ReadFull(p)
		if err != nil {
			// LOOP stream was interrupted before we received all of the
			// expected packets.
			Warn.Printf("Loop stream %d/%d read interrupted: %s, aborting",
				loopNum, numLoops, err.Error())
			break
		}

		err = l.UnmarshalBinary(p)
		if err != nil {
			// Most likely a CRC error.  We are probably out of sync with the
			// steam of 99-byte LOOP packets so the safest action is to abort.
			Error.Printf("Loop stream %d/%d decode error: %s, aborting",
				loopNum, numLoops, err.Error())
			break
		}

		// We have a valid decoded packet
		Trace.Println("Valid loop")
		Trace.Printf("Packet\n%s", hex.Dump(p))
		Trace.Printf("Decoded\n%s", Sdump(l))

		// Since our Loop is combiation of LOOP1&2 don't start emitting until we have
		// at least one of each or some values will still be zeroed resulting in
		// inaccurate data.
		if loopNum > 0 {
			select {
			case ec <- l:
			default:
				Warn.Println("Event channel is full, discarding latest loop")
			}
		}

		// A LOOP1 decode includes the next archive record indicator and if it changes
		// a new archive record is ready to be read.
		if nextArcRec < 0 {
			nextArcRec = l.NextArcRec
		} else if nextArcRec != l.NextArcRec {
			Debug.Printf("New archive record is available (%d->%d)", nextArcRec, l.NextArcRec)
			c.NewArcRec = true
			return
		}

		// Loops are low priority so if something else is waiting to run then exit.
		if len(c.Q) > 0 {
			Debug.Println("Command queue is not empty, cancelling get loops")
			c.softReset()
			break
		}
	}

	return
}
