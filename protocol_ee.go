// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

// GetEEPROM retrieves the entire EEPROM configuration.
func (c *Conn) GetEEPROM(ec chan<- interface{}) error {
	p, err := c.writeCmd([]byte("GETEE\n"), 4098)
	if err != nil {
		return err
	}

	ee := EEPROM{}
	err = ee.FromPacket(p)
	if err != nil {
		return err
	}

	ec <- ee

	return nil
}
