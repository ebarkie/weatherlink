// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

// Common binary packet encoding logic for packets.
//
// Refer to Vantage ProTM, Vantage Pro2TM and Vantage VueTM Serial
// Communication Reference Manual, section X. Data Formats.

import "time"

// setCrc sets the last 2-bytes of a given packet to the proper
// CRC value based on the rest of content.
func (p Packet) setCrc() {
	c := crc(p[0 : len(p)-2])

	p[len(p)-2] = uint8(c >> 8)
	p[len(p)-1] = uint8(c & 0xff)
}

// setDateTimeBig sets a 6-byte date and time like the console
// uses.
func (p Packet) setDateTimeBig(t time.Time) {
	p[0] = uint8(t.Second())
	p[1] = uint8(t.Minute())
	p[2] = uint8(t.Hour())
	p[3] = uint8(t.Day())
	p[4] = uint8(t.Month())
	p[5] = uint8(t.Year() - 1900)
}

// setDateTimeSmall sets a 4-byte date and time like in archive
// records.
func (p Packet) setDateTimeSmall(t time.Time) {
	// The date is stored in the first two bytes as:
	//
	//  YYYY YYYM MMMD DDDD
	// 15       8         0
	date := t.Day() + int(t.Month())*0x20 + (t.Year()-2000)*0x200
	p[0] = uint8(date & 0xff)
	p[1] = uint8(date >> 8)

	// The time is stored in second two bytes stored as: hour * 100 + min
	hour := 100*t.Hour() + t.Minute()
	p[2] = uint8(hour & 0xff)
	p[3] = uint8(hour >> 8)
}
