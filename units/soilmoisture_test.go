// Copyright 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package units

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSoilMoisture(t *testing.T) {
	assert.Equal(t, 100, SoilMoisture(0*Centibars).Percent(Loam), "Completely saturated")
	assert.Equal(t, 50, SoilMoisture(65*Centibars).Percent(Loam), "Half saturated")
	assert.Equal(t, 0, SoilMoisture(130*Centibars).Percent(Loam), "No available moisture")
}
