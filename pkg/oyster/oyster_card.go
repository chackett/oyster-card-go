package oyster

import (
	"errors"
	uuid "github.com/satori/go.uuid"
)

// Card payment instrument
type Card struct {
	// ID unique identifier of card
	ID string
	// Journeys stores the journeys for each card. For simplicity this is in the card struct, but wouldn't usually
	// reside there (from a design perspective)
	Journeys []Journey
	Balance  int
}

// NewCard returns a new Oyster Card with zero balance
func NewCard() Card {
	id, _ := uuid.NewV4()
	return Card{
		ID:       id.String(),
		Balance:  0,
		Journeys: make([]Journey, 0),
	}
}

// TopUp can be used to add increase the balance of the card. Value is currency minor units (pence).
func (card *Card) TopUp(amount int) error {
	if amount < 1 {
		return errors.New("minimum top up is 1")
	}
	card.updateBalance(amount)
	return nil
}

// updateBalance is a simplistic setter for balance. Adds the provide amount to the existing balance, be it a positive
// or negative number.
// Not thread safe.
func (card *Card) updateBalance(amount int) {
	card.Balance += amount
}
