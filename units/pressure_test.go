// Copyright 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package units

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMillibars(t *testing.T) {
	assert.Equal(t, 846.593815, Pressure(25.0*Inches).Millibars(), "Inches to Millibars")
}
