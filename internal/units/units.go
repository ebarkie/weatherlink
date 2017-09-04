// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// Package units implements very simple and lightweight unit
// conversion functions.
package units

// Ft converts Meters to Feet.
func Ft(m float64) float64 {
	return m * 3.28084
}
