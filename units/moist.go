// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package units

type SoilType uint

const (
	Sand SoilType = iota
	SandyLoam
	Loam
	ClayLoam
)

// SoilMoisture converts soil moisture tension in centibars to a percentage.
func SoilMoisture(t SoilType, cb int) int {
	// This uses a simple linear scale based on depletion of plant
	// available water for each soil type.

	// Scheduling Irrigations: When and How Much Water to Apply.
	// Division of Agriculture and Natural Resources Publication 3396.
	// University of California Irrigation Program.
	// University of California, Davis. pp. 106.
	var depletedCbs = []int{
		30,  // Sand/Loamy Sand
		50,  // Sandy Loam
		130, // Loam/Silt Loam
		170, // Clay Loam/Clay
	}

	dcb := depletedCbs[t]
	if cb > dcb {
		return 0
	}
	return 100 - int((float32(cb)/float32(dcb))*100)
}
