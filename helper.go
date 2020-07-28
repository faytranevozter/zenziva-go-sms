package zenziva

func inArray(s string, haystack []string) bool {
	for _, v := range haystack {
		if v == s {
			return true
		}
	}
	return false
}
