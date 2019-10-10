// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package data

import (
	"testing"

	"gitlab.com/ebarkie/weatherlink/packet"

	"github.com/stretchr/testify/assert"
)

func BenchmarkCRC(b *testing.B) {
	for n := 0; n < b.N; n++ {
		packet.Crc(testLoopPackets["1NoRain"])
	}
}

func TestCRC(t *testing.T) {
	a := assert.New(t)
	a.Zero(packet.Crc(testLoopPackets["1Rain"]), "Loop1 CRC check")
	a.NotZero(packet.Crc(testLoopPackets["2BadCrc"]), "Loop1 bad CRC check")
}
