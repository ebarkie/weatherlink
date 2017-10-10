// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

import (
	"encoding/hex"
	"fmt"
)

// getLoops starts a stream of loop packets and sends them to the
// event channel. It exits when either numLoops is hit, an archive
// record was written, or a command is pending.
func (w *Weatherlink) getLoops(ec chan interface{}) (err error) {
	// If the archive period is 5m then this will never get exhausted
	// because the next archive record will change first, triggering an
	// exit.
	const numLoops = 165 // 2 seconds each * 165 = ~5m30s

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
			// LOOP stream was interrupted before we received all of the
			// expected packets.
			Warn.Printf("Loop stream %d/%d read interrupted: %s, aborting",
				loopNum, numLoops, err.Error())
			break
		}

		err = l.FromPacket(p)
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
		Trace.Printf("Decoded\n%+v", l)

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
		//
		// Only preemt the get loops cycle if the last archive time was set manually or
		// as the result of a (user initiated) CmdGetDmps command, indicating the user is
		// interested in archive records.
		if nextArchRec < 0 {
			nextArchRec = l.nextArchRec
		} else if !w.LastDmpTime.IsZero() && nextArchRec != l.nextArchRec {
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
