package common

import "testing"

func TestConvertFlipToRotation(t *testing.T) {
	tests := []struct {
		name     string
		input    uint32
		expected float32
	}{
		{
			name:     "no flipping should be 0º rotation",
			input:    0,
			expected: 0,
		},
		{
			name:     "flipping diagonally should be 90º rotation",
			input:    2684354560,
			expected: 90,
		},
		{
			name:     "flipping horizontally and vertically should be 180º rotation",
			input:    3221225472,
			expected: 180,
		},
		{
			name:     "flipping diagonally, horizontally and vertically should be 270º rotation",
			input:    1610612736,
			expected: 270,
		},
	}

	for _, test := range tests {
		actual := convertFlipToRotation(test.input)

		if actual != test.expected {
			t.Errorf("%s expected=%f ; got=%f", test.name, test.expected, actual)
		}
	}
}
