// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

import (
	"io"
	"time"

	"github.com/pkg/term"
)

// Serial represents a Weatherlink serial or USB device.
type Serial struct {
	*term.Term
}

// DialSerial establishes a serial port connection with a Weatherlink device.
func DialSerial(dev string, timeout ...time.Duration) (s Serial, err error) {
	t := 6 * time.Second
	if len(timeout) > 0 {
		t = timeout[0]
	}
	s.Term, err = term.Open(dev, term.Speed(19200), term.ReadTimeout(t), term.RawMode)

	return
}

// ReadFull reads the full size of the provided byte buffer from the
// Weatherlink device.  It blocks until the entire buffer is filled
// or the timeout triggers.
func (s Serial) ReadFull(b []byte) (int, error) {
	return io.ReadFull(s.Term, b)
}
