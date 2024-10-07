package filter

import (
	"testing"
)

func TestMatchesPartial(t *testing.T) {
	tests := []struct {
		fieldValue  string
		filterValue string
		expected    bool
	}{
		{"Restrooms Available", "restrooms", true},
		{"Allow Picnic", "picnic", true},
		{"Fishing Allowed", "fishing", true},
	}

	for _, test := range tests {
		result := matchesPartial(test.fieldValue, test.filterValue)
		if result != test.expected {
			t.Errorf("matchesPartial(%q, %q) = %v; want %v", test.fieldValue, test.filterValue, result, test.expected)
		}
	}
}

func TestFilterTrailsParallel(t *testing.T) {
	trails := []Trail{
		{AccessName: "Trail A", RESTROOMS: "Yes"},
		{AccessName: "Trail B", RESTROOMS: "No"},
	}
	filters := map[string]string{"RESTROOMS": "Yes"}

	filtered := FilterTrailsParallel(trails, filters)
	if len(filtered) != 1 || filtered[0].AccessName != "Trail A" {
		t.Errorf("FilterTrailsParallel() = %v; want [Trail A]", filtered)
	}
}
