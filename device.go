// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

import "io"

// device is an interface for the protocol to use to perform basic I/O
// operations with different Weatherlink devices.
type device interface {
	io.ReadWriteCloser
	Dial(addr string) error
	Flush() error
	ReadFull(buf []byte) (n int, err error)
}
