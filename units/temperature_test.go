// Copyright 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package units

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCelsius(t *testing.T) {
	assert.Equal(t, 0.0, Fahrenheit(32.0).Celsius(), "Fahrenheit to Celsius")
}

func TestFahrenheit(t *testing.T) {
	assert.Equal(t, 32.0, Celsius(0.0).Fahrenheit(), "Celsius to Fahrenheit")
}
