// Copyright (c) 2020 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

import (
	"time"

	"github.com/ebarkie/weatherlink/data"
)

// GetFirmBuildTime gets the firmware build time.
func (c Conn) GetFirmBuildTime() (time.Time, error) {
	p, err := c.writeCmd([]byte("VER\n"), []byte("\n\rOK\n\r"), 13)
	if err != nil {
		return time.Time{}, err
	}

	var ft data.FirmTime
	err = ft.UnmarshalText(p)
	return time.Time(ft), err
}

// GetFirmVer gets the firmware version number.
func (c Conn) GetFirmVer() (string, error) {
	p, err := c.writeCmd([]byte("NVER\n"), []byte("\n\rOK\n\r"), 6)
	if err != nil {
		return "", err
	}

	var fv data.FirmVer
	err = fv.UnmarshalText(p)
	return string(fv), err
}
