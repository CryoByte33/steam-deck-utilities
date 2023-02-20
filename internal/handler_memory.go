package internal

func getHugePagesStatus() bool {
	status, err := getUnitStatus("hugepages")
	if err != nil {
		CryoUtils.ErrorLog.Println("Unable to get current hugepages value")
		return false
	}
	if status == RecommendedHugePages {
		return true
	}
	return false
}

func getCompactionProactivenessStatus() bool {
	status, err := getUnitStatus("compaction_proactiveness")
	if err != nil {
		CryoUtils.ErrorLog.Println("Unable to get current compaction_proactiveness")
		return false
	}
	if status == RecommendedCompactionProactiveness {
		return true
	}
	return false
}

func getPageLockUnfairnessStatus() bool {
	status, err := getUnitStatus("page_lock_unfairness")
	if err != nil {
		CryoUtils.ErrorLog.Println("Unable to get current page_lock_unfairness")
		return false
	}
	if status == RecommendedPageLockUnfairness {
		return true
	}
	return false
}

func getShMemStatus() bool {
	status, err := getUnitStatus("shmem_enabled")
	if err != nil {
		CryoUtils.ErrorLog.Println("Unable to get current shmem_enabled")
		return false
	}
	if status == RecommendedShMem {
		return true
	}
	return false
}

func getDefragStatus() bool {
	status, err := getUnitStatus("defrag")
	if err != nil {
		CryoUtils.ErrorLog.Println("Unable to get current defrag")
		return false
	}
	if status == RecommendedHugePageDefrag {
		return true
	}
	return false
}

// ToggleHugePages Simple one-function toggle for the button to use
func ToggleHugePages() error {
	if getHugePagesStatus() {
		err := RevertHugePages()
		if err != nil {
			return err
		}
	} else {
		err := SetHugePages()
		if err != nil {
			return err
		}
	}
	return nil
}

// ToggleShMem Simple one-function toggle for the button to use
func ToggleShMem() error {
	if getShMemStatus() {
		err := RevertShMem()
		if err != nil {
			return err
		}
	} else {
		err := SetShMem()
		if err != nil {
			return err
		}
	}
	return nil
}

// ToggleCompactionProactiveness Simple one-function toggle for the button to use
func ToggleCompactionProactiveness() error {
	if getCompactionProactivenessStatus() {
		err := RevertCompactionProactiveness()
		if err != nil {
			return err
		}
	} else {
		err := SetCompactionProactiveness()
		if err != nil {
			return err
		}
	}
	return nil
}

// ToggleDefrag Simple one-function toggle for the button to use
func ToggleDefrag() error {
	if getDefragStatus() {
		err := RevertDefrag()
		if err != nil {
			return err
		}
	} else {
		err := SetDefrag()
		if err != nil {
			return err
		}
	}
	return nil
}

// TogglePageLockUnfairness Simple one-function toggle for the button to use
func TogglePageLockUnfairness() error {
	if getPageLockUnfairnessStatus() {
		err := RevertPageLockUnfairness()
		if err != nil {
			return err
		}
	} else {
		err := SetPageLockUnfairness()
		if err != nil {
			return err
		}
	}
	return nil
}

func SetHugePages() error {
	CryoUtils.InfoLog.Println("Enabling hugepages...")
	// Remove a file accidentally included in a beta for testing
	_ = removeFile(NHPTestingFile)
	err := setUnitValue("hugepages", RecommendedHugePages)
	if err != nil {
		return err
	}
	err = writeUnitFile("hugepages", RecommendedHugePages)
	if err != nil {
		return err
	}
	return nil
}

func RevertHugePages() error {
	CryoUtils.InfoLog.Println("Disabling hugepages...")
	err := setUnitValue("hugepages", DefaultHugePages)
	if err != nil {
		return err
	}
	err = removeUnitFile("hugepages")
	if err != nil {
		return err
	}
	return nil
}

func SetCompactionProactiveness() error {
	CryoUtils.InfoLog.Println("Setting compaction_proactiveness...")
	err := setUnitValue("compaction_proactiveness", RecommendedCompactionProactiveness)
	if err != nil {
		return err
	}
	err = writeUnitFile("compaction_proactiveness", RecommendedCompactionProactiveness)
	if err != nil {
		return err
	}
	return nil
}

func RevertCompactionProactiveness() error {
	CryoUtils.InfoLog.Println("Disabling compaction_proactiveness...")
	err := setUnitValue("compaction_proactiveness", DefaultCompactionProactiveness)
	if err != nil {
		return err
	}
	err = removeUnitFile("compaction_proactiveness")
	if err != nil {
		return err
	}
	return nil
}

func SetPageLockUnfairness() error {
	CryoUtils.InfoLog.Println("Enabling page_lock_unfairness...")
	err := setUnitValue("page_lock_unfairness", RecommendedPageLockUnfairness)
	if err != nil {
		return err
	}
	err = writeUnitFile("page_lock_unfairness", RecommendedPageLockUnfairness)
	if err != nil {
		return err
	}
	return nil
}

func RevertPageLockUnfairness() error {
	CryoUtils.InfoLog.Println("Disabling page_lock_unfairness...")
	err := setUnitValue("page_lock_unfairness", DefaultPageLockUnfairness)
	if err != nil {
		return err
	}
	err = removeUnitFile("page_lock_unfairness")
	if err != nil {
		return err
	}
	return nil
}

func SetShMem() error {
	CryoUtils.InfoLog.Println("Enabling shmem_enabled...")
	err := setUnitValue("shmem_enabled", RecommendedShMem)
	if err != nil {
		return err
	}
	err = writeUnitFile("shmem_enabled", RecommendedShMem)
	if err != nil {
		return err
	}
	return nil
}

func RevertShMem() error {
	CryoUtils.InfoLog.Println("Disabling shmem_enabled...")
	err := setUnitValue("shmem_enabled", DefaultShMem)
	if err != nil {
		return err
	}
	err = removeUnitFile("shmem_enabled")
	if err != nil {
		return err
	}
	return nil
}

func SetDefrag() error {
	CryoUtils.InfoLog.Println("Enabling shmem_enabled...")
	err := setUnitValue("defrag", RecommendedHugePageDefrag)
	if err != nil {
		return err
	}
	err = writeUnitFile("defrag", RecommendedHugePageDefrag)
	if err != nil {
		return err
	}
	return nil
}

func RevertDefrag() error {
	CryoUtils.InfoLog.Println("Disabling shmem_enabled...")
	err := setUnitValue("defrag", DefaultHugePageDefrag)
	if err != nil {
		return err
	}
	err = removeUnitFile("defrag")
	if err != nil {
		return err
	}
	return nil
}
