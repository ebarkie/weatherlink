# weatherlink

```go
import "github.com/ebarkie/weatherlink"
```

Package weatherlink implements the Davis Instruments serial, USB, and TCP/IP
communication protocol.

## Usage

```go
const (
	CmdGetDmps cmd = iota
	CmdGetLoops
	CmdStop
	CmdSyncConsTime
)
```
Commands that can be requested.

```go
var (
	Trace = log.New(ioutil.Discard, "[TRCE]", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	Debug = log.New(ioutil.Discard, "[DBUG]", log.LstdFlags|log.Lshortfile)
	Info  = log.New(ioutil.Discard, "[INFO]", log.LstdFlags)
	Warn  = log.New(ioutil.Discard, "[WARN]", log.LstdFlags|log.Lshortfile)
	Error = log.New(ioutil.Discard, "[ERRO]", log.LstdFlags|log.Lshortfile)
)
```
Loggers

```go
var (
	ErrNotDmp  = errors.New("Not a dmp metadata packet")
	ErrNotDmpB = errors.New("Not a revision b dmp packet")
)
```
Errors.

```go
var (
	ErrNotLoop     = errors.New("Not a loop packet")
	ErrUnknownLoop = errors.New("Loop packet type is unknown")
)
```
Errors.

```go
var (
	ErrProtoCmdFailed = errors.New("Protocol command failed")
	ErrNoLoopChan     = errors.New("Can't start command broker without a Loop channel")
)
```
Errors.

```go
var (
	ErrBadCRC = errors.New("CRC check failed")
)
```
Errors.

#### type Archive

```go
type Archive struct {
	Bar            float64   `json:"barometer"`
	ET             float64   `json:"ET"`
	ExtraHumidity  [2]*int   `json:"extraHumidity"`
	ExtraTemp      [3]*int   `json:"extraTemperature"`
	Forecast       string    `json:"forecast"`
	InsideHumidity int       `json:"insideHumidity"`
	InsideTemp     float64   `json:"insideTemperature"`
	LeafTemp       [2]*int   `json:"leafTemperature"`
	LeafWetness    [2]*int   `json:"leafWetness"`
	OutHumidity    int       `json:"outsideHumidity"`
	OutTemp        float64   `json:"outsideTemperature"`
	OutTempHi      float64   `json:"outsideTemperatureHigh"`
	OutTempLow     float64   `json:"outsideTemperatureLow"`
	RainAccum      float64   `json:"rainAccumulation"`
	RainRateHi     float64   `json:"rainRateHigh"`
	SoilMoist      [4]*int   `json:"soilMoisture"`
	SoilTemp       [4]*int   `json:"soilTemperature"`
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

Archive represents all of the data in a revision B DMP archive record.

#### type ConsTime

```go
type ConsTime time.Time
```

ConsTime is the console current time.

#### func (*ConsTime) FromPacket

```go
func (ct *ConsTime) FromPacket(p Packet) (err error)
```
FromPacket unpacks the data from an 8-byte GETTIME response packet into a
console timestamp.

#### func (ConsTime) ToPacket

```go
func (ct ConsTime) ToPacket() (p Packet)
```
ToPacket packs the console timestamp into an 8-byte packet suitable for the
SETTIME command.

#### type Device

```go
type Device interface {
	io.ReadWriteCloser
	Flush() error
	ReadFull(buf []byte) (int, error)
}
```

Device is an interface for the protocol to use to perform basic I/O operations
with different Weatherlink devices.

#### type Dmp

```go
type Dmp [5]Archive
```

Dmp is a revision B DMP archive page consisting of 5 archive records.

#### func (*Dmp) FromPacket

```go
func (d *Dmp) FromPacket(p Packet) (err error)
```
FromPacket unpacks the data from a 267-byte DMP revision B archive page packet
into a Dmp array of 5 Archive records.

#### type DmpAft

```go
type DmpAft time.Time
```

DmpAft is a timestamp appropriate for the "DMP after" command.

#### func (DmpAft) ToPacket

```go
func (da DmpAft) ToPacket() (p Packet)
```
ToPacket packs the data from the DmpAft struct into a 6-byte packet appropriate
for use with the DMPAFT command.

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

#### func (*DmpMeta) FromPacket

```go
func (dm *DmpMeta) FromPacket(p Packet) (err error)
```
FromPacket unpacks the data from a 6-byte DMP metadata packet.

#### type IP

```go
type IP struct {
	Timeout time.Duration
}
```

IP represents a Weatherlink IP.

#### func  DialIP

```go
func DialIP(address string, timeout ...time.Duration) (i IP, err error)
```
DialIP establishes a TCP connection with a Weatherlink IP.

#### func (IP) Close

```go
func (i IP) Close() error
```
Close closes the TCP connection of the Weatherlink IP.

#### func (IP) Flush

```go
func (i IP) Flush() error
```
Flush flushes the input buffers of the Weatherlink IP.

#### func (IP) Read

```go
func (i IP) Read(b []byte) (int, error)
```
Read reads up to the size of the provided byte buffer from the Weatherlink IP.
It blocks until at least one byte is read or the timeout triggers. In practice,
exactly how much it reads beyond one byte seems unpredictable.

#### func (IP) ReadFull

```go
func (i IP) ReadFull(b []byte) (int, error)
```
ReadFull reads the full size of the provided byte buffer from the Weatherlink
IP. It blocks until the entire buffer is filled or the timeout triggers.

#### func (IP) Write

```go
func (i IP) Write(b []byte) (int, error)
```
Write writes the byte buffer to the Weatherlink IP.

#### type Loop

```go
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
}
```

Loop is a combined struct representation of the union of loop1 and loop2
packets. They have a lot of overlap but the precision is sometimes different and
they complement each other.

During the protocol loop polling with the LPS command the two versions are
interleaved.

#### func (*Loop) FromPacket

```go
func (l *Loop) FromPacket(p Packet) (err error)
```
FromPacket unpacks the data from a 99-byte loop 1 or 2 packet into the Loop
struct.

#### func (*Loop) ToPacket

```go
func (l *Loop) ToPacket(t int) (p Packet, err error)
```
ToPacket packs the data from the Loop struct into a 99-byte loop 1 or two
packet.

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
	TransStatus    int     `json:"transmitterStatus"`
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
	Accum          LoopRainAccum `json:"accumulation"`
	Rate           float64       `json:"rate"`
	StormStartDate time.Time     `json:"stormStartDate"`
}
```

LoopRain is the rain sensor related readings for a Loop struct.

#### type LoopRainAccum

```go
type LoopRainAccum struct {
	Last15Min   float64 `json:"last15Minutes"`
	LastHour    float64 `json:"lastHour"`
	Last24Hours float64 `json:"last24Hours"`
	Today       float64 `json:"today"`
	LastMonth   float64 `json:"lastMonth"`
	LastYear    float64 `json:"lastYear"`
	Storm       float64 `json:"storm"`
}
```

LoopRainAccum is the rain accumulation related readings for a LoopRain struct.

#### type LoopWind

```go
type LoopWind struct {
	Avg  LoopWindAvgs  `json:"average"`
	Cur  LoopWindCur   `json:"current"`
	Gust LoopWindGusts `json:"gust"`
}
```

LoopWind is the wind related readings for a Loop struct.

#### type LoopWindAvgs

```go
type LoopWindAvgs struct {
	Last2MinSpeed  float64 `json:"last2MinutesSpeed"`
	Last10MinSpeed float64 `json:"last10MinutesSpeed"`
}
```

LoopWindAvgs is the average wind speed related readings for a LoopWind struct.

#### type LoopWindCur

```go
type LoopWindCur struct {
	Dir   int `json:"direction"`
	Speed int `json:"speed"`
}
```

LoopWindCur is the current wind direction and speed for a LoopWind struct.

#### type LoopWindGusts

```go
type LoopWindGusts struct {
	Last10MinDir   int     `json:"last10MinutesDirection"`
	Last10MinSpeed float64 `json:"last10MinutesSpeed"`
}
```

LoopWindGusts is the wind gust related readings for a LoopWind struct.

#### type Packet

```go
type Packet []byte
```

Packet is a binary data packet.

#### type Serial

```go
type Serial struct {
	*term.Term
}
```

Serial represents a Weatherlink serial or USB device.

#### func  DialSerial

```go
func DialSerial(dev string, timeout ...time.Duration) (s Serial, err error)
```
DialSerial establishes a serial port connection with a Weatherlink device.

#### func (Serial) ReadFull

```go
func (s Serial) ReadFull(b []byte) (int, error)
```
ReadFull reads the full size of the provided byte buffer from the Weatherlink
device. It blocks until the entire buffer is filled or the timeout triggers.

#### type Weatherlink

```go
type Weatherlink struct {
	Archive     chan Archive
	LastDmpTime time.Time
	Loops       chan Loop
	CmdQ        chan cmd
}
```

Weatherlink is used to track the Weatherlink device.

#### func  Dial

```go
func Dial(dev string) (w Weatherlink, err error)
```
Dial opens the connection to the Weatherlink.

#### func (*Weatherlink) Close

```go
func (w *Weatherlink) Close() error
```
Close closes the connection to the Weatherlink.

#### func (*Weatherlink) Start

```go
func (w *Weatherlink) Start() (err error)
```
Start starts the command broker. It attempts to intelligently select what
commands should be run but also accepts commands via the CmdQ channel. The
channel is especially useful for building multiplexing services.

#### func (Weatherlink) Stop

```go
func (w Weatherlink) Stop()
```
Stop stops the command broker.
