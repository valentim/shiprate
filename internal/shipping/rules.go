package shipping

import "fmt"

// Charge is one line of a quote.
type Charge struct {
	Label  string
	Amount Money
}

// Rule returns the charge for one aspect of a parcel.
type Rule interface {
	Apply(p Parcel) (Charge, error)
}

// WeightTier is the rate for parcels up to and including MaxKg.
type WeightTier struct {
	MaxKg float64
	Rate  Money
}

// WeightTiers is the rule for the base rate by weight. Tiers must be listed in
// ascending order of MaxKg. Parcels heavier than the last tier pay Above.
type WeightTiers struct {
	Tiers []WeightTier
	Above Money
}

// Apply returns the rate of the first tier the parcel fits in.
func (r WeightTiers) Apply(p Parcel) (Charge, error) {
	for _, t := range r.Tiers {
		if p.WeightKg <= t.MaxKg {
			return Charge{Label: "weight tier", Amount: t.Rate}, nil
		}
	}
	return Charge{Label: "weight tier", Amount: r.Above}, nil
}

// ZoneSurcharges is the rule for the surcharge of each zone.
type ZoneSurcharges map[Zone]Money

// Apply returns the surcharge for the zone of the parcel.
func (r ZoneSurcharges) Apply(p Parcel) (Charge, error) {
	amount, ok := r[p.Zone]
	if !ok {
		return Charge{}, fmt.Errorf("no surcharge for zone %q", p.Zone)
	}
	return Charge{Label: "zone surcharge", Amount: amount}, nil
}
