// Copyright (c) 2020 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package data

// Packet coding logic for console NVER and VER commands.
//
// Refer to Vantage Pro™, Vantage Pro2™ and Vantage Vue™ Serial
// Communication Reference Manual, section VIII. Command Summary,
// subsection 1. Testing Commands.

import (
	"strings"
	"time"
)

// FirmTime is the firmware build time.
type FirmTime time.Time

// MarshalText encodes the firmware build time into a 13-byte packet suitable
// for the VER command.
func (ft FirmTime) MarshalText() ([]byte, error) {
	return time.Time(ft).AppendFormat([]byte{}, "Jan 02 2006\n\r"), nil
}

// UnmarshalText decodes a 13-byte firmware build time response packet into
// the FirmTime struct.
func (ft *FirmTime) UnmarshalText(p []byte) error {
	t, err := time.Parse("Jan 02 2006\n\r", string(p))
	*ft = FirmTime(t)

	return err
}

// FirmVer is the firmware version number.
type FirmVer string

// MarshalText encodes the firmware version into a 6-byte packet suitable
// for the NVER command.
func (fv FirmVer) MarshalText() ([]byte, error) {
	if len(fv) != 4 {
		return []byte("0.00"), ErrBadFirmVer
	}

	return []byte(fv + "\n\r"), nil
}

// UnmarshalText decodes a 6-byte firmware version response packet into the
// FirmVer struct.
func (fv *FirmVer) UnmarshalText(p []byte) error {
	s := strings.TrimRight(string(p), "\n\r")
	*fv = FirmVer(s)

	return nil
}
