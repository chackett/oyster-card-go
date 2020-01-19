package oyster

import "fmt"

// FormatCurrencyMinor formats a minor amount to a string suitable for printing to users
func FormatCurrencyMinor(minor int) string {
	fBal := float64(minor) / 100.0
	return fmt.Sprintf("£%.2f", fBal)
}
