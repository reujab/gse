package gse

// Enable enables an extension by UUID.
func Enable(uuid string) error {
	enabled, err := getEnabledExtensions()
	if err != nil {
		return err
	}

	alreadyEnabled := false
	for _, extension := range enabled {
		if extension == uuid {
			alreadyEnabled = true
			break
		}
	}

	if !alreadyEnabled {
		enabled = append(enabled, uuid)
		err = setEnabledExtensions(enabled)
		if err != nil {
			return err
		}
	}

	return nil
}
