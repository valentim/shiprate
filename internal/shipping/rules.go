package shipping

import "fmt"

type Charge struct {
	Label  string
	Amount Money
}

type Rule interface {
	Apply(p Parcel) (Charge, error)
}

type WeightTier struct {
	MaxKg float64
	Rate  Money
}

type WeightTiers struct {
	Tiers []WeightTier
	Above Money
}

func (r WeightTiers) Apply(p Parcel) (Charge, error) {
	for _, t := range r.Tiers {
		if p.WeightKg <= t.MaxKg {
			return Charge{Label: "weight tier", Amount: t.Rate}, nil
		}
	}
	return Charge{Label: "weight tier", Amount: r.Above}, nil
}

type ZoneSurcharges map[Zone]Money

func (r ZoneSurcharges) Apply(p Parcel) (Charge, error) {
	amount, ok := r[p.Zone]
	if !ok {
		return Charge{}, fmt.Errorf("no surcharge for zone %q", p.Zone)
	}
	return Charge{Label: "zone surcharge", Amount: amount}, nil
}
