// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

import "time"

// getConsTime reads the console time and returns it as a time.
func (w *Weatherlink) getConsTime() (t time.Time, err error) {
	var p Packet
	p, err = w.sendCommand([]byte("GETTIME\n"), 8)
	if err != nil {
		return
	}

	ct := ConsTime{}
	err = ct.FromPacket(p)
	if err != nil {
		return
	}

	t = time.Time(ct)
	return
}

// setConsTime sets the console time to the specified time.
func (w *Weatherlink) setConsTime(t time.Time) (err error) {
	_, err = w.sendCommand([]byte("SETTIME\n"), 0)
	if err != nil {
		return
	}
	_, err = w.sendCommand(ConsTime(t).ToPacket(), 0)

	return
}

// syncConsTime uses getConsTime and compares it to the system time.
// If it exceeds the offset threshold then setConsTime is called to
// get them in sync.
func (w *Weatherlink) syncConsTime() (err error) {
	const maxOffset = 10 * time.Second

	var t time.Time
	t, err = w.getConsTime()
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
		err = w.setConsTime(time.Now())
		if err != nil {
			Error.Println(err.Error())
		}
	}

	return
}
