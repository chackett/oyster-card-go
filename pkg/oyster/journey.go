package oyster

import "fmt"

// Journey represents a journey taken by a passenger, and the applied fare.
type Journey struct {
	// Start of journey
	Start Station
	// End of journey
	End Station
	// Fare applied for journey - currency minor units (pence)
	Fare int
}

// Describe returns a string with a simple description of the fare, useful for statements.
func (j Journey) Describe() string {
	fareHuman := FormatCurrencyMinor(j.Fare)
	if j.Start.Typ == typeBus {
		return fmt.Sprintf("(%s) %s", fareHuman, j.Start.Name)
	}
	return fmt.Sprintf("(%s) %s -> %s", fareHuman, j.Start.Name, j.End.Name)
}
