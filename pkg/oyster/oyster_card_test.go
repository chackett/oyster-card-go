package oyster

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOysterCard_updateBalance(t *testing.T) {
	testTable := []struct {
		card            *Card
		updateBalance   int
		expectedBalance int
	}{
		{
			card:            &Card{Balance: 100},
			updateBalance:   100,
			expectedBalance: 200,
		},
		{
			card:            &Card{Balance: 0},
			updateBalance:   0,
			expectedBalance: 0,
		},
		{
			card:            &Card{Balance: 100},
			updateBalance:   -100,
			expectedBalance: 0,
		},
		{
			card:            &Card{Balance: 0},
			updateBalance:   -100,
			expectedBalance: -100,
		},
	}

	for _, testItem := range testTable {
		testItem.card.updateBalance(testItem.updateBalance)
		assert.Equal(t, testItem.expectedBalance, testItem.card.Balance)
	}
}
