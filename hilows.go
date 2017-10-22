// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

import "github.com/ebarkie/weatherlink/data"

// GetHiLows retrieves the record high and lows.
func (c Conn) GetHiLows(ec chan<- interface{}) error {
	p, err := c.writeCmd([]byte("HILOWS\n"), []byte{ack}, 438)
	if err != nil {
		return err
	}

	hl := data.HiLows{}
	err = hl.UnmarshalBinary(p)
	if err != nil {
		return err
	}

	ec <- hl

	return nil
}
