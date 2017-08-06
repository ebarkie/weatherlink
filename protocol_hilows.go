// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

// getHiLows retrieves the record high and lows.
func (w *Weatherlink) getHiLows(ec chan interface{}) error {
	p, err := w.sendCommand([]byte("HILOWS\n"), 438)
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
