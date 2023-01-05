package util

// CompareSlices determines if two slices are equal. Slices are equal if they are of the same length, contain the
// same elements in the same order.
func CompareSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, _ := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// IndexOf returns the index of the element in the given slice. If it does not exist, -1 is returned.
func IndexOf(e string, slice []string) int {
	for i, el := range slice {
		if el == e {
			return i
		}
	}

	return -1
}