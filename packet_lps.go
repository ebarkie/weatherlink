// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

// Packet coding logic for loop revision 1 and 2 packets.
//
// Refer to Vantage ProTM, Vantage Pro2TM and Vantage VueTM Serial
// Communication Reference Manual, section X. Data Formats,
// subsection 1. LOOP data format.

import (
	"errors"
	"time"
)

// Errors.
var (
	ErrNotLoop     = errors.New("Not a loop packet")
	ErrUnknownLoop = errors.New("Loop packet type is unknown")
)

// Loop is a combined struct representation of the union of loop1
// and loop2 packets.  They have a lot of overlap but the precision
// is sometimes different and they complement each other.
//
// During the protocol loop polling with the LPS command the two
// versions are interleaved.
type Loop struct {
	Bar           LoopBar   `json:"barometer"`
	Bat           LoopBat   `json:"battery"`
	DewPoint      float64   `json:"dewPoint"`
	ET            LoopET    `json:"ET"`
	ExtraHumidity [7]*int   `json:"extraHumidity"`
	ExtraTemp     [7]*int   `json:"extraTemperature"`
	Forecast      string    `json:"forecast"`
	HeatIndex     float64   `json:"heatIndex"`
	Icons         []string  `json:"icons"`
	InHumidity    int       `json:"insideHumidity"`
	InTemp        float64   `json:"insideTemperature"`
	LeafTemp      [4]*int   `json:"leafTemperature"`
	LeafWetness   [4]*int   `json:"leafWetness"`
	OutHumidity   int       `json:"outsideHumidity"`
	OutTemp       float64   `json:"outsideTemperature"`
	Rain          LoopRain  `json:"rain"`
	SoilMoist     [4]*int   `json:"soilMoisture"`
	SoilTemp      [4]*int   `json:"soilTemperature"`
	SolarRad      int       `json:"solarRadiation"`
	Sunrise       time.Time `json:"sunrise"`
	Sunset        time.Time `json:"sunset"`
	THSWIndex     float64   `json:"THSWIndex"`
	UVIndex       float64   `json:"UVIndex"`
	Wind          LoopWind  `json:"wind"`
	WindChill     float64   `json:"windChill"`
	nextArchRec   int
}

// LoopBar is the barometer related readings for a Loop struct.
type LoopBar struct {
	Altimeter float64 `json:"altimeter"`
	SeaLevel  float64 `json:"seaLevel"`
	Station   float64 `json:"station"`
	Trend     string  `json:"trend"`
}

// LoopBat is the console and transmitter battery readings for a Loop struct.
type LoopBat struct {
	ConsoleVoltage float64 `json:"consoleVoltage"`
	TransStatus    int     `json:"transmitterStatus"`
}

// LoopET is the evapotranspiration related readings for a Loop struct.
type LoopET struct {
	Today     float64 `json:"today"`
	LastMonth float64 `json:"lastMonth"`
	LastYear  float64 `json:"lastYear"`
}

// LoopRain is the rain sensor related readings for a Loop struct.
type LoopRain struct {
	Accum          LoopRainAccum `json:"accumulation"`
	Rate           float64       `json:"rate"`
	StormStartDate time.Time     `json:"stormStartDate"`
}

// LoopRainAccum is the rain accumulation related readings for a LoopRain struct.
type LoopRainAccum struct {
	Last15Min   float64 `json:"last15Minutes"`
	LastHour    float64 `json:"lastHour"`
	Last24Hours float64 `json:"last24Hours"`
	Today       float64 `json:"today"`
	LastMonth   float64 `json:"lastMonth"`
	LastYear    float64 `json:"lastYear"`
	Storm       float64 `json:"storm"`
}

// LoopWind is the wind related readings for a Loop struct.
type LoopWind struct {
	Avg  LoopWindAvgs  `json:"average"`
	Cur  LoopWindCur   `json:"current"`
	Gust LoopWindGusts `json:"gust"`
}

// LoopWindAvgs is the average wind speed related readings for a LoopWind struct.
type LoopWindAvgs struct {
	Last2MinSpeed  float64 `json:"last2MinutesSpeed"`
	Last10MinSpeed float64 `json:"last10MinutesSpeed"`
}

// LoopWindCur is the current wind direction and speed for a LoopWind struct.
type LoopWindCur struct {
	Dir   int `json:"direction"`
	Speed int `json:"speed"`
}

// LoopWindGusts is the wind gust related readings for a LoopWind struct.
type LoopWindGusts struct {
	Last10MinDir   int     `json:"last10MinutesDirection"`
	Last10MinSpeed float64 `json:"last10MinutesSpeed"`
}

// FromPacket unpacks the data from a 99-byte loop 1 or 2 packet into
// the Loop struct.
func (l *Loop) FromPacket(p Packet) error {
	if crc(p) != 0 {
		return ErrBadCRC
	}

	t := p.getLoopType()
	switch t {
	case -1:
		// Packet length or header didn't make sense.
		return ErrNotLoop
	case 1:
		// Loop1
		l.Bar.SeaLevel = p.getPressure(7)
		l.Bar.Trend = p.getBarTrend(3)
		l.Bat.ConsoleVoltage = p.getVoltage(87)
		l.Bat.TransStatus = p.get1ByteInt(86)
		l.ET.Today = p.get2ByteFloat(56) / 1000
		l.ET.LastMonth = p.get2ByteFloat(58) / 100
		l.ET.LastYear = p.get2ByteFloat(60) / 100
		for i := uint(0); i < 7; i++ {
			if v := p.get1ByteInt(34 + i); v != 255 {
				l.ExtraHumidity[i] = &v
			}
			if v := p.get1ByteTemp(18 + i); v != 165 {
				l.ExtraTemp[i] = &v
			}
		}
		l.Forecast = p.getForecast(90)
		l.Icons = p.getForecastIcons(89)
		l.InHumidity = p.get1ByteInt(11)
		l.InTemp = p.get2ByteTemp(9)
		for i := uint(0); i < 4; i++ {
			if v := p.get1ByteTemp(29 + i); v != 165 {
				l.LeafTemp[i] = &v
			}
			if v := p.get1ByteInt(66 + i); v != 255 {
				// There's a bug in my Davis firmware where the last leaf
				// wetness sensor returns 0 when it should be returning the
				// dash value.  This hack corrects it but could nil out a
				// valid value of zero.
				if (i == 3) && (v == 0) {
					continue
				}
				l.LeafWetness[i] = &v
			}
		}
		l.OutHumidity = p.get1ByteInt(33)
		l.OutTemp = p.get2ByteTemp(12)
		l.Rain.Accum.Today = p.getRainClicks(50)
		l.Rain.Accum.LastMonth = p.getRainClicks(52)
		l.Rain.Accum.LastYear = p.getRainClicks(54)
		l.Rain.Accum.Storm = p.getRainClicks(46)
		l.Rain.Rate = p.getRainClicks(41)
		l.Rain.StormStartDate = p.get2ByteDate(48)
		for i := uint(0); i < 4; i++ {
			if v := p.get1ByteInt(62 + i); v != 255 {
				l.SoilMoist[i] = &v
			}
			if v := p.get1ByteTemp(25 + i); v != 165 {
				l.SoilTemp[i] = &v
			}
		}
		l.SolarRad = p.get2ByteInt(44)
		l.Sunrise = p.get2ByteTime(91)
		l.Sunset = p.get2ByteTime(93)
		l.UVIndex = p.getUVIndex(43)
		l.Wind.Cur.Dir = p.get2ByteInt(16)
		l.Wind.Cur.Speed = p.get1ByteMPH(14)
		// Intentionally skip l.Wind.Avg.Last10MinSpeed because
		// the loop2 decode is more precise.
		// l.Wind.Avg.Last10MinSpeed = p.get1ByteMPH(15)

		l.nextArchRec = p.get2ByteInt(5)
	case 2:
		// Loop2
		l.Bar.Altimeter = p.getPressure(69)
		l.Bar.SeaLevel = p.getPressure(7)
		l.Bar.Station = p.getPressure(65)
		l.Bar.Trend = p.getBarTrend(3)
		// Most calculated values are in whole degrees so the division done with
		// get2ByteTemp needs to be reversed.
		l.DewPoint = p.get2ByteTemp(30) * 10
		l.ET.Today = p.get2ByteFloat(56) / 1000
		l.HeatIndex = p.get2ByteTemp(35) * 10
		l.InHumidity = p.get1ByteInt(11)
		l.InTemp = p.get2ByteTemp(9)
		l.OutHumidity = p.get1ByteInt(33)
		l.OutTemp = p.get2ByteTemp(12)
		l.Rain.Accum.Last15Min = p.getRainClicks(52)
		l.Rain.Accum.LastHour = p.getRainClicks(54)
		l.Rain.Accum.Last24Hours = p.getRainClicks(58)
		l.Rain.Accum.Today = p.getRainClicks(50)
		l.Rain.Accum.Storm = p.getRainClicks(46)
		l.Rain.Rate = p.getRainClicks(41)
		l.SolarRad = p.get2ByteInt(44)
		l.THSWIndex = p.get2ByteTemp(39) * 10
		l.UVIndex = p.getUVIndex(43)
		l.Wind.Cur.Dir = p.get2ByteInt(16)
		l.Wind.Cur.Speed = p.get1ByteMPH(14)
		l.Wind.Avg.Last2MinSpeed = p.get2ByteMPH(20)
		l.Wind.Avg.Last10MinSpeed = p.get2ByteMPH(18)
		l.Wind.Gust.Last10MinDir = p.get2ByteInt(24)
		l.Wind.Gust.Last10MinSpeed = p.get2ByteMPH(22)
		l.WindChill = p.get2ByteTemp(37) * 10
	default:
		// Valid loop but a newer version than we know about.  This
		// should never happen since the protocol LPS loop request bit
		// mask only calls for the above versions.
		return ErrUnknownLoop
	}

	return nil
}

// ToPacket packs the data from the Loop struct into a 99-byte loop 1
// or 2 packet.
func (l *Loop) ToPacket(t int) (p Packet, err error) {
	p = make(Packet, 99)

	switch t {
	case 1:
		// TODO
	case 2:
		// TODO
	default:
		err = ErrUnknownLoop
	}

	p.setCrc()

	return
}

// getLoopType returns the loop packet numeric type or -1 if it
// is not a valid loop packet.
func (p Packet) getLoopType() int {
	if (len(p) == 99) &&
		(p[0] == 0x4c) && (p[1] == 0x4f) && (p[2] == 0x4f) { // LOO
		return p.get1ByteInt(4) + 1
	}

	return -1
}
