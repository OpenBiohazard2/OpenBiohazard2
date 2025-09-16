package render

import (
	"math"
	"testing"
)

func TestConvertToScreenX(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float32
	}{
		{
			name:     "Left_edge",
			input:    0.0,
			expected: -1.0,
		},
		{
			name:     "Center",
			input:    160.0, // 320/2
			expected: 0.0,
		},
		{
			name:     "Right_edge",
			input:    320.0,
			expected: 1.0,
		},
		{
			name:     "Quarter_left",
			input:    80.0, // 320/4
			expected: -0.5,
		},
		{
			name:     "Three_quarters_right",
			input:    240.0, // 320*3/4
			expected: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToScreenX(tt.input)
			if math.Abs(float64(result-tt.expected)) > 1e-6 {
				t.Errorf("ConvertToScreenX(%f) = %f, expected %f", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvertToScreenY(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float32
	}{
		{
			name:     "Top_edge",
			input:    0.0,
			expected: 1.0, // Note: Y is flipped
		},
		{
			name:     "Center",
			input:    120.0, // 240/2
			expected: 0.0,
		},
		{
			name:     "Bottom_edge",
			input:    240.0,
			expected: -1.0, // Note: Y is flipped
		},
		{
			name:     "Quarter_top",
			input:    60.0, // 240/4
			expected: 0.5,  // Note: Y is flipped
		},
		{
			name:     "Three_quarters_bottom",
			input:    180.0, // 240*3/4
			expected: -0.5,  // Note: Y is flipped
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToScreenY(tt.input)
			if math.Abs(float64(result-tt.expected)) > 1e-6 {
				t.Errorf("ConvertToScreenY(%f) = %f, expected %f", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvertToTextureU(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float32
	}{
		{
			name:     "Left_edge",
			input:    0.0,
			expected: 0.0,
		},
		{
			name:     "Center",
			input:    160.0, // 320/2
			expected: 0.5,
		},
		{
			name:     "Right_edge",
			input:    320.0,
			expected: 1.0,
		},
		{
			name:     "Quarter_left",
			input:    80.0, // 320/4
			expected: 0.25,
		},
		{
			name:     "Three_quarters_right",
			input:    240.0, // 320*3/4
			expected: 0.75,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToTextureU(tt.input)
			if math.Abs(float64(result-tt.expected)) > 1e-6 {
				t.Errorf("ConvertToTextureU(%f) = %f, expected %f", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvertToTextureV(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float32
	}{
		{
			name:     "Top_edge",
			input:    0.0,
			expected: 0.0,
		},
		{
			name:     "Center",
			input:    120.0, // 240/2
			expected: 0.5,
		},
		{
			name:     "Bottom_edge",
			input:    240.0,
			expected: 1.0,
		},
		{
			name:     "Quarter_top",
			input:    60.0, // 240/4
			expected: 0.25,
		},
		{
			name:     "Three_quarters_bottom",
			input:    180.0, // 240*3/4
			expected: 0.75,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToTextureV(tt.input)
			if math.Abs(float64(result-tt.expected)) > 1e-6 {
				t.Errorf("ConvertToTextureV(%f) = %f, expected %f", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvertToScreenX_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float32
	}{
		{
			name:     "Negative_input",
			input:    -10.0,
			expected: -1.0625, // 2.0*(-10/320) - 1.0
		},
		{
			name:     "Beyond_right_edge",
			input:    400.0,
			expected: 1.5, // 2.0*(400/320) - 1.0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToScreenX(tt.input)
			if math.Abs(float64(result-tt.expected)) > 1e-6 {
				t.Errorf("ConvertToScreenX(%f) = %f, expected %f", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvertToScreenY_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float32
	}{
		{
			name:     "Negative_input",
			input:    -10.0,
			expected: 1.083333, // -1.0 * (2.0*(-10/240) - 1.0)
		},
		{
			name:     "Beyond_bottom_edge",
			input:    300.0,
			expected: -1.5, // -1.0 * (2.0*(300/240) - 1.0)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToScreenY(tt.input)
			if math.Abs(float64(result-tt.expected)) > 1e-6 {
				t.Errorf("ConvertToScreenY(%f) = %f, expected %f", tt.input, result, tt.expected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkConvertToScreenX(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ConvertToScreenX(160.0)
	}
}

func BenchmarkConvertToScreenY(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ConvertToScreenY(120.0)
	}
}

func BenchmarkConvertToTextureU(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ConvertToTextureU(160.0)
	}
}

func BenchmarkConvertToTextureV(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ConvertToTextureV(120.0)
	}
}
