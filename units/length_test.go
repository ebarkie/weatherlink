// Copyright 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package units

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFeet(t *testing.T) {
	assert.Equal(t, 1.0, Length(12.0*Inches).Feet(), "Inches to Feet")
	assert.Equal(t, 0.025399999187200026, Length(1.0*Inches).Meters(), "Inches to Meters")

	assert.Equal(t, 3.28084, Length(1.0*Meters).Feet(), "Meters to Feet")
	assert.Equal(t, 1.0, Length(1.0*Meters).Meters(), "Meters to Meters")
}
