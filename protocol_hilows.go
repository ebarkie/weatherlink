// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

// GetHiLows retrieves the record high and lows.
func (c Conn) GetHiLows(ec chan<- interface{}) error {
	p, err := c.writeCmd([]byte("HILOWS\n"), []byte{ack}, 438)
	if err != nil {
		return err
	}

	hl := HiLows{}
	err = hl.FromPacket(p)
	if err != nil {
		return err
	}

	ec <- hl

	return nil
}
