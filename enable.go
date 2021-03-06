package gse

// Enable enables an extension by UUID.
func Enable(uuid string) error {
	enabled, err := Enabled()
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
		err = SetEnabled(enabled)
		if err != nil {
			return err
		}
	}

	return nil
}
