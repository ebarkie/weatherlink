// Copyright 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package units

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKnots(t *testing.T) {
	assert.Equal(t, 47.784, Speed(55.0*MPH).Knots(), "Miles Per Hour (MPH) to Knots")
}
