package oyster

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStation_MultiZone(t *testing.T) {
	multi := Station{Zones: []int{1, 2}}
	notMulti := Station{Zones: []int{1}}

	assert.True(t, multi.MultiZone())
	assert.False(t, notMulti.MultiZone())
}

func TestStation_InZone(t *testing.T) {
	testTable := []struct {
		testZone int
		inZone   bool
		st       Station
	}{
		{
			st:       Station{Zones: []int{1, 2}},
			testZone: 1,
			inZone:   true,
		},
		{
			st:       Station{Zones: []int{1, 2}},
			testZone: 5,
			inZone:   false,
		},
		{
			st:       Station{Zones: []int{1, 2}},
			testZone: 0,
			inZone:   false,
		},
		{
			st:       Station{Zones: []int{1, 2}},
			testZone: -1,
			inZone:   false,
		},
		{
			st:       Station{Zones: []int{}},
			testZone: 1,
			inZone:   false,
		},
		{
			testZone: 1,
			inZone:   false,
		},
	}

	for _, testItem := range testTable {
		actual := testItem.st.InZone(testItem.testZone)
		assert.Equal(t, testItem.inZone, actual)
	}
}
