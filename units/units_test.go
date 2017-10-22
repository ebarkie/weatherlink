// Copyright 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package units

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestC(t *testing.T) {
	assert.Equal(t, 0.0, C(32.0), "Fahrenheit to Celsius")
}

func TestF(t *testing.T) {
	assert.Equal(t, 32.0, F(0.0), "Celsius to Fahrenheit")
}

func TestFt(t *testing.T) {
	assert.Equal(t, 3.28084, Ft(1.0), "Meters to Feet")
}

func TestKn(t *testing.T) {
	assert.Equal(t, 47.784, Kn(55.0), "Miles Per Hour (MPH) to Knots")
}
