// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

import (
	"time"

	"github.com/ebarkie/weatherlink/data"
)

// GetConsTime gets the console time.
func (c Conn) GetConsTime() (t time.Time, err error) {
	var p []byte
	p, err = c.writeCmd([]byte("GETTIME\n"), []byte{ack}, 8)
	if err != nil {
		return
	}

	ct := data.ConsTime{}
	err = ct.UnmarshalBinary(p)
	if err != nil {
		return
	}
	t = time.Time(ct)

	return
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
func (c Conn) SyncConsTime() (err error) {
	const maxOffset = 10 * time.Second

	var t time.Time
	t, err = c.GetConsTime()
	if err != nil {
		return
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

	return
}
