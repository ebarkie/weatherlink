// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package packet

import "time"

// SetCrc sets the last 2-bytes of a given packet to the proper
// CRC value based on the rest of content.
func SetCrc(p *[]byte) {
	c := Crc((*p)[0 : len(*p)-2])
	(*p)[len(*p)-2] = byte(c >> 8)
	(*p)[len(*p)-1] = byte(c)
}

// SetDateTime32 sets a 4-byte date and time value in a given packet
// at the specified index.
func SetDateTime32(p *[]byte, i uint, t time.Time) {
	// The date is stored in the first two bytes as:
	//
	//  YYYY YYYM MMMD DDDD
	// 15       8         0
	date := t.Day() + int(t.Month())*0x20 + (t.Year()-2000)*0x200
	(*p)[i] = byte(date)
	(*p)[i+1] = byte(date >> 8)

	// The time is stored in second two bytes stored as: hour * 100 + min
	hour := 100*t.Hour() + t.Minute()
	(*p)[i+2] = byte(hour)
	(*p)[i+3] = byte(hour >> 8)
}

// SetDateTime48 sets a 6-byte date and time value in a given packet
// at the specified index.
func SetDateTime48(p *[]byte, i uint, t time.Time) {
	(*p)[i] = byte(t.Second())
	(*p)[i+1] = byte(t.Minute())
	(*p)[i+2] = byte(t.Hour())
	(*p)[i+3] = byte(t.Day())
	(*p)[i+4] = byte(t.Month())
	(*p)[i+5] = byte(t.Year() - 1900)
}

// SetFloat16 sets a 2-byte signed two's complement float value in
// a given packet at the specified index.
func SetFloat16(p *[]byte, i uint, v float64) {
	(*p)[i] = byte(v)
	(*p)[i+1] = byte(uint16(v) >> 8)
}

// SetFloat16_10 sets a 2-byte signed two's complement float value
// in tenths in a given packet at the specified index.
func SetFloat16_10(p *[]byte, i uint, v float64) {
	SetFloat16(p, i, v*10.0)
}

// SetInt8 sets a 1-byte integer value in a given packet at the
// specified index.
func SetInt8(p *[]byte, i uint, v int) {
	(*p)[i] = byte(v)
}

// SetInt16 sets a 2-byte integer value in a given packet at the
// specified index.
func SetInt16(p *[]byte, i uint, v int) {
	(*p)[i] = byte(v)
	(*p)[i+1] = byte(uint16(v) >> 8)
}

// SetMPH8 sets a 1-byte MPH value in a given packet at the specified
// index.
func SetMPH8(p *[]byte, i uint, v int) {
	SetInt8(p, i, v)
}

// SetMPH16 sets a 2-byte MPH value in a given packet at the
// specified index.
func SetMPH16(p *[]byte, i uint, v float64) {
	SetFloat16(p, i, v*10.0)
}

// SetPressure sets a pressure value in a given packet at the
// specified index.
func SetPressure(p *[]byte, i uint, v float64) {
	SetFloat16(p, i, v*1000.0)
}

// SetRainClicks sets a rain rate or accumulation value
// in a given packet at the specified index.
func SetRainClicks(p *[]byte, i uint, v float64) {
	SetFloat16(p, i, v*100.0)
}

// SetTemp8 sets a 1-byte temprature value in a given packet at
// the specified index.
func SetTemp8(p *[]byte, i uint, v int) {
	SetInt8(p, i, v+90)
}

// SetVoltage sets a battery voltage value in a given packet
// at the specified index.
func SetVoltage(p *[]byte, i uint, v float64) {
	SetFloat16(p, i, v*300.0*512.0*100.0)
}
