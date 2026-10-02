package shipping

// Quote is the result for a parcel: the charges and their total.
type Quote struct {
	Charges []Charge
	Total   Money
}

// Calculator prices parcels with a list of rules.
type Calculator struct {
	rules []Rule
}

// NewCalculator returns a Calculator that applies the given rules in order.
func NewCalculator(rules ...Rule) *Calculator {
	return &Calculator{rules: rules}
}

// Quote validates the parcel and returns the charge of each rule and the total.
func (c *Calculator) Quote(p Parcel) (Quote, error) {
	if err := p.Validate(); err != nil {
		return Quote{}, err
	}
	var q Quote
	for _, r := range c.rules {
		charge, err := r.Apply(p)
		if err != nil {
			return Quote{}, err
		}
		q.Charges = append(q.Charges, charge)
		q.Total += charge.Amount
	}
	return q, nil
}
