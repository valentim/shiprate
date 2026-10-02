package shipping

func StandardRates() *Calculator {
	return NewCalculator(
		WeightTiers{
			Tiers: []WeightTier{
				{MaxKg: 1, Rate: EUR(5)},
				{MaxKg: 5, Rate: EUR(10)},
			},
			Above: EUR(20),
		},
		ZoneSurcharges{
			ZoneDomestic:      EUR(0),
			ZoneEU:            EUR(8),
			ZoneInternational: EUR(15),
		},
	)
}
