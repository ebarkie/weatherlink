// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// Package weatherlink implements the Davis Instruments serial, USB, and
// TCP/IP communication protocol.
package weatherlink

import (
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

const (
	cr  = 0x0d // Carriage return
	lf  = 0x0a // Line Feed
	ack = 0x06 // Acknowledge
	nak = 0x15 // Not acknowledge
	esc = 0x1b // Escape
)

type cmd uint8

// Commands that can be requested.
const (
	CmdGetDmps cmd = iota
	CmdGetLoops
	CmdStop
	CmdSyncConsTime
)

// Errors.
var (
	ErrProtoCmdFailed = errors.New("Protocol command failed")
	ErrNoLoopChan     = errors.New("Can't start command broker without a Loop channel")
)

// Weatherlink is used to track the Weatherlink device.
type Weatherlink struct {
	dev string // Device name is saved to re-connect if a hard reset is necessary
	d   Device // Device interface is either IP or Serial

	Archive     chan Archive
	LastDmpTime time.Time
	Loops       chan Loop
	CmdQ        chan cmd
}

// Dial opens the connection to the Weatherlink.
func Dial(dev string) (w Weatherlink, err error) {
	w.CmdQ = make(chan cmd, 16) // XXX What's the right command buffer size?

	w.dev = dev
	err = w.open()

	return
}

// open connects to the Weatherlink.  It's split from Dial so it can be used
// for hard-resets without losing the previous state.
func (w *Weatherlink) open() (err error) {
	const rwTimeout = 6 * time.Second

	if strings.HasPrefix(w.dev, "/dev/") {
		w.d, err = DialSerial(w.dev, rwTimeout)
	} else {
		w.d, err = DialIP(w.dev, rwTimeout)
	}

	return
}

// Close closes the connection to the Weatherlink.
func (w *Weatherlink) Close() error {
	return w.d.Close()
}

// reset tries to get the Weatherlink device in a state where it's responding
// to commands.   It's usually used to interrupt a LPS or DMPAFT command.
func (w *Weatherlink) reset() {
	const flushTime = 1 * time.Second

	w.d.Write([]byte{lf})
	time.Sleep(flushTime)
	w.d.Flush()
}

// sendCommand is used to send commands to the Weatherlink and check to make sure
// the command is ACKnowledged.  Optionally it will do a readFull on ps
// bytes, if ps is greater than 0.
func (w *Weatherlink) sendCommand(c []byte, ps int) (p Packet, err error) {
	const retries = 3

	// Determine what to print when showing the command in debug mode.  If
	// it ends with a line-feed it usually means it's printable.
	printableCommand := string(c[0 : len(c)-1])
	if c[len(c)-1] != lf {
		printableCommand = "[bytes]"
	}

	response := make(Packet, 1)
	acked := false
	for tryNum := 0; tryNum < retries; tryNum++ {
		w.d.Write(c)

		w.d.Read(response)
		if (len(response) > 0) && (response[0] == ack) {
			acked = true
			break
		} else {
			Warn.Printf("Command '%s' bad response, retrying (%d/%d)",
				printableCommand, tryNum+1, retries)
			w.reset()
		}
	}
	if !acked {
		Error.Printf("Command '%s' bad response after repeated attempts", printableCommand)
		err = ErrProtoCmdFailed
		return
	}

	Debug.Printf("Command '%s' successful", printableCommand)

	// If the dataSize is 0 we are just validating the ACK and leaving
	// the rest of the response to be read elsewhere (e.g. DMP* and LPS commands).
	if ps < 1 {
		return nil, nil
	}

	p = make(Packet, ps)
	_, err = w.d.ReadFull(p)
	Trace.Printf("Hex\n%s", hex.Dump(p))

	return
}

// Start starts the command broker.  It attempts to intelligently select what
// commands should be run but also accepts commands via the CmdQ
// channel.  The channel is especially useful for building multiplexing
// services.
func (w *Weatherlink) Start() (err error) {
	const syncConsTimerFreq = 24 * time.Hour

	if w.Loops == nil {
		err = ErrNoLoopChan
		return
	}

	// Send an dmp command on startup before doing anything else.
	w.CmdQ <- CmdGetDmps

	// Send a console time sync command on startup and every syncConsTimerFreq.
	syncConsTimer := time.NewTimer(0)

	for {
		// Before we do anything make sure we're in a non-error state.
		if err != nil {
			// Try a soft-reset first.
			//
			// There's a TEST command however it's a lot more convenient to use
			// a command that follows the ACK/NAK response flow and GETTIME fits.
			Warn.Printf("%s, trying soft-reset", err.Error())
			_, err = w.getConsTime()
			// Hard-reset if we're still in an error state.
			if err != nil {
				Error.Printf("%s, trying hard-reset", err.Error())
				w.Close()
				err = w.open()
				continue
			}
		}

		// Process command queue channel.
		Debug.Printf("%d command(s) in queue", len(w.CmdQ))
		select {
		case c := <-w.CmdQ:
			switch c {
			case CmdGetDmps:
				w.LastDmpTime, err = w.getDmps(w.LastDmpTime)
			case CmdGetLoops:
				err = w.getLoops()
			case CmdStop:
				return
			case CmdSyncConsTime:
				err = w.syncConsTime()
			default:
				// Should never happen unless new commands Cmd*'s are added and
				// not defined here.
				Error.Printf("Unhandled command")
			}
		case <-syncConsTimer.C:
			err = w.syncConsTime()
			syncConsTimer.Reset(syncConsTimerFreq)
		default:
			// If there's nothing in the command queue then poll loops.
			err = w.getLoops()
		}
	}
}

// Stop stops the command broker.
func (w Weatherlink) Stop() {
	// Drain the command queue and then send a stop command.
	for range w.CmdQ {
	}
	w.CmdQ <- CmdStop
}
