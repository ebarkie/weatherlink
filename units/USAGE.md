# units

```go
import "github.com/ebarkie/weatherlink/units"
```

Package units implements simple unit conversion functions.

## Usage

#### func  C

```go
func C(f float64) float64
```
C converts Celsius to Fahrenheit.

#### func  F

```go
func F(c float64) float64
```
F converts Fahrenheit to Celsius.

#### func  Ft

```go
func Ft(m float64) float64
```
Ft converts Meters to Feet.

#### func  Kn

```go
func Kn(mph float64) float64
```
Kn converts Miles Per Hour (MPH) to Knots.
