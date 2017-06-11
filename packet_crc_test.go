// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCRC(t *testing.T) {
	a := assert.New(t)
	a.Zero(crc(testLoopPackets["1Rain"]), "Loop1 CRC check")
	a.NotZero(crc(testLoopPackets["2BadCrc"]), "Loop1 bad CRC check")
}
