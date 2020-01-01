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
