// Copyright 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package units

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestC(t *testing.T) {
	assert.Equal(t, 0.0, C(32.0), "Fahrenheit to Celsius")
}

func TestF(t *testing.T) {
	assert.Equal(t, 32.0, F(0.0), "Celsius to Fahrenheit")
}
