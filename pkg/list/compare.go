package list

import "sort"

// CompareListStr compares two string lists and
// returns true if there is any difference between
// the lists
func CompareListStr(base, new []string) (hasDiff bool) {
	if len(base) != len(new) {
		return true
	}

	sort.Strings(base)
	sort.Strings(new)

	for i, b := range base {
		if b != new[i] {
			return true
		}
	}
	return false
}
