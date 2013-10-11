package main

import (
	"github.com/PreetamJinka/orderedlist"
)

type ComparableString string

func (cs ComparableString) Compare(c orderedlist.Comparable) int {
	if cs > c.(ComparableString) {
		return 1
	}
	if cs < c.(ComparableString) {
		return -1
	}
	return 0
}
