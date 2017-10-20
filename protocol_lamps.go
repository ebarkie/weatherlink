// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

// SetLamps sets the console lamps state.
func (c Conn) SetLamps(on bool) (err error) {
	state := "0"
	if on {
		state = "1"
	}
	_, err = c.writeCmd([]byte("LAMPS "+state+"\n"), []byte{lf, cr, 'O', 'K', lf, cr}, 0)

	return
}
