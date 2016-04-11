package billprocessor

import (
	"github.com/luketower/roomies/field"
	"testing"
)

// Test Helpers

func compareStrings(result, expected, methodName string, t *testing.T) {
	if result != expected {
		t.Errorf("\nGot: %q\nExpected: %q",
			methodName, result, expected)
	}
}

// Tests

func TestHasValid(t *testing.T) {
	validArgs := []string{"date", "03/2015", "gas", "34.56",
		"--", "bob", "55", "susan", "45"}
	wrongNumOfArgs := []string{"date", "03/2015", "gas", "34.56", "electric",
		"--", "bob", "55", "susan", "45"}
	noDateArgs := []string{"03/2015", "gas", "34.56",
		"--", "bob", "55", "susan", "45"}
	noBillsArgs := []string{"date", "03/2015", "--", "bob", "55", "susan", "45"}

	if !HasValid(validArgs) {
		t.Errorf("\nArgs are valid but code says NOT: %q", validArgs)
	}

	if HasValid(wrongNumOfArgs) {
		t.Errorf("\nWrong number of args: %q", wrongNumOfArgs)
	}

	if HasValid(noDateArgs) {
		t.Errorf("\nImproper date format in args: %q", noDateArgs)
	}

	if HasValid(noBillsArgs) {
		t.Errorf("\nNo bills in args: %q", noBillsArgs)
	}
}

func TestMakeHeader(t *testing.T) {
	dateArgs := []string{"date", "03/2015", "gas", "34.56",
		"--", "bob", "55", "susan", "45"}
	monthArgs := []string{"month", "03/2015", "gas", "34.56",
		"--", "bob", "55", "susan", "45"}
	customArgs := []string{"header", "03/2015", "gas", "34.56",
		"--", "bob", "55", "susan", "45"}

	expectedDate := "March 2015"
	expectedMonth := "March 2015"
	expectedCustom := "03/2015"

	if dateHeader := makeHeader("03/2015", dateArgs, 1); dateHeader != expectedDate {
		t.Errorf("\nGot: %q\nExpected: %q\nargs: %q", dateHeader,
			expectedDate, dateArgs)
	}
	if monthHeader := makeHeader("03/2015", monthArgs, 1); monthHeader != expectedMonth {
		t.Errorf("\nGot: %q\nExpected: %q\nargs: %q", monthHeader,
			expectedMonth, monthArgs)
	}
	if customHeader := makeHeader("03/2015", customArgs, 1); customHeader != expectedCustom {
		t.Errorf("\nGot: %q\nExpected: %q\nargs: %q", customHeader,
			expectedCustom, customArgs)
	}
}

func TestToString(t *testing.T) {
	f := field.Field{"Gas", 45.67, false}
	result := f.ToString(0)
	expected := "Gas: $45.67\n"
	compareStrings(result, expected, "billToString(name, amount, 0)", t)
}

func TestCalcShare(t *testing.T) {
	result := calcShare("55", 1000)
	expected := 550.00
	if result != expected {
		t.Error("\nGot: %q\nExpected: %q", result, expected)
	}
}
