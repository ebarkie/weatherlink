// Copyright 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package units

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSoilMoisture(t *testing.T) {
	assert.Equal(t, 100, FromCB(0).P(Loam), "Completely saturated")
	assert.Equal(t, 50, FromCB(65).P(Loam), "Half saturated")
	assert.Equal(t, 0, FromCB(130).P(Loam), "No available moisture")
}
