// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package packet

// Common binary packet decoding logic for packets.
//
// Refer to Vantage Pro™, Vantage Pro2™ and Vantage Vue™ Serial
// Communication Reference Manual, section X. Data Formats.

import (
	"strings"
	"time"
)

// Barometer trends.
const (
	Dash         = "-"
	FallingRapid = "Falling Rapidly"
	FallingSlow  = "Falling Slowly"
	Steady       = "Steady"
	RisingSlow   = "Rising Slowly"
	RisingRapid  = "Rising Rapidly"
)

// GetBarTrend gets a barometer trend from a given packet at
// the specified index.
func GetBarTrend(p []byte, i uint) string {
	switch GetUInt8(p, i) {
	case -60:
		return FallingRapid
	case -20:
		return FallingSlow
	case 0:
		return Steady
	case 20:
		return RisingSlow
	case 60:
		return RisingRapid
	default:
		return Dash
	}
}

// GetDate16 gets a 2-byte date (no time) value from a given packet
// at the specified index.
func GetDate16(p []byte, i uint) time.Time {
	// If unitialized then return a zero Time.
	d := GetUInt16(p, i)
	if d == 0xffff {
		return time.Time{}
	}

	// The date is stored in the two bytes as:
	//
	//  MMMM DDDD DYYY YYYY
	// 15       8         0
	year := 2000 + d&0x007f
	day := d & 0x0f80 >> 7
	month := d & 0xf000 >> 12

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
}

// GetDateTime32 gets a 4-byte date and time value from a given packet
// at the specified index.
func GetDateTime32(p []byte, i uint) time.Time {
	// The date is stored in the first two bytes as:
	//
	//  YYYY YYYM MMMD DDDD
	// 15       8         0
	d := GetUInt16(p, i)
	day := d & 0x001f
	month := (d & 0x01e0) >> 5
	year := 2000 + (d&0xfe00)>>9

	// The time is stored in second two bytes stored as: hour * 100 + min
	t := GetUInt16(p, i+2)
	hour := t / 100
	minute := t % 100

	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.Local)
}

// GetDateTime48 gets a 6-byte date and time value from a given packet
// at the specified index.
func GetDateTime48(p []byte, i uint) time.Time {
	second := GetUInt8(p, i)
	minute := GetUInt8(p, i+1)
	hour := GetUInt8(p, i+2)
	day := GetUInt8(p, i+3)
	month := GetUInt8(p, i+4)
	year := 1900 + GetUInt8(p, i+5)

	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
}

// GetFloat16 gets a 2-byte signed two's complement float value from
// a given packet at the specified index.
func GetFloat16(p []byte, i uint) float64 {
	return float64(int16(uint16(p[i+1])<<8 | uint16(p[i])))
}

// GetFloat16_10 gets a 2-byte signed two's complement float value
// in tenths in a given packet at the specified index.
func GetFloat16_10(p []byte, i uint) float64 {
	return GetFloat16(p, i) / 10.0
}

// GetForecast gets a forecast string from a given packet at the
// specified index.
func GetForecast(p []byte, i uint) string {
	// There are 51 unique forecast messages and each rule can combine
	// up to 3 of them.

	// Rule to message mapping.
	var rules = [][]int{
		{0},
		{1},
		{2},
		{3},
		{1},
		{47},
		{48},
		{48},
		{4},
		{48},
		{48},
		{1},
		{11, 40},
		{48},
		{1},
		{12, 39},
		{1},
		{48},
		{1},
		{12, 34},
		{1},
		{48},
		{1},
		{11, 39},
		{4, 43},
		{48},
		{1},
		{11, 34, 43},
		{4, 43},
		{11},
		{48},
		{1},
		{11, 34, 43},
		{4, 43},
		{11},
		{48},
		{1},
		{11, 34, 43},
		{48},
		{1},
		{4, 41},
		{4},
		{48},
		{1},
		{12, 40},
		{12},
		{48},
		{1},
		{11, 38},
		{48},
		{1},
		{11, 38, 44},
		{48},
		{1},
		{11, 38, 44},
		{48},
		{1},
		{11, 33},
		{48},
		{1},
		{11, 33, 44},
		{48},
		{1},
		{11, 38, 44},
		{48},
		{1},
		{11, 34},
		{48},
		{1},
		{11, 23},
		{14, 25},
		{48},
		{14, 25},
		{1},
		{14, 25},
		{47},
		{48},
		{0},
		{14, 25},
		{1},
		{14, 25},
		{0},
		{48},
		{1},
		{12, 39},
		{18, 21},
		{48},
		{1},
		{18, 23},
		{19, 21},
		{19, 23},
		{48},
		{1},
		{13, 30},
		{12, 30},
		{18, 21, 43},
		{48},
		{1},
		{18, 23, 43},
		{19, 21, 43},
		{19, 23, 43},
		{48},
		{1},
		{13, 38, 45},
		{12, 38, 45},
		{48},
		{1},
		{13, 29, 45},
		{12, 29, 45},
		{18, 26, 45},
		{18, 45},
		{19, 26, 45},
		{19, 45},
		{18, 26, 45},
		{48},
		{1},
		{18, 39, 45},
		{19, 26, 45},
		{19, 39, 45},
		{15, 25},
		{15},
		{18, 25, 46},
		{18, 46},
		{15},
		{48},
		{1},
		{19, 34, 44},
		{48},
		{1},
		{13, 35, 44},
		{18, 25, 44},
		{48},
		{1},
		{18, 34, 44},
		{18, 28},
		{18},
		{18, 22, 44},
		{48},
		{1},
		{18, 33, 44},
		{19, 22, 44},
		{48},
		{1},
		{19, 33, 44},
		{48},
		{1},
		{12, 35, 44},
		{18, 44},
		{18, 22, 44},
		{48},
		{1},
		{18, 24, 44},
		{19, 22, 44},
		{19, 24, 44},
		{48},
		{1},
		{13, 29, 44},
		{12, 29, 44},
		{18, 21, 46},
		{48},
		{1},
		{18, 23, 46},
		{19, 21, 46},
		{19, 23, 46},
		{13, 29, 46},
		{48},
		{1},
		{13, 29, 45},
		{12, 29, 46},
		{12, 29, 45},
		{48},
		{1},
		{13, 29, 46},
		{12, 29, 46},
		{48},
		{1},
		{13, 38, 46},
		{12, 38, 46},
		{18, 27, 46},
		{48},
		{1},
		{18, 32, 46},
		{19, 26, 46},
		{19, 32, 46},
		{18, 21},
		{48},
		{1},
		{18, 23, 46},
		{19, 21},
		{19, 23},
		{48},
		{1},
		{18, 35, 44},
		{20}, // 193
	}

	var msg = [51]string{
		"Mostly clear and cooler.",
		"Mostly clear with little temperature change.",
		"Mostly clear for 12 hours with little temperature change.",
		"Mostly clear for 12 to 24 hours and cooler.",
		"Mostly clear and warmer.",
		"Mostly clear for 6 to 12 hours with little temperature change.",
		"Mostly clear for 12 to 24 hours with little temperature change.",
		"",
		"",
		"",
		"",
		"Increasing clouds and warmer.",
		"Increasing clouds with little temperature change.",
		"Increasing clouds and cooler.",
		"Clearing and cooler.",
		"Clearing, cooler and windy.",
		"",
		"",
		"Mostly cloudy and cooler.",
		"Mostly cloudy with little temperature change.",
		"FORECAST REQUIRES 3 HOURS OF RECENT DATA.",
		"Precipitation continuing.",
		"Precipitation continuing, possibly heavy at times.",
		"Precipitation likely.",
		"Precipitation likely, possibly heavy at times.",
		"Precipitation ending within 6 hours.",
		"Precipitation ending within 12 hours.",
		"Precipitation possibly heavy at times and ending within 12 hours.",
		"Precipitation ending in 12 to 24 hours.",
		"Precipitation possible within 6 hours.",
		"Precipitation possible and windy within 6 hours.",
		"",
		"Precipitation possible within 6 to 12 hours, possibly heavy at times.",
		"Precipitation possible within 6 to 12 hours.",
		"Precipitation possible within 12 hours.",
		"Precipitation possible within 12 hours, possibly heavy at times.",
		"",
		"",
		"Precipitation possible within 12 to 24 hours.",
		"Precipitation possible within 24 hours.",
		"Precipitation possible within 24 to 48 hours.",
		"Precipitation possible within 48 hours.",
		"",
		"Increasing winds.",
		"Windy.",
		"Possible wind shift to the W, NW, or N.",
		"Windy with possible wind shift to the W, NW, or N.",
		"Partly cloudy and cooler.",
		"Partly cloudy with little temperature change.",
		"Possible wind shift to the W, SW, or S. ",
		"Windy with possible wind shift to the W, SW, or S.",
	}

	r := GetUInt8(p, i)

	// If forecast rule is not within the bounds of our table then
	// return the dash value.
	if r < 0 || r >= len(rules) {
		return Dash
	}

	// Return all of the messages for the rule.
	msgs := make([]string, len(rules[r]))
	for j := 0; j < len(rules[r]); j++ {
		msgs[j] = msg[rules[r][j]]
	}

	return strings.Join(msgs, " ")
}

// GetForecastIcons gets a forecast icon bit map from a given packet at
// the specified index.
func GetForecastIcons(p []byte, i uint) (icons []string) {
	var iconBits = []string{ // Bit
		"Rain",          // 0
		"Cloud",         // 1
		"Partly Cloudy", // 2
		"Sun",           // 3
		"Snow",          // 4
	}

	for j := 0; j < len(iconBits); j++ {
		if GetUInt8(p, i)&(1<<uint(j)) != 0 {
			icons = append(icons, iconBits[j])
		}
	}

	return
}

// GetMPH8 gets a 1-byte MPH value from a given packet at the
// specified index.
func GetMPH8(p []byte, i uint) int {
	return GetUInt8(p, i)
}

// GetMPH16 gets a 2-byte MPH value from a given packet at the
// specified index.
func GetMPH16(p []byte, i uint) float64 {
	return GetFloat16(p, i) / 10.0
}

// GetPressure gets a pressure value from a given packet at
// the specified index.
func GetPressure(p []byte, i uint) float64 {
	return GetFloat16(p, i) / 1000.0
}

// GetRain gets a rain rate or accumulation value from a given
// packet at the specified index.
func GetRain(p []byte, i uint) float64 {
	return GetFloat16(p, i) / 100.0
}

// GetTemp8 gets a 1-byte temperature value in a given packet
// at the specified index.
func GetTemp8(p []byte, i uint) int {
	return GetUInt8(p, i) - 90
}

// GetTime16 gets a 2-byte time (no date) value in a given packet
// at the specified index.
func GetTime16(p []byte, i uint) time.Time {
	// If uninitialized then return a zero Time.
	t := GetUInt16(p, i)
	if t == 0xffff {
		return time.Time{}
	}

	// The time is stored as: hour * 100 + min
	hour := t / 100
	minute := t % 100

	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, time.Local)
}

// GetTransStatus gets the transmitter status from the given packet
// at the specified index and returns a slice of the ID's/channels
// that have low battery indicators.
func GetTransStatus(p []byte, i uint) (low []int) {
	for j := uint(0); j < 8; j++ {
		if GetUInt8(p, i)&(1<<j) != 0 {
			low = append(low, int(j)+1)
		}
	}

	return
}

// GetUFloat8 gets a 1-byte unsigned float value from a given packet
// at the specified index.
func GetUFloat8(p []byte, i uint) float64 {
	return float64(p[i])
}

// GetUInt8 gets a 1-byte unsigned integer value from a given packet
// at the specified index.
func GetUInt8(p []byte, i uint) int {
	return int(p[i])
}

// GetUInt16 gets a 2-byte unsigned integer value from a given packet
// at the specified index.
func GetUInt16(p []byte, i uint) int {
	return int(p[i+1])<<8 | int(p[i])
}

// GetUVIndex gets a Ultraviolet index value from a given packet
// at the specified index.
func GetUVIndex(p []byte, i uint) float64 {
	return GetUFloat8(p, i) / 10.0
}

// GetVoltage gets a battery voltage value from a given packet
// at the specified index.
func GetVoltage(p []byte, i uint) float64 {
	return GetFloat16(p, i) * 300.0 / 512.0 / 100.0
}

// GetWindDir gets a wind direction value in degrees from a
// given packet at the specified index.
func GetWindDir(p []byte, i uint) int {
	c := GetUInt8(p, i)
	if c < 0 || c > 15 {
		return 0
	}

	return int(float64(c)*22.5 + 0.5)
}
