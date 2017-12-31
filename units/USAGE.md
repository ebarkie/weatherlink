# units
```
import "github.com/ebarkie/weatherlink/units"
```

Package units implements simple unit conversion functions.

## Usage

#### type Length

```go
type Length struct {
}
```

Length is a length stored in Inches.

#### func  FromFt

```go
func FromFt(ft float64) Length
```
FromFt returns a length from a value in Feet.

#### func  FromIn

```go
func FromIn(in float64) Length
```
FromIn returns a length from a value in Inches.

#### func  FromM

```go
func FromM(m float64) Length
```
FromM returns a length from a value in Meters.

#### func (Length) Ft

```go
func (l Length) Ft() float64
```
Ft returns the length in feet.

#### func (Length) In

```go
func (l Length) In() float64
```
In returns the length in inches.

#### func (Length) M

```go
func (l Length) M() float64
```
M returns the length in meters.

#### func (Length) Mm

```go
func (l Length) Mm() float64
```
Mm returns the length in Millimeters.

#### type Moisture

```go
type Moisture struct {
}
```

Moisture is moisture in centibars of tension.

#### func  FromCB

```go
func FromCB(cb int) Moisture
```
FromCB returns a moisure level stored in centibars of tension.

#### func (Moisture) P

```go
func (m Moisture) P(t SoilType) int
```
P returns the soil moisure as a percentage.

#### type Pressure

```go
type Pressure struct {
}
```

Pressure is a barometric pressure stored in Inches.

#### func  FromMercuryIn

```go
func FromMercuryIn(in float64) Pressure
```
FromMercuryIn returns a pressure from a value in Inches.

#### func (Pressure) Hpa

```go
func (p Pressure) Hpa() float64
```
Hpa returns the pressure in Hectopascals.

#### func (Pressure) In

```go
func (p Pressure) In() float64
```
In returns the pressure in Inches.

#### func (Pressure) Mb

```go
func (p Pressure) Mb() float64
```
Mb returns the pressure in Millibars.

#### type SoilType

```go
type SoilType uint
```

SoilType is the soil type used for calculating suction.

```go
const (
	Sand SoilType = iota
	SandyLoam
	Loam
	Clay
)
```
Soil types ranging from sand to clay.

#### type Speed

```go
type Speed struct {
}
```

Speed is a speed.

#### func  FromMPH

```go
func FromMPH(mph float64) Speed
```
FromMPH returns a speed from a value in Miles Per Hour.

#### func (Speed) Kn

```go
func (s Speed) Kn() float64
```
Kn returns the speed in Knots.

#### func (Speed) MPH

```go
func (s Speed) MPH() float64
```
MPH returns the speed in Miles per Hour.

#### func (Speed) MPS

```go
func (s Speed) MPS() float64
```
MPS returns the speed in Meters per Second.

#### type Temperature

```go
type Temperature struct {
}
```

Temperature is a temperature stored in Fahrenheit.

#### func  FromC

```go
func FromC(c float64) Temperature
```
FromC returns a temperature from a value in Celsius.

#### func  FromF

```go
func FromF(f float64) Temperature
```
FromF returns a temperature from a value in Fahrenheit.

#### func (Temperature) C

```go
func (t Temperature) C() float64
```
C returns the temperature in Celsius.

#### func (Temperature) F

```go
func (t Temperature) F() float64
```
F returns the temperature in Fahrenheit.
