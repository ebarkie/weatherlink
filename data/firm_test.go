// Copyright (c) 2020 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package data

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFirmTimeUnmarshalText(t *testing.T) {
	a := assert.New(t)

	ft := FirmTime{}
	err := ft.UnmarshalText([]byte("Apr 24 2002\n\r"))
	a.Nil(err, "UnmarshalText FirmTime")

	a.Equal(time.Date(2002, time.April, 24, 0, 0, 0, 0, time.UTC),
		time.Time(ft), "Firmware build time")
}

func TestFirmVerUnmarshalText(t *testing.T) {
	a := assert.New(t)

	var fv FirmVer
	err := fv.UnmarshalText([]byte("1.73\n\r"))
	a.Nil(err, "UnmarshalText FirmVer")

	a.Equal("1.73", string(fv), "Firmware version number")
}
