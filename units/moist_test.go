// Copyright 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package units

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSoilMoisture(t *testing.T) {
	assert.Equal(t, 100, SoilMoisture(Loam, 0), "Completely saturated")
	assert.Equal(t, 50, SoilMoisture(Loam, 65), "Half saturated")
	assert.Equal(t, 0, SoilMoisture(Loam, 130), "No available moisture")
}
