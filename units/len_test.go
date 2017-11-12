// Copyright 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package units

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFt(t *testing.T) {
	assert.Equal(t, 3.28084, Ft(1.0), "Meters to Feet")
}
