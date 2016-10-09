package utils

// StringInSlice is self explanatory.  Return true or false.
func StringInSlice(s string, slice []string) bool {
	for _, item := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func AllStringsInSlice(strings []string, slice []string) bool {
	for _, string := range strings {
		for i, item := range slice {
			if string == item {
				break
			}
			if i == len(slice)-1 {
				return false
			}
		}
	}
	return true
}
