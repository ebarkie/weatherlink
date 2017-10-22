// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package packet

// Common binary packet decoding logic for packets.
//
// Refer to Vantage ProTM, Vantage Pro2TM and Vantage VueTM Serial
// Communication Reference Manual, section X. Data Formats.

import (
	"math"
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
	switch GetInt8(p, i) {
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
	d := GetInt16(p, i)
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
	d := GetInt16(p, i)
	day := d & 0x001f
	month := (d & 0x01e0) >> 5
	year := 2000 + (d&0xfe00)>>9

	// The time is stored in second two bytes stored as: hour * 100 + min
	t := GetInt16(p, i+2)
	hour := t / 100
	minute := t % 100

	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.Local)
}

// GetDateTime48 gets a 6-byte date and time value from a given packet
// at the specified index.
func GetDateTime48(p []byte, i uint) time.Time {
	second := GetInt8(p, i)
	minute := GetInt8(p, i+1)
	hour := GetInt8(p, i+2)
	day := GetInt8(p, i+3)
	month := GetInt8(p, i+4)
	year := 1900 + GetInt8(p, i+5)

	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
}

// GetFloat8 gets a 1-byte float value from a given packet at
// the specified index.
func GetFloat8(p []byte, i uint) float64 {
	return float64(p[i])
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
	var rules = []string{
		"Mostly clear and cooler.",
		"Mostly clear with little temperature change.",
		"Mostly clear for 12 hrs. with little temperature change.",
		"Mostly clear for 12 to 24 hrs. and cooler.",
		"Mostly clear with little temperature change.",
		"Partly cloudy and cooler.",
		"Partly cloudy with little temperature change.",
		"Partly cloudy with little temperature change.",
		"Mostly clear and warmer.",
		"Partly cloudy with little temperature change.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds and warmer. Precipitation possible within 24 to 48 hrs.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds with little temperature change. Precipitation possible within 24 hrs.",
		"Mostly clear with little temperature change.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds with little temperature change. Precipitation possible within 12 hrs.",
		"Mostly clear with little temperature change.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds and warmer. Precipitation possible within 24 hrs.",
		"Mostly clear and warmer. Increasing winds.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds and warmer. Precipitation possible within 12 hrs. Increasing winds.",
		"Mostly clear and warmer. Increasing winds.",
		"Increasing clouds and warmer.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds and warmer. Precipitation possible within 12 hrs. Increasing winds.",
		"Mostly clear and warmer. Increasing winds.",
		"Increasing clouds and warmer.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds and warmer. Precipitation possible within 12 hrs. Increasing winds.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Mostly clear and warmer. Precipitation possible within 48 hrs.",
		"Mostly clear and warmer.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds with little temperature change. Precipitation possible within 24 to 48 hrs.",
		"Increasing clouds with little temperature change.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds and warmer. Precipitation possible within 12 to 24 hrs.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds and warmer. Precipitation possible within 12 to 24 hrs. Windy.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds and warmer. Precipitation possible within 12 to 24 hrs. Windy.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds and warmer. Precipitation possible within 6 to 12 hrs.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds and warmer. Precipitation possible within 6 to 12 hrs. Windy.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds and warmer. Precipitation possible within 12 to 24 hrs. Windy.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds and warmer. Precipitation possible within 12 hrs.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds and warmer. Precipitation likely.",
		"clearing and cooler. Precipitation ending within 6 hrs.",
		"Partly cloudy with little temperature change.",
		"clearing and cooler. Precipitation ending within 6 hrs.",
		"Mostly clear with little temperature change.",
		"Clearing and cooler. Precipitation ending within 6 hrs.",
		"Partly cloudy and cooler.",
		"Partly cloudy with little temperature change.",
		"Mostly clear and cooler.",
		"clearing and cooler. Precipitation ending within 6 hrs.",
		"Mostly clear with little temperature change.",
		"Clearing and cooler. Precipitation ending within 6 hrs.",
		"Mostly clear and cooler.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds with little temperature change. Precipitation possible within 24 hrs.",
		"Mostly cloudy and cooler. Precipitation continuing.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Mostly cloudy and cooler. Precipitation likely.",
		"Mostly cloudy with little temperature change. Precipitation continuing.",
		"Mostly cloudy with little temperature change. Precipitation likely.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds and cooler. Precipitation possible and windy within 6 hrs.",
		"Increasing clouds with little temperature change. Precipitation possible and windy within 6 hrs.",
		"Mostly cloudy and cooler. Precipitation continuing. Increasing winds.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Mostly cloudy and cooler. Precipitation likely. Increasing winds.",
		"Mostly cloudy with little temperature change. Precipitation continuing. Increasing winds.",
		"Mostly cloudy with little temperature change. Precipitation likely. Increasing winds.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds and cooler. Precipitation possible within 12 to 24 hrs. Possible wind shift to the W, NW, or N.",
		"Increasing clouds with little temperature change. Precipitation possible within 12 to 24 hrs. Possible wind shift to the W, NW, or N.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds and cooler. Precipitation possible within 6 hrs. Possible wind shift to the W, NW, or N.",
		"Increasing clouds with little temperature change. Precipitation possible within 6 hrs. Possible wind shift to the W, NW, or N.",
		"Mostly cloudy and cooler. Precipitation ending within 12 hrs. Possible wind shift to the W, NW, or N.",
		"Mostly cloudy and cooler. Possible wind shift to the W, NW, or N.",
		"Mostly cloudy with little temperature change. Precipitation ending within 12 hrs. Possible wind shift to the W, NW, or N.",
		"Mostly cloudy with little temperature change. Possible wind shift to the W, NW, or N.",
		"Mostly cloudy and cooler. Precipitation ending within 12 hrs. Possible wind shift to the W, NW, or N.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Mostly cloudy and cooler. Precipitation possible within 24 hrs. Possible wind shift to the W, NW, or N.",
		"Mostly cloudy with little temperature change. Precipitation ending within 12 hrs. Possible wind shift to the W, NW, or N.",
		"Mostly cloudy with little temperature change. Precipitation possible within 24 hrs. Possible wind shift to the W, NW, or N.",
		"clearing, cooler and windy. Precipitation ending within 6 hrs.",
		"clearing, cooler and windy.",
		"Mostly cloudy and cooler. Precipitation ending within 6 hrs. Windy with possible wind shift to the W, NW, or N.",
		"Mostly cloudy and cooler. Windy with possible wind shift to the W, NW, or N.",
		"clearing, cooler and windy.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Mostly cloudy with little temperature change. Precipitation possible within 12 hrs. Windy.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds and cooler. Precipitation possible within 12 hrs., possibly heavy at times. Windy.",
		"Mostly cloudy and cooler. Precipitation ending within 6 hrs. Windy.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Mostly cloudy and cooler. Precipitation possible within 12 hrs. Windy.",
		"Mostly cloudy and cooler. Precipitation ending in 12 to 24 hrs.",
		"Mostly cloudy and cooler.",
		"Mostly cloudy and cooler. Precipitation continuing, possible heavy at times. Windy.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Mostly cloudy and cooler. Precipitation possible within 6 to 12 hrs. Windy.",
		"Mostly cloudy with little temperature change. Precipitation continuing, possibly heavy at times. Windy.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Mostly cloudy with little temperature change. Precipitation possible within 6 to 12 hrs. Windy.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds with little temperature change. Precipitation possible within 12 hrs., possibly heavy at times. Windy.",
		"Mostly cloudy and cooler. Windy.",
		"Mostly cloudy and cooler. Precipitation continuing, possibly heavy at times. Windy.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Mostly cloudy and cooler. Precipitation likely, possibly heavy at times. Windy.",
		"Mostly cloudy with little temperature change. Precipitation continuing, possibly heavy at times. Windy.",
		"Mostly cloudy with little temperature change. Precipitation likely, possibly heavy at times. Windy.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds and cooler. Precipitation possible within 6 hrs. Windy.",
		"Increasing clouds with little temperature change. Precipitation possible within 6 hrs. windy",
		"Increasing clouds and cooler. Precipitation continuing. Windy with possible wind shift to the W, NW, or N.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Mostly cloudy and cooler. Precipitation likely. Windy with possible wind shift to the W, NW, or N.",
		"Mostly cloudy with little temperature change. Precipitation continuing. Windy with possible wind shift to the W, NW, or N.",
		"Mostly cloudy with little temperature change. Precipitation likely. Windy with possible wind shift to the W, NW, or N.",
		"Increasing clouds and cooler. Precipitation possible within 6 hrs. Windy with possible wind shift to the W, NW, or N.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds and cooler. Precipitation possible within 6 hrs. Possible wind shift to the W, NW, or N.",
		"Increasing clouds with little temperature change. Precipitation possible within 6 hrs. Windy with possible wind shift to the W, NW, or N.",
		"Increasing clouds with little temperature change. Precipitation possible within 6 hrs. Possible wind shift to the W, NW, or N.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds and cooler. Precipitation possible within 6 hrs. Windy with possible wind shift to the W, NW, or N.",
		"Increasing clouds with little temperature change. Precipitation possible within 6 hrs. Windy with possible wind shift to the W, NW, or N.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Increasing clouds and cooler. Precipitation possible within 12 to 24 hrs. Windy with possible wind shift to the W, NW, or N.",
		"Increasing clouds with little temperature change. Precipitation possible within 12 to 24 hrs. Windy with possible wind shift to the W, NW, or N.",
		"Mostly cloudy and cooler. Precipitation possibly heavy at times and ending within 12 hrs. Windy with possible wind shift to the W, NW, or N.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Mostly cloudy and cooler. Precipitation possible within 6 to 12 hrs., possibly heavy at times. Windy with possible wind shift to the W, NW, or N.",
		"Mostly cloudy with little temperature change. Precipitation ending within 12 hrs. Windy with possible wind shift to the W, NW, or N.",
		"Mostly cloudy with little temperature change. Precipitation possible within 6 to 12 hrs., possibly heavy at times. Windy with possible wind shift to the W, NW, or N.",
		"Mostly cloudy and cooler. Precipitation continuing.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Mostly cloudy and cooler. Precipitation likely, windy with possible wind shift to the W, NW, or N.",
		"Mostly cloudy with little temperature change. Precipitation continuing.",
		"Mostly cloudy with little temperature change. Precipitation likely.",
		"Partly cloudy with little temperature change.",
		"Mostly clear with little temperature change.",
		"Mostly cloudy and cooler. Precipitation possible within 12 hours, possibly heavy at times. Windy.",
		"FORECAST REQUIRES 3 HOURS OF RECENT DATA",
		"Mostly clear and cooler.",
		"Mostly clear and cooler.",
		"Mostly clear and cooler.",
	}

	r := GetInt8(p, i)
	if r > 0 && r <= len(rules) {
		return rules[r]
	}

	return Dash
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
		if GetInt8(p, i)&int(math.Pow(2, float64(j))) != 0 {
			icons = append(icons, iconBits[j])
		}
	}

	return
}

// GetInt8 gets a 1-byte integer value from a given packet at
// the specified index.
func GetInt8(p []byte, i uint) int {
	return int(p[i])
}

// GetInt16 gets a 2-byte integer value from a given packet at
// the specified index.
func GetInt16(p []byte, i uint) int {
	return int(p[i+1])<<8 | int(p[i])
}

// GetMPH8 gets a 1-byte MPH value from a given packet at the
// specified index.
func GetMPH8(p []byte, i uint) int {
	return GetInt8(p, i)
}

// GetMPH16 gets a 2-byte MPH value from a given packet at the
// specified index.
func GetMPH16(p []byte, i uint) float64 {
	return GetFloat16(p, i) / 10.0
}

// GetTemp8 gets a 1-byte temperature value in a given packet
// at the specified index.
func GetTemp8(p []byte, i uint) int {
	return GetInt8(p, i) - 90
}

// GetTime16 gets a 2-byte time (no date) value in a given packet
// at the specified index.
func GetTime16(p []byte, i uint) time.Time {
	// If uninitialized then return a zero Time.
	t := GetInt16(p, i)
	if t == 0xffff {
		return time.Time{}
	}

	// The time is stored as: hour * 100 + min
	hour := t / 100
	minute := t % 100

	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, time.Local)
}

// GetPressure gets a pressure value from a given packet at
// the specified index.
func GetPressure(p []byte, i uint) float64 {
	return GetFloat16(p, i) / 1000.0
}

// GetRainClicks gets a rain rate or accumulation value from
// a given packet at the specified index.
func GetRainClicks(p []byte, i uint) float64 {
	return GetFloat16(p, i) / 100.0
}

// GetUVIndex gets a Ultraviolet index value from a given packet
// at the specified index.
func GetUVIndex(p []byte, i uint) float64 {
	return GetFloat8(p, i) / 10.0
}

// GetVoltage gets a battery voltage value from a given packet
// at the specified index.
func GetVoltage(p []byte, i uint) float64 {
	return GetFloat16(p, i) * 300.0 / 512.0 / 100.0
}

// GetWindDir gets a wind direction value in degrees from a
// given packet at the specified index.
func GetWindDir(p []byte, i uint) int {
	c := GetInt8(p, i)
	if c < 0 || c > 15 {
		return 0
	}

	return int(float64(c)*22.5 + 0.5)
}
