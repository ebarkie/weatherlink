// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

import "gitlab.com/ebarkie/weatherlink/data"

// GetEEPROM retrieves the entire EEPROM configuration.
func (c Conn) GetEEPROM(ec chan<- interface{}) error {
	p, err := c.writeCmd([]byte("GETEE\n"), []byte{ack}, 4098)
	if err != nil {
		return err
	}

	ee := data.EEPROM{}
	err = ee.UnmarshalBinary(p)
	if err != nil {
		return err
	}

	ec <- ee

	return nil
}
