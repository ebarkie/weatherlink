# data
```go
    import "github.com/ebarkie/weatherlink/data"
```

Package data implements encoding and decoding of Davis Instruments binary data
types.

## Usage

```go
var (
	ErrNotArcB     = errors.New("not a revision B archive record")
	ErrBadCRC      = errors.New("CRC check failed")
	ErrBadFirmVer  = errors.New("firmware version is not valid")
	ErrBadLocation = errors.New("location is inconsistent")
	ErrNotDmp      = errors.New("not a download memory page")
	ErrNotDmpMeta  = errors.New("not a download memory page metadata packet")
	ErrNotLoop     = errors.New("not a loop packet")
	ErrUnknownLoop = errors.New("unknown loop packet type")
)
```
Errors.

#### type Archive

```go
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
```

Archive represents all of the data in a revision B archive record.

#### func (*Archive) UnmarshalBinary

```go
func (a *Archive) UnmarshalBinary(p []byte) error
```
UnmarshalBinary decodes a 52-byte revision B archive record.

#### type ConsTime

```go
type ConsTime time.Time
```

ConsTime is the console current time.

#### func (ConsTime) MarshalBinary

```go
func (ct ConsTime) MarshalBinary() (p []byte, err error)
```
MarshalBinary encodes the console time into an 8-byte packet suitable for the
SETTIME command.

#### func (*ConsTime) UnmarshalBinary

```go
func (ct *ConsTime) UnmarshalBinary(p []byte) error
```
UnmarshalBinary decodes an 8-byte console time response packet into the ConsTime
struct.

#### type Dmp

```go
type Dmp [5]Archive
```

Dmp is a download memory page which contains 5 archive records.

#### func (*Dmp) UnmarshalBinary

```go
func (d *Dmp) UnmarshalBinary(p []byte) error
```
UnmarshalBinary decodes a 267-byte download memory page into an array of 5
Archive records.

#### type DmpAft

```go
type DmpAft time.Time
```

DmpAft is a timestamp appropriate for the "DMP after" command.

#### func (DmpAft) MarshalBinary

```go
func (da DmpAft) MarshalBinary() (p []byte, err error)
```
MarshalBinary encodes the data from the DmpAft struct into a 6-byte packet
appropriate for use with the DMPAFT command.

#### type DmpMeta

```go
type DmpMeta struct {
	Pages           int // Number of pages to download
	FirstPageOffset int // Offset of the first record to read within the first page
}
```

DmpMeta is the DMP metadata sent after the DMPAFT command is issued. It informs
the downloader how much data to expect and where the first record is within the
first page.

#### func (*DmpMeta) UnmarshalBinary

```go
func (dm *DmpMeta) UnmarshalBinary(p []byte) error
```
UnmarshalBinary decodes a 6-byte DMP metadata packet into the DmpMeta stuct.

#### type EEPROM

```go
type EEPROM struct {
	ArchivePeriod int           `json:"archivePeriod"`
	Elev          int           `json:"elevation"`
	Lat           float64       `json:"latitude"`
	Lon           float64       `json:"longitude"`
	TimeOffset    time.Duration `json:"timeOffset"`
}
```

EEPROM represents the configuration settings.

#### func (*EEPROM) UnmarshalBinary

```go
func (ee *EEPROM) UnmarshalBinary(p []byte) error
```
UnmarshalBinary decodes a 4096-byte EEPROM packet into the EEPROM struct.

#### type FirmTime

```go
type FirmTime time.Time
```

FirmTime is the firmware build time.

#### func (FirmTime) MarshalText

```go
func (ft FirmTime) MarshalText() ([]byte, error)
```
MarshalText encodes the firmware build time into a 13-byte packet suitable for
the VER command.

#### func (*FirmTime) UnmarshalText

```go
func (ft *FirmTime) UnmarshalText(p []byte) error
```
UnmarshalText decodes a 13-byte firmware build time response packet into the
FirmTime struct.

#### type FirmVer

```go
type FirmVer string
```

FirmVer is the firmware version number.

#### func (FirmVer) MarshalText

```go
func (fv FirmVer) MarshalText() ([]byte, error)
```
MarshalText encodes the firmware version into a 6-byte packet suitable for the
NVER command.

#### func (*FirmVer) UnmarshalText

```go
func (fv *FirmVer) UnmarshalText(p []byte) error
```
UnmarshalText decodes a 6-byte firmware version response packet into the FirmVer
struct.

#### type HiHeatIndex

```go
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
```

HiHeatIndex is the record high heat index readings.

#### type HiLowBar

```go
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
```

HiLowBar is the record high and low barometer readings.

#### type HiLowExtraTemp

```go
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
```

HiLowExtraTemp is the record high and low extra temperature readings.

#### type HiLowHumidity

```go
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
```

HiLowHumidity is the record high and low humidity readings.

#### type HiLowLeafWetness

```go
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
```

HiLowLeafWetness is the record high and low leaf wetness readings.

#### type HiLowSoilMoist

```go
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
```

HiLowSoilMoist is the record high and low soil moisture readings.

#### type HiLowTemp

```go
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
```

HiLowTemp is the record high and low temperature readings and dew point
calculations.

#### type HiLows

```go
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
```

HiLows represents all of the record high and lows by day, month, and year. The
day also includes the time(s) when the record occurred.

#### func (*HiLows) UnmarshalBinary

```go
func (hl *HiLows) UnmarshalBinary(p []byte) error
```
UnmarshalBinary decodes a 438-byte high and lows packet into the HiLows struct.

#### type HiRainRate

```go
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
```

HiRainRate is the record high rain rate readings.

#### type HiSolarRad

```go
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
```

HiSolarRad is the record high solar radiation readings.

#### type HiTHSWIndex

```go
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
```

HiTHSWIndex is the record high THSW index calculations.

#### type HiUVIndex

```go
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
```

HiUVIndex is the record high UltraViolet index readings.

#### type HiWindSpeed

```go
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
```

HiWindSpeed is the record high wind speed readings.

#### type Loop

```go
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
	LeafWet       [4]*int   `json:"leafWetness,omitempty"`
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
```

Loop is a combined struct representation of the union of loop1 and loop2
packets. They have a lot of overlap but the precision is sometimes different and
they complement each other.

During the protocol loop polling with the LPS command the two versions are
interleaved.

#### func (*Loop) MarshalBinary

```go
func (l *Loop) MarshalBinary() (p []byte, err error)
```
MarshalBinary encodes the data from the Loop struct into a 99-byte loop 1 or 2
packet.

#### func (*Loop) UnmarshalBinary

```go
func (l *Loop) UnmarshalBinary(p []byte) error
```
UnmarshalBinary decodes a 99-byte loop 1 or 2 packet into the Loop struct.

#### type LoopBar

```go
type LoopBar struct {
	Altimeter float64 `json:"altimeter"`
	SeaLevel  float64 `json:"seaLevel"`
	Station   float64 `json:"station"`
	Trend     string  `json:"trend"`
}
```

LoopBar is the barometer related readings for a Loop struct.

#### type LoopBat

```go
type LoopBat struct {
	ConsoleVoltage float64 `json:"consoleVoltage"`
	TransLow       []int   `json:"transmittersLow"`
}
```

LoopBat is the console and transmitter battery readings for a Loop struct.

#### type LoopET

```go
type LoopET struct {
	Today     float64 `json:"today"`
	LastMonth float64 `json:"lastMonth"`
	LastYear  float64 `json:"lastYear"`
}
```

LoopET is the evapotranspiration related readings for a Loop struct.

#### type LoopRain

```go
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
```

LoopRain is the rain sensor related readings for a Loop struct.

#### type LoopWind

```go
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
```

LoopWind is the wind related readings for a Loop struct.

#### type LowWindChill

```go
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
```

LowWindChill is the record low wind chill calculations.
