// Copyright 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package units

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFt(t *testing.T) {
	assert.Equal(t, 1.0, FromIn(12.0).Ft(), "Inches to Feet")
	assert.Equal(t, 0.025399999187200026, FromIn(1.0).M(), "Inches to Meters")

	assert.Equal(t, 3.28084, FromM(1.0).Ft(), "Meters to Feet")
	assert.Equal(t, 1.0, FromM(1.0).M(), "Meters to Meters")
}
