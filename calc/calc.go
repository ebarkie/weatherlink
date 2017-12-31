// Copyright 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// Package calc implements weather calculations.
package calc

import (
	"math"

	"github.com/ebarkie/weatherlink/units"
)

// DewPoint takes a temperature in Fahrenheit and humidity and
// returns the dew point in Fahrenheit.  It uses Magnus-Tetens
// formula.
func DewPoint(tf float64, h int) float64 {
	const (
		a = 17.27
		b = 237.7
	)

	tc := units.FromF(tf).C()
	x := a*tc/(b+tc) + math.Log(float64(h)/100.0)
	dpc := b * x / (a - x)

	return units.FromC(dpc).F()
}
