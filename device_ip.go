// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

import (
	"io"
	"net"
	"time"
)

// IP represents a Weatherlink IP.
type IP struct {
	conn    net.Conn
	Timeout time.Duration
}

// DialIP establishes a TCP connection with a Weatherlink IP.
func DialIP(address string, timeout ...time.Duration) (i IP, err error) {
	t := 6 * time.Second
	if len(timeout) > 0 {
		t = timeout[0]
	}
	i.Timeout = t
	i.conn, err = net.Dial("tcp", address)

	return
}

// Close closes the TCP connection of the Weatherlink IP.
func (i IP) Close() error {
	return i.conn.Close()
}

// Flush flushes the input buffers of the Weatherlink IP.
func (i IP) Flush() error {
	// No lower level flush is available so allocate an absurdly
	// large buffer and read everything we can, expecting the
	// timeout to kick in.
	b := make([]byte, 8*1024)
	i.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	i.ReadFull(b)
	return nil
}

// Read reads up to the size of the provided byte buffer from the
// Weatherlink IP.  It blocks until at least one byte is read
// or the timeout triggers.  In practice, exactly how much it
// reads beyond one byte seems unpredictable.
func (i IP) Read(b []byte) (int, error) {
	i.conn.SetReadDeadline(time.Now().Add(i.Timeout))
	return i.conn.Read(b)
}

// ReadFull reads the full size of the provided byte buffer from the
// Weatherlink IP.  It blocks until the entire buffer is filled
// or the timeout triggers.
func (i IP) ReadFull(b []byte) (int, error) {
	i.conn.SetReadDeadline(time.Now().Add(i.Timeout))
	return io.ReadFull(i.conn, b)
}

// Write writes the byte buffer to the Weatherlink IP.
func (i IP) Write(b []byte) (int, error) {
	i.conn.SetWriteDeadline(time.Now().Add(i.Timeout))
	return i.conn.Write(b)
}
