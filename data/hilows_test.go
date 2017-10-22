// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of thighs source code is governed by the MIT license
// that can be found in the LICENSE file.

package data

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testHiLowsPackets = map[string][]byte{
	"std": {
		0x68, 0x75, 0xe1, 0x75, 0x3a, 0x74, 0xe1, 0x75,
		0x11, 0x72, 0x92, 0x78, 0x1c, 0x07, 0xfc, 0x03,
		0x0c, 0x8d, 0x00, 0x0d, 0x1b, 0x20, 0x03, 0x05,
		0x03, 0x0c, 0x00, 0xae, 0x03, 0xea, 0x02, 0x2d,
		0x03, 0x47, 0x02, 0x30, 0x03, 0x2b, 0x26, 0x7c,
		0x02, 0x31, 0x00, 0x33, 0x25, 0x3a, 0x15, 0xce,
		0x02, 0x79, 0x03, 0xfe, 0x01, 0x36, 0x05, 0xf5,
		0x03, 0xa8, 0x02, 0xf5, 0x03, 0x5c, 0x00, 0x42,
		0x00, 0x49, 0x00, 0x96, 0x01, 0x0f, 0x00, 0x52,
		0x00, 0x41, 0x00, 0x52, 0x00, 0x00, 0x00, 0x48,
		0x00, 0xca, 0x01, 0x44, 0x00, 0x09, 0x00, 0x60,
		0x00, 0x34, 0x05, 0x77, 0x00, 0x77, 0x00, 0x6d,
		0x00, 0x8e, 0x05, 0x84, 0x00, 0x84, 0x00, 0x88,
		0x04, 0x21, 0x05, 0xf2, 0x04, 0x46, 0x05, 0x3e,
		0xeb, 0x04, 0x52, 0x5d, 0x00, 0x00, 0xff, 0xff,
		0x00, 0x00, 0xec, 0x01, 0x20, 0x1c, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xaa, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xad, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xf1, 0x02, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0x50, 0x06, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xaf,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa7, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xaf, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0x81, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0x34, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0x53, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0x3b, 0x05, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0x05, 0x00, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0x62, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0x23, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0x63, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0x0f, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0x1d, 0xff, 0xff, 0xff,
		0xcc, 0x06, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0x17, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x01, 0xff, 0xff, 0xff,
		0xc4, 0xff, 0xff, 0xff, 0x01, 0xff, 0xff, 0xff,
		0xc4, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0x85, 0x9a,
	},
}

func BenchmarkHiLowsUnmarshalBinary(b *testing.B) {
	for n := 0; n < b.N; n++ {
		hl := HiLows{}
		hl.UnmarshalBinary(testHiLowsPackets["std"])
	}
}

func TestHiLowsUnmarshalBinary(t *testing.T) {
	a := assert.New(t)

	hl := HiLows{}
	err := hl.UnmarshalBinary(testHiLowsPackets["std"])
	a.Nil(err, "UnmarshalBinary hilows")

	// Barometer
	a.Equal(30.056, hl.Bar.Day.Low, "Barometer day low")
	a.Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 18, 20, 0, 0, time.Local),
		hl.Bar.Day.LowTime, "Dew point day low time")
	a.Equal(30.177, hl.Bar.Day.Hi, "Barometer day high")
	a.Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 10, 20, 0, 0, time.Local),
		hl.Bar.Day.HiTime, "Dew point day high time")
	a.Equal(29.754, hl.Bar.Month.Low, "Barometer month low")
	a.Equal(30.177, hl.Bar.Month.Hi, "Barometer month high")
	a.Equal(29.201, hl.Bar.Year.Low, "Barometer year low")
	a.Equal(30.866, hl.Bar.Year.Hi, "Barometer year high")

	// Dew point
	a.Equal(66.0, hl.DewPoint.Day.Low, "Dew point day low")
	a.Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 4, 6, 0, 0, time.Local),
		hl.DewPoint.Day.LowTime, "Dew point day low time")
	a.Equal(73.0, hl.DewPoint.Day.Hi, "Dew point day high")
	a.Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 15, 0, 0, time.Local),
		hl.DewPoint.Day.HiTime, "Dew point day high time")
	a.Equal(65.0, hl.DewPoint.Month.Low, "Dew point month low")
	a.Equal(82.0, hl.DewPoint.Month.Hi, "Dew point month high")
	a.Equal(0.0, hl.DewPoint.Year.Low, "Dew point year low")
	a.Equal(82.0, hl.DewPoint.Year.Hi, "Dew point year high")

	// Extra humidity and temperatures
	for i := 0; i < 7; i++ {
		a.Nil(hl.ExtraHumidity[i], fmt.Sprintf("Extra humidity %d", i))
		a.Nil(hl.ExtraTemp[i], fmt.Sprintf("Extra temperature %d", i))
	}

	// Heat index
	a.Equal(96.0, hl.HeatIndex.Day.Hi, "Heat index day high")
	a.Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 13, 32, 0, 0, time.Local),
		hl.HeatIndex.Day.HiTime, "Heat index day high time")
	a.Equal(119.0, hl.HeatIndex.Month.Hi, "Heat index month high")
	a.Equal(119.0, hl.HeatIndex.Year.Hi, "Heat index year high")

	// Inside humidity
	a.Equal(38, hl.InHumidity.Day.Low, "Inside humidity day low")
	a.Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 49, 0, 0, time.Local),
		hl.InHumidity.Day.LowTime, "Inside humidity day low time")
	a.Equal(43, hl.InHumidity.Day.Hi, "Inside humidity day high")
	a.Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 6, 36, 0, 0, time.Local),
		hl.InHumidity.Day.HiTime, "Inside humidity day high time")
	a.Equal(37, hl.InHumidity.Month.Low, "Inside humidity month low")
	a.Equal(51, hl.InHumidity.Month.Hi, "Inside humidity month high")
	a.Equal(21, hl.InHumidity.Year.Low, "Inside humidity year low")
	a.Equal(58, hl.InHumidity.Year.Hi, "Inside humidity year high")

	// Inside temperature
	a.Equal(77.3, hl.InTemp.Day.Low, "Inside temperature day low")
	a.Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 9, 42, 0, 0, time.Local),
		hl.InTemp.Day.LowTime, "Inside temperature day low time")
	a.Equal(80.0, hl.InTemp.Day.Hi, "Inside temperature day high")
	a.Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 12, 0, 0, time.Local),
		hl.InTemp.Day.HiTime, "Inside temperature day high time")
	a.Equal(74.6, hl.InTemp.Month.Low, "Inside temperature month low")
	a.Equal(81.3, hl.InTemp.Month.Hi, "Inside temperature month high")
	a.Equal(58.3, hl.InTemp.Year.Low, "Inside temperature year low")
	a.Equal(81.6, hl.InTemp.Year.Hi, "Inside temperature year high")

	// Leaf temperature and wetness
	for i := 0; i < 4; i++ {
		a.Nil(hl.LeafTemp[i], fmt.Sprintf("Leaf temperature %d", i))
		a.Nil(hl.LeafWetness[i], fmt.Sprintf("Leaf wetness %d", i))
	}

	// Outside humidity
	a.Equal(52, hl.OutHumidity.Day.Low, "Outside humidity day low")
	a.Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 13, 39, 0, 0, time.Local),
		hl.OutHumidity.Day.LowTime, "Outside humidity day low time")
	a.Equal(83, hl.OutHumidity.Day.Hi, "Outside humidity day high")
	a.Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 5, 0, 0, time.Local),
		hl.OutHumidity.Day.HiTime, "Outside humidity day high time")
	a.Equal(35, hl.OutHumidity.Month.Low, "Outside humidity month low")
	a.Equal(98, hl.OutHumidity.Month.Hi, "Outside humidity month high")
	a.Equal(15, hl.OutHumidity.Year.Low, "Outside humidity year low")
	a.Equal(99, hl.OutHumidity.Year.Hi, "Outside humidity year high")

	// Outside temperature
	a.Equal(71.8, hl.OutTemp.Day.Low, "Outside temperature day low")
	a.Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 5, 10, 0, 0, time.Local),
		hl.OutTemp.Day.LowTime, "Outside temperature day low time")
	a.Equal(88.9, hl.OutTemp.Day.Hi, "Outside temperature day high")
	a.Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 13, 34, 0, 0, time.Local),
		hl.OutTemp.Day.HiTime, "Outside temperature day high time")
	a.Equal(68.0, hl.OutTemp.Month.Low, "Outside temperature month low")
	a.Equal(101.3, hl.OutTemp.Month.Hi, "Outside temperature month high")
	a.Equal(9.2, hl.OutTemp.Year.Low, "Outside temperature year low")
	a.Equal(101.3, hl.OutTemp.Year.Hi, "Outside temperature year high")

	// Rain rate
	a.Equal(0.0, hl.RainRate.Hour.Hi, "Rain rate hour high")
	a.Equal(0.0, hl.RainRate.Day.Hi, "Rain rate day high")
	a.Equal(time.Time{}, hl.RainRate.Day.HiTime, "Rain rate day high time")
	a.Equal(4.92, hl.RainRate.Month.Hi, "Rain rate month high")
	a.Equal(72.0, hl.RainRate.Year.Hi, "Rain rate year high")

	// Soil moisture and temperature
	a.Equal(23, hl.SoilMoist[0].Day.Low, "Soil moisture day low")
	a.Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local),
		hl.SoilMoist[0].Day.LowTime, "Soil moisture day low time")
	a.Equal(29, hl.SoilMoist[0].Day.Hi, "Soil moisture day high")
	a.Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 17, 40, 0, 0, time.Local),
		hl.SoilMoist[0].Day.HiTime, "Soil moisture day high time")
	a.Equal(1, hl.SoilMoist[0].Month.Low, "Soil moisture month low")
	a.Equal(196, hl.SoilMoist[0].Month.Hi, "Soil moisture month high")
	a.Equal(1, hl.SoilMoist[0].Year.Low, "Soil moisture year low")
	a.Equal(196, hl.SoilMoist[0].Year.Hi, "Soil moisture year high")

	a.Equal(80, hl.SoilTemp[0].Day.Low, "Soil temperature day low")
	a.Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 7, 53, 0, 0, time.Local),
		hl.SoilTemp[0].Day.LowTime, "Soil temperature day low time")
	a.Equal(83, hl.SoilTemp[0].Day.Hi, "Soil temperature day high")
	a.Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 16, 16, 0, 0, time.Local),
		hl.SoilTemp[0].Day.HiTime, "Soil temperature day high time")
	a.Equal(77, hl.SoilTemp[0].Month.Low, "Soil temperature month low")
	a.Equal(85, hl.SoilTemp[0].Month.Hi, "Soil temperature month high")
	a.Equal(39, hl.SoilTemp[0].Year.Low, "Soil temperature year low")
	a.Equal(85, hl.SoilTemp[0].Year.Hi, "Soil temperature year high")

	for i := 1; i < 4; i++ {
		a.Nil(hl.SoilMoist[i], fmt.Sprintf("Soil moisture %d", i))
		a.Nil(hl.SoilTemp[i], fmt.Sprintf("Soil temperature %d", i))
	}

	// Solar radiation
	a.Equal(1160, hl.SolarRad.Day.Hi, "Solar radiation day high")
	a.Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 13, 13, 0, 0, time.Local),
		hl.SolarRad.Day.HiTime, "Solar radiation day high time")
	a.Equal(1266, hl.SolarRad.Month.Hi, "Solar radiation month high")
	a.Equal(1350, hl.SolarRad.Year.Hi, "Solar radiation year high")

	// THSW index
	a.Equal(109.0, hl.THSWIndex.Day.Hi, "THSW index day high")
	a.Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 14, 22, 0, 0, time.Local),
		hl.THSWIndex.Day.HiTime, "THSW index day high time")
	a.Equal(132.0, hl.THSWIndex.Month.Hi, "THSW index month high")
	a.Equal(132.0, hl.THSWIndex.Year.Hi, "THSW index year high")

	// UltraViolet index
	a.Equal(6.2, hl.UVIndex.Day.Hi, "UltraViolet index day high")
	a.Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 12, 59, 0, 0, time.Local),
		hl.UVIndex.Day.HiTime, "UltraViolet index day high time")
	a.Equal(8.2, hl.UVIndex.Month.Hi, "UltraViolet index month high")
	a.Equal(9.3, hl.UVIndex.Year.Hi, "UltraViolet index year high")

	// Wind speed
	a.Equal(12, hl.WindSpeed.Day.Hi, "Wind speed day high")
	a.Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 1, 41, 0, 0, time.Local),
		hl.WindSpeed.Day.HiTime, "Wind speed day high time")
	a.Equal(13, hl.WindSpeed.Month.Hi, "Wind speed month high")
	a.Equal(27, hl.WindSpeed.Year.Hi, "Wind speed year high")

	// Wind chill
	a.Equal(72.0, hl.WindChill.Day.Low, "Wind chill day low")
	a.Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 4, 58, 0, 0, time.Local),
		hl.WindChill.Day.LowTime, "Wind chill day low time")
	a.Equal(68.0, hl.WindChill.Month.Low, "Wind chill month low")
	a.Equal(9.0, hl.WindChill.Year.Low, "Wind chill year low")
}
