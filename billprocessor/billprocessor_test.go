package billprocessor

import (
	"os"
	"testing"
)

// Test-helpers

func compareStrings(result, expected, methodName string, t *testing.T) {
	if result != expected {
		t.Errorf("\nGot: %q\nExpected: %q",
			methodName, result, expected)
	}
}

// Tests

func TestValid(t *testing.T) {
	validArgs := []string{"date", "03/2015", "gas", "34.56",
		"--", "bob", "55", "susan", "45"}

	if !Valid(validArgs) {
		t.Errorf("\nArgs are valid but code says NOT: %q", validArgs)
	}

	wrongNumOfArgs := []string{"date", "03/2015", "gas", "34.56", "electric",
		"--", "bob", "55", "susan", "45"}
	if Valid(wrongNumOfArgs) {
		t.Errorf("\nWrong number of args: %q", wrongNumOfArgs)
	}

	noDateArgs := []string{"03/2015", "gas", "34.56", "electric",
		"--", "bob", "55", "susan", "45"}
	if Valid(noDateArgs) {
		t.Errorf("\nImproper date format in args: %q", noDateArgs)
	}

	noBillsArgs := []string{"date", "03/2015", "--", "bob", "55", "susan", "45"}
	if Valid(noBillsArgs) {
		t.Errorf("\nNo bills in args: %q", noBillsArgs)
	}
}

func TestMonthHeader(t *testing.T) {
	args := []string{"date", "03/2015", "gas", "34.56",
		"--", "bob", "55", "susan", "45"}
	header := monthHeader(args)
	expected := "March 2015"
	if header != expected {
		t.Errorf("\nGot: %q\nExpected: %q\nargs: %q", header, expected, args)
	}
}

func TestBillsMap(t *testing.T) {
	os.Args = []string{"gobills", "date", "03/2015", "gas", "34.56",
		"--", "bob", "55", "susan", "45"}
	args := os.Args[1:]
	b := map[string]string{}
	billsMap(args, b)
	expected := map[string]string{"gas": "34.56"}

	if b["gas"] != expected["gas"] {
		t.Errorf("\nGot: %q\nExpected: %q\nargs: %q", b, expected, args)
	}
}

func TestBillToString(t *testing.T) {
	name, amount := "Gas", "45.67"
	result := billToString(name, amount)
	expected := "Gas:\t\t$45.67\n"
	compareStrings(result, expected, "billToString(name, amount)", t)
}

func TestIndividualShares(t *testing.T) {
	total := "1000"
	args := []string{"date", "03/2015", "gas", "34.56",
		"--", "bob", "45", "susan", "55"}
	result := individualShares(total, args)
	expected := "Bob's Total:\t$450.00\nSusan's Total:\t$550.00\n"
	compareStrings(result, expected, "individualShares(total, args)", t)
}

func TestCalcShare(t *testing.T) {
	result := calcShare("55", "1000")
	expected := "550.00"
	compareStrings(result, expected, "calcShare(\"55\", \"1000\"", t)
}

func TestIncludeIn(t *testing.T) {
	args := []string{"date", "03/2015", "gas", "34.56", "--",
		"bob", "55", "Susan", "45"}
	if includeIn(args, "--") != true {
		t.Errorf("\nGot: false\nExpected: true\nargs = %q", args)
	}
}

func TestOrderKeys(t *testing.T) {
	result := orderKeys(map[string]string{"bob": "50.00", "andy": "40.00",
		"alfred": "67.00"})
	expected := []string{"alfred", "andy", "bob"}
	for i, name := range result {
		if name != expected[i] {
			t.Errorf("\nGot: %q\nExpected: %q", result, expected)
		}
	}
}
