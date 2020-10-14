// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package device

// A Weatherlink device is simulated by guessing what commands were
// requested based on the packet sizes.  It's not perfect but is a
// convenient way to allow low level protocol testing.

import (
	"io"
	"math/rand"
	"time"

	"github.com/ebarkie/weatherlink/data"
)

// Sim represents a simulted Weatherlink device.
type Sim struct {
	l            data.Loop // Current loop packet state
	nextLoopType int       // Loop type to send next (so they are interleaved)

	// lastWrite and readsSinceWrite are used by ReadFull() to determine
	// what's expected to be read.  This is simple and avoids implementing
	// a state machine.
	lastWrite       []byte
	readsSinceWrite int
}

// Dial initializes the state of a simulated Weatherlink device.
func (s *Sim) Dial(addr string) error {
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
func (s *Sim) Close() error {
	s.l = data.Loop{}
	s.nextLoopType = 0

	return nil
}

// Flush flushes the input buffers of the simulated Weatherlink device.
func (Sim) Flush() error {
	return nil
}

// Read reads up to the size of the provided byte buffer from the
// simulated Weatherlink device.
func (Sim) Read([]byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

// ReadFull reads the full size of the provided byte buffer from the
// simulted Weatherlink device.
func (s *Sim) ReadFull(b []byte) (n int, err error) {
	const ack = 0x06 // Acknowledge

	s.readsSinceWrite++

	var p []byte
	switch {
	case len(b) == 1: // Command ack
		p = []byte{ack}
	case len(b) == 6 && s.readsSinceWrite < 2: // Command OK
		p = []byte("\n\rOK\n\r")
	case string(s.lastWrite) == "GETTIME\n":
		ct := data.ConsTime(time.Now())
		p, err = ct.MarshalBinary()
	case string(s.lastWrite) == "NVER\n":
		fv := data.FirmVer("1.73")
		p, err = fv.MarshalText()
	case string(s.lastWrite) == "TEST\n":
		p = []byte("\n\rTEST\n\r")
	case string(s.lastWrite) == "VER\n":
		ft := data.FirmTime(time.Date(2002, time.April, 24, 0, 0, 0, 0, time.UTC))
		p, err = ft.MarshalText()
	case len(b) == 99: // LPS 3 x
		// Interleave loop types.
		s.l.LoopType = s.nextLoopType + 1
		s.nextLoopType = (s.nextLoopType + 1) % 2

		// Make observation values wander around like they would on a
		// real station.
		s.l.Bar.Altimeter = wander(s.l.Bar.Altimeter, 0.01)
		s.l.Bar.SeaLevel = wander(s.l.Bar.SeaLevel, 0.01)
		s.l.Bar.Station = wander(s.l.Bar.Station, 0.01)
		s.l.OutHumidity = int(wander(float64(s.l.OutHumidity), 1))
		s.l.OutTemp = wander(s.l.OutTemp, 0.5)
		s.l.Wind.Cur.Speed = int(wander(float64(s.l.Wind.Cur.Speed), 1))

		s.l.LoopType = s.nextLoopType + 1
		s.nextLoopType = (s.nextLoopType + 1) % 2

		p, err = s.l.MarshalBinary()

		// Create 2s delay between packets.
		time.Sleep(2 * time.Second)
	default:
		return 0, io.ErrUnexpectedEOF
	}

	n = copy(b, p)
	return
}

// Write simulates a write of the byte buffer.
func (s *Sim) Write(b []byte) (int, error) {
	s.lastWrite = b
	s.readsSinceWrite = 0

	return len(b), nil
}

// wander takes a value and randomly adds +/- step or zero.
func wander(v, step float64) float64 {
	rand.Seed(int64(time.Now().Nanosecond()))
	return v + float64(rand.Intn(3)-1)*step
}
