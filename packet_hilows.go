// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

// Packet coding logic for HILOWS packets.
//
// Refer to Vantage ProTM, Vantage Pro2TM and Vantage VueTM Serial
// Communication Reference Manual, section X. Data Formats,
// subsection 3. HILOW data format.

import "time"

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

// FromPacket unpacks a 438-byte high and lows packet into the
// HiLows struct.
func (hl *HiLows) FromPacket(p Packet) error {
	if crc(p) != 0 {
		return ErrBadCRC
	}

	// It would have been nice if the decoding was more universal
	// but the ordering of the fields is not consistent, making it
	// difficult.

	// Barometer
	hl.Bar.Day.Low = p.getPressure(0)
	hl.Bar.Day.LowTime = p.get2ByteTime(12)
	hl.Bar.Day.Hi = p.getPressure(2)
	hl.Bar.Day.HiTime = p.get2ByteTime(14)
	hl.Bar.Month.Low = p.getPressure(4)
	hl.Bar.Month.Hi = p.getPressure(6)
	hl.Bar.Year.Low = p.getPressure(8)
	hl.Bar.Year.Hi = p.getPressure(10)

	// Dew point
	hl.DewPoint.Day.Low = p.get2ByteFloat(63)
	hl.DewPoint.Day.LowTime = p.get2ByteTime(67)
	hl.DewPoint.Day.Hi = p.get2ByteFloat(65)
	hl.DewPoint.Day.HiTime = p.get2ByteTime(69)
	hl.DewPoint.Month.Low = p.get2ByteFloat(73)
	hl.DewPoint.Month.Hi = p.get2ByteFloat(71)
	hl.DewPoint.Year.Low = p.get2ByteFloat(77)
	hl.DewPoint.Year.Hi = p.get2ByteFloat(75)

	// Extra humidity and temperatures
	extraHumidity := func(i uint) (h HiLowHumidity) {
		h.Day.Low = p.get1ByteInt(276 + i)
		h.Day.LowTime = p.get2ByteTime(292 + i*2)
		h.Day.Hi = p.get1ByteInt(284 + i)
		h.Day.HiTime = p.get2ByteTime(308 + i*2)
		h.Month.Low = p.get1ByteInt(332 + i)
		h.Month.Hi = p.get1ByteInt(324 + i)
		h.Year.Low = p.get1ByteInt(348 + i)
		h.Year.Hi = p.get1ByteInt(340 + i)

		return
	}
	extraTemp := func(i uint) (et HiLowExtraTemp) {
		et.Day.Low = p.get1ByteTemp(126 + i)
		et.Day.LowTime = p.get2ByteTime(156 + i*2)
		et.Day.Hi = p.get1ByteTemp(141 + i)
		et.Day.HiTime = p.get2ByteTime(186 + i*2)
		et.Month.Low = p.get1ByteTemp(231 + i)
		et.Month.Hi = p.get1ByteTemp(216 + i)
		et.Year.Low = p.get1ByteTemp(261 + i)
		et.Year.Hi = p.get1ByteTemp(246 + i)

		return
	}
	for i := uint(0); i < 7; i++ {
		if eh := extraHumidity(1 + i); eh.Day.Low != 255 {
			hl.ExtraHumidity[i] = &eh
		}
		if et := extraTemp(i); et.Day.Low != 165 {
			hl.ExtraTemp[i] = &et
		}
	}

	// Heat index
	hl.HeatIndex.Day.Hi = p.get2ByteFloat(87)
	hl.HeatIndex.Day.HiTime = p.get2ByteTime(89)
	hl.HeatIndex.Month.Hi = p.get2ByteFloat(91)
	hl.HeatIndex.Year.Hi = p.get2ByteFloat(93)

	// Inside humidity
	hl.InHumidity.Day.Low = p.get1ByteInt(38)
	hl.InHumidity.Day.LowTime = p.get2ByteTime(41)
	hl.InHumidity.Day.Hi = p.get1ByteInt(37)
	hl.InHumidity.Day.HiTime = p.get2ByteTime(39)
	hl.InHumidity.Month.Low = p.get1ByteInt(44)
	hl.InHumidity.Month.Hi = p.get1ByteInt(43)
	hl.InHumidity.Year.Low = p.get1ByteInt(46)
	hl.InHumidity.Year.Hi = p.get1ByteInt(45)

	// Inside temperature
	hl.InTemp.Day.Low = p.get2ByteFloat10(23)
	hl.InTemp.Day.LowTime = p.get2ByteTime(27)
	hl.InTemp.Day.Hi = p.get2ByteFloat10(21)
	hl.InTemp.Day.HiTime = p.get2ByteTime(25)
	hl.InTemp.Month.Low = p.get2ByteFloat10(29)
	hl.InTemp.Month.Hi = p.get2ByteFloat10(31)
	hl.InTemp.Year.Low = p.get2ByteFloat10(33)
	hl.InTemp.Year.Hi = p.get2ByteFloat10(35)

	// Leaf temperature and wetness
	for i := uint(0); i < 4; i++ {
		if et := extraTemp(11 + i); et.Day.Low != 165 {
			hl.LeafTemp[i] = &et
		}

		if low := p.get1ByteInt(408 + i); low != 255 {
			lw := HiLowLeafWetness{}
			lw.Day.Low = low
			lw.Day.LowTime = p.get2ByteTime(412 + i*2)
			lw.Day.Hi = p.get1ByteInt(396 + i)
			lw.Day.HiTime = p.get2ByteTime(400 + i*2)
			lw.Month.Low = p.get1ByteInt(420 + i)
			lw.Month.Hi = p.get1ByteInt(424 + i)
			lw.Year.Low = p.get1ByteInt(428 + i)
			lw.Year.Hi = p.get1ByteInt(432 + i)
			hl.LeafWetness[i] = &lw
		}
	}

	// Outside humidity
	hl.OutHumidity = extraHumidity(0)

	// Outside temperature
	hl.OutTemp.Day.Low = p.get2ByteFloat10(47)
	hl.OutTemp.Day.LowTime = p.get2ByteTime(51)
	hl.OutTemp.Day.Hi = p.get2ByteFloat10(49)
	hl.OutTemp.Day.HiTime = p.get2ByteTime(53)
	hl.OutTemp.Month.Low = p.get2ByteFloat10(57)
	hl.OutTemp.Month.Hi = p.get2ByteFloat10(55)
	hl.OutTemp.Year.Low = p.get2ByteFloat10(61)
	hl.OutTemp.Year.Hi = p.get2ByteFloat10(59)

	// Rain rate
	hl.RainRate.Hour.Hi = p.getRainClicks(120)
	hl.RainRate.Day.Hi = p.getRainClicks(116)
	hl.RainRate.Day.HiTime = p.get2ByteTime(118)
	hl.RainRate.Month.Hi = p.getRainClicks(122)
	hl.RainRate.Year.Hi = p.getRainClicks(124)

	// Soil moisture and temperature
	for i := uint(0); i < 4; i++ {
		if low := p.get1ByteInt(368 + i); low != 255 {
			sm := HiLowSoilMoist{}
			sm.Day.Low = low
			sm.Day.LowTime = p.get2ByteTime(372 + i*2)
			sm.Day.Hi = p.get1ByteInt(356 + i)
			sm.Day.HiTime = p.get2ByteTime(360 + i*2)
			sm.Month.Low = p.get1ByteInt(380 + i)
			sm.Month.Hi = p.get1ByteInt(384 + i)
			sm.Year.Low = p.get1ByteInt(388 + i)
			sm.Year.Hi = p.get1ByteInt(392 + i)
			hl.SoilMoist[i] = &sm
		}

		if et := extraTemp(7 + i); et.Day.Low != 165 {
			hl.SoilTemp[i] = &et
		}
	}

	// Solar radiation
	hl.SolarRad.Day.Hi = p.get2ByteInt(103)
	hl.SolarRad.Day.HiTime = p.get2ByteTime(105)
	hl.SolarRad.Month.Hi = p.get2ByteInt(107)
	hl.SolarRad.Year.Hi = p.get2ByteInt(109)

	// THSW index
	hl.THSWIndex.Day.Hi = p.get2ByteFloat(95)
	hl.THSWIndex.Day.HiTime = p.get2ByteTime(97)
	hl.THSWIndex.Month.Hi = p.get2ByteFloat(99)
	hl.THSWIndex.Year.Hi = p.get2ByteFloat(101)

	// UltraViolet index
	hl.UVIndex.Day.Hi = p.getUVIndex(111)
	hl.UVIndex.Day.HiTime = p.get2ByteTime(112)
	hl.UVIndex.Month.Hi = p.getUVIndex(114)
	hl.UVIndex.Year.Hi = p.getUVIndex(115)

	// Wind speed
	hl.WindSpeed.Day.Hi = p.get1ByteMPH(16)
	hl.WindSpeed.Day.HiTime = p.get2ByteTime(17)
	hl.WindSpeed.Month.Hi = p.get1ByteMPH(19)
	hl.WindSpeed.Year.Hi = p.get1ByteMPH(20)

	// Wind chill
	hl.WindChill.Day.Low = p.get2ByteFloat(79)
	hl.WindChill.Day.LowTime = p.get2ByteTime(81)
	hl.WindChill.Month.Low = p.get2ByteFloat(83)
	hl.WindChill.Year.Low = p.get2ByteFloat(85)

	return nil
}
