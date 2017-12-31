// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package units

// Pressure is a barometric pressure stored in Inches.
type Pressure struct {
	in float64
}

// FromMercuryIn returns a pressure from a value in Inches.
func FromMercuryIn(in float64) Pressure {
	return Pressure{in: in}
}

// Hpa returns the pressure in Hectopascals.
func (p Pressure) Hpa() float64 {
	return p.Mb()
}

// In returns the pressure in Inches.
func (p Pressure) In() float64 {
	return p.in
}

// Mb returns the pressure in Millibars.
func (p Pressure) Mb() float64 {
	return p.in * 33.8637526
}
