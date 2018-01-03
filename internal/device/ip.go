// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package device

import (
	"io"
	"net"
	"time"
)

// IP represents a Weatherlink IP device.
type IP struct {
	net.Conn
	Timeout time.Duration
}

// Dial establishes a TCP/IP connection with a Weatherlink IP.
func (i *IP) Dial(addr string) (err error) {
	i.Conn, err = net.Dial("tcp", addr)

	return
}

// Flush flushes the input buffers of the Weatherlink IP.
func (i IP) Flush() error {
	// No lower level flush is available so allocate an absurdly
	// large buffer and read what we can.
	b := make([]byte, 8*1024)
	i.Conn.SetReadDeadline(time.Now().Add(i.Timeout))
	i.Read(b)

	return nil
}

// Read reads up to the size of the provided byte buffer from the
// Weatherlink IP.  It blocks until at least one byte is read
// or the timeout triggers.  In practice, exactly how much it
// reads beyond one byte seems unpredictable.
func (i IP) Read(b []byte) (int, error) {
	i.Conn.SetReadDeadline(time.Now().Add(i.Timeout))
	return i.Conn.Read(b)
}

// ReadFull reads the full size of the provided byte buffer from the
// Weatherlink IP.  It blocks until the entire buffer is filled
// or the timeout triggers.
func (i IP) ReadFull(b []byte) (int, error) {
	i.Conn.SetReadDeadline(time.Now().Add(i.Timeout))
	return io.ReadFull(i.Conn, b)
}

// Write writes the byte buffer to the Weatherlink IP.
func (i IP) Write(b []byte) (int, error) {
	i.Conn.SetWriteDeadline(time.Now().Add(i.Timeout))
	return i.Conn.Write(b)
}
