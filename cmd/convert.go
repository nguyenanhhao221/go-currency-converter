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
	"strconv"

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
	Use:     "convert [amount]",
	Short:   "Convert currency from one unit to another",
	Long:    `Convert a specified amount from one currency to another using exchange rates.`,
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		fromCurrency, err := cmd.Flags().GetString("form")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}

		toCurrency, err := cmd.Flags().GetString("to")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		// MAYBE Need to validate error when accessing with index, however cobra.ExactArgs also ensure that the args will has at least 1
		amountString := args[0]
		amount, err := strconv.ParseFloat(amountString, 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

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

		result, err := convertAction(exChangeRateApi.ConversionRates, fromCurrency, toCurrency, amount)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		if _, err := printConverResult(os.Stdout, amount, "USD", toCurrency, result); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().StringP("form", "f", "USD", "From currency")
	convertCmd.Flags().StringP("to", "t", "VND", "To currency")
}

func convertAction(conversionRates map[string]float64, fromCurrency string, toCurrency string, amount float64) (float64, error) {
	// Get the exchange rates
	fromRate, fromOk := conversionRates[fromCurrency]
	toRate, toOk := conversionRates[toCurrency]
	if !fromOk || !toOk {
		fmt.Printf("Currency code not found. From: %v, To: %v\n", fromOk, toOk)
		return 0.0, ErrCurrencyNotFound
	}

	// Perform the conversion
	amountInUSD := amount / fromRate
	convertedAmount := amountInUSD * toRate
	return convertedAmount, nil
}

func printConverResult(w io.Writer, amount float64, fromCurrency string, toCurrency string, result float64) (int, error) {
	return fmt.Fprintf(w, "From: %2.f %s to %s\nResult:  %.2f %s\n", amount, fromCurrency, toCurrency, result, toCurrency)
}
