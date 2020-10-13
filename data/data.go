// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// Package data implements encoding and decoding of Davis
// Instruments binary data types.
package data

import "errors"

// Errors.
var (
	ErrNotArcB     = errors.New("not a revision B archive record")
	ErrBadCRC      = errors.New("CRC check failed")
	ErrBadFirmVer  = errors.New("firmware version is not valid")
	ErrBadLocation = errors.New("location is inconsistent")
	ErrNotDmp      = errors.New("not a download memory page")
	ErrNotDmpMeta  = errors.New("not a download memory page metadata packet")
	ErrNotLoop     = errors.New("not a loop packet")
	ErrUnknownLoop = errors.New("unknown loop packet type")
)
