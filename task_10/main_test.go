package main

import (
	"testing"
)

func checkPasswordLength(t *testing.T, length int, good bool) {
	pass, err := GeneratePassword(length, true, false, false, false)
	if good {
		if len(pass) != length {
			t.Errorf("Expected password length %d, got %d", length, len(pass))
		}
	} else {
		if err == nil {
			t.Errorf("Expected error, got %s", pass)
		}
	}
}

func TestPasswordLength(t *testing.T) {
	checkPasswordLength(t, 0, true)
	checkPasswordLength(t, -2, false)
	checkPasswordLength(t, 10, true)
}

func TestOnlyNumbers(t *testing.T) {
	pass, err := GeneratePassword(1000, true, false, false, false)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	for _, c := range pass {
		if c < 48 || c > 57 {
			t.Errorf("Expected only numbers, got %c", c)
		}
	}
}

func TestOnlySymbols(t *testing.T) {
	pass, err := GeneratePassword(1000, false, true, false, false)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	for _, c := range pass {
		if !((c >= 33 && c <= 47) || (c >= 58 && c <= 64) || (c >= 91 && c <= 96) || (c >= 123 && c <= 126)) {
			t.Errorf("Expected only symbols, got %c", c)
		}
	}
}

func TestOnlyUpperCases(t *testing.T) {
	pass, err := GeneratePassword(1000, false, false, true, false)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	for _, c := range pass {
		if c < 65 || c > 90 {
			t.Errorf("Expected only upper cases, got %c", c)
		}
	}
}

func TestOnlyLowerCases(t *testing.T) {
	pass, err := GeneratePassword(1000, false, false, false, true)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	for _, c := range pass {
		if c < 97 || c > 122 {
			t.Errorf("Expected only lower cases, got %c", c)
		}
	}
}
