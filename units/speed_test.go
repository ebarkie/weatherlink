// Copyright 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package units

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKn(t *testing.T) {
	assert.Equal(t, 47.784, Kn(55.0), "Miles Per Hour (MPH) to Knots")
}
