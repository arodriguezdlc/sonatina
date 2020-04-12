package utils

import (
	"reflect"
	"sort"
	"testing"
)

func TestFindStringWithExistingString(t *testing.T) {
	testSlice := []string{"foo", "bar", "test"}
	expectedIndex := 1
	expectedResult := true
	obtainedIndex, obtainedResult := FindString(testSlice, "bar")

	if !reflect.DeepEqual(expectedIndex, obtainedIndex) {
		t.Errorf("Incorrect Index, expected: %v, obtained: %v", expectedIndex, obtainedIndex)
	}

	if !reflect.DeepEqual(expectedResult, obtainedResult) {
		t.Errorf("Incorrect Result, expected: %v, obtained: %v", expectedResult, obtainedResult)
	}
}

func TestFindStringWithUnexistingString(t *testing.T) {
	testSlice := []string{"foo", "bar", "test"}
	expectedIndex := -1
	expectedResult := false
	obtainedIndex, obtainedResult := FindString(testSlice, "unexisting")

	if !reflect.DeepEqual(expectedIndex, obtainedIndex) {
		t.Errorf("Incorrect Index, expected: %v, obtained: %v", expectedIndex, obtainedIndex)
	}

	if !reflect.DeepEqual(expectedResult, obtainedResult) {
		t.Errorf("Incorrect Result, expected: %v, obtained: %v", expectedResult, obtainedResult)
	}
}

func TestRemoveDuplicatedStrings(t *testing.T) {
	testSlice := []string{"foo", "bar", "test", "foo"}
	expectedResult := []string{"foo", "bar", "test"}
	sort.Strings(expectedResult) // RemoveDuplicatedStrings doesn't guaranteed the order

	obtainedResult := RemoveDuplicatedStrings(testSlice)
	sort.Strings(obtainedResult) // RemoveDuplicatedStrings doesn't guaranteed the order

	if !reflect.DeepEqual(expectedResult, obtainedResult) {
		t.Errorf("Incorrect Result, expected: %v, obtained: %v", expectedResult, obtainedResult)
	}
}
