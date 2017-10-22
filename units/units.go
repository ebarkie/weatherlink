// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// Package units implements simple unit conversion functions.
package units

// C converts Fahrenheit to Celsius.
func C(f float64) float64 {
	return (f - 32.0) * 5.0 / 9.0
}

// F converts Celsius to Fahrenheit.
func F(c float64) float64 {
	return c*1.8 + 32.0
}

// Ft converts Meters to Feet.
func Ft(m float64) float64 {
	return m * 3.28084
}

// Kn converts Miles Per Hour (MPH) to Knots.
func Kn(mph float64) float64 {
	return mph * 0.8688
}
