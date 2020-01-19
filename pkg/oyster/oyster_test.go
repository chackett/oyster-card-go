package oyster

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOyster_CountZones(t *testing.T) {
	testTable := []struct {
		journey       Journey
		expectedZones int
		err           error
	}{
		{
			expectedZones: 1,
			journey: Journey{
				Start: Station{Zones: []int{1}},
				End:   Station{Zones: []int{1}},
			},
		},
		{
			expectedZones: 2,
			journey: Journey{
				Start: Station{Zones: []int{2}},
				End:   Station{Zones: []int{1}},
			},
		},
		{
			expectedZones: 3,
			journey: Journey{
				Start: Station{Zones: []int{1}},
				End:   Station{Zones: []int{3}},
			},
		},
		{
			expectedZones: 1,
			journey: Journey{
				Start: Station{Zones: []int{3}},
				End:   Station{Zones: []int{3}},
			},
		},
		{
			err: errors.New("both start and end stations being multi zone is not supported"),
			journey: Journey{
				Start: Station{Zones: []int{3, 2}},
				End:   Station{Zones: []int{4, 5}},
			},
		},
	}

	oyster := Oyster{}

	for _, testItem := range testTable {
		zones, err := oyster.countZones(testItem.journey)
		if testItem.err != nil {
			assert.Error(t, err)
			assert.EqualError(t, err, testItem.err.Error())
			continue
		}
		assert.Equal(t, testItem.expectedZones, zones)
	}
}

func TestOyster_zoneDiff(t *testing.T) {
	testTable := []struct {
		z1       int
		z2       int
		expected int
	}{
		{
			z1:       1,
			z2:       2,
			expected: 2,
		},
		{
			z1:       2,
			z2:       2,
			expected: 1,
		},
		{
			z1:       0,
			z2:       0,
			expected: -1,
		},
		{
			z1:       -1,
			z2:       2,
			expected: -1,
		},
		{
			z1:       2,
			z2:       -1,
			expected: -1,
		},
	}

	oyster := Oyster{}

	for i, testItem := range testTable {
		label := fmt.Sprintf("Test: %d", i)
		t.Run(label, func(t *testing.T) {
			actual := oyster.zoneDiff(testItem.z1, testItem.z2)
			assert.EqualValues(t, testItem.expected, actual)
		})
	}
}

func TestOyster_countZones(t *testing.T) {
	testTable := []struct {
		journey  Journey
		expected int
		err      error
	}{
		{
			journey: Journey{
				Start: Station{Zones: []int{1}},
				End:   Station{Zones: []int{1}},
			},
			expected: 1,
		},
		{
			journey: Journey{
				Start: Station{Zones: []int{1, 2}},
				End:   Station{Zones: []int{1, 2}},
			},
			err: errors.New("both start and end stations being multi zone is not supported"),
		},
		{
			journey: Journey{
				Start: Station{Zones: []int{1, 2, 3}},
				End:   Station{Zones: []int{4}},
			},
			expected: 2,
		},
		{
			journey: Journey{
				Start: Station{Zones: []int{5}},
				End:   Station{Zones: []int{4, 5, 6}},
			},
			expected: 1,
		},
		{
			journey: Journey{
				Start: Station{Zones: []int{1, 2, 3}},
				End:   Station{Zones: []int{1}},
			},
			expected: 1,
		},
		{
			journey: Journey{
				Start: Station{Zones: []int{0}},
				End:   Station{Zones: []int{0}},
			},
			expected: -1,
		},
		{
			journey: Journey{
				Start: Station{},
				End:   Station{Zones: []int{1}},
			},
			err: errors.New("Station zone empty"),
		},
		{
			journey: Journey{
				Start: Station{Zones: nil},
				End:   Station{Zones: []int{1}},
			},
			err: errors.New("Station zone empty"),
		},
	}

	oyster := Oyster{}

	for i, testItem := range testTable {
		label := fmt.Sprintf("Test: %d", i)
		t.Run(label, func(t *testing.T) {
			actual, err := oyster.countZones(testItem.journey)
			if testItem.err != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, testItem.err.Error())
			}
			assert.Equal(t, testItem.expected, actual)
		})
	}
}

func TestOyster_EnterBus(t *testing.T) {
	testTable := []struct {
		Card            Card
		StopID          string
		ExpectedBalance int
		ExpectedJourney Journey
		Error           error
	}{
		{
			Card: Card{
				ID:       "card-id",
				Balance:  1000,
				Journeys: []Journey{},
			},
			StopID: BusStopEarlsCourt,
			ExpectedJourney: Journey{
				Start: Station{
					Typ:  typeBus,
					Name: BusStopEarlsCourt,
				},
				Fare: busFixedFare,
			},
			ExpectedBalance: 1000 - busFixedFare,
		},
	}

	for i, ti := range testTable {
		label := fmt.Sprintf("Test: %d", i)
		t.Run(label, func(t *testing.T) {
			oyster := NewOyster([]*Card{&ti.Card})

			err := oyster.EnterBus(ti.StopID, ti.Card.ID)
			if ti.Error != nil {
				assert.Error(t, err)
				assert.Equal(t, ti.Error.Error(), err.Error())
			}

			assert.Equal(t, ti.ExpectedBalance, ti.Card.Balance)
		})
	}
}

func TestOyster_EnterExitTube(t *testing.T) {
	testTable := []struct {
		Card            Card
		StartStationID  string
		EndStationID    string
		ExpectedBalance int
		ExpectedJourney Journey
		Error           error
	}{
		{
			Card: Card{
				ID:       "card-id",
				Balance:  1000,
				Journeys: []Journey{},
			},
			StartStationID: StationEarlsCourt,
			EndStationID:   StationHolborn,
			ExpectedJourney: Journey{
				Start: Station{
					Typ:   typeTube,
					Name:  StationEarlsCourt,
					Zones: []int{1, 2},
				},
				End: Station{
					Typ:   typeTube,
					Name:  StationHolborn,
					Zones: []int{1},
				},

				Fare: 250,
			},
			ExpectedBalance: 750,
		},
	}

	for i, ti := range testTable {
		label := fmt.Sprintf("Test: %d", i)
		t.Run(label, func(t *testing.T) {
			oyster := NewOyster([]*Card{&ti.Card})

			err := oyster.EnterTube(ti.StartStationID, ti.Card.ID)
			if ti.Error != nil {
				assert.Error(t, err)
				assert.Equal(t, ti.Error.Error(), err.Error())
			}
			err = oyster.ExitTube(ti.EndStationID, ti.Card.ID)
			if ti.Error != nil {
				assert.Error(t, err)
				assert.Equal(t, ti.Error.Error(), err.Error())
			}

			assert.Equal(t, ti.ExpectedBalance, ti.Card.Balance)
			assert.Equal(t, ti.ExpectedJourney, ti.Card.Journeys[0])
		})
	}
}

//
//func TestOyster_(t *testing.T) {
//	testTable := []struct {
//	}{
//		{
//
//		},
//	}
//
//	oyster := Oyster{}
//
//	for _, testItem := range testTable {
//
//	}
//}
//
//func TestOyster_(t *testing.T) {
//	testTable := []struct {
//	}{
//		{
//
//		},
//	}
//
//	oyster := Oyster{}
//
//	for _, testItem := range testTable {
//
//	}
//}
//
