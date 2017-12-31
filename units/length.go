// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package units

// Length is a length stored in Inches.
type Length struct {
	in float64
}

// FromFt returns a length from a value in Feet.
func FromFt(ft float64) Length {
	return Length{in: ft * 12.0}
}

// FromIn returns a length from a value in Inches.
func FromIn(in float64) Length {
	return Length{in: in}
}

// FromM returns a length from a value in Meters.
func FromM(m float64) Length {
	return Length{in: m * 39.37008}
}

// Ft returns the length in feet.
func (l Length) Ft() float64 {
	return l.in / 12.0
}

// In returns the length in inches.
func (l Length) In() float64 {
	return l.in
}

// M returns the length in meters.
func (l Length) M() float64 {
	return l.in * 0.025399999187200026
}

// Mm returns the length in Millimeters.
func (l Length) Mm() float64 {
	return l.M() * 1000.0
}
