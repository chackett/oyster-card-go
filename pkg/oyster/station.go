package oyster

const BusStopEarlsCourt = "Earls Court Bus Stop"
const StationHolborn = "Holborn"
const StationEarlsCourt = "Earl's Court"
const StationWimbledon = "Wimbledon"
const StationHammersmith = "Hammersmith"

type typ int

const (
	typeBus  typ = iota
	typeTube typ = iota
)

type Station struct {
	Name  string
	Zones []int
	Typ   typ
}

// InZone returns true if the specified zone is assigned to the Station, false if not.
func (s *Station) InZone(zone int) bool {
	for _, z := range s.Zones {
		if z == zone {
			return true
		}
	}
	return false
}

// MultiZone returns true if the Station has more than one assigned to it and false for just one zone.
func (s *Station) MultiZone() bool {
	return len(s.Zones) > 1
}
