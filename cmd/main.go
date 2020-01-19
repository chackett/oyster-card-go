package main

import (
	"fmt"
	"os"

	"github.com/chackett/oyster-test/pkg/oyster"
)

func main() {
	sampleCard := oyster.NewCard()

	systemCards := []*oyster.Card{&sampleCard}

	// Create a new Oyster instance with batch of cards
	oy := oyster.NewOyster(systemCards)

	topUpAmt := 3000
	err := sampleCard.TopUp(topUpAmt)
	printErrorAndQuit(err)

	err = oy.EnterTube(oyster.StationHolborn, sampleCard.ID)
	printErrorAndQuit(err)

	err = oy.ExitTube(oyster.StationEarlsCourt, sampleCard.ID)
	printErrorAndQuit(err)

	err = oy.EnterBus(oyster.BusStopEarlsCourt, sampleCard.ID)

	err = oy.EnterTube(oyster.StationEarlsCourt, sampleCard.ID)
	printErrorAndQuit(err)

	err = oy.ExitTube(oyster.StationHammersmith, sampleCard.ID)
	printErrorAndQuit(err)

	fmt.Printf("Card balance: %s\n", oyster.FormatCurrencyMinor(sampleCard.Balance))

	fmt.Println("Travel Summary:")
	for _, journey := range oy.Cards[sampleCard.ID].Journeys {
		fmt.Println(journey.Describe())
	}
}

func printErrorAndQuit(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
