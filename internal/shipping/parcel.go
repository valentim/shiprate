package shipping

import (
	"fmt"
	"math"
	"slices"
)

// Zone is the destination zone of a parcel.
type Zone string

// The zones a parcel can be sent to.
const (
	ZoneDomestic      Zone = "domestic"
	ZoneEU            Zone = "eu"
	ZoneInternational Zone = "international"
)

var zones = []Zone{ZoneDomestic, ZoneEU, ZoneInternational}

// Parcel is a parcel to be priced.
type Parcel struct {
	WeightKg float64
	LengthCm float64
	WidthCm  float64
	HeightCm float64
	Zone     Zone
}

// Validate returns an error if a measurement is not a positive number or the
// zone is unknown.
func (p Parcel) Validate() error {
	measurements := []struct {
		name  string
		value float64
	}{
		{"weight", p.WeightKg},
		{"length", p.LengthCm},
		{"width", p.WidthCm},
		{"height", p.HeightCm},
	}
	for _, m := range measurements {
		if math.IsNaN(m.value) || math.IsInf(m.value, 0) || m.value <= 0 {
			return fmt.Errorf("%s must be a positive number, got %v", m.name, m.value)
		}
	}
	if !slices.Contains(zones, p.Zone) {
		return fmt.Errorf("unknown zone %q, valid zones are %v", p.Zone, zones)
	}
	return nil
}
