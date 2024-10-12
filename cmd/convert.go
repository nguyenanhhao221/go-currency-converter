/*
Copyright © 2024 Hao Nguyen <hao@haonguyen.tech>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

type ExchangeRateAPIResponse struct {
	Result             string             `json:"result"`
	Documentation      string             `json:"documentation"`
	TermsOfUse         string             `json:"terms_of_use"`
	TimeLastUpdateUnix int64              `json:"time_last_update_unix"`
	TimeLastUpdateUTC  string             `json:"time_last_update_utc"`
	TimeNextUpdateUnix int64              `json:"time_next_update_unix"`
	TimeNextUpdateUTC  string             `json:"time_next_update_utc"`
	BaseCode           string             `json:"base_code"`
	ConversionRates    map[string]float64 `json:"conversion_rates"`
}

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:     "convert",
	Short:   "Convert currency from one unit to another",
	Long:    `Convert a specified amount from one currency to another using exchange rates.`,
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Open("./testdata/mockResponse.json")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening convert json file %v", err)
			return
		}
		defer f.Close()

		exChangeRateApi := ExchangeRateAPIResponse{}
		err = json.NewDecoder(f).Decode(&exChangeRateApi)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		result, err := convertAction(exChangeRateApi, 100, "VND")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		if _, err := printConverResult(os.Stdout, "USD", "VND", result); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
}

func convertAction(exChangeRateApi ExchangeRateAPIResponse, amount float64, toCurrency string) (float64, error) {
	exChangeRate, ok := exChangeRateApi.ConversionRates[toCurrency]
	if !ok {
		return 0.0, ErrCurrencyNotFound
	}
	return exChangeRate * amount, nil
}

func printConverResult(w io.Writer, fromCurrency string, toCurrency string, result float64) (int, error) {
	return fmt.Fprintf(w, "From: %s to %s\nResult:  %.2f\n", fromCurrency, toCurrency, result)
}
