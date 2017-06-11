// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

// Common binary packet decoding logic for packets.
//
// Refer to Vantage ProTM, Vantage Pro2TM and Vantage VueTM Serial
// Communication Reference Manual, section X. Data Formats.

import (
	"fmt"
	"math"
	"time"
)

func (p Packet) get1ByteInt(i uint) int {
	return int(p[i])
}

func (p Packet) get2ByteFloat(i uint) float64 {
	return float64(p.get2ByteInt(i))
}

func (p Packet) get2ByteInt(i uint) int {
	return int(p[i+1])<<8 | int(p[i])
}

// get1ByteTemp gets a 1-byte integer temprature like
// extra sensors.
func (p Packet) get1ByteTemp(i uint) int {
	return p.get1ByteInt(i) - 90
}

// get2ByteTemp gets a 2-byte float temperature like
// inside/outside sensors.
func (p Packet) get2ByteTemp(i uint) float64 {
	// Decode signed two's complement.
	return float64(int16(uint16(p[i+1])<<8|uint16(p[i]))) / 10
}

// get2ByteDate gets a 2-byte date (no time) like rain storm
// start date.
func (p Packet) get2ByteDate(i uint) (t time.Time) {
	// Date is stored in the two bytes as:
	//
	//  MMMM DDDD DYYY YYYY
	// 15       8         0
	year := 2000 + p.get2ByteInt(i)&0x007f
	day := (p.get2ByteInt(i) & 0x0f80) >> 7
	month := (p.get2ByteInt(i) & 0xf000) >> 12

	// If no date has been set (e.g. no storm in progress),  we'll get
	// back 0xffff which decodes to 2127-15-31.  In this circumstance
	// return 1980-01-01 since it's a cool throwback to DOS.
	if month == 15 {
		//t = time.Unix(315550800, 0)
	} else {
		t, _ = time.ParseInLocation("2006-01-02",
			fmt.Sprintf("%04d-%02d-%02d", year, month, day),
			time.Now().Location())
	}

	return
}

// get2ByteTime gets a 2-byte time (no date) like sunrise and sunset.
// The date will be set to today.
func (p Packet) get2ByteTime(i uint) time.Time {
	// Time is stored as: hour * 100 + min
	hour := p.get2ByteInt(i) / 100
	minute := p.get2ByteInt(i) % 100

	t, _ := time.ParseInLocation("2006-01-02 15:04",
		fmt.Sprintf("%04d-%02d-%02d %02d:%02d",
			time.Now().Year(),
			time.Now().Month(),
			time.Now().Day(),
			hour,
			minute),
		time.Now().Location())
	return t
}

// get4ByteDateTime gets a 4-byte date and time like in archive
// records.
func (p Packet) get4ByteDateTime(i uint) (t time.Time) {
	// Date is stored in the first two bytes as:
	//
	//  YYYY YYYM MMMD DDDD
	// 15       8         0
	day := p.get2ByteInt(i) & 0x001f
	month := (p.get2ByteInt(i) & 0x01e0) >> 5
	year := 2000 + (p.get2ByteInt(i)&0xfe00)>>9

	// Time is stored in second two bytes stored as: hour * 100 + min
	hour := p.get2ByteInt(i+2) / 100
	minute := p.get2ByteInt(i+2) % 100

	t, _ = time.ParseInLocation("2006-01-02T15:04",
		fmt.Sprintf("%04d-%02d-%02dT%02d:%02d", year, month, day, hour, minute),
		time.Now().Location())

	return
}

// get6ByteDateTime gets a 6-byte date and time like the console.
func (p Packet) get6ByteDateTime(i uint) (t time.Time) {
	second := p.get1ByteInt(i + 0)
	minute := p.get1ByteInt(i + 1)
	hour := p.get1ByteInt(i + 2)
	day := p.get1ByteInt(i + 3)
	month := p.get1ByteInt(i + 4)
	year := 1900 + p.get1ByteInt(i+5)

	t, _ = time.ParseInLocation("2006-01-02T15:04:05",
		fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d", year, month, day, hour, minute, second),
		time.Now().Location())

	return
}

// get1ByteMPH gets a 1-byte MPH which is all speed values except for
// the 2 and 10 minute values in a loop2 packet.
func (p Packet) get1ByteMPH(i uint) int {
	return p.get1ByteInt(i)
}

// get2ByteMPH gets a 2-byte MPH like 2 and 10 minute values.
func (p Packet) get2ByteMPH(i uint) float64 {
	return p.get2ByteFloat(i) / 10
}

// getBarTrend converts a barometer trend code to a string.
func (p Packet) getBarTrend(i uint) string {
	switch p.get1ByteInt(i) {
	case -60:
		return "Falling Rapidly"
	case -20:
		return "Falling Slowly"
	case 0:
		return "Steady"
	case 20:
		return "Rising Slowly"
	case 60:
		return "Rising Rapidly"
	default:
		return "-"
	}
}

// getForecast converts a forecast rule index to a string.
func (p Packet) getForecast(i uint) string {
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

	if (p.get1ByteInt(i) > 0) && (p.get1ByteInt(i) <= len(rules)) {
		return rules[p.get1ByteInt(i)]
	}

	return "-"
}

// getForecastIcons converts a forecast icon bit map to a slice of icon
// names.
func (p Packet) getForecastIcons(i uint) (icons []string) {
	var iconBits = []string{ // Bit
		"Rain",          // 0
		"Cloud",         // 1
		"Partly Cloudy", // 2
		"Sun",           // 3
		"Snow",          // 4
	}

	for j := 0; j < len(iconBits); j++ {
		if p.get1ByteInt(i)&int(math.Pow(2, float64(j))) != 0 {
			icons = append(icons, iconBits[j])
		}
	}

	return
}

func (p Packet) getRainClicks(i uint) float64 {
	return p.get2ByteFloat(i) / 100
}

func (p Packet) getPressure(i uint) float64 {
	return p.get2ByteFloat(i) / 1000
}

func (p Packet) getUVIndex(i uint) float64 {
	return float64(p[i]) / 10
}

func (p Packet) getVoltage(i uint) float64 {
	return p.get2ByteFloat(i) * 300 / 512 / 100.0
}

// getWindDir converts an archive record wind direction code
// to a direction in degrees.
func (p Packet) getWindDir(i uint) int {
	c := float64(p.get1ByteInt(i))
	if (c < 0) || (c > 15) {
		return 0
	}

	return int((c * 22.5) + 0.5)
}
