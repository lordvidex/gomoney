package gomoney

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurrency(t *testing.T) {
	testCases := []struct {
		input  Currency
		output string
	}{
		{USD, "USD"},
		{RUB, "RUB"},
		{NGN, "NGN"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.output, func(t *testing.T) {
			stringValue := string(testCase.input)
			assert.Equal(t, testCase.output, stringValue)
		})
	}
}

func TestCurrencyMarshal(t *testing.T) {
	testCases := []struct {
		input    Currency
		expected string
	}{
		{USD, `"USD"`},
		{NGN, `"NGN"`},
		{RUB, `"RUB"`},
	}
	for _, tc := range testCases {
		t.Run(string(tc.input), func(t *testing.T) {
			bytes, err := json.Marshal(tc.input)
			if err != nil {
				t.Errorf("Error marshaling currency")
				return
			}
			got := string(bytes)
			if got != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, got)
				return
			}
		})
	}
}
