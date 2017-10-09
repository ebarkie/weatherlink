// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

// Common binary packet encoding logic for packets.
//
// Refer to Vantage ProTM, Vantage Pro2TM and Vantage VueTM Serial
// Communication Reference Manual, section X. Data Formats.

import "time"

func (p Packet) set1ByteInt(i uint, v int) {
	p[i] = byte(v)
}

func (p Packet) set2ByteFloat(i uint, v float64) {
	// Encode signed two's complement.
	p[i] = byte(v)
	p[i+1] = byte(uint16(v) >> 8)
}

// setCrc sets the last 2-bytes of a given packet to the proper
// CRC value based on the rest of content.
func (p Packet) setCrc() {
	c := crc(p[0 : len(p)-2])

	p[len(p)-2] = byte(c >> 8)
	p[len(p)-1] = byte(c)
}

// set1ByteMPH sets a 1-byte MPH which is all speed values except for
// the 2 and 10 minute values in a loop2 packet.
func (p Packet) set1ByteMPH(i uint, v int) {
	p.set1ByteInt(i, v)
}

// set1ByteTemp sets a 1-byte integer temprature like extra
// sensors.
func (p Packet) set1ByteTemp(i uint, v int) {
	p.set1ByteInt(i, v+90)
}

// setDateTimeBig sets a 6-byte date and time like the console
// uses.
func (p Packet) setDateTimeBig(i uint, t time.Time) {
	p[i] = byte(t.Second())
	p[i+1] = byte(t.Minute())
	p[i+2] = byte(t.Hour())
	p[i+3] = byte(t.Day())
	p[i+4] = byte(t.Month())
	p[i+5] = byte(t.Year() - 1900)
}

// setDateTimeSmall sets a 4-byte date and time like in archive
// records.
func (p Packet) setDateTimeSmall(i uint, t time.Time) {
	// The date is stored in the first two bytes as:
	//
	//  YYYY YYYM MMMD DDDD
	// 15       8         0
	date := t.Day() + int(t.Month())*0x20 + (t.Year()-2000)*0x200
	p[i] = byte(date)
	p[i+1] = byte(date >> 8)

	// The time is stored in second two bytes stored as: hour * 100 + min
	hour := 100*t.Hour() + t.Minute()
	p[i+2] = byte(hour)
	p[i+3] = byte(hour >> 8)
}

func (p Packet) setPressure(i uint, v float64) {
	p.set2ByteFloat(i, v*1000.0)
}
