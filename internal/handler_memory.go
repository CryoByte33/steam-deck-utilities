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
