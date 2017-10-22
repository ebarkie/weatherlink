// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// Package weatherlink implements the Davis Instruments serial, USB, and
// TCP/IP communication protocol.
package weatherlink

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/ebarkie/weatherlink/internal/device"
)

const (
	cr  = 0x0d // Carriage return
	lf  = 0x0a // Line Feed
	ack = 0x06 // Acknowledge
	nak = 0x21 // Not acknowledge
	esc = 0x18 // Escape
)

type cmd uint8

// Commands.
const (
	GetDmps cmd = iota
	GetEEPROM
	GetHiLows
	GetLoops
	LampsOff
	LampsOn
	Stop
	SyncConsTime
)

// Errors.
var (
	ErrCmdFailed = errors.New("Command failed")
)

// Tunables.
var (
	ConsTimeSyncFreq = 24 * time.Hour
)

// dev is an interface for the protocol to use to perform basic I/O
// operations with different Weatherlink devices.
type dev interface {
	io.ReadWriteCloser
	Dial(addr string) error
	Flush() error
	ReadFull(buf []byte) (n int, err error)
}

// Conn holds the weatherlink connnection context.
type Conn struct {
	addr string // Device address
	d    dev    // Device interface (IP, serial(/USB), or simulator)

	LastDmp   time.Time // Time of the last downloaded archive record
	NewArcRec bool      // Indicates a new archive record is available

	Q chan cmd // Command queue
}

// Dial establishes the weatherlink connection.
func Dial(addr string) (c Conn, err error) {
	c.Q = make(chan cmd, 1)

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
		c.d = &device.Sim{}
	case strings.HasPrefix(c.addr, "/dev/"):
		c.d = &device.Serial{Timeout: timeout}
	default:
		c.d = &device.IP{Timeout: timeout}
	}
	err = c.d.Dial(c.addr)

	return
}

// Close closes the weatherlink connection.
func (c Conn) Close() error {
	Trace.Printf("Closing device %s", c.addr)
	return c.d.Close()
}

// softReset tries to get the weatherlink device to abort the current command
// and get into a ready state.  It's usually used to interrupt LPS or DMPAFT
// commands.
func (c Conn) softReset() {
	const flushTime = 1 * time.Second

	c.d.Write([]byte{lf})
	time.Sleep(flushTime)
	c.d.Flush()
}

// test sends a test command.
func (c Conn) test() (err error) {
	_, err = c.writeCmd([]byte("TEST\n"), []byte{lf, cr, 'T', 'E', 'S', 'T', lf, cr}, 0)

	return
}

// writeCmd runs a command and requires an acknowledgement response.  If n > 0
// then a Packet of that length will be read after the acknowledgement.
func (c Conn) writeCmd(cmd []byte, cmdAck []byte, n int) (p []byte, err error) {
	const retries = 3

	// Determine what to print when showing the command in debug mode.  If it
	// ends with a line feed it's probably printable.
	cmdStr := string(cmd[0 : len(cmd)-1])
	if cmd[len(cmd)-1] != lf {
		cmdStr = "[bytes]"
	}

	resp := make([]byte, len(cmdAck))
	acked := false
	for tryNum := 0; tryNum < retries; tryNum++ {
		Trace.Printf("Command\n%s", hex.Dump(cmd))
		c.d.Write(cmd)

		c.d.ReadFull(resp)
		if bytes.Equal(cmdAck, resp) {
			acked = true
			break
		} else {
			Trace.Printf("Expected ack\n%s", hex.Dump(cmdAck))
			Trace.Printf("Actual ack\n%s", hex.Dump(resp))
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

	p = make([]byte, n)
	_, err = c.d.ReadFull(p)
	Trace.Printf("Packet\n%s", hex.Dump(p))

	return
}

// Idler is the idle function the command broker executes when
// there are no pending commands in the queue.
type Idler func(*Conn, chan<- interface{}) error

// StdIdle is the standard idler which reads loop packets and new
// archive records when they're available.
func StdIdle(c *Conn, ec chan<- interface{}) (err error) {
	if c.NewArcRec {
		c.NewArcRec = false
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
				Warn.Printf("%s, trying soft-reset", err.Error())
				err = c.test()
				// Hard-reset if we're still in an error state.
				if err != nil {
					Error.Printf("%s, trying hard-reset", err.Error())
					c.Close()
					err = c.open()
					continue
				}
			}

			// Process command queue channel.
			Debug.Printf("%d command(s) in queue", len(c.Q))
			select {
			case cmd := <-c.Q:
				switch cmd {
				case GetEEPROM:
					err = c.GetEEPROM(ec)
				case GetDmps:
					c.LastDmp, err = c.GetDmps(ec, c.LastDmp)
				case GetHiLows:
					err = c.GetHiLows(ec)
				case GetLoops:
					err = c.GetLoops(ec)
				case LampsOff:
					err = c.SetLamps(false)
				case LampsOn:
					err = c.SetLamps(true)
				case Stop:
					return
				case SyncConsTime:
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
		case <-c.Q:
		default:
			c.Q <- Stop
			return
		}
	}
}
