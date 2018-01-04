# packet

```go
import "github.com/ebarkie/weatherlink/packet"
```

Package packet implements getters and setters for Davis Instruments packets.

## Usage

```go
const (
	Dash         = "-"
	FallingRapid = "Falling Rapidly"
	FallingSlow  = "Falling Slowly"
	Steady       = "Steady"
	RisingSlow   = "Rising Slowly"
	RisingRapid  = "Rising Rapidly"
)
```
Barometer trends.

#### func  Crc

```go
func Crc(p []byte) (c uint16)
```
Crc calculates the 16-bit CRC for the Packet. When decoding the result should be
zero.

#### func  GetBarTrend

```go
func GetBarTrend(p []byte, i uint) string
```
GetBarTrend gets a barometer trend from a given packet at the specified index.

#### func  GetDate16

```go
func GetDate16(p []byte, i uint) time.Time
```
GetDate16 gets a 2-byte date (no time) value from a given packet at the
specified index.

#### func  GetDateTime32

```go
func GetDateTime32(p []byte, i uint) time.Time
```
GetDateTime32 gets a 4-byte date and time value from a given packet at the
specified index.

#### func  GetDateTime48

```go
func GetDateTime48(p []byte, i uint) time.Time
```
GetDateTime48 gets a 6-byte date and time value from a given packet at the
specified index.

#### func  GetFloat16

```go
func GetFloat16(p []byte, i uint) float64
```
GetFloat16 gets a 2-byte signed two's complement float value from a given packet
at the specified index.

#### func  GetFloat16_10

```go
func GetFloat16_10(p []byte, i uint) float64
```
GetFloat16_10 gets a 2-byte signed two's complement float value in tenths in a
given packet at the specified index.

#### func  GetForecast

```go
func GetForecast(p []byte, i uint) string
```
GetForecast gets a forecast string from a given packet at the specified index.

#### func  GetForecastIcons

```go
func GetForecastIcons(p []byte, i uint) (icons []string)
```
GetForecastIcons gets a forecast icon bit map from a given packet at the
specified index.

#### func  GetMPH16

```go
func GetMPH16(p []byte, i uint) float64
```
GetMPH16 gets a 2-byte MPH value from a given packet at the specified index.

#### func  GetMPH8

```go
func GetMPH8(p []byte, i uint) int
```
GetMPH8 gets a 1-byte MPH value from a given packet at the specified index.

#### func  GetPressure

```go
func GetPressure(p []byte, i uint) float64
```
GetPressure gets a pressure value from a given packet at the specified index.

#### func  GetRain

```go
func GetRain(p []byte, i uint) float64
```
GetRain gets a rain rate or accumulation value from a given packet at the
specified index.

#### func  GetTemp8

```go
func GetTemp8(p []byte, i uint) int
```
GetTemp8 gets a 1-byte temperature value in a given packet at the specified
index.

#### func  GetTime16

```go
func GetTime16(p []byte, i uint) time.Time
```
GetTime16 gets a 2-byte time (no date) value in a given packet at the specified
index.

#### func  GetTransStatus

```go
func GetTransStatus(p []byte, i uint) (low []int)
```
GetTransStatus gets the transmitter status from the given packet at the
specified index and returns a slice of the ID's/channels that have low battery
indicators.

#### func  GetUFloat8

```go
func GetUFloat8(p []byte, i uint) float64
```
GetUFloat8 gets a 1-byte unsigned float value from a given packet at the
specified index.

#### func  GetUInt16

```go
func GetUInt16(p []byte, i uint) int
```
GetUInt16 gets a 2-byte unsigned integer value from a given packet at the
specified index.

#### func  GetUInt8

```go
func GetUInt8(p []byte, i uint) int
```
GetUInt8 gets a 1-byte unsigned integer value from a given packet at the
specified index.

#### func  GetUVIndex

```go
func GetUVIndex(p []byte, i uint) float64
```
GetUVIndex gets a Ultraviolet index value from a given packet at the specified
index.

#### func  GetVoltage

```go
func GetVoltage(p []byte, i uint) float64
```
GetVoltage gets a battery voltage value from a given packet at the specified
index.

#### func  GetWindDir

```go
func GetWindDir(p []byte, i uint) int
```
GetWindDir gets a wind direction value in degrees from a given packet at the
specified index.

#### func  SetCrc

```go
func SetCrc(p *[]byte)
```
SetCrc sets the last 2-bytes of a given packet to the proper CRC value based on
the rest of content.

#### func  SetDateTime32

```go
func SetDateTime32(p *[]byte, i uint, t time.Time)
```
SetDateTime32 sets a 4-byte date and time value in a given packet at the
specified index.

#### func  SetDateTime48

```go
func SetDateTime48(p *[]byte, i uint, t time.Time)
```
SetDateTime48 sets a 6-byte date and time value in a given packet at the
specified index.

#### func  SetFloat16

```go
func SetFloat16(p *[]byte, i uint, v float64)
```
SetFloat16 sets a 2-byte signed two's complement float value in a given packet
at the specified index.

#### func  SetFloat16_10

```go
func SetFloat16_10(p *[]byte, i uint, v float64)
```
SetFloat16_10 sets a 2-byte signed two's complement float value in tenths in a
given packet at the specified index.

#### func  SetMPH16

```go
func SetMPH16(p *[]byte, i uint, v float64)
```
SetMPH16 sets a 2-byte MPH value in a given packet at the specified index.

#### func  SetMPH8

```go
func SetMPH8(p *[]byte, i uint, v int)
```
SetMPH8 sets a 1-byte MPH value in a given packet at the specified index.

#### func  SetPressure

```go
func SetPressure(p *[]byte, i uint, v float64)
```
SetPressure sets a pressure value in a given packet at the specified index.

#### func  SetRainClicks

```go
func SetRainClicks(p *[]byte, i uint, v float64)
```
SetRainClicks sets a rain rate or accumulation value in a given packet at the
specified index.

#### func  SetTemp8

```go
func SetTemp8(p *[]byte, i uint, v int)
```
SetTemp8 sets a 1-byte temprature value in a given packet at the specified
index.

#### func  SetUInt16

```go
func SetUInt16(p *[]byte, i uint, v int)
```
SetUInt16 sets a 2-byte unsigned integer value in a given packet at the
specified index.

#### func  SetUInt8

```go
func SetUInt8(p *[]byte, i uint, v int)
```
SetUInt8 sets a 1-byte unsigned integer value in a given packet at the specified
index.

#### func  SetVoltage

```go
func SetVoltage(p *[]byte, i uint, v float64)
```
SetVoltage sets a battery voltage value in a given packet at the specified
index.
