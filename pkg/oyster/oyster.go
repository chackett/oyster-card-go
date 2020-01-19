package oyster

import (
	"fmt"
	"github.com/pkg/errors"
	"math"
)

const maxFare = 320
const busFixedFare = 180

// Oyster travel system to calculate fares, track cards and journeys.
type Oyster struct {
	// Cards the Cards issued to passengers
	Cards    map[string]*Card
	Stations map[string]Station
}

// NewOyster returns a new instance of Oyster for processing journeys etc
func NewOyster(cards []*Card) Oyster {
	result := Oyster{
		Cards:    make(map[string]*Card),
		Stations: make(map[string]Station),
	}

	// Setup the stations
	result.Stations[StationHolborn] = Station{Name: StationHolborn, Zones: []int{1}, Typ: typeTube}
	result.Stations[StationEarlsCourt] = Station{Name: StationEarlsCourt, Zones: []int{1, 2}, Typ: typeTube}
	result.Stations[StationWimbledon] = Station{Name: StationWimbledon, Zones: []int{3}, Typ: typeTube} // Unused
	result.Stations[StationHammersmith] = Station{Name: StationHammersmith, Zones: []int{2}, Typ: typeTube}
	result.Stations[BusStopEarlsCourt] = Station{Name: BusStopEarlsCourt, Zones: []int{}, Typ: typeBus} // Bus is a special case Station that exists separate of zones. Use -1.

	// Need to index the Cards for easy lookup later
	for _, card := range cards {
		result.Cards[card.ID] = card
	}

	return result
}

// EnterBus allow user to tap into bus
func (oyster *Oyster) EnterBus(stopID string, cardID string) error {
	stBegin, ok := oyster.Stations[stopID]
	if !ok {
		return fmt.Errorf("bus stop `%s` not found", stopID)
	}

	j := Journey{
		Start: stBegin,
		Fare:  busFixedFare,
	}

	err := oyster.addJourney(cardID, j)
	if err != nil {
		return errors.Wrap(err, "add journey")
	}
	// Deduct the fare from card
	oyster.Cards[cardID].updateBalance(j.Fare * -1)

	return nil
}

// EnterTube allow user to tap into a Station
func (oyster *Oyster) EnterTube(stationID string, cardID string) error {
	stBegin, ok := oyster.Stations[stationID]
	if !ok {
		return fmt.Errorf("Station `%s` not found", stationID)
	}
	if stBegin.Typ != typeTube {
		return errors.New("cannot enter tube at non tube Station")
	}

	j := Journey{
		Start: stBegin,
		Fare:  maxFare,
	}

	err := oyster.addJourney(cardID, j)
	if err != nil {
		return errors.Wrap(err, "add journey")
	}
	// Deduct the fare from card
	oyster.Cards[cardID].updateBalance(j.Fare * -1)

	return nil
}

// ExitTube allow user to tap out of a Station
func (oyster *Oyster) ExitTube(stationID string, cardID string) error {
	stEnd, ok := oyster.Stations[stationID]
	if !ok {
		return fmt.Errorf("Station `%s` not found", stationID)
	}

	journey, err := oyster.getLastJourney(cardID)
	if err != nil {
		return errors.Wrap(err, "get last journey")
	}
	if journey.Start.Typ != typeTube {
		return errors.New("cannot exit tube for a journey starting at non tube Station")
	}
	journey.End = stEnd

	// Refund the initial fare
	oyster.Cards[cardID].updateBalance(journey.Fare)

	f, err := oyster.calculateFare(journey)
	if err != nil {
		return errors.Wrap(err, "calculate fare")
	}
	journey.Fare = f

	err = oyster.updateLastJourney(cardID, journey)
	if err != nil {
		return errors.Wrap(err, "update last journey")
	}

	oyster.Cards[cardID].updateBalance(journey.Fare * -1)

	return nil
}

func (oyster *Oyster) calculateFare(j Journey) (int, error) {
	// 2 Special cases
	if j.Start.Typ == typeBus {
		return busFixedFare, nil
	}

	if j.End.Name == "" {
		return maxFare, nil
	}

	// Edge case: Tapped in and out of same Station. Assume no travel for simplicity.
	// In real life, timestamps should be considered.
	if j.Start.Name == j.End.Name {
		return 0, nil
	}

	zones, err := oyster.countZones(j)
	if err != nil {
		return 0, errors.Wrap(err, "count zones")
	}

	// There is probably a more elegant way to determine the fares. Consider a structure to map these and
	// define rules.

	// Anywhere in Zone 1
	if zones == 1 && j.Start.InZone(1) && j.End.InZone(1) { // Belt and braces checking
		return 250, nil
	} else if zones == 1 { // Any one zone outside zone 1
		return 200, nil
	}
	// Any two zones including Zone 1
	if zones == 2 && (j.Start.InZone(1) || j.End.InZone(1)) {
		return 300, nil
	} else if zones == 2 { // Any two zones excluding Zone 1
		return 225, nil
	}
	// Any three zones
	if zones == 3 {
		return 320, nil
	}

	return 0, nil
}

// countZones counts the number of zones visited during travel. Stations that are designated in multiple zones
// are also considered, favouring the customer. To favour the customer in multi zoned stations, compare all zones
// until the lowest count is found, because the cost is directly proportional to zone count
func (oyster *Oyster) countZones(j Journey) (int, error) {
	if len(j.Start.Zones) == 0 || len(j.End.Zones) == 0 {
		return 0, errors.New("Station zone empty")
	}

	if j.Start.MultiZone() && j.End.MultiZone() {
		return 0, errors.New("both start and end stations being multi zone is not supported")
	}

	// For simplicity, ignoring possibility of start and end Station being multi zone
	if j.Start.MultiZone() {
		var min = math.MaxInt32
		for _, startZone := range j.Start.Zones {
			diff := oyster.zoneDiff(j.End.Zones[0], startZone)
			min = int(math.Min(float64(diff), float64(min)))
		}
		return min, nil

	} else if j.End.MultiZone() {
		var min = math.MaxInt32
		for _, endZone := range j.End.Zones {
			diff := oyster.zoneDiff(j.End.Zones[0], endZone)
			min = int(math.Min(float64(diff), float64(min)))
		}
		return min, nil
	}

	// Final case, neither start or end are multi zone stations
	return oyster.zoneDiff(j.Start.Zones[0], j.End.Zones[0]), nil
}

// zoneDiff calculates the difference in two specified zones. i.e. How many zones were visited in a journey.
// Must add one because simple subtraction is exclusive. i.e. difference of zones 1-1 is zero. But results is 1
func (oyster *Oyster) zoneDiff(z1, z2 int) int {
	if z1 < 1 || z2 < 1 {
		return -1
	}
	return int(math.Abs(float64(z1-z2)) + 1)
}

func (oyster *Oyster) getLastJourney(cardID string) (Journey, error) {
	if len(oyster.Cards[cardID].Journeys) == 0 {
		return Journey{}, errors.New("no journeys")
	}
	return oyster.Cards[cardID].Journeys[len(oyster.Cards[cardID].Journeys)-1], nil

}

func (oyster *Oyster) updateLastJourney(cardID string, j Journey) error {
	if len(oyster.Cards[cardID].Journeys) == 0 {
		return errors.New("No journeys")
	}
	idxLast := len(oyster.Cards[cardID].Journeys) - 1

	oyster.Cards[cardID].Journeys[idxLast] = j

	return nil
}

func (oyster *Oyster) addJourney(cardID string, j Journey) error {
	oyster.Cards[cardID].Journeys = append(oyster.Cards[cardID].Journeys, j)
	return nil
}

// Edge case, tapping "in" twice, before tapping out of the first instance.
