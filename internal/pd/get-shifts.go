// Copyright © 2020 The Homeport Team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package pd

import (
	"strconv"
	"time"
)

// Requirement: Shifts must be sorted in .pd.yml file

// Shift specifies time range and name of a shift, start and end times are saved in minutes
type Shift struct {
	Start int
	End   int
	Name  string
}

// GetNextShift returns information about the next shift and time until it starts
func GetNextShift(shifts []Shift, shiftPos int) (Shift, int, error) {

	timeInUTC := time.Now().UTC()
	currentTime := timeInUTC.Hour()*60 + timeInUTC.Minute()

	nextShift := shifts[(shiftPos+1)%len(shifts)]

	timeUntilNextShift := nextShift.Start - currentTime
	if currentTime > nextShift.Start { // prevents timeUntilNextShift from being negative if the next shift starts during the next day
		timeUntilNextShift += 24 * 60
	}

	return nextShift, timeUntilNextShift, nil
}

// GetCurrentShift returns all shifts in a slice and the position of the current shift
func GetCurrentShift() ([]Shift, int, error) {

	timeInUTC := time.Now().UTC()
	currentTime := timeInUTC.Hour()*60 + timeInUTC.Minute()

	shifts, err := LoadShifts()
	if err != nil {
		return []Shift{}, 0, err
	}

	shiftPos := -1
	for i, shift := range shifts {
		if shift.Start < shift.End { // shift starts and ends during the same day
			if currentTime >= shift.Start && currentTime < shift.End {
				shiftPos = i
			}
		} else {
			if currentTime >= shift.Start || currentTime < shift.End {
				shiftPos = i
			}
		}
	}

	return shifts, shiftPos, nil
}

// LoadShifts loads shifts out of the .pd.yml file
func LoadShifts() ([]Shift, error) {

	var err error
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}

	finalShifts := make([]Shift, len(config.ShiftTimes))
	var min int
	for i, shift := range config.ShiftTimes {
		finalShifts[i] = Shift{}
		finalShifts[i].Start, err = strconv.Atoi(shift.Start[:2])
		if err != nil {
			return nil, err
		}
		finalShifts[i].Start *= 60
		min, err = strconv.Atoi(shift.Start[3:])
		if err != nil {
			return nil, err
		}

		finalShifts[i].End, err = strconv.Atoi(shift.End[:2])
		if err != nil {
			return nil, err
		}
		finalShifts[i].End *= 60
		min, err = strconv.Atoi(shift.End[3:])
		if err != nil {
			return nil, err
		}
		finalShifts[i].End += min
		finalShifts[i].Name = shift.Name
	}

	return finalShifts, nil
}