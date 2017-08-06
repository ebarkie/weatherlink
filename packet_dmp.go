// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

// Packet coding logic for DMP revision B packets. Unlike loop
// packets where multiple versions are still in use, DMP switched
// from revision A to B in April 2002.  Since it is 2016 when
// implementing this it made no sense to support revision A.
//
// Refer to Vantage ProTM, Vantage Pro2TM and Vantage VueTM Serial
// Communication Reference Manual, section X. Data Formats,
// subsection 4. DMP and DMPAFT data format.

import "time"

// Dmp is a revision B DMP archive page consisting of 5 archive
// records.
type Dmp [5]Archive

// Archive represents all of the data in a revision B DMP archive
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

// FromPacket unpacks the data from a 267-byte DMP revision B archive page
// packet into a Dmp array of 5 Archive records.
func (d *Dmp) FromPacket(p Packet) error {
	if crc(p) != 0 {
		return ErrBadCRC
	}

	// Each individual record within the archive page contains a
	// revision marker but they're all going to be the same so
	// it's only necessary to check the first one.
	if p.getDmpType() != "b" {
		return ErrNotDmpB
	}

	// Break apart the page of 5 52-byte archive records and process
	// each one.  There are 4 unused bytes at the end.
	for i := 0; i < 5; i++ {
		offset := (52 * i) + 1
		pr := p[offset : offset+52]

		d[i].Bar = pr.getPressure(14)
		d[i].ET = float64(pr[29]) / 1000
		// There are 2 extra humidity sensors and 3 extra temperature
		// sensors.  Usually they match but not for archive.
		for j := uint(0); j < 2; j++ {
			if v := pr.get1ByteInt(43 + j); v != 255 {
				d[i].ExtraHumidity[j] = &v
			}
		}
		for j := uint(0); j < 3; j++ {
			if v := pr.get1ByteTemp(45 + j); v != 165 {
				d[i].ExtraTemp[j] = &v
			}
		}
		d[i].Forecast = pr.getForecast(33)
		d[i].InHumidity = pr.get1ByteInt(22)
		d[i].InTemp = pr.get2ByteTemp10(20)
		for j := uint(0); j < 2; j++ {
			if v := pr.get1ByteTemp(34 + j); v != 165 {
				d[i].LeafTemp[j] = &v
			}
			if v := pr.get1ByteInt(36 + j); v != 255 {
				d[i].LeafWetness[j] = &v
			}
		}
		d[i].OutHumidity = pr.get1ByteInt(23)
		d[i].OutTemp = pr.get2ByteTemp10(4)
		d[i].OutTempHi = pr.get2ByteTemp10(6)
		d[i].OutTempLow = pr.get2ByteTemp10(8)
		d[i].RainAccum = pr.getRainClicks(10)
		d[i].RainRateHi = pr.getRainClicks(12)
		for j := uint(0); j < 4; j++ {
			if v := pr.get1ByteInt(48 + j); v != 255 {
				d[i].SoilMoist[j] = &v
			}
			if v := pr.get1ByteTemp(38 + j); v != 165 {
				d[i].SoilTemp[j] = &v
			}
		}
		d[i].SolarRad = pr.get2ByteInt(16)
		d[i].SolarRadHi = pr.get2ByteInt(30)
		d[i].Timestamp = pr.get4ByteDateTime(0)
		d[i].UVIndexAvg = pr.getUVIndex(28)
		d[i].UVIndexHi = pr.getUVIndex(32)
		d[i].WindDirHi = pr.getWindDir(26)
		d[i].WindDirPrevail = pr.getWindDir(27)
		d[i].WindSamples = pr.get2ByteInt(18)
		d[i].WindSpeedAvg = pr.get1ByteMPH(24)
		d[i].WindSpeedHi = pr.get1ByteMPH(25)
	}

	return nil
}

// Refer to Vantage ProTM, Vantage Pro2TM and Vantage VueTM Serial
// Communication Reference Manual, section XI. Download Protocol.

// DmpAft is a timestamp appropriate for the "DMP after" command.
type DmpAft time.Time

// ToPacket packs the data from the DmpAft struct into a 6-byte packet
// appropriate for use with the DMPAFT command.
func (da DmpAft) ToPacket() (p Packet) {
	// 4-bytes for the time and 2-bytes for the CRC.
	p = make(Packet, 6)
	p.setDateTimeSmall(time.Time(da))
	p.setCrc()

	return
}

// DmpMeta is the DMP metadata sent after the DMPAFT command is issued.  It
// informs the downloader how much data to expect and where the first record
// is within the first page.
type DmpMeta struct {
	Pages           int // Number of pages to download
	FirstPageOffset int // Offset of the first record to read within the first page
}

// FromPacket unpacks the data from a 6-byte DMP metadata packet.
func (dm *DmpMeta) FromPacket(p Packet) (err error) {
	if crc(p) != 0 {
		err = ErrBadCRC
		return
	}

	if len(p) != 6 {
		err = ErrNotDmp
		return
	}

	dm.Pages = p.get2ByteInt(0)
	dm.FirstPageOffset = p.get2ByteInt(2)

	return
}

// getDmpType returns the Dmp packet revision or empty string
// if it's not a valid archive packet.
func (p Packet) getDmpType() (t string) {
	switch p[43] {
	case 0xff:
		t = "a"
	case 0x00:
		t = "b"
	}

	return
}
