package shipping

import (
	"math"
	"reflect"
	"strings"
	"testing"
)

func parcel(weightKg float64, zone Zone) Parcel {
	return Parcel{WeightKg: weightKg, LengthCm: 10, WidthCm: 10, HeightCm: 10, Zone: zone}
}

func TestStandardRatesTotal(t *testing.T) {
	tests := []struct {
		name     string
		weightKg float64
		zone     Zone
		want     Money
	}{
		{"light domestic", 0.5, ZoneDomestic, EUR(5)},
		{"light eu", 0.5, ZoneEU, EUR(13)},
		{"light international", 0.5, ZoneInternational, EUR(20)},
		{"medium domestic", 3, ZoneDomestic, EUR(10)},
		{"medium eu (example from the brief)", 3, ZoneEU, EUR(18)},
		{"medium international", 3, ZoneInternational, EUR(25)},
		{"heavy domestic", 10, ZoneDomestic, EUR(20)},
		{"heavy eu", 10, ZoneEU, EUR(28)},
		{"heavy international", 10, ZoneInternational, EUR(35)},

		{"exactly 1 kg", 1, ZoneDomestic, EUR(5)},
		{"just above 1 kg", 1.01, ZoneDomestic, EUR(10)},
		{"exactly 5 kg", 5, ZoneDomestic, EUR(10)},
		{"just above 5 kg", 5.01, ZoneDomestic, EUR(20)},
	}
	calc := StandardRates()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quote, err := calc.Quote(parcel(tt.weightKg, tt.zone))
			if err != nil {
				t.Fatalf("Quote() error = %v", err)
			}
			if quote.Total != tt.want {
				t.Errorf("Total = %v, want %v", quote.Total, tt.want)
			}
		})
	}
}

func TestStandardRatesPricesEveryZone(t *testing.T) {
	calc := StandardRates()
	for _, zone := range zones {
		if _, err := calc.Quote(parcel(1, zone)); err != nil {
			t.Errorf("Quote() for zone %q error = %v", zone, err)
		}
	}
}

func TestStandardRatesBreakdown(t *testing.T) {
	quote, err := StandardRates().Quote(parcel(3, ZoneEU))
	if err != nil {
		t.Fatalf("Quote() error = %v", err)
	}
	want := []Charge{
		{Label: "weight tier", Amount: EUR(10)},
		{Label: "zone surcharge", Amount: EUR(8)},
	}
	if !reflect.DeepEqual(quote.Charges, want) {
		t.Errorf("Charges = %v, want %v", quote.Charges, want)
	}
}

func TestQuoteRejectsInvalidParcel(t *testing.T) {
	tests := []struct {
		name    string
		parcel  Parcel
		wantErr string
	}{
		{"zero weight", parcel(0, ZoneEU), "weight"},
		{"negative weight", parcel(-1, ZoneEU), "weight"},
		{"NaN weight", parcel(math.NaN(), ZoneEU), "weight"},
		{"infinite weight", parcel(math.Inf(1), ZoneEU), "weight"},
		{"zero length", Parcel{WeightKg: 1, LengthCm: 0, WidthCm: 20, HeightCm: 10, Zone: ZoneEU}, "length"},
		{"negative width", Parcel{WeightKg: 1, LengthCm: 30, WidthCm: -20, HeightCm: 10, Zone: ZoneEU}, "width"},
		{"missing height", Parcel{WeightKg: 1, LengthCm: 30, WidthCm: 20, Zone: ZoneEU}, "height"},
		{"unknown zone", parcel(1, "mars"), `unknown zone "mars", valid zones are [domestic eu international]`},
		{"empty zone", parcel(1, ""), "unknown zone"},
		{"zone in upper case", parcel(1, "EU"), "unknown zone"},
	}
	calc := StandardRates()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := calc.Quote(tt.parcel)
			if err == nil {
				t.Fatal("Quote() error = nil, want an error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Quote() error = %q, want it to mention %q", err, tt.wantErr)
			}
		})
	}
}

func TestQuoteRejectsZoneWithoutSurcharge(t *testing.T) {
	calc := NewCalculator(ZoneSurcharges{ZoneDomestic: EUR(0)})
	_, err := calc.Quote(parcel(1, ZoneEU))
	if err == nil {
		t.Fatal("Quote() error = nil, want an error")
	}
	if want := `no surcharge for zone "eu"`; !strings.Contains(err.Error(), want) {
		t.Errorf("Quote() error = %q, want it to mention %q", err, want)
	}
}

type flatFee Money

func (f flatFee) Apply(Parcel) (Charge, error) {
	return Charge{Label: "handling", Amount: Money(f)}, nil
}

func TestCalculatorSumsEveryRule(t *testing.T) {
	calc := NewCalculator(
		WeightTiers{Tiers: []WeightTier{{MaxKg: 2, Rate: EUR(3)}}, Above: EUR(9)},
		flatFee(150),
	)
	quote, err := calc.Quote(parcel(2, ZoneDomestic))
	if err != nil {
		t.Fatalf("Quote() error = %v", err)
	}
	if want := Money(450); quote.Total != want {
		t.Errorf("Total = %v, want %v", quote.Total, want)
	}
	if len(quote.Charges) != 2 {
		t.Errorf("got %d charges, want 2", len(quote.Charges))
	}
}
