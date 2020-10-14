// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

import (
	"time"

	"github.com/ebarkie/weatherlink/data"
)

// GetConsTime gets the console time.
func (c Conn) GetConsTime() (time.Time, error) {
	p, err := c.writeCmd([]byte("GETTIME\n"), []byte{ack}, 8)
	if err != nil {
		return time.Time{}, err
	}

	var ct data.ConsTime
	err = ct.UnmarshalBinary(p)
	return time.Time(ct), err
}

// setConsTime sets the console time.
func (c Conn) setConsTime(t time.Time) (err error) {
	_, err = c.writeCmd([]byte("SETTIME\n"), []byte{ack}, 0)
	if err != nil {
		return
	}

	var p []byte
	p, err = data.ConsTime(t).MarshalBinary()
	if err != nil {
		return
	}

	_, err = c.writeCmd(p, []byte{ack}, 0)

	return
}

// SyncConsTime synchronizes the console time with the local
// system time if the offset exceeds 10 seconds.
func (c Conn) SyncConsTime() error {
	const maxOffset = 10 * time.Second

	t, err := c.GetConsTime()
	if err != nil {
		return err
	}

	offset := time.Since(t)
	if offset < 0 {
		offset *= -1
	}
	Debug.Printf("Console time is %s, offset is %s", t, offset)

	if offset > maxOffset {
		Info.Printf("Console time is off by %s, syncing", offset)
		err = c.setConsTime(time.Now())
		if err != nil {
			Error.Println(err.Error())
		}
	}

	return err
}

// SetLamps sets the console lamps state.
func (c Conn) SetLamps(on bool) (err error) {
	state := "0"
	if on {
		state = "1"
	}
	_, err = c.writeCmd([]byte("LAMPS "+state+"\n"), []byte("\n\rOK\n\r"), 0)

	return
}
