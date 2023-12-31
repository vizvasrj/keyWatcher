package keywatcher

import "testing"

func TestCheckKeyCombination(t *testing.T) {
	testCases := []struct {
		name          string
		currentKeys   []string
		expectedKeys  []Key
		expectedMatch bool
	}{
		{
			name:          "Matched keys",
			currentKeys:   []string{"Ctrl", "Enter"},
			expectedKeys:  []Key{{KeyString: "Ctrl"}, {KeyString: "Enter"}},
			expectedMatch: true,
		},
		{
			name:          "Mismatched keys",
			currentKeys:   []string{"Ctrl", "Shift"},
			expectedKeys:  []Key{{KeyString: "Ctrl"}, {KeyString: "Enter"}},
			expectedMatch: false,
		},
		{
			name:          "Different lengths",
			currentKeys:   []string{"Ctrl", "Enter"},
			expectedKeys:  []Key{{KeyString: "Ctrl"}},
			expectedMatch: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := checkKeyCombination(tc.currentKeys, tc.expectedKeys)
			if result != tc.expectedMatch {
				t.Errorf("Expected %v, got %v", tc.expectedMatch, result)
			}
		})
	}
}
