package gse

// Disable disables an extension by UUID.
func Disable(uuid string) error {
	enabled, err := Enabled()
	if err != nil {
		return err
	}

	modified := false
	for i := len(enabled) - 1; i >= 0; i-- {
		if enabled[i] == uuid {
			enabled = append(enabled[:i], enabled[i+1:]...)
			modified = true
		}
	}

	if modified {
		err = setEnabledExtensions(enabled)
		if err != nil {
			return err
		}
	}

	return nil
}
