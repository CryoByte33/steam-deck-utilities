// CryoUtilities
// Copyright (C) 2023 CryoByte33 and contributors to the CryoUtilities project

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package internal

func isTweakEnabled(unit string) bool {
	status, err := getUnitStatus(unit)
	if err != nil {
		CryoUtils.ErrorLog.Println("Unable to get current", unit, "value")
		return false
	}
	return status == TweakList[unit].Recommended
}

func EnableTweak(unit string) error {
	CryoUtils.InfoLog.Println("Enabling", unit, "...")
	err := setUnitValue(unit, TweakList[unit].Recommended)
	if err != nil {
		CryoUtils.ErrorLog.Println("Unable to set", unit, "to", TweakList[unit].Recommended)
		return err
	}
	err = writeUnitFile(unit, TweakList[unit].Recommended)
	if err != nil {
		CryoUtils.ErrorLog.Println("Unable to write unit file for", unit)
		return err
	}
	return nil
}

func RevertTweak(unit string) error {
	CryoUtils.InfoLog.Println("Disabling", unit, "...")
	err := setUnitValue(unit, TweakList[unit].Default)
	if err != nil {
		CryoUtils.ErrorLog.Println("Unable to revert", unit, "to", TweakList[unit].Default)
		return err
	}
	err = removeUnitFile(unit)
	if err != nil {
		CryoUtils.ErrorLog.Println("Unable to remove unit file for", unit)
		return err
	}
	return nil
}

func ToggleTweak(unit string) error {
	if isTweakEnabled(unit) {
		return RevertTweak(unit)
	}
	return EnableTweak(unit)
}
