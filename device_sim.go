// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

// A Weatherlink device is simulated by guessing what commands were
// requested based on the packet sizes.  It's not perfect but is a
// convenient way to allow low level protocol testing.

import (
	"io"
	"math/rand"
	"time"
)

// sim represents a simulted Weatherlink device.
type sim struct {
	l            Loop // Current loop packet state
	nextLoopType int  // Loop type to send next (so they are interleaved)
}

// Dial initializes the state of a simulated Weatherlink device.
func (s *sim) Dial(addr string) error {
	// Starting loop values which will pass typical QC processes.
	s.l.Bar.Altimeter = 29.0
	s.l.Bar.SeaLevel = 29.0
	s.l.Bar.Station = 29.0
	s.l.OutHumidity = 50
	s.l.OutTemp = 65.0
	s.l.Wind.Cur.Speed = 3

	return nil
}

// Close closes the simulated Weatherlink device.
func (s *sim) Close() error {
	s.l = Loop{}
	s.nextLoopType = 0

	return nil
}

// Flush flushes the input buffers of the simulated Weatherlink device.
func (sim) Flush() error {
	return nil
}

// Read reads up to the size of the provided byte buffer from the
// simulated Weatherlink device.
func (sim) Read(b []byte) (int, error) {
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
func (s *sim) ReadFull(b []byte) (n int, err error) {
	switch len(b) {
	case 8:
		// GETTIME
		ct := ConsTime(time.Now())
		n = copy(b, ct.ToPacket())
	case 99:
		// LPS 3 x

		// Make observation values wander around like they would on a
		// real station.
		s.l.Bar.Altimeter = wander(s.l.Bar.Altimeter, 0.01)
		s.l.Bar.SeaLevel = wander(s.l.Bar.SeaLevel, 0.01)
		s.l.Bar.Station = wander(s.l.Bar.Station, 0.01)
		s.l.OutHumidity = int(wander(float64(s.l.OutHumidity), 1))
		s.l.OutTemp = wander(s.l.OutTemp, 0.5)
		s.l.Wind.Cur.Speed = int(wander(float64(s.l.Wind.Cur.Speed), 1))

		// Create 2s delay between packets and interleave loop
		// types.
		time.Sleep(2 * time.Second)
		var p Packet
		p, err = s.l.ToPacket(s.nextLoopType + 1)
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
func (sim) Write(b []byte) (int, error) {
	return len(b), nil
}

// wander takes a value and randomly adds +/- step or zero.
func wander(v, step float64) float64 {
	rand.Seed(int64(time.Now().Nanosecond()))
	return v + float64(rand.Intn(3)-1)*step
}
