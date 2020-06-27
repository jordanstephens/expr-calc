package main

import "testing"

func TestMul(t *testing.T) {
	result, err := run([]string{"1", "*", "3"})
	if err != nil {
		t.Errorf("Err: %w", err)
	}
	if result != 3 {
		t.Errorf("1 * 3 = %f, want 3", result)
	}
}

func TestPrecedence(t *testing.T) {
	result, err := run([]string{"1", "+", "2", "*", "3"})
	if err != nil {
		t.Errorf("Err: %w", err)
	}
	if result != 7 {
		t.Errorf("1 + 2 * 3 = %f, want 7", result)
	}
}
