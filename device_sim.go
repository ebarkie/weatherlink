// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

// A Weatherlink device is simulated by guessing what commands were
// requested based on the packet sizes.  It's not perfect but is a
// is a convenient way to allow low level protocol testing.

import (
	"io"
	"math/rand"
	"time"
)

// Sim represents a simulted Weatherlink device.
type Sim struct {
	nextLoopType int
}

// Close closes the simulated Weatherlink device.
func (Sim) Close() error {
	return nil
}

// Flush flushes the input buffers of the simulated Weatherlink device.
func (Sim) Flush() error {
	return nil
}

// Read reads up to the size of the provided byte buffer from the
// simulated Weatherlink device.
func (Sim) Read(b []byte) (int, error) {
	switch len(b) {
	case 1:
		b[0] = ack
		return 1, nil
	default:
		Debug.Printf("Unhandled simulated read %d bytes", len(b))
		return 0, io.ErrUnexpectedEOF
	}
}

// ReadFull reads the full size of the provided byte buffer from the
// simulted Weatherlink device.
func (s *Sim) ReadFull(b []byte) (n int, err error) {
	switch len(b) {
	case 8:
		// GETTIME
		ct := ConsTime(time.Now())
		n = copy(b, ct.ToPacket())
	case 99:
		// LPS 3 x

		// Set minimal data to be useful for testing and pass any
		// QC processes.
		l := Loop{}
		l.Bar.Altimeter = 6.8
		l.Bar.SeaLevel = 25.0
		l.Bar.Station = 6.8
		l.Wind.Cur.Speed = rand.Intn(10)

		// Simulate delay between packets and interleave loop
		// types.
		time.Sleep(2 * time.Second)
		var p Packet
		p, err = l.ToPacket(s.nextLoopType + 1)
		s.nextLoopType = (s.nextLoopType + 1) % 2
		n = copy(b, p)
	default:
		Debug.Printf("Unhandled simulated read full %d bytes", len(b))
		err = io.ErrUnexpectedEOF
		return
	}

	return
}

// Write simulates a write of the byte buffer.
func (Sim) Write(b []byte) (int, error) {
	return len(b), nil
}
