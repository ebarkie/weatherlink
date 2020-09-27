// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package data

// Packet coding logic for EEPROM packets.
//
// Refer to Vantage Pro™, Vantage Pro2™ and Vantage Vue™ Serial
// Communication Reference Manual, section XIII. EEPROM
// configuration settings.

import (
	"time"

	"github.com/ebarkie/weatherlink/packet"
	"github.com/ebarkie/weatherlink/units"
)

// EEPROM represents the configuration settings.
type EEPROM struct {
	ArchivePeriod int           `json:"archivePeriod"`
	Elev          int           `json:"elevation"`
	Lat           float64       `json:"latitude"`
	Lon           float64       `json:"longitude"`
	TimeOffset    time.Duration `json:"timeOffset"`
}

// UnmarshalBinary decodes a 4096-byte EEPROM packet into the
// EEPROM struct.
func (ee *EEPROM) UnmarshalBinary(p []byte) error {
	if packet.Crc(p) != 0 {
		return ErrBadCRC
	}

	// Setup bit breakdown:
	//
	// Bit  7     | 6     | 5      4      | 3         | 2             | 1        | 0
	//     -------+-------+---------------+-----------+---------------+----------+-----------
	//      Lon   | Lat   | Rain Coll     | Wind Cup  | Month/Day     | Is AM/PM | Time mode
	//     -------+-------+---------------+-----------+---------------+----------+-----------
	//      0 = W | 0 = S | 0 = 0.01in    | 0 = Small | 0 = Month/Day | 0 = PM   | 0 = AM/PM
	//      1 = E | 1 = N | 1 = 0.2mm     | 1 = Large | 1 = Day/Month | 1 = AM   | 1 = 24hr
	//            |       | 2 = 0.1mm     |
	setup := packet.GetUInt8(p, 43)

	// Unit bit breakdown:
	//
	// Bit  7    6    | 5      | 4         | 3      2      | 1    0
	//     -----------+--------+-----------+---------------+-----------
	//      Wind      | Rain   | Elevation | Temperature   | Barometer
	//     -----------+--------+-----------+---------------+-----------
	//      0 = mph   | 0 = in | 0 = ft    | 0 = F (whole) | 0 = in
	//      1 = m/s   | 1 = mm | 1 = m     | 1 = F (tenth) | 1 = mm
	//      2 = km/h  |        |           | 2 = C (whole) | 2 = hpa
	//      3 = knots |        |           | 3 = C (tenth) | 3 = mb
	unit := packet.GetUInt8(p, 41)

	ee.ArchivePeriod = packet.GetUInt8(p, 45)

	// Location
	ee.Elev = packet.GetUInt16(p, 15)
	if ft := unit&0x10 == 0; !ft {
		// Elevation is in meters so convert to feet
		ee.Elev = int(units.Length(float64(ee.Elev) * units.Meters).Feet())
	}
	ee.Lat = packet.GetFloat16_10(p, 11)
	if north := setup&0x40 != 0; (north && ee.Lat < 0.0) || (!north && ee.Lat > 0.0) {
		// Equator hemisphere setting and latitude do not agree
		return ErrBadLocation
	}
	ee.Lon = packet.GetFloat16_10(p, 13)
	if east := setup&0x80 != 0; (east && ee.Lon < 0.0) || (!east && ee.Lon > 0.0) {
		// Prime meridian hemisphere setting and longitude do not agree
		return ErrBadLocation
	}

	ee.TimeOffset = time.Duration(packet.GetFloat16(p, 20)/100.0) * time.Hour

	return nil
}
