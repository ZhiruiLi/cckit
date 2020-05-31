package slices

func FindString(s []string, v string) (int, bool) {
	for i, item := range s {
		if item == v {
			return i, true
		}
	}
	return -1, false
}
