// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package units

// Temperature is a temperature stored in Fahrenheit.
type Temperature struct {
	f float64
}

// FromC returns a temperature from a value in Celsius.
func FromC(c float64) Temperature {
	return Temperature{f: c*1.8 + 32.0}
}

// FromF returns a temperature from a value in Fahrenheit.
func FromF(f float64) Temperature {
	return Temperature{f: f}
}

// C returns the temperature in Celsius.
func (t Temperature) C() float64 {
	return (t.f - 32.0) * 5.0 / 9.0
}

// F returns the temperature in Fahrenheit.
func (t Temperature) F() float64 {
	return t.f
}
