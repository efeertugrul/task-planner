package utility

import (
	"testing"
	"time"
)

func TestToPointer(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
	}{
		{
			name:  "string pointer",
			value: "test string",
		},
		{
			name:  "int pointer",
			value: 42,
		},
		{
			name:  "float pointer",
			value: 3.14,
		},
		{
			name:  "bool pointer",
			value: true,
		},
		{
			name:  "time pointer",
			value: time.Now(),
		},
		{
			name:  "empty string pointer",
			value: "",
		},
		{
			name:  "zero int pointer",
			value: 0,
		},
		{
			name:  "zero float pointer",
			value: 0.0,
		},
		{
			name:  "false bool pointer",
			value: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Create a new variable for the closure
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Test string pointer
			if str, ok := tt.value.(string); ok {
				ptr := ToPointer(str)
				if ptr == nil {
					t.Error("Expected non-nil pointer for string")
					return
				}
				if *ptr != str {
					t.Errorf("Expected %v, got %v", str, *ptr)
				}
				return
			}

			// Test int pointer
			if i, ok := tt.value.(int); ok {
				ptr := ToPointer(i)
				if ptr == nil {
					t.Error("Expected non-nil pointer for int")
					return
				}
				if *ptr != i {
					t.Errorf("Expected %v, got %v", i, *ptr)
				}
				return
			}

			// Test float pointer
			if f, ok := tt.value.(float64); ok {
				ptr := ToPointer(f)
				if ptr == nil {
					t.Error("Expected non-nil pointer for float")
					return
				}
				if *ptr != f {
					t.Errorf("Expected %v, got %v", f, *ptr)
				}
				return
			}

			// Test bool pointer
			if b, ok := tt.value.(bool); ok {
				ptr := ToPointer(b)
				if ptr == nil {
					t.Error("Expected non-nil pointer for bool")
					return
				}
				if *ptr != b {
					t.Errorf("Expected %v, got %v", b, *ptr)
				}
				return
			}

			// Test time pointer
			if tm, ok := tt.value.(time.Time); ok {
				ptr := ToPointer(tm)
				if ptr == nil {
					t.Error("Expected non-nil pointer for time")
					return
				}
				if !ptr.Equal(tm) {
					t.Errorf("Expected %v, got %v", tm, *ptr)
				}
				return
			}
		})
	}
}
