package main

type ComparableString string

func (cs ComparableString) Compare(c interface{}) int {
	if cs > c.(ComparableString) {
		return 1
	}
	if cs < c.(ComparableString) {
		return -1
	}
	return 0
}
