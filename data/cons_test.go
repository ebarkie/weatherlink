// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package data

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testConsTimePackets = map[string][]byte{
	"std": {
		0x02, 0x2c, 0x0f, 0x1e, 0x06, 0x74, 0x10, 0xe6,
	},
}

func TestConsTimeMarshalBinary(t *testing.T) {
	a := assert.New(t)

	ct := ConsTime(time.Date(2016, time.June, 30, 15, 44, 2, 0, time.Local))
	p, err := ct.MarshalBinary()
	a.Nil(err, "MarshalBinary ConsTime")

	a.Equal(testConsTimePackets["std"], p, "Console time")
}

func BenchmarkConsTimeUnmarshalBinary(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ct := ConsTime{}
		ct.UnmarshalBinary(testConsTimePackets["std"])
	}
}

func TestConsTimeUnmarshalBinary(t *testing.T) {
	a := assert.New(t)

	ct := ConsTime{}
	err := ct.UnmarshalBinary(testConsTimePackets["std"])
	a.Nil(err, "UnmarshalBinary ConsTime")

	a.Equal(time.Date(2016, time.June, 30, 15, 44, 2, 0, time.Local),
		time.Time(ct), "Console time")
}
