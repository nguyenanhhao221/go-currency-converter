package cmd

import (
	"errors"
	"testing"
)

func getMockExchangeRateApi() ExchangeRateAPIResponse {
	mockProvider := map[string]float64{
		"USD": 1.0,
		"VND": 24822.3006,
	}
	return ExchangeRateAPIResponse{
		ConversionRates: mockProvider,
	}
}

func TestConvertAction(t *testing.T) {
	testCases := []struct {
		name       string
		expErr     error
		expRes     float64
		amount     float64
		toCurrency string
	}{
		{name: "Valid USD to VND", amount: 100, toCurrency: "VND", expErr: nil, expRes: 2482230.060000},
		{name: "Invalid currency", toCurrency: "foo", expErr: ErrCurrencyNotFound},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := convertAction(getMockExchangeRateApi(), tc.amount, tc.toCurrency)
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

			if result != tc.expRes {
				t.Errorf("Expect %f got %f", tc.expRes, result)
			}
		})
	}
}
