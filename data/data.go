// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// Package data implements encoding and decoding of Davis
// Instruments binary data types.
package data

import "errors"

// Errors.
var (
	ErrBadCRC      = errors.New("CRC check failed")
	ErrBadLocation = errors.New("Location is inconsistent")
	ErrNotArchive  = errors.New("Not a revision B archive record")
	ErrNotDmp      = errors.New("Not a download memory page")
	ErrNotDmpMeta  = errors.New("Not a download memory page metadata packet")
	ErrNotLoop     = errors.New("Not a loop packet")
	ErrUnknownLoop = errors.New("Unknown loop packet type")
)
