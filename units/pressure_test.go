// Copyright 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package units

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMb(t *testing.T) {
	assert.Equal(t, 846.593815, FromMercuryIn(25.0).Mb(), "Inches to Millibars")
}
