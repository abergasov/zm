package utils_test

import (
	"reflect"
	"testing"
	"zm/internal/utils"
)

type testPerson struct {
	ID   int
	Name string
}

func TestStringsFromObjectSlice(t *testing.T) {
	// Example data
	personSlice := []testPerson{
		{ID: 1, Name: "John"},
		{ID: 2, Name: "Jane"},
		{ID: 3, Name: "Doe"},
	}

	// Extractor function for Person
	extractor := func(p testPerson) string {
		return p.Name
	}

	// Expected result
	expected := []string{"John", "Jane", "Doe"}

	// Test the function
	result := utils.StringsFromObjectSlice(personSlice, extractor)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
