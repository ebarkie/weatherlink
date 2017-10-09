// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

// Packet coding logic for console GETTIME and SETTIME commands.
//
// Refer to Vantage ProTM, Vantage Pro2TM and Vantage VueTM Serial
// Communication Reference Manual, section VIII. Command Summary,
// subsection 7. Configuration Commands.

import "time"

// ConsTime is the console current time.
type ConsTime time.Time

// FromPacket unpacks an 8-byte console time response packet into
// the ConsTime struct.
func (ct *ConsTime) FromPacket(p Packet) error {
	if crc(p) != 0 {
		return ErrBadCRC
	}

	*ct = ConsTime(p.get6ByteDateTime(0))

	return nil
}

// ToPacket packs the console timestamp into an 8-byte packet suitable
// for the SETTIME command.
func (ct ConsTime) ToPacket() (p Packet) {
	p = make(Packet, 8)
	p.setDateTimeBig(0, time.Time(ct))
	p.setCrc()

	return
}
