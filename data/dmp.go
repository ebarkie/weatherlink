// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package data

// Packet coding logic for DMP revision B packets. Unlike loop
// packets where multiple versions are still in use, DMP switched
// from revision A to B in April 2002.  Since it is 2016 when
// implementing this it made no sense to support revision A.
//
// Refer to Vantage ProTM, Vantage Pro2TM and Vantage VueTM Serial
// Communication Reference Manual, section X. Data Formats,
// subsection 4. DMP and DMPAFT data format.

import (
	"time"

	"github.com/ebarkie/weatherlink/packet"
)

// Archive represents all of the data in a revision B archive
// record.
type Archive struct {
	Bar            float64   `json:"barometer"`
	ET             float64   `json:"ET"`
	ExtraHumidity  [2]*int   `json:"extraHumidity,omitempty"`
	ExtraTemp      [3]*int   `json:"extraTemperature,omitempty"`
	Forecast       string    `json:"forecast"`
	InHumidity     int       `json:"insideHumidity"`
	InTemp         float64   `json:"insideTemperature"`
	LeafTemp       [2]*int   `json:"leafTemperature,omitempty"`
	LeafWetness    [2]*int   `json:"leafWetness,omitempty"`
	OutHumidity    int       `json:"outsideHumidity"`
	OutTemp        float64   `json:"outsideTemperature"`
	OutTempHi      float64   `json:"outsideTemperatureHigh"`
	OutTempLow     float64   `json:"outsideTemperatureLow"`
	RainAccum      float64   `json:"rainAccumulation"`
	RainRateHi     float64   `json:"rainRateHigh"`
	SoilMoist      [4]*int   `json:"soilMoisture,omitempty"`
	SoilTemp       [4]*int   `json:"soilTemperature,omitempty"`
	SolarRad       int       `json:"solarRadiation"`
	SolarRadHi     int       `json:"solarRadiationHigh"`
	Timestamp      time.Time `json:"timestamp"`
	UVIndexAvg     float64   `json:"UVIndexAverage"`
	UVIndexHi      float64   `json:"UVIndexHigh"`
	WindDirHi      int       `json:"windDirectionHigh"`
	WindDirPrevail int       `json:"windDirectionPrevailing"`
	WindSamples    int       `json:"windSamples"`
	WindSpeedAvg   int       `json:"windSpeedAverage"`
	WindSpeedHi    int       `json:"windSpeedHigh"`
}

// UnmarshalBinary decodes a 52-byte revision B archive record.
func (a *Archive) UnmarshalBinary(p []byte) error {
	if getArcType(p) != "b" {
		return ErrNotArchive
	}

	a.Bar = packet.GetPressure(p, 14)
	a.ET = packet.GetFloat8(p, 29) / 1000
	// There are 2 extra humidity sensors and 3 extra temperature
	// sensors.  Usually the quantities match but not for archive
	// records.
	for j := uint(0); j < 2; j++ {
		if v := packet.GetInt8(p, 43+j); v != 255 {
			a.ExtraHumidity[j] = &v
		}
	}
	for j := uint(0); j < 3; j++ {
		if v := packet.GetTemp8(p, 45+j); v != 165 {
			a.ExtraTemp[j] = &v
		}
	}
	a.Forecast = packet.GetForecast(p, 33)
	a.InHumidity = packet.GetInt8(p, 22)
	a.InTemp = packet.GetFloat16_10(p, 20)
	for j := uint(0); j < 2; j++ {
		if v := packet.GetTemp8(p, 34+j); v != 165 {
			a.LeafTemp[j] = &v
		}
		if v := packet.GetInt8(p, 36+j); v != 255 {
			a.LeafWetness[j] = &v
		}
	}
	a.OutHumidity = packet.GetInt8(p, 23)
	a.OutTemp = packet.GetFloat16_10(p, 4)
	a.OutTempHi = packet.GetFloat16_10(p, 6)
	a.OutTempLow = packet.GetFloat16_10(p, 8)
	a.RainAccum = packet.GetRainClicks(p, 10)
	a.RainRateHi = packet.GetRainClicks(p, 12)
	for j := uint(0); j < 4; j++ {
		if v := packet.GetInt8(p, 48+j); v != 255 {
			a.SoilMoist[j] = &v
		}
		if v := packet.GetTemp8(p, 38+j); v != 165 {
			a.SoilTemp[j] = &v
		}
	}
	a.SolarRad = packet.GetInt16(p, 16)
	a.SolarRadHi = packet.GetInt16(p, 30)
	a.Timestamp = packet.GetDateTime32(p, 0)
	a.UVIndexAvg = packet.GetUVIndex(p, 28)
	a.UVIndexHi = packet.GetUVIndex(p, 32)
	a.WindDirHi = packet.GetWindDir(p, 26)
	a.WindDirPrevail = packet.GetWindDir(p, 27)
	a.WindSamples = packet.GetInt16(p, 18)
	a.WindSpeedAvg = packet.GetMPH8(p, 24)
	a.WindSpeedHi = packet.GetMPH8(p, 25)

	return nil
}

// Dmp is a download memory page which contains 5 archive
// records.
type Dmp [5]Archive

// UnmarshalBinary decodes a 267-byte download memory page into an
// array of 5 Archive records.
func (d *Dmp) UnmarshalBinary(p []byte) error {
	if packet.Crc(p) != 0 {
		return ErrBadCRC
	} else if len(p) != 267 {
		return ErrNotDmp
	}

	// Break apart the page of 5 52-byte archive records and process
	// each one.  There are 4 unused bytes at the end.
	for i := 0; i < 5; i++ {
		offset := 1 + (52 * i)
		err := d[i].UnmarshalBinary(p[offset : offset+52])
		if err != nil {
			return err
		}
	}

	return nil
}

// Refer to Vantage ProTM, Vantage Pro2TM and Vantage VueTM Serial
// Communication Reference Manual, section XI. Download Protocol.

// DmpAft is a timestamp appropriate for the "DMP after" command.
type DmpAft time.Time

// MarshalBinary encodes the data from the DmpAft struct into a 6-byte packet
// appropriate for use with the DMPAFT command.
func (da DmpAft) MarshalBinary() (p []byte, err error) {
	// 4-bytes for the time and 2-bytes for the CRC.
	p = make([]byte, 6)
	packet.SetDateTime32(&p, 0, time.Time(da))
	packet.SetCrc(&p)

	return
}

// DmpMeta is the DMP metadata sent after the DMPAFT command is issued.  It
// informs the downloader how much data to expect and where the first record
// is within the first page.
type DmpMeta struct {
	Pages           int // Number of pages to download
	FirstPageOffset int // Offset of the first record to read within the first page
}

// UnmarshalBinary decodes a 6-byte DMP metadata packet into the
// DmpMeta stuct.
func (dm *DmpMeta) UnmarshalBinary(p []byte) error {
	if packet.Crc(p) != 0 {
		return ErrBadCRC
	}

	if len(p) != 6 {
		return ErrNotDmpMeta
	}

	dm.Pages = packet.GetInt16(p, 0)
	dm.FirstPageOffset = packet.GetInt16(p, 2)

	return nil
}

// getArcType returns the Dmp packet revision or empty string
// if it's not a valid archive packet.
func getArcType(p []byte) (t string) {
	if len(p) != 52 {
		return
	}

	switch p[42] {
	case 0xff:
		t = "a"
	case 0x0:
		t = "b"
	}

	return
}
