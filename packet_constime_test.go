// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testConsTimePackets = map[string]Packet{
	"standard": {
		0x02, 0x2c, 0x0f, 0x1e, 0x06, 0x74, 0x10, 0xe6,
	},
}

func BenchmarkConsTimeFromPacket(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ct := ConsTime{}
		ct.FromPacket(testConsTimePackets["standard"])
	}
}

func TestConsTimeFromPacket(t *testing.T) {
	ct := ConsTime{}
	ct.FromPacket(testConsTimePackets["standard"])

	a := assert.New(t)
	a.Equal(time.Date(2016, time.June, 30, 15, 44, 2, 0, time.Local),
		time.Time(ct), "Console time")
}

func TestConsTimeToPacket(t *testing.T) {
	ct := ConsTime(time.Date(2016, time.June, 30, 15, 44, 2, 0, time.Local))

	a := assert.New(t)
	a.Equal(testConsTimePackets["standard"], ct.ToPacket(), "Console time")
}
