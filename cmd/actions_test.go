package cmd

import (
	"bytes"
	"errors"
	"testing"
)

func getMockExchangeRateApi() ExchangeRateAPIResponse {
	mockProvider := map[string]float64{
		"USD": 1.0,
		"VND": 24822.3006,
		"EUR": 0.9142,
	}
	return ExchangeRateAPIResponse{
		ConversionRates: mockProvider,
	}
}

func TestConvertAction(t *testing.T) {
	testCases := []struct {
		name         string
		expErr       error
		expRes       float64
		amount       float64
		toCurrency   string
		fromCurrency string
	}{
		{name: "USD to VND", amount: 100, fromCurrency: "USD", toCurrency: "VND", expErr: nil, expRes: 2482230.060000},
		{name: "EUR to VND", amount: 100, fromCurrency: "EUR", toCurrency: "VND", expErr: nil, expRes: 2715193.677532},
		{name: "Invalid currency", fromCurrency: "bar", toCurrency: "foo", expErr: ErrCurrencyNotFound},
		{name: "Invalid currency 2", fromCurrency: "bar", toCurrency: "VND", expErr: ErrCurrencyNotFound},
	}
	conversionRates := getMockExchangeRateApi().ConversionRates
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := convertAction(conversionRates, tc.fromCurrency, tc.toCurrency, tc.amount)
			if tc.expErr != nil {
				if err == nil {
					t.Error("Expect to have error, got nil")
				}
				if !errors.Is(tc.expErr, err) {
					t.Errorf("Expect error %v got %v", tc.expErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("Expect no error, got %v", err)
			}

			if int(result) != int(tc.expRes) {
				t.Errorf("Expect %f got %f", tc.expRes, result)
			}
		})
	}
}

func TestPrintCoverResult(t *testing.T) {
	conversionRates := getMockExchangeRateApi().ConversionRates
	var out bytes.Buffer
	result, err := convertAction(conversionRates, "USD", "VND", 100)
	if err != nil {
		t.Fatal(err)
	}

	_, err = printConverResult(&out, 100, "USD", "VND", result)
	if err != nil {
		t.Fatal(err)
	}

	exp := "From: 100 USD to VND\nResult:  2,482,230 VND\n"
	if exp != out.String() {
		t.Errorf("Expect: %s, got %s\n", exp, out.String())
	}
}
