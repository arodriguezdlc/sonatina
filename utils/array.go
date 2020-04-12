package utils

// FindString takes a slice and looks for a string in it. If found it will
// return it's key, otherwise it will return nil and a bool of false.
func FindString(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// RemoveDuplicatedStrings returns a string slice without duplicated elements.
// Order is NOT preserved
func RemoveDuplicatedStrings(slice []string) []string {
	result := []string{}

	check := make(map[string]bool)
	for _, element := range slice {
		check[element] = true
	}

	for key := range check {
		result = append(result, key)
	}

	return result
}
