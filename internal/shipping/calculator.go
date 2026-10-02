package shipping

type Quote struct {
	Charges []Charge
	Total   Money
}

type Calculator struct {
	rules []Rule
}

func NewCalculator(rules ...Rule) *Calculator {
	return &Calculator{rules: rules}
}

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
