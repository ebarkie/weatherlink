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
	CmdGetEEPROM
	CmdGetHiLows
	CmdGetLoops
	CmdStop
	CmdSyncConsTime
)

// Errors.
var (
	ErrBadCRC      = errors.New("CRC check failed")
	ErrBadLocation = errors.New("Location is inconsistent")
	ErrNotDmp      = errors.New("Not a DMP metadata packet")
	ErrNotDmpB     = errors.New("Not a revision B DMP packet")
	ErrNotLoop     = errors.New("Not a loop packet")
	ErrUnknownLoop = errors.New("Loop packet type is unknown")
	ErrCmdFailed   = errors.New("Protocol command failed")
)

// Tunables.
var (
	ConsTimeSyncFreq = 24 * time.Hour
)

// Conn holds the weatherlink connnection context.
type Conn struct {
	addr string // Device address
	dev  device // Device interface (IP, serial(/USB), or simulator)

	CmdQ chan cmd // Broker command queue

	LatestArchRec time.Time // Time of the latest archive record available
	LastDmp       time.Time // Time of the last downloaded archive record
}

// Dial establishes the weatherlink connection.
func Dial(addr string) (c Conn, err error) {
	c.CmdQ = make(chan cmd)

	c.addr = addr
	err = c.open()

	return
}

// open makes the connection to the weatherlink device.  It is separate
// from Dial() so it can be used as a reconnect during hard resets without
// losing state.
func (c *Conn) open() (err error) {
	const timeout = 6 * time.Second

	Trace.Printf("Opening device %s with a %s timeout", c.addr, timeout)
	switch {
	case c.addr == "/dev/null":
		c.dev = &sim{}
	case strings.HasPrefix(c.addr, "/dev/"):
		c.dev = &serial{Timeout: timeout}
	default:
		c.dev = &ip{Timeout: timeout}
	}
	err = c.dev.Dial(c.addr)

	return
}

// Close closes the weatherlink connection.
func (c *Conn) Close() error {
	Trace.Printf("Closing device %s", c.addr)
	return c.dev.Close()
}

// softReset tries to get the weatherlink device to abort the current command
// and get into a ready state.  It's usually used to interrupt LPS or DMPAFT
// commands.
func (c *Conn) softReset() {
	const flushTime = 1 * time.Second

	c.dev.Write([]byte{lf})
	time.Sleep(flushTime)
	c.dev.Flush()
}

// writeCmd runs commands with acknowledgement. It reads a response
// packet of size n, which can be zero.
func (c *Conn) writeCmd(cmd []byte, n int) (p Packet, err error) {
	const retries = 3

	// Determine what to print when showing the command in debug mode.  If it
	// ends with a line feed it's probably printable.
	cmdStr := string(cmd[0 : len(cmd)-1])
	if cmd[len(cmd)-1] != lf {
		cmdStr = "[bytes]"
	}

	resp := make(Packet, 1)
	acked := false
	for tryNum := 0; tryNum < retries; tryNum++ {
		c.dev.Write(cmd)

		c.dev.Read(resp)
		if len(resp) > 0 && resp[0] == ack {
			acked = true
			break
		} else {
			Warn.Printf("Command '%s' bad response, retrying (%d/%d)",
				cmdStr, tryNum+1, retries)
			c.softReset()
		}
	}
	if !acked {
		Error.Printf("Command '%s' bad response after repeated attempts", cmdStr)
		err = ErrCmdFailed
		return
	}

	Debug.Printf("Command '%s' successful", cmdStr)

	// If the dataSize is 0 we are just validating the ACK and leaving
	// the rest of the response to be read elsewhere (e.g. DMP* and LPS commands).
	if n < 1 {
		return nil, nil
	}

	p = make(Packet, n)
	_, err = c.dev.ReadFull(p)
	Trace.Printf("Hex\n%s", hex.Dump(p))

	return
}

// Idler is the idle function the command broker executes when
// there are no pending commands in the queue.
type Idler func(*Conn, chan<- interface{}) error

// StdIdle is the standard idler which reads loop packets and new
// archive records when they're available.
func StdIdle(c *Conn, ec chan<- interface{}) (err error) {
	if c.LatestArchRec.After(c.LastDmp) {
		c.LastDmp, err = c.GetDmps(ec, c.LastDmp)
	} else {
		err = c.GetLoops(ec)
	}

	return
}

// Start starts the command broker.  If no commands are pending it runs
// the idler.
func (c *Conn) Start(idle Idler) <-chan interface{} {
	// Buffer the event channel to the maximum records a Vantage
	// Pro 2 console can hold in memory.  This can speed up large
	// downloads when the receiver is I/O bound with database writes.
	ec := make(chan interface{}, 5*512)

	go func() (err error) {
		defer close(ec)

		// Send a console time sync command on startup and every ConsTimeSyncFreq.
		syncConsTime := time.NewTimer(0)

		for {
			// Before we do anything make sure we're in a non-error state.
			if err != nil {
				// Try a soft-reset first.
				//
				// There's a TEST command however it's a lot more convenient to use
				// a command that follows the ACK/NAK response flow and GETTIME fits.
				Warn.Printf("%s, trying soft-reset", err.Error())
				_, err = c.GetConsTime()
				// Hard-reset if we're still in an error state.
				if err != nil {
					Error.Printf("%s, trying hard-reset", err.Error())
					c.Close()
					err = c.open()
					continue
				}
			}

			// Process command queue channel.
			Debug.Printf("%d command(s) in queue", len(c.CmdQ))
			select {
			case cmd := <-c.CmdQ:
				switch cmd {
				case CmdGetEEPROM:
					err = c.GetEEPROM(ec)
				case CmdGetDmps:
					c.LastDmp, err = c.GetDmps(ec, c.LastDmp)
				case CmdGetHiLows:
					err = c.GetHiLows(ec)
				case CmdGetLoops:
					err = c.GetLoops(ec)
				case CmdStop:
					return
				case CmdSyncConsTime:
					err = c.SyncConsTime()
				default:
					// Should never happen unless new commands Cmd*'s are added and
					// not defined here.
					Error.Printf("Unhandled command: %d", cmd)
					err = ErrCmdFailed
				}
			case <-syncConsTime.C:
				err = c.SyncConsTime()
				if err != nil {
					syncConsTime.Reset(0)
				} else {
					syncConsTime.Reset(ConsTimeSyncFreq)
				}
			default:
				err = idle(c, ec)
			}
		}
	}()

	return ec
}

// Stop stops the command broker.
func (c Conn) Stop() {
	Trace.Println("Stopping command broker by request")
	// Drain the command queue and send a stop command.
	for {
		select {
		case <-c.CmdQ:
		default:
			c.CmdQ <- CmdStop
			return
		}
	}
}
