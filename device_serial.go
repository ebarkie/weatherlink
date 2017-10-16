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
	Timeout time.Duration
}

// Dial opens a serial port connection with a weatherlink device.
func (s *Serial) Dial(addr string) (err error) {
	s.Term, err = term.Open(addr, term.Speed(19200), term.ReadTimeout(s.Timeout), term.RawMode)

	return
}

// ReadFull reads the full size of the provided byte buffer from the
// Weatherlink device.  It blocks until the entire buffer is filled
// or the timeout triggers.
func (s Serial) ReadFull(b []byte) (int, error) {
	return io.ReadFull(s.Term, b)
}
