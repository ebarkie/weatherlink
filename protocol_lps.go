// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

import (
	"encoding/hex"
	"fmt"
)

// getLoops starts a stream of loop packets and sends them to the
// loops channel. It exits when either numLoops is hit or an archive
// record was written.
func (w *Weatherlink) getLoops() (err error) {
	const numLoops = 165 // (2 seconds each * 165 = ~5m30s)

	Info.Printf("Retrieving %d loop packets", numLoops)

	// Start a stream of LOOP1&2 packets, loop through, decode, and
	// send each one to the loops channel.
	_, err = w.sendCommand([]byte(fmt.Sprintf("LPS 3 %d\n", numLoops)), 0)
	if err != nil {
		Error.Printf("LPS command error: %s, aborting", err.Error())
		return
	}

	p := make(Packet, 99)
	var l Loop
	nextArchRec := -1
	for loopNum := 0; loopNum < numLoops; loopNum++ {
		_, err = w.d.ReadFull(p)
		if err != nil {
			// LOOP stream was interrupted before we received all of the packets
			// that we expected.
			Warn.Printf("Loop stream %d/%d read interrupted: %s, aborting",
				loopNum, numLoops, err.Error())
			break
		}

		err = l.FromPacket(p)
		if err != nil {
			// Most likely a CRC error.  We are likely out of sync with the stream
			// of 99-byte LOOP packets so the safest action is to abort.
			Error.Printf("Loop stream %d/%d decode error: %s, aborting",
				loopNum, numLoops, err.Error())
			break
		}

		// We have a valid decoded packet
		Trace.Println("Valid loop")
		Trace.Printf("Packet\n%s", hex.Dump(p))
		Trace.Printf("Decoded\n%+v", l)

		// Since our Loop is combiation of LOOP1&2 don't start emitting until we have
		// at least one of each or some values will still be zeroed resulting in
		// inaccure data.
		if loopNum > 0 {
			select {
			case w.Loops <- l:
			default:
				Warn.Println("Loop channel is full, discarding latest")
			}
		}

		// A LOOP1 decode includes the next archive record indicator and if it changes we
		// want to read it immediately.
		if nextArchRec < 0 {
			nextArchRec = l.nextArchRec
		} else if nextArchRec != l.nextArchRec {
			Trace.Printf("New archive record is available (%d->%d)", nextArchRec, l.nextArchRec)
			// Send a dmp but make sure it doesn't block or we'll be deadlocked.
			select {
			case w.CmdQ <- CmdGetDmps:
			default:
				// If it's blocked then just drop it, nextArchRec is not
				// updated so this will keep retrying.
			}
		}

		// Loops are low priority so if something else is waiting to run then exit.
		if len(w.CmdQ) > 0 {
			Debug.Println("Command queue is not empty, cancelling get loops")
			w.reset()
			break
		}
	}

	return
}
