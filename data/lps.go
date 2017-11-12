// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package data

import (
	"time"

	"github.com/ebarkie/weatherlink/packet"
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
	ExtraHumidity [7]*int   `json:"extraHumidity,omitempty"`
	ExtraTemp     [7]*int   `json:"extraTemperature,omitempty"`
	Forecast      string    `json:"forecast"`
	HeatIndex     float64   `json:"heatIndex"`
	Icons         []string  `json:"icons"`
	InHumidity    int       `json:"insideHumidity"`
	InTemp        float64   `json:"insideTemperature"`
	LeafTemp      [4]*int   `json:"leafTemperature,omitempty"`
	LeafWetness   [4]*int   `json:"leafWetness,omitempty"`
	OutHumidity   int       `json:"outsideHumidity"`
	OutTemp       float64   `json:"outsideTemperature"`
	Rain          LoopRain  `json:"rain"`
	SoilMoist     [4]*int   `json:"soilMoisture,omitempty"`
	SoilTemp      [4]*int   `json:"soilTemperature,omitempty"`
	SolarRad      int       `json:"solarRadiation"`
	Sunrise       time.Time `json:"sunrise,omitempty"`
	Sunset        time.Time `json:"sunset,omitempty"`
	THSWIndex     float64   `json:"THSWIndex"`
	UVIndex       float64   `json:"UVIndex"`
	Wind          LoopWind  `json:"wind"`
	WindChill     float64   `json:"windChill"`

	LoopType   int `json:"-"`
	NextArcRec int `json:"-"`
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
	Accum struct {
		Last15Min   float64 `json:"last15Minutes"`
		LastHour    float64 `json:"lastHour"`
		Last24Hours float64 `json:"last24Hours"`
		Today       float64 `json:"today"`
		LastMonth   float64 `json:"lastMonth"`
		LastYear    float64 `json:"lastYear"`
		Storm       float64 `json:"storm"`
	} `json:"accumulation"`
	Rate           float64   `json:"rate"`
	StormStartDate time.Time `json:"stormStartDate,omitempty"`
}

// LoopWind is the wind related readings for a Loop struct.
type LoopWind struct {
	Avg struct {
		Last2MinSpeed  float64 `json:"last2MinutesSpeed"`
		Last10MinSpeed float64 `json:"last10MinutesSpeed"`
	} `json:"average"`
	Cur struct {
		Dir   int `json:"direction"`
		Speed int `json:"speed"`
	} `json:"current"`
	Gust struct {
		Last10MinDir   int     `json:"last10MinutesDirection"`
		Last10MinSpeed float64 `json:"last10MinutesSpeed"`
	} `json:"gust"`
}

// UnmarshalBinary decodes a 99-byte loop 1 or 2 packet into the
// Loop struct.
func (l *Loop) UnmarshalBinary(p []byte) error {
	if packet.Crc(p) != 0 {
		return ErrBadCRC
	}

	l.LoopType = getLoopType(p)
	switch l.LoopType {
	case -1:
		// Packet length or header didn't make sense.
		return ErrNotLoop
	case 1:
		// Loop1
		l.Bar.SeaLevel = packet.GetPressure(p, 7)
		l.Bar.Trend = packet.GetBarTrend(p, 3)
		l.Bat.ConsoleVoltage = packet.GetVoltage(p, 87)
		l.Bat.TransStatus = packet.GetUInt8(p, 86)
		l.ET.Today = packet.GetFloat16(p, 56) / 1000.0
		l.ET.LastMonth = packet.GetFloat16(p, 58) / 100.0
		l.ET.LastYear = packet.GetFloat16(p, 60) / 100.0
		for i := uint(0); i < 7; i++ {
			if v := packet.GetUInt8(p, 34+i); v != 255 {
				l.ExtraHumidity[i] = &v
			}
			if v := packet.GetTemp8(p, 18+i); v != 165 {
				l.ExtraTemp[i] = &v
			}
		}
		l.Forecast = packet.GetForecast(p, 90)
		l.Icons = packet.GetForecastIcons(p, 89)
		l.InHumidity = packet.GetUInt8(p, 11)
		l.InTemp = packet.GetFloat16_10(p, 9)
		for i := uint(0); i < 4; i++ {
			if v := packet.GetTemp8(p, 29+i); v != 165 {
				l.LeafTemp[i] = &v
			}
			if v := packet.GetUInt8(p, 66+i); v != 255 {
				// There's a bug in my Davis firmware where the last leaf
				// wetness sensor returns 0 when it should be returning the
				// dash value.  This hack corrects it but could nil out a
				// valid value of zero.
				if i == 3 && v == 0 {
					continue
				}
				l.LeafWetness[i] = &v
			}
		}
		l.OutHumidity = packet.GetUInt8(p, 33)
		l.OutTemp = packet.GetFloat16_10(p, 12)
		l.Rain.Accum.Today = packet.GetRainClicks(p, 50)
		l.Rain.Accum.LastMonth = packet.GetRainClicks(p, 52)
		l.Rain.Accum.LastYear = packet.GetRainClicks(p, 54)
		l.Rain.Accum.Storm = packet.GetRainClicks(p, 46)
		l.Rain.Rate = packet.GetRainClicks(p, 41)
		l.Rain.StormStartDate = packet.GetDate16(p, 48)
		for i := uint(0); i < 4; i++ {
			if v := packet.GetUInt8(p, 62+i); v != 255 {
				l.SoilMoist[i] = &v
			}
			if v := packet.GetTemp8(p, 25+i); v != 165 {
				l.SoilTemp[i] = &v
			}
		}
		l.SolarRad = packet.GetUInt16(p, 44)
		l.Sunrise = packet.GetTime16(p, 91)
		l.Sunset = packet.GetTime16(p, 93)
		l.UVIndex = packet.GetUVIndex(p, 43)
		l.Wind.Cur.Dir = packet.GetUInt16(p, 16)
		l.Wind.Cur.Speed = packet.GetMPH8(p, 14)
		// Intentionally skip l.Wind.Avg.Last10MinSpeed because
		// the loop2 decode is more precise.
		// l.Wind.Avg.Last10MinSpeed = packet.GetMPH8(p, 15)

		l.NextArcRec = packet.GetUInt16(p, 5)
	case 2:
		// Loop2
		l.Bar.Altimeter = packet.GetPressure(p, 69)
		l.Bar.SeaLevel = packet.GetPressure(p, 7)
		l.Bar.Station = packet.GetPressure(p, 65)
		l.Bar.Trend = packet.GetBarTrend(p, 3)
		l.DewPoint = packet.GetFloat16(p, 30)
		l.ET.Today = packet.GetFloat16(p, 56) / 1000.0
		l.HeatIndex = packet.GetFloat16(p, 35)
		l.InHumidity = packet.GetUInt8(p, 11)
		l.InTemp = packet.GetFloat16_10(p, 9)
		l.OutHumidity = packet.GetUInt8(p, 33)
		l.OutTemp = packet.GetFloat16_10(p, 12)
		l.Rain.Accum.Last15Min = packet.GetRainClicks(p, 52)
		l.Rain.Accum.LastHour = packet.GetRainClicks(p, 54)
		l.Rain.Accum.Last24Hours = packet.GetRainClicks(p, 58)
		l.Rain.Accum.Today = packet.GetRainClicks(p, 50)
		l.Rain.Accum.Storm = packet.GetRainClicks(p, 46)
		l.Rain.Rate = packet.GetRainClicks(p, 41)
		l.SolarRad = packet.GetUInt16(p, 44)
		l.THSWIndex = packet.GetFloat16(p, 39)
		l.UVIndex = packet.GetUVIndex(p, 43)
		l.Wind.Cur.Dir = packet.GetUInt16(p, 16)
		l.Wind.Cur.Speed = packet.GetMPH8(p, 14)
		l.Wind.Avg.Last2MinSpeed = packet.GetMPH16(p, 20)
		l.Wind.Avg.Last10MinSpeed = packet.GetMPH16(p, 18)
		l.Wind.Gust.Last10MinDir = packet.GetUInt16(p, 24)
		l.Wind.Gust.Last10MinSpeed = packet.GetMPH16(p, 22)
		l.WindChill = packet.GetFloat16(p, 37)
	default:
		// Valid loop but a newer version than we know about.  This
		// should never happen since the protocol LPS loop request bit
		// mask only calls for the above versions.
		return ErrUnknownLoop
	}

	return nil
}

// MarshalBinary encodes the data from the Loop struct into a 99-byte loop 1
// or 2 packet.
func (l *Loop) MarshalBinary() (p []byte, err error) {
	p = make([]byte, 99)

	switch l.LoopType {
	case 1:
		// Loop1
		packet.SetPressure(&p, 7, l.Bar.SeaLevel)
		packet.SetFloat16(&p, 56, l.ET.Today*1000.0)
		packet.SetFloat16(&p, 58, l.ET.LastMonth*100.0)
		packet.SetFloat16(&p, 60, l.ET.LastYear*100.0)
		for i := uint(0); i < 7; i++ {
			if l.ExtraHumidity[i] != nil {
				packet.SetInt8(&p, 34+i, *l.ExtraHumidity[i])
			} else {
				packet.SetInt8(&p, 34+i, 255)
			}
			if l.ExtraTemp[i] != nil {
				packet.SetTemp8(&p, 18+i, *l.ExtraTemp[i])
			} else {
				packet.SetTemp8(&p, 18+i, 165)
			}
		}
		packet.SetInt8(&p, 11, l.InHumidity)
		packet.SetFloat16_10(&p, 9, l.InTemp)
		packet.SetInt8(&p, 33, l.OutHumidity)
		packet.SetFloat16_10(&p, 12, l.OutTemp)
		packet.SetRainClicks(&p, 50, l.Rain.Accum.Today)
		packet.SetRainClicks(&p, 52, l.Rain.Accum.LastMonth)
		packet.SetRainClicks(&p, 54, l.Rain.Accum.LastYear)
		packet.SetRainClicks(&p, 46, l.Rain.Accum.Storm)
		packet.SetRainClicks(&p, 41, l.Rain.Rate)
		for i := uint(0); i < 4; i++ {
			if l.SoilMoist[i] != nil {
				packet.SetInt8(&p, 62+i, *l.SoilMoist[i])
			} else {
				packet.SetInt8(&p, 62+i, 255)
			}
			if l.SoilTemp[i] != nil {
				packet.SetTemp8(&p, 25+i, *l.SoilTemp[i])
			} else {
				packet.SetTemp8(&p, 25+i, 165)
			}
		}
		packet.SetMPH8(&p, 14, l.Wind.Cur.Speed)

		packet.SetInt16(&p, 5, l.NextArcRec)
	case 2:
		// Loop2
		packet.SetPressure(&p, 69, l.Bar.Altimeter)
		packet.SetPressure(&p, 7, l.Bar.SeaLevel)
		packet.SetPressure(&p, 65, l.Bar.Station)
		packet.SetFloat16(&p, 30, l.DewPoint)
		packet.SetFloat16(&p, 56, l.ET.Today*1000.0)
		packet.SetFloat16(&p, 35, l.HeatIndex)
		packet.SetInt8(&p, 11, l.InHumidity)
		packet.SetFloat16_10(&p, 9, l.InTemp)
		packet.SetInt8(&p, 33, l.OutHumidity)
		packet.SetFloat16_10(&p, 12, l.OutTemp)
		packet.SetRainClicks(&p, 52, l.Rain.Accum.Last15Min)
		packet.SetRainClicks(&p, 54, l.Rain.Accum.LastHour)
		packet.SetRainClicks(&p, 58, l.Rain.Accum.Last24Hours)
		packet.SetRainClicks(&p, 50, l.Rain.Accum.Today)
		packet.SetRainClicks(&p, 46, l.Rain.Accum.Storm)
		packet.SetRainClicks(&p, 41, l.Rain.Rate)
		packet.SetMPH8(&p, 14, l.Wind.Cur.Speed)
	default:
		err = ErrUnknownLoop
	}

	setLoopType(&p, l.LoopType)

	packet.SetCrc(&p)

	return
}

// getLoopType returns the loop packet numeric type or -1 if it
// is not a valid loop packet.
func getLoopType(p []byte) int {
	if len(p) == 99 &&
		p[0] == 0x4c &&
		p[1] == 0x4f &&
		p[2] == 0x4f { // LOO
		return packet.GetUInt8(p, 4) + 1
	}

	return -1
}

func setLoopType(p *[]byte, t int) {
	(*p)[0] = 0x4c
	(*p)[1] = 0x4f
	(*p)[2] = 0x4f // LOO

	(*p)[4] = byte(t - 1)
}
