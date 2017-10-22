// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package data

// Packet coding logic for console GETTIME and SETTIME commands.
//
// Refer to Vantage ProTM, Vantage Pro2TM and Vantage VueTM Serial
// Communication Reference Manual, section VIII. Command Summary,
// subsection 7. Configuration Commands.

import (
	"time"

	"github.com/ebarkie/weatherlink/packet"
)

// ConsTime is the console current time.
type ConsTime time.Time

// UnmarshalBinary decodes an 8-byte console time response packet into
// the ConsTime struct.
func (ct *ConsTime) UnmarshalBinary(p []byte) error {
	if packet.Crc(p) != 0 {
		return ErrBadCRC
	}

	*ct = ConsTime(packet.GetDateTime48(p, 0))

	return nil
}

// MarshalBinary encodes the console timestamp into an 8-byte packet suitable
// for the SETTIME command.
func (ct ConsTime) MarshalBinary() (p []byte, err error) {
	p = make([]byte, 8)
	packet.SetDateTime48(&p, 0, time.Time(ct))
	packet.SetCrc(&p)

	return
}