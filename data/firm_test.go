// Copyright (c) 2020 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package data

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFirmTimeMarshalText(t *testing.T) {
	a := assert.New(t)

	ft := FirmTime(time.Date(2002, time.April, 24, 0, 0, 0, 0, time.UTC))
	p, err := ft.MarshalText()
	a.Nil(err, "MarshalText FirmTime")

	a.Equal([]byte("Apr 24 2002\n\r"), p, "Firmware build time")
}

func TestFirmTimeUnmarshalText(t *testing.T) {
	a := assert.New(t)

	ft := FirmTime{}
	err := ft.UnmarshalText([]byte("Apr 24 2002\n\r"))
	a.Nil(err, "UnmarshalText FirmTime")

	a.Equal(FirmTime(time.Date(2002, time.April, 24, 0, 0, 0, 0, time.UTC)),
		ft, "Firmware build time")
}

func TestFirmVerMarshalText(t *testing.T) {
	a := assert.New(t)

	fv := FirmVer("1.73")
	p, err := fv.MarshalText()
	a.Nil(err, "MarshalText FirmVer")

	a.Equal([]byte("1.73\n\r"), p, "Firmware version number")
}

func TestFirmVerUnmarshalText(t *testing.T) {
	a := assert.New(t)

	var fv FirmVer
	err := fv.UnmarshalText([]byte("1.73\n\r"))
	a.Nil(err, "UnmarshalText FirmVer")

	a.Equal(FirmVer("1.73"), fv, "Firmware version number")
}
