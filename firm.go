// Copyright (c) 2020 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

import (
	"github.com/ebarkie/weatherlink/data"
)

// GetFirmTime gets the firmware build time.
func (c Conn) GetFirmTime() (ft data.FirmTime, err error) {
	var p []byte
	p, err = c.writeCmd([]byte("VER\n"), []byte("\n\rOK\n\r"), 13)
	if err != nil {
		return
	}

	err = ft.UnmarshalText(p)

	return
}

// GetFirmVer gets the firmware version number.
func (c Conn) GetFirmVer() (fv data.FirmVer, err error) {
	var p []byte
	p, err = c.writeCmd([]byte("NVER\n"), []byte("\n\rOK\n\r"), 6)
	if err != nil {
		return
	}

	err = fv.UnmarshalText(p)

	return
}
