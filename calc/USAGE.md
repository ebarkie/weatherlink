# calc

```go
import "github.com/ebarkie/weatherlink/calc"
```

Package calc implements weather calculations.

## Usage

#### func  DewPoint

```go
func DewPoint(tf float64, h int) float64
```
DewPoint takes a temperature in Fahrenheit and humidity and returns the dew
point in Fahrenheit. It uses Magnus-Tetens formula.
