// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package data

// Packet coding logic for HILOWS packets.
//
// Refer to Vantage ProTM, Vantage Pro2TM and Vantage VueTM Serial
// Communication Reference Manual, section X. Data Formats,
// subsection 3. HILOW data format.

import (
	"time"

	"github.com/ebarkie/weatherlink/packet"
)

// HiLows represents all of the record high and lows by day, month, and
// year.  The day also includes the time(s) when the record occurred.
type HiLows struct {
	Bar           HiLowBar             `json:"barometer"`
	DewPoint      HiLowTemp            `json:"dewPoint"`
	ExtraHumidity [7]*HiLowHumidity    `json:"extraHumidity,omitempty"`
	ExtraTemp     [7]*HiLowExtraTemp   `json:"extraTemperature,omitempty"`
	HeatIndex     HiHeatIndex          `json:"heatIndex"`
	InHumidity    HiLowHumidity        `json:"insideHumidity"`
	InTemp        HiLowTemp            `json:"insideTemperature"`
	LeafTemp      [4]*HiLowExtraTemp   `json:"leafTemperature,omitempty"`
	LeafWetness   [4]*HiLowLeafWetness `json:"leafWetness,omitempty"`
	OutHumidity   HiLowHumidity        `json:"outsideHumidity"`
	OutTemp       HiLowTemp            `json:"outsideTemperature"`
	RainRate      HiRainRate           `json:"rainRate"`
	SoilMoist     [4]*HiLowSoilMoist   `json:"soilMoisture,omitempty"`
	SoilTemp      [4]*HiLowExtraTemp   `json:"soilTemperature,omitempty"`
	SolarRad      HiSolarRad           `json:"solarRadiation"`
	THSWIndex     HiTHSWIndex          `json:"THSWIndex"`
	UVIndex       HiUVIndex            `json:"UVIndex"`
	WindSpeed     HiWindSpeed          `json:"windSpeed"`
	WindChill     LowWindChill         `json:"windChill"`
}

// HiLowBar is the record high and low barometer readings.
type HiLowBar struct {
	Day struct {
		Hi      float64   `json:"hi"`
		HiTime  time.Time `json:"hiTime,omitempty"`
		Low     float64   `json:"low"`
		LowTime time.Time `json:"lowTime,omitempty"`
	} `json:"day"`
	Month struct {
		Hi  float64 `json:"hi"`
		Low float64 `json:"low"`
	} `json:"month"`
	Year struct {
		Hi  float64 `json:"hi"`
		Low float64 `json:"low"`
	} `json:"year"`
}

// HiLowExtraTemp is the record high and low extra temperature readings.
type HiLowExtraTemp struct {
	Day struct {
		Hi      int       `json:"hi"`
		HiTime  time.Time `json:"hiTime,omitempty"`
		Low     int       `json:"low"`
		LowTime time.Time `json:"lowTime,omitempty"`
	} `json:"day"`
	Month struct {
		Hi  int `json:"hi"`
		Low int `json:"low"`
	} `json:"month"`
	Year struct {
		Hi  int `json:"hi"`
		Low int `json:"low"`
	} `json:"year"`
}

// HiHeatIndex is the record high heat index readings.
type HiHeatIndex struct {
	Day struct {
		Hi     float64   `json:"hi"`
		HiTime time.Time `json:"hiTime,omitempty"`
	} `json:"day"`
	Month struct {
		Hi float64 `json:"hi"`
	} `json:"month"`
	Year struct {
		Hi float64 `json:"hi"`
	} `json:"year"`
}

// HiLowHumidity is the record high and low humidity readings.
type HiLowHumidity struct {
	Day struct {
		Hi      int       `json:"hi"`
		HiTime  time.Time `json:"hiTime,omitempty"`
		Low     int       `json:"low"`
		LowTime time.Time `json:"lowTime,omitempty"`
	} `json:"day"`
	Month struct {
		Hi  int `json:"hi"`
		Low int `json:"low"`
	} `json:"month"`
	Year struct {
		Hi  int `json:"hi"`
		Low int `json:"low"`
	} `json:"year"`
}

// HiLowLeafWetness is the record high and low leaf wetness readings.
type HiLowLeafWetness struct {
	Day struct {
		Hi      int       `json:"hi"`
		HiTime  time.Time `json:"hiTime,omitempty"`
		Low     int       `json:"low"`
		LowTime time.Time `json:"lowTime,omitempty"`
	} `json:"day"`
	Month struct {
		Hi  int `json:"hi"`
		Low int `json:"low"`
	} `json:"month"`
	Year struct {
		Hi  int `json:"hi"`
		Low int `json:"low"`
	} `json:"year"`
}

// HiLowTemp is the record high and low temperature readings and dew point
// calculations.
type HiLowTemp struct {
	Day struct {
		Hi      float64   `json:"hi"`
		HiTime  time.Time `json:"hiTime,omitempty"`
		Low     float64   `json:"low"`
		LowTime time.Time `json:"lowTime,omitempty"`
	} `json:"day"`
	Month struct {
		Hi  float64 `json:"hi"`
		Low float64 `json:"low"`
	} `json:"month"`
	Year struct {
		Hi  float64 `json:"hi"`
		Low float64 `json:"low"`
	} `json:"year"`
}

// HiRainRate is the record high rain rate readings.
type HiRainRate struct {
	Hour struct {
		Hi float64 `json:"hi"`
	} `json:"hour"`
	Day struct {
		Hi     float64   `json:"hi"`
		HiTime time.Time `json:"hiTime,omitempty"`
	} `json:"day"`
	Month struct {
		Hi float64 `json:"hi"`
	} `json:"month"`
	Year struct {
		Hi float64 `json:"hi"`
	} `json:"year"`
}

// HiLowSoilMoist is the record high and low soil moisture readings.
type HiLowSoilMoist struct {
	Day struct {
		Hi      int       `json:"hi"`
		HiTime  time.Time `json:"hiTime,omitempty"`
		Low     int       `json:"low"`
		LowTime time.Time `json:"lowTime,omitempty"`
	} `json:"day"`
	Month struct {
		Hi  int `json:"hi"`
		Low int `json:"low"`
	} `json:"month"`
	Year struct {
		Hi  int `json:"hi"`
		Low int `json:"low"`
	} `json:"year"`
}

// HiSolarRad is the record high solar radiation readings.
type HiSolarRad struct {
	Day struct {
		Hi     int       `json:"hi"`
		HiTime time.Time `json:"hiTime,omitempty"`
	} `json:"day"`
	Month struct {
		Hi int `json:"hi"`
	} `json:"month"`
	Year struct {
		Hi int `json:"hi"`
	} `json:"year"`
}

// HiTHSWIndex is the record high THSW index calculations.
type HiTHSWIndex struct {
	Day struct {
		Hi     float64   `json:"hi"`
		HiTime time.Time `json:"hiTime,omitempty"`
	} `json:"day"`
	Month struct {
		Hi float64 `json:"hi"`
	} `json:"month"`
	Year struct {
		Hi float64 `json:"hi"`
	} `json:"year"`
}

// HiUVIndex is the record high UltraViolet index readings.
type HiUVIndex struct {
	Day struct {
		Hi     float64   `json:"hi"`
		HiTime time.Time `json:"hiTime,omitempty"`
	} `json:"day"`
	Month struct {
		Hi float64 `json:"hi"`
	} `json:"month"`
	Year struct {
		Hi float64 `json:"hi"`
	} `json:"year"`
}

// HiWindSpeed is the record high wind speed readings.
type HiWindSpeed struct {
	Day struct {
		Hi     int       `json:"hi"`
		HiTime time.Time `json:"hiTime,omitempty"`
	} `json:"day"`
	Month struct {
		Hi int `json:"hi"`
	} `json:"month"`
	Year struct {
		Hi int `json:"hi"`
	} `json:"year"`
}

// LowWindChill is the record low wind chill calculations.
type LowWindChill struct {
	Day struct {
		Low     float64   `json:"low"`
		LowTime time.Time `json:"lowTime,omitempty"`
	} `json:"day"`
	Month struct {
		Low float64 `json:"low"`
	} `json:"month"`
	Year struct {
		Low float64 `json:"low"`
	} `json:"year"`
}

// UnmarshalBinary decodes a 438-byte high and lows packet into the
// HiLows struct.
func (hl *HiLows) UnmarshalBinary(p []byte) error {
	if packet.Crc(p) != 0 {
		return ErrBadCRC
	}

	// It would have been nice if the decoding was more universal
	// but the ordering of the fields is not consistent, making it
	// difficult.

	// Barometer
	hl.Bar.Day.Low = packet.GetPressure(p, 0)
	hl.Bar.Day.LowTime = packet.GetTime16(p, 12)
	hl.Bar.Day.Hi = packet.GetPressure(p, 2)
	hl.Bar.Day.HiTime = packet.GetTime16(p, 14)
	hl.Bar.Month.Low = packet.GetPressure(p, 4)
	hl.Bar.Month.Hi = packet.GetPressure(p, 6)
	hl.Bar.Year.Low = packet.GetPressure(p, 8)
	hl.Bar.Year.Hi = packet.GetPressure(p, 10)

	// Dew point
	hl.DewPoint.Day.Low = packet.GetFloat16(p, 63)
	hl.DewPoint.Day.LowTime = packet.GetTime16(p, 67)
	hl.DewPoint.Day.Hi = packet.GetFloat16(p, 65)
	hl.DewPoint.Day.HiTime = packet.GetTime16(p, 69)
	hl.DewPoint.Month.Low = packet.GetFloat16(p, 73)
	hl.DewPoint.Month.Hi = packet.GetFloat16(p, 71)
	hl.DewPoint.Year.Low = packet.GetFloat16(p, 77)
	hl.DewPoint.Year.Hi = packet.GetFloat16(p, 75)

	// Extra humidity and temperatures
	extraHumidity := func(p []byte, i uint) (h HiLowHumidity) {
		h.Day.Low = packet.GetUInt8(p, 276+i)
		h.Day.LowTime = packet.GetTime16(p, 292+i*2)
		h.Day.Hi = packet.GetUInt8(p, 284+i)
		h.Day.HiTime = packet.GetTime16(p, 308+i*2)
		h.Month.Low = packet.GetUInt8(p, 332+i)
		h.Month.Hi = packet.GetUInt8(p, 324+i)
		h.Year.Low = packet.GetUInt8(p, 348+i)
		h.Year.Hi = packet.GetUInt8(p, 340+i)

		return
	}
	extraTemp := func(p []byte, i uint) (et HiLowExtraTemp) {
		et.Day.Low = packet.GetTemp8(p, 126+i)
		et.Day.LowTime = packet.GetTime16(p, 156+i*2)
		et.Day.Hi = packet.GetTemp8(p, 141+i)
		et.Day.HiTime = packet.GetTime16(p, 186+i*2)
		et.Month.Low = packet.GetTemp8(p, 231+i)
		et.Month.Hi = packet.GetTemp8(p, 216+i)
		et.Year.Low = packet.GetTemp8(p, 261+i)
		et.Year.Hi = packet.GetTemp8(p, 246+i)

		return
	}
	for i := uint(0); i < 7; i++ {
		if eh := extraHumidity(p, 1+i); eh.Day.Low != 255 {
			hl.ExtraHumidity[i] = &eh
		}
		if et := extraTemp(p, i); et.Day.Low != 165 {
			hl.ExtraTemp[i] = &et
		}
	}

	// Heat index
	hl.HeatIndex.Day.Hi = packet.GetFloat16(p, 87)
	hl.HeatIndex.Day.HiTime = packet.GetTime16(p, 89)
	hl.HeatIndex.Month.Hi = packet.GetFloat16(p, 91)
	hl.HeatIndex.Year.Hi = packet.GetFloat16(p, 93)

	// Inside humidity
	hl.InHumidity.Day.Low = packet.GetUInt8(p, 38)
	hl.InHumidity.Day.LowTime = packet.GetTime16(p, 41)
	hl.InHumidity.Day.Hi = packet.GetUInt8(p, 37)
	hl.InHumidity.Day.HiTime = packet.GetTime16(p, 39)
	hl.InHumidity.Month.Low = packet.GetUInt8(p, 44)
	hl.InHumidity.Month.Hi = packet.GetUInt8(p, 43)
	hl.InHumidity.Year.Low = packet.GetUInt8(p, 46)
	hl.InHumidity.Year.Hi = packet.GetUInt8(p, 45)

	// Inside temperature
	hl.InTemp.Day.Low = packet.GetFloat16_10(p, 23)
	hl.InTemp.Day.LowTime = packet.GetTime16(p, 27)
	hl.InTemp.Day.Hi = packet.GetFloat16_10(p, 21)
	hl.InTemp.Day.HiTime = packet.GetTime16(p, 25)
	hl.InTemp.Month.Low = packet.GetFloat16_10(p, 29)
	hl.InTemp.Month.Hi = packet.GetFloat16_10(p, 31)
	hl.InTemp.Year.Low = packet.GetFloat16_10(p, 33)
	hl.InTemp.Year.Hi = packet.GetFloat16_10(p, 35)

	// Leaf temperature and wetness
	for i := uint(0); i < 4; i++ {
		if et := extraTemp(p, 11+i); et.Day.Low != 165 {
			hl.LeafTemp[i] = &et
		}

		if low := packet.GetUInt8(p, 408+i); low != 255 {
			lw := HiLowLeafWetness{}
			lw.Day.Low = low
			lw.Day.LowTime = packet.GetTime16(p, 412+i*2)
			lw.Day.Hi = packet.GetUInt8(p, 396+i)
			lw.Day.HiTime = packet.GetTime16(p, 400+i*2)
			lw.Month.Low = packet.GetUInt8(p, 420+i)
			lw.Month.Hi = packet.GetUInt8(p, 424+i)
			lw.Year.Low = packet.GetUInt8(p, 428+i)
			lw.Year.Hi = packet.GetUInt8(p, 432+i)
			hl.LeafWetness[i] = &lw
		}
	}

	// Outside humidity
	hl.OutHumidity = extraHumidity(p, 0)

	// Outside temperature
	hl.OutTemp.Day.Low = packet.GetFloat16_10(p, 47)
	hl.OutTemp.Day.LowTime = packet.GetTime16(p, 51)
	hl.OutTemp.Day.Hi = packet.GetFloat16_10(p, 49)
	hl.OutTemp.Day.HiTime = packet.GetTime16(p, 53)
	hl.OutTemp.Month.Low = packet.GetFloat16_10(p, 57)
	hl.OutTemp.Month.Hi = packet.GetFloat16_10(p, 55)
	hl.OutTemp.Year.Low = packet.GetFloat16_10(p, 61)
	hl.OutTemp.Year.Hi = packet.GetFloat16_10(p, 59)

	// Rain rate
	hl.RainRate.Hour.Hi = packet.GetRainClicks(p, 120)
	hl.RainRate.Day.Hi = packet.GetRainClicks(p, 116)
	hl.RainRate.Day.HiTime = packet.GetTime16(p, 118)
	hl.RainRate.Month.Hi = packet.GetRainClicks(p, 122)
	hl.RainRate.Year.Hi = packet.GetRainClicks(p, 124)

	// Soil moisture and temperature
	for i := uint(0); i < 4; i++ {
		if low := packet.GetUInt8(p, 368+i); low != 255 {
			sm := HiLowSoilMoist{}
			sm.Day.Low = low
			sm.Day.LowTime = packet.GetTime16(p, 372+i*2)
			sm.Day.Hi = packet.GetUInt8(p, 356+i)
			sm.Day.HiTime = packet.GetTime16(p, 360+i*2)
			sm.Month.Low = packet.GetUInt8(p, 380+i)
			sm.Month.Hi = packet.GetUInt8(p, 384+i)
			sm.Year.Low = packet.GetUInt8(p, 388+i)
			sm.Year.Hi = packet.GetUInt8(p, 392+i)
			hl.SoilMoist[i] = &sm
		}

		if et := extraTemp(p, 7+i); et.Day.Low != 165 {
			hl.SoilTemp[i] = &et
		}
	}

	// Solar radiation
	hl.SolarRad.Day.Hi = packet.GetUInt16(p, 103)
	hl.SolarRad.Day.HiTime = packet.GetTime16(p, 105)
	hl.SolarRad.Month.Hi = packet.GetUInt16(p, 107)
	hl.SolarRad.Year.Hi = packet.GetUInt16(p, 109)

	// THSW index
	hl.THSWIndex.Day.Hi = packet.GetFloat16(p, 95)
	hl.THSWIndex.Day.HiTime = packet.GetTime16(p, 97)
	hl.THSWIndex.Month.Hi = packet.GetFloat16(p, 99)
	hl.THSWIndex.Year.Hi = packet.GetFloat16(p, 101)

	// UltraViolet index
	hl.UVIndex.Day.Hi = packet.GetUVIndex(p, 111)
	hl.UVIndex.Day.HiTime = packet.GetTime16(p, 112)
	hl.UVIndex.Month.Hi = packet.GetUVIndex(p, 114)
	hl.UVIndex.Year.Hi = packet.GetUVIndex(p, 115)

	// Wind speed
	hl.WindSpeed.Day.Hi = packet.GetMPH8(p, 16)
	hl.WindSpeed.Day.HiTime = packet.GetTime16(p, 17)
	hl.WindSpeed.Month.Hi = packet.GetMPH8(p, 19)
	hl.WindSpeed.Year.Hi = packet.GetMPH8(p, 20)

	// Wind chill
	hl.WindChill.Day.Low = packet.GetFloat16(p, 79)
	hl.WindChill.Day.LowTime = packet.GetTime16(p, 81)
	hl.WindChill.Month.Low = packet.GetFloat16(p, 83)
	hl.WindChill.Year.Low = packet.GetFloat16(p, 85)

	return nil
}
