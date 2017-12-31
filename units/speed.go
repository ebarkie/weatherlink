// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package units

// Speed is a speed.
type Speed struct {
	mph float64 // Miles per Hour
}

// FromMPH returns a speed from a value in Miles Per Hour.
func FromMPH(mph float64) Speed {
	return Speed{mph: mph}
}

// Kn returns the speed in Knots.
func (s Speed) Kn() float64 {
	return s.mph * 0.8688
}

// MPH returns the speed in Miles per Hour.
func (s Speed) MPH() float64 {
	return s.mph
}

// MPS returns the speed in Meters per Second.
func (s Speed) MPS() float64 {
	return s.mph * 0.44704
}
